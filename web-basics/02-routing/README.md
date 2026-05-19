# 02 – Routing với ServeMux (Go 1.22+)

Từ Go 1.22, `http.ServeMux` được nâng cấp lớn: hỗ trợ **HTTP method** và **path variables**. Trước đó phải dùng third-party (gorilla/mux, chi...) — giờ không cần nữa.

## Chạy

```bash
go run ./web-basics/02-routing
```

Test:
```bash
curl http://localhost:8080/                       # Home
curl http://localhost:8080/users/42                # Get user 42
curl -X POST http://localhost:8080/users           # Create
curl -X DELETE http://localhost:8080/users/42      # Delete
curl http://localhost:8080/users/1/posts/99        # Nested
curl http://localhost:8080/files/a/b/c.txt         # Wildcard
```

## Pattern syntax (Go 1.22+)

```
[METHOD] [HOST]/[PATH]
```

| Pattern | Match |
|---------|-------|
| `/users` | GET hoặc bất kỳ method, đúng path `/users` |
| `GET /users` | Chỉ GET `/users` |
| `POST /users` | Chỉ POST `/users` |
| `GET /users/{id}` | GET `/users/X` với X là segment bất kỳ. Lấy bằng `r.PathValue("id")` |
| `GET /users/{id}/posts/{postID}` | Nested variable |
| `GET /files/{path...}` | Wildcard — `path` chứa toàn bộ remainder (kể cả `/`) |
| `api.example.com/v1/ping` | Host-specific |

### Lấy path variable

```go
mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")  // string
    // chuyển sang int:
    n, err := strconv.Atoi(id)
    // ...
})
```

`PathValue` trả `""` nếu key không có. Không panic.

## Method routing — lợi gì?

Trước Go 1.22:
```go
mux.HandleFunc("/users", func(w, r) {
    switch r.Method {
    case "GET":   ...
    case "POST":  ...
    default: http.Error(w, "method not allowed", 405)
    }
})
```

Sau:
```go
mux.HandleFunc("GET /users",  listHandler)
mux.HandleFunc("POST /users", createHandler)
// Method khác → tự trả 405 Method Not Allowed
```

Sạch hơn nhiều, ít bug.

## Precedence — pattern nào match trước?

Khi nhiều pattern khớp, **pattern cụ thể hơn thắng**:

| Pattern | Specificity |
|---------|-------------|
| `/users/me` | Cao nhất (literal segment) |
| `/users/{id}` | Trung bình (1 variable) |
| `/users/{path...}` | Thấp nhất (wildcard) |

Conflict không giải quyết được → panic khi `ListenAndServe`.

## Trailing slash

```
mux.HandleFunc("/users/",  listHandler)   // /users/anything (prefix match)
mux.HandleFunc("/users",   listHandler)   // ĐÚNG /users
```

Pattern có `/` cuối = prefix match. Không có = exact match. Cẩn thận khi định nghĩa.

## So với router của framework

Go 1.22 ServeMux đủ cho 90% use case. Tuy nhiên KHÔNG có:

| Feature | ServeMux | Chi / Gorilla / Echo |
|---------|----------|-----------------------|
| Method + path variable | ✅ | ✅ |
| Wildcard | ✅ | ✅ |
| Regex constraint | ❌ | ✅ (vd `/{id:[0-9]+}`) |
| Subrouter / group | ❌ | ✅ |
| Middleware chain | Phải tự | ✅ built-in |
| Param parsing helper | ❌ | ✅ |

Nếu cần group + middleware nhiều, framework như `chi` (vẫn rất minimal) tốt hơn.

## Pitfall

### 1. Path không match → 404

```go
mux.HandleFunc("GET /users/{id}", h)
// curl /users/      → 404 (không match {id} rỗng)
// curl /users/1/    → 404 (trailing slash)
```

### 2. PathValue trên key không tồn tại

`r.PathValue("wrong")` trả `""`, không panic. Phải validate.

### 3. Nhầm `HandleFunc` vs `Handle`

- `mux.HandleFunc(pattern, func)` — nhận function.
- `mux.Handle(pattern, http.Handler)` — nhận implement Handler interface.

Cả 2 đều OK, chỉ là syntax tiện.

### 4. Đăng ký 2 lần

```go
mux.HandleFunc("GET /a", h1)
mux.HandleFunc("GET /a", h2)  // panic: pattern "/a" conflicts
```

## Bài tập

- Thêm route `GET /api/v{version}/users` — bắt version number.
- Implement static file server với wildcard: `GET /static/{path...}` → đọc file.
- Versioned API: `/v1/users` vs `/v2/users` với handler khác.
