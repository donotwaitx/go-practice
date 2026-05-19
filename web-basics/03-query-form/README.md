# 03 – Query string & Form data

3 cách client gửi data: query string, form body, headers. Đây là bài về cách đọc.

## Chạy

```bash
go run ./web-basics/03-query-form
```

Test:
```bash
# Query string
curl 'http://localhost:8080/search?q=golang&page=2&tags=web&tags=server'

# Form (POST body)
curl -X POST http://localhost:8080/login \
  -d 'username=alice&password=secret123'

# Header
curl -H 'Authorization: Bearer abc123' http://localhost:8080/headers
```

## Query string

```
GET /search?q=golang&page=2&tags=web&tags=server
                ▲       ▲       ▲          ▲
                key=value pairs, tags lặp lại
```

### API

```go
q := r.URL.Query()        // url.Values = map[string][]string

q.Get("q")                // string — chỉ value đầu tiên, "" nếu không có
q["tags"]                 // []string — tất cả values của key
q.Has("filter")           // bool — chỉ check key có tồn tại
```

### Parse số

```go
page, err := strconv.Atoi(q.Get("page"))
if err != nil { page = 1 } // default
```

Hoặc dùng `strconv.ParseInt`, `ParseBool`, `ParseFloat` tùy type.

## Form data

3 content-types phổ biến gửi data:

| Content-Type | Use case | Đọc bằng |
|--------------|----------|----------|
| `application/x-www-form-urlencoded` | HTML form đơn giản | `r.ParseForm()` |
| `multipart/form-data` | Upload file | `r.ParseMultipartForm(maxBytes)` |
| `application/json` | API hiện đại | `json.NewDecoder(r.Body).Decode()` (bài 04) |

### Form-urlencoded

```go
if err := r.ParseForm(); err != nil { ... }

// Sau ParseForm, dùng:
r.FormValue(key)       // query string + body, body ưu tiên
r.PostFormValue(key)   // CHỈ body POST/PUT
```

### Upload file (multipart)

```go
r.ParseMultipartForm(10 << 20)  // max 10MB
file, header, err := r.FormFile("avatar")
defer file.Close()

// header.Filename, header.Size, header.Header...
io.Copy(dst, file)
```

## Headers

```go
r.Header.Get("Authorization")    // chỉ value đầu tiên
r.Header.Values("Accept")        // tất cả values
r.UserAgent()                    // r.Header.Get("User-Agent")
r.Referer()                      // r.Header.Get("Referer")
r.Host                           // r.Header.Get("Host") nhưng đặc biệt
```

Custom header convention: `X-Foo-Bar`. Một số "magic" header phổ biến:
- `Authorization`, `Content-Type`, `Accept`
- `X-Request-ID`, `X-Forwarded-For`, `X-Real-IP`

## Pitfall

### 1. Quên `ParseForm`

```go
// Sai
v := r.PostForm.Get("user")  // ❌ panic — r.PostForm = nil

// Đúng
r.ParseForm()
v := r.PostForm.Get("user")

// Hoặc dùng PostFormValue (tự ParseForm bên trong)
v := r.PostFormValue("user")
```

### 2. Tin tưởng client input

```go
n, _ := strconv.Atoi(r.URL.Query().Get("page"))
// Nếu client gửi "page=abc" → n = 0 (vì bỏ err)
// Nếu client gửi "page=-99999999999" → overflow
```

→ Validate range, default sensible.

### 3. Sensitive data trong query string

Query string xuất hiện trong:
- Browser history
- Server log access.log
- Referer header

→ KHÔNG bao giờ truyền password, token qua query. Dùng POST body hoặc header.

### 4. URL decoding

Go tự decode percent-encoding. `?q=hello%20world` → `q.Get("q")` = `"hello world"`. Không cần decode tay.

### 5. Body size

```go
r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // limit 1MB
```

Không limit → client gửi 1GB → DoS.

## Bài tập

- Endpoint `/sort?by=name&order=desc` — validate `by` chỉ trong whitelist `[name, age, email]`.
- Endpoint `/upload` nhận file ảnh, lưu vào disk, trả URL.
- Endpoint `/echo` trả lại mọi query + form + header dưới dạng JSON.
