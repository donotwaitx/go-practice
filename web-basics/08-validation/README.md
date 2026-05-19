# 08 – Validation & Structured Error Responses

Input validation và return error JSON có cấu trúc — pattern then chốt của REST API chất lượng.

## Chạy

```bash
go run ./web-basics/08-validation
```

Test:
```bash
# Hợp lệ
curl -X POST localhost:8080/users \
  -H 'Content-Type: application/json' \
  -d '{"name":"Alice","email":"a@x.com","password":"secret123","age":30}'

# Lỗi validation (nhiều field cùng lúc)
curl -X POST localhost:8080/users \
  -H 'Content-Type: application/json' \
  -d '{"name":"","email":"bad","password":"123","age":-1}'
```

Response lỗi:

```json
{
  "code": "validation_failed",
  "message": "request validation failed",
  "fields": [
    {"field":"name","message":"required"},
    {"field":"email","message":"invalid email"},
    {"field":"password","message":"min 8 chars"},
    {"field":"age","message":"must be 0..150"}
  ]
}
```

## Tại sao validate?

- **Security** — input untrusted có thể là attack (injection, XSS, oversized payload).
- **UX** — trả tất cả lỗi cùng lúc tốt hơn từng cái một.
- **Data integrity** — sai data sẽ corrupt DB.

## Pattern: ValidationErrors

```go
type FieldError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

type ValidationErrors []FieldError

func (v ValidationErrors) Error() string { ... }
```

`ValidationErrors` **implement `error`** → trả về từ function bình thường:

```go
func (r CreateUserReq) Validate() error {
    var errs ValidationErrors
    if r.Name == "" {
        errs = append(errs, FieldError{"name", "required"})
    }
    // ... các check khác
    if len(errs) > 0 { return errs }
    return nil
}
```

**Quan trọng:** thu thập TẤT CẢ lỗi rồi mới return, không return sớm. → Client thấy hết các vấn đề trong 1 lần.

## APIError — chuẩn hóa error response

```go
type APIError struct {
    Status  int    `json:"-"`        // không xuất hiện trong JSON
    Code    string `json:"code"`     // máy đọc được: "user_not_found"
    Message string `json:"message"`  // người đọc được
    Fields  any    `json:"fields,omitempty"`
}
```

- **Code** — stable identifier, FE biết cách handle (vd `user_not_found` → redirect login).
- **Message** — human-readable, có thể i18n.
- **Fields** — chi tiết field-level lỗi (cho form).

## Centralized error writer

```go
func writeErr(w http.ResponseWriter, err error) {
    var apiErr *APIError
    if errors.As(err, &apiErr) { writeJSON(w, apiErr.Status, apiErr); return }
    
    var vErr ValidationErrors
    if errors.As(err, &vErr) {
        writeJSON(w, 400, &APIError{Code: "validation_failed", Fields: vErr})
        return
    }
    
    // fallback
    writeJSON(w, 500, &APIError{Code: "internal_error", Message: err.Error()})
}
```

Handler chỉ cần `return err`, không lo response format:

```go
func createUser(w, r) {
    var req CreateUserReq
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeErr(w, &APIError{Status: 400, Code: "invalid_json", Message: err.Error()})
        return
    }
    if err := req.Validate(); err != nil { writeErr(w, err); return }
    // ...
}
```

## Validation rules phổ biến

| Rule | Cách |
|------|------|
| Required | `s == ""` (sau `strings.TrimSpace`) |
| Min/Max length | `utf8.RuneCountInString(s) <= max` (đếm rune, không byte) |
| Email | `mail.ParseAddress` (đủ tốt) hoặc regex |
| URL | `url.Parse` + check scheme |
| Range | `n >= min && n <= max` |
| Regex | `regexp.MustCompile(...).MatchString` |
| Whitelist | check trong `[]string{"a","b","c"}` |
| Format date | `time.Parse` |
| Unique (DB) | query trong handler (race condition cần handle ở DB level) |

## Vì sao KHÔNG dùng library?

`go-playground/validator` rất phổ biến với struct tags:

```go
type User struct {
    Name  string `validate:"required,max=50"`
    Email string `validate:"required,email"`
    Age   int    `validate:"gte=0,lte=150"`
}
validator.Struct(u)
```

**Pro:** ngắn gọn.
**Con:** dùng reflection (chậm), error message tệ, khó customize.

Bài này dùng tay vì:
1. Học pattern rõ ràng hơn.
2. Validate phức tạp (cross-field, async DB check) cần code thật.

Production: hybrid — library cho simple rule, custom cho complex.

## Pitfall

### 1. Return sớm — mất các lỗi khác

```go
// ❌ Client phải submit form 4 lần để biết hết lỗi
if r.Name == "" { return errors.New("name required") }
if r.Email == "" { return errors.New("email required") }
// ...

// ✅ Gom tất cả
var errs ValidationErrors
// append từng cái
return errs
```

### 2. Validate quá lỏng

`len(s) > 0` ≠ `valid name`. Cần check:
- Whitespace only (`strings.TrimSpace`)
- Max length (chống DoS payload)
- Allowed characters (chống injection)

### 3. Error message lộ internal

```go
http.Error(w, err.Error(), 500)
// → "pq: relation users does not exist"  ⚠️ lộ DB info
```

Sanitize error trước khi trả client. Log chi tiết → server log, response → generic message.

### 4. Lỗi của lỗi

```go
writeJSON(w, 500, ...)
json.NewEncoder(w).Encode(v)  // có thể fail nhưng đã WriteHeader rồi → không trả được
```

→ Log lỗi encode, không thể "retry" response.

### 5. Status code sai semantically

| Code | Đúng | Sai |
|------|------|-----|
| 400 | JSON sai, format lỗi | DB unavailable (→ 503) |
| 404 | Resource không có | Method không hỗ trợ (→ 405) |
| 422 | Validation business fail | JSON parse fail (→ 400) |
| 500 | Internal bug | Client gửi sai (→ 400) |

## Bài tập

- Implement validator `nameRule(min, max)` reusable.
- Custom error code mapping table: `code → status code`.
- Async validation: check email không trùng trong DB.
- I18n: error message dựa vào `Accept-Language` header.
