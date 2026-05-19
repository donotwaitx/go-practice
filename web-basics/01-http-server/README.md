# 01 – HTTP Server cơ bản

Server HTTP đơn giản nhất với `net/http`. Không framework, chỉ standard library.

## Chạy

```bash
go run ./web-basics/01-http-server
```

Test:
```bash
curl -i http://localhost:8080/
curl -i http://localhost:8080/about
curl -i http://localhost:8080/teapot
```

## Khái niệm

### `http.HandlerFunc`

Adapter biến một function `func(w http.ResponseWriter, r *http.Request)` thành `http.Handler`:

```go
func hello(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Hello")
}

http.HandleFunc("/", hello)  // đăng ký vào DefaultServeMux
```

### Hai actor

| Object | Vai trò |
|--------|---------|
| `http.ResponseWriter` | Ghi response: status, header, body |
| `*http.Request` | Đọc request: method, URL, header, body |

### Thứ tự ghi response

Khi ghi response, có thứ tự bắt buộc:

```go
1. w.Header().Set(...)      // thêm header trước
2. w.WriteHeader(status)    // gửi status code (chỉ 1 lần)
3. fmt.Fprintln(w, body)    // ghi body — implicit WriteHeader(200) nếu chưa gọi
```

Gọi `WriteHeader` 2 lần → warning, header thứ 2 bị bỏ. Set header sau khi `WriteHeader` → bỏ qua.

### `ListenAndServe`

```go
http.ListenAndServe(":8080", nil)
```

- `:8080` — bind mọi interface, port 8080. `"localhost:8080"` chỉ bind loopback.
- `nil` — dùng `DefaultServeMux` (mux global). Truyền mux riêng để tránh global state (xem bài 02).
- **Block forever** đến khi error. Wrap với `log.Fatal` để in lỗi rồi exit.

## Status codes (HTTP)

| Constant | Code | Khi nào |
|----------|------|---------|
| `StatusOK` | 200 | Default — request thành công |
| `StatusCreated` | 201 | POST tạo resource thành công |
| `StatusNoContent` | 204 | DELETE thành công, không body |
| `StatusBadRequest` | 400 | Client gửi sai (validation fail) |
| `StatusUnauthorized` | 401 | Cần đăng nhập |
| `StatusForbidden` | 403 | Đăng nhập rồi nhưng không có quyền |
| `StatusNotFound` | 404 | Resource không tồn tại |
| `StatusInternalServerError` | 500 | Server bug |

## Pitfall

### 1. `nil` route → 404

Nếu path không match handler nào, ServeMux trả 404. Trong demo này, `/` match TẤT CẢ path (catch-all) vì đăng ký prefix matching. Vd: `/random` cũng đi vào `helloHandler`. Bài 02 sẽ fix điều này.

### 2. Sử dụng `http.DefaultServeMux`

`http.HandleFunc(path, h)` đăng ký vào DefaultServeMux toàn cục. Trong app lớn, code khác có thể đăng ký bậy vào → conflict. **Best practice:** tự `http.NewServeMux()` rồi truyền vào `ListenAndServe`.

### 3. Server không graceful shutdown

`ListenAndServe` không xử lý SIGTERM. Production cần:

```go
server := &http.Server{Addr: ":8080", Handler: mux}
go server.ListenAndServe()

// Listen SIGTERM/SIGINT...
server.Shutdown(ctx)  // dừng graceful, chờ in-flight requests xong
```

### 4. Quên `Content-Type`

Nếu không set, Go đoán theo body. Cho JSON cần tự set:
```go
w.Header().Set("Content-Type", "application/json")
```

## Bài tập

- Thêm route `/headers` in tất cả header client gửi (`for k, v := range r.Header`).
- Thêm route `/redirect` chuyển hướng sang `/about` bằng `http.Redirect`.
- Đo thời gian xử lý: in `time.Since(start)` cuối handler.
