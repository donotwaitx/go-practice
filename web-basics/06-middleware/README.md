# 06 – Middleware

Middleware là layer chạy TRƯỚC/SAU handler chính, thêm logic chung (log, recover, auth, metrics...). Là pattern then chốt của Go web.

## Chạy

```bash
go run ./web-basics/06-middleware
```

Test:
```bash
curl http://localhost:8080/hello
curl http://localhost:8080/panic       # → 500 nhưng server không crash
curl -i http://localhost:8080/hello    # xem X-Request-ID header
```

## Pattern

```go
func MyMW(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // BEFORE: code chạy trước handler
        next.ServeHTTP(w, r)
        // AFTER: code chạy sau handler
    })
}
```

Middleware = **wrapper**. Nó NHẬN một handler, TRẢ VỀ handler khác (có wrapping logic).

## Compose

Chain nhiều middleware:

```go
handler = Logger(Recover(Auth(realHandler)))
//        ▲outer                     inner▲
//
// Request đi vào theo thứ tự: Logger → Recover → Auth → realHandler
// Response đi ra theo thứ tự ngược: realHandler → Auth → Recover → Logger
```

Helper `chain`:

```go
func chain(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
    for i := len(mws) - 1; i >= 0; i-- {
        h = mws[i](h)
    }
    return h
}

handler := chain(mux, Logger, Recover, Auth)  // Logger outer nhất
```

## Middleware kinh điển

### 1. Logger

```go
start := time.Now()
sw := &statusWriter{ResponseWriter: w, status: 200}
next.ServeHTTP(sw, r)
log.Printf("%s %s -> %d (%s)", r.Method, r.URL.Path, sw.status, time.Since(start))
```

Cần wrap ResponseWriter để capture status code (Go không expose sẵn).

### 2. Recover

Bắt panic, trả 500 thay vì crash server:

```go
defer func() {
    if rec := recover(); rec != nil {
        log.Printf("panic: %v\n%s", rec, debug.Stack())
        http.Error(w, "internal server error", 500)
    }
}()
next.ServeHTTP(w, r)
```

⚠️ **Quan trọng:** Go server **mỗi request 1 goroutine**, panic 1 goroutine không kill server. Nhưng trang đó treo → user nhận response trống. Recover middleware giải quyết.

### 3. Request ID

Gán ID duy nhất cho mỗi request → tracing log:

```go
id := r.Header.Get("X-Request-ID")
if id == "" {
    id = generateID()
}
w.Header().Set("X-Request-ID", id)
ctx := context.WithValue(r.Context(), requestIDKey, id)
next.ServeHTTP(w, r.WithContext(ctx))
```

### 4. Auth (API key)

```go
if r.Header.Get("X-API-Key") != expectedKey {
    http.Error(w, "unauthorized", 401)
    return  // ⚠️ KHÔNG gọi next
}
next.ServeHTTP(w, r)
```

### 5. CORS

```go
w.Header().Set("Access-Control-Allow-Origin", "*")
w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
if r.Method == "OPTIONS" {
    w.WriteHeader(204)
    return
}
next.ServeHTTP(w, r)
```

### 6. Rate limiting

Đếm request per IP, từ chối nếu quá. Stdlib không có — dùng `golang.org/x/time/rate`.

## Per-route middleware

Một số middleware chỉ áp cho một số route:

```go
mux.Handle("GET /public",         publicHandler)
mux.Handle("GET /admin",   authMW(adminHandler))
mux.Handle("POST /admin",  authMW(adminCreate))
```

Hoặc nhóm — nhưng `ServeMux` chuẩn không có group. Cần `chi` hoặc tự viết.

## Pitfall

### 1. Quên gọi `next.ServeHTTP`

```go
func brokenMW(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println("before")
        // ❌ quên gọi next → handler chính không chạy
    })
}
```

### 2. Gọi `next` sau khi đã `WriteHeader`

```go
w.WriteHeader(200)
next.ServeHTTP(w, r)  // ⚠️ handler trong vẫn ghi được, nhưng header đã sent
```

Trừ khi cố tình "short-circuit", thường KHÔNG gọi `next` nếu đã `WriteHeader`.

### 3. Modify Request without `WithContext`

```go
r.Context() = ctx  // ❌ không compile
```

Request có context immutable. Tạo request mới:
```go
r = r.WithContext(ctx)
next.ServeHTTP(w, r)
```

### 4. Context value key dùng string thay vì custom type

```go
ctx := context.WithValue(r.Context(), "userID", id)  // ❌
```

→ Conflict với code khác cũng dùng `"userID"`. Dùng custom unexported type:

```go
type ctxKey string
const userIDKey ctxKey = "userID"
ctx := context.WithValue(r.Context(), userIDKey, id)  // ✅
```

### 5. statusWriter wrap không đầy đủ

Một số middleware như compression, streaming cần wrap thêm method:

```go
type fancyWriter struct {
    http.ResponseWriter
    status int
    bytes  int
}

func (w *fancyWriter) WriteHeader(s int) { w.status = s; w.ResponseWriter.WriteHeader(s) }
func (w *fancyWriter) Write(b []byte) (int, error) {
    n, err := w.ResponseWriter.Write(b)
    w.bytes += n
    return n, err
}
```

## Bài tập

- Implement `metricsMW` đếm số request / status code, expose qua `/metrics`.
- Implement `rateLimitMW` — max 10 req/giây/IP.
- Implement `cacheMW` cho GET, cache 60 giây trong memory.
- Implement `gzipMW` compress response nếu client gửi `Accept-Encoding: gzip`.
