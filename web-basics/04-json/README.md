# 04 – JSON encode/decode

JSON là lingua franca của web API. `encoding/json` của Go cực mạnh và hoàn toàn standard.

## Chạy

```bash
go run ./web-basics/04-json
```

Test:
```bash
curl http://localhost:8080/user
curl http://localhost:8080/users

curl -X POST http://localhost:8080/user \
  -H 'Content-Type: application/json' \
  -d '{"id":42,"name":"Alice","email":"a@x.com"}'
```

## Struct tags

```go
type User struct {
    ID       int    `json:"id"`              // rename field
    Name     string `json:"name"`
    Email    string `json:"email,omitempty"` // bỏ field nếu zero value
    password string                          // ❌ không export → không có trong JSON
    Secret   string `json:"-"`               // export nhưng vẫn ẩn
}
```

### Các option phổ biến

| Tag | Ý nghĩa |
|-----|---------|
| `json:"name"` | Đặt tên field trong JSON |
| `json:"name,omitempty"` | Bỏ field nếu zero value |
| `json:"-"` | KHÔNG include trong JSON (cho secret) |
| `json:",string"` | Encode int/bool thành string |
| `json:"-,"` | Tên field thật là `-` (hiếm dùng) |

### Zero value để omitempty hoạt động

- `int = 0`, `string = ""`, `bool = false`, `slice/map = nil`, pointer = nil

Muốn phân biệt "0 thật sự" vs "không gửi" → dùng pointer:

```go
Score *int `json:"score,omitempty"`
// nil = không gửi, &0 = gửi 0
```

## Encode (Go → JSON)

### Stream qua Writer (tốt cho HTTP response)

```go
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(user)  // ghi thẳng vào ResponseWriter
```

→ Không cần allocate string trung gian.

### Marshal sang `[]byte`

```go
b, err := json.Marshal(user)            // compact: {"id":1,"name":"A"}
b, err := json.MarshalIndent(u, "", "  ") // pretty:
// {
//   "id": 1,
//   "name": "A"
// }
```

## Decode (JSON → Go)

### Stream từ Reader

```go
var u User
err := json.NewDecoder(r.Body).Decode(&u)
```

### Unmarshal từ `[]byte`

```go
err := json.Unmarshal(data, &u)
```

### Strict mode

```go
dec := json.NewDecoder(r.Body)
dec.DisallowUnknownFields() // fail nếu JSON có field không có trong struct
```

Hữu ích để bắt typo (vd client gửi `"emial"` thay vì `"email"`).

## Dynamic JSON: `any` / `map[string]any`

Khi schema không cố định:

```go
var data any
json.NewDecoder(r.Body).Decode(&data)
// data có thể là map, slice, string, float64, bool, nil
```

Hoặc cụ thể hơn:

```go
var m map[string]any
json.NewDecoder(r.Body).Decode(&m)
m["count"]   // any — phải type assert
```

⚠️ Số trong JSON parse thành `float64` mặc định — kể cả integer. Cần convert:
```go
n := int(m["count"].(float64))
```

Hoặc dùng `json.Number`:
```go
dec.UseNumber()
n, _ := m["count"].(json.Number).Int64()
```

## Custom Marshal/Unmarshal

Implement interface khi cần format đặc biệt:

```go
type Date struct{ time.Time }

func (d Date) MarshalJSON() ([]byte, error) {
    return []byte(`"` + d.Format("2006-01-02") + `"`), nil
}

func (d *Date) UnmarshalJSON(b []byte) error {
    t, err := time.Parse(`"2006-01-02"`, string(b))
    d.Time = t
    return err
}
```

## Pitfall

### 1. Field không export → không Marshal

```go
type User struct {
    name string  // ❌ JSON output: {}
    Name string  // ✅ JSON output: {"Name":"..."}
}
```

### 2. Pass value thay vì pointer khi Decode

```go
json.Unmarshal(data, u)   // ❌ u là copy, không update
json.Unmarshal(data, &u)  // ✅
```

### 3. Thiếu Content-Type

Browser/client có thể tự đoán content type từ body. Set tường minh:
```go
w.Header().Set("Content-Type", "application/json")
```

### 4. Encoder add newline

`json.NewEncoder(w).Encode(v)` thêm `\n` cuối. Không phải bug — đây là behavior đã document.

### 5. Time format

Time mặc định JSON là RFC 3339:
```json
{"created_at":"2025-01-15T10:30:00Z"}
```

Muốn format khác → custom MarshalJSON.

### 6. Big number precision loss

JSON number = `float64` → mất precision với `int64` > 2^53. Dùng `json.Number` hoặc encode thành string.

## So với JSON ở các ngôn ngữ khác

| | Go | JavaScript | Python |
|---|----|------------|--------|
| Encode | `json.Marshal` | `JSON.stringify` | `json.dumps` |
| Decode | `json.Unmarshal(&v)` | `JSON.parse` | `json.loads` |
| Field renaming | Struct tag | (manual) | (manual) |
| Type checking | Compile-time với struct | Runtime | Runtime |

Go's compile-time check là điểm mạnh — bug shape JSON detect được sớm.

## Bài tập

- Custom `Marshaler` cho type `Money` in dạng `"$10.00"`.
- Endpoint accept dynamic JSON (`map[string]any`), in summary các field.
- Endpoint nhận date string `"2025-01-15"` thành `time.Time`.
- Test với `DisallowUnknownFields` — gửi field thừa và check error.
