# 07 – Context: cancellation, timeout, request-scoped values

`context.Context` là cách Go truyền tín hiệu (cancel, deadline) và giá trị (request-scoped) xuống call chain. Trong HTTP server: mỗi request đã có sẵn `r.Context()`.

## Chạy

```bash
go run ./web-basics/07-context
```

Test:
```bash
# Slow query — Ctrl+C giữa chừng để xem cancel
curl 'http://localhost:8080/slow?ms=5000'

# Timeout — request mất >3s, handler cancel sau 1s
curl http://localhost:8080/timeout

# Pass value qua context
curl -H 'X-User: alice' http://localhost:8080/whoami
```

## Tại sao cần Context?

3 vấn đề Context giải quyết:

1. **Cancellation cascade** — client disconnect → server biết, dừng DB query/HTTP call/computation đang chạy.
2. **Deadline** — đặt thời gian tối đa cho chuỗi operation.
3. **Request-scoped values** — truyền user info, request ID xuống các layer mà không cần parameter hardcode.

## API chính

### Tạo context

```go
ctx := context.Background()         // root context, dùng làm parent
ctx := context.TODO()               // placeholder khi chưa biết dùng gì

// Derive với deadline
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
ctx, cancel := context.WithDeadline(parent, deadline)

// Derive với cancel
ctx, cancel := context.WithCancel(parent)

// Derive với value
ctx := context.WithValue(parent, key, value)

defer cancel()  // ⚠️ LUÔN gọi để giải phóng resource
```

### Đọc context

```go
ctx.Done()        // <-chan struct{} — đóng khi cancelled/deadline
ctx.Err()         // nil khi active, context.Canceled / DeadlineExceeded khi xong
ctx.Deadline()    // (time.Time, bool)
ctx.Value(key)    // any
```

## Pattern: respect context trong long operation

```go
func slowQuery(ctx context.Context, ms int) (string, error) {
    select {
    case <-time.After(time.Duration(ms) * time.Millisecond):
        return "done", nil
    case <-ctx.Done():
        return "", ctx.Err()
    }
}
```

Mọi operation block (DB query, HTTP call, sleep) nên check ctx.

## HTTP server tự manage context

Khi client disconnect (đóng tab, network drop), Go **tự cancel** `r.Context()`:

```go
func handler(w http.ResponseWriter, r *http.Request) {
    result, err := slowQuery(r.Context(), 10000)
    // Nếu client đóng kết nối → err = context.Canceled
}
```

Đây là "magic" của Go HTTP — tận dụng cho mọi I/O downstream.

## Context-aware libraries

Stdlib chuẩn:

```go
db.QueryContext(ctx, "SELECT ...")          // database/sql
http.NewRequestWithContext(ctx, ...)        // outgoing HTTP
cmd.CommandContext(ctx, "ls")               // os/exec
```

Library tốt luôn có `XxxContext` version. Khi gặp lib không có → là red flag.

## Request-scoped values

```go
type ctxKey string  // unexported type tránh collision
const userKey ctxKey = "user"

// SET (thường trong middleware)
ctx := context.WithValue(r.Context(), userKey, user)
next.ServeHTTP(w, r.WithContext(ctx))

// GET (trong handler)
user, _ := r.Context().Value(userKey).(string)
```

⚠️ **Best practice với Value:**

- Chỉ dùng cho **request-scoped data**: user ID, request ID, trace ID, logger có request context.
- KHÔNG dùng thay parameter function bình thường.
- Type assertion luôn dùng `, ok` pattern phòng nil.

```go
// BAD — config nên là parameter, không phải value
ctx := context.WithValue(ctx, "config", cfg)

// GOOD
func process(ctx context.Context, cfg *Config) { ... }
```

## Pitfall

### 1. Quên `defer cancel()`

```go
ctx, cancel := context.WithTimeout(ctx, time.Second)
// ⚠️ thiếu defer cancel — context bị leak (vẫn timer chạy)
```

Linter `govet` warn cái này.

### 2. Lưu Context vào struct

```go
type Server struct {
    ctx context.Context  // ❌ anti-pattern
}
```

Context là **request-scoped**, không phải config. Truyền qua parameter.

```go
func (s *Server) Handle(ctx context.Context, ...) { ... }  // ✅
```

### 3. Pass nil context

```go
db.QueryContext(nil, "SELECT")  // ❌ panic
```

Dùng `context.Background()` hoặc `context.TODO()`.

### 4. Type assertion fail im lặng

```go
user := r.Context().Value(userKey).(string)  // ❌ panic nếu nil
user, ok := r.Context().Value(userKey).(string)  // ✅
```

### 5. Cancel sau khi đã có error path

```go
ctx, cancel := context.WithTimeout(...)
defer cancel()
result, err := doStuff(ctx)
cancel()  // ⚠️ thừa defer + explicit — không hại nhưng confused
```

Chọn 1 trong 2: `defer cancel()` HOẶC `cancel()` ngay sau khi xong.

### 6. Context không xuyên goroutine

```go
go func() {
    // ⚠️ goroutine này KHÔNG có context tự động
    // → khi request kết thúc, goroutine vẫn chạy → leak
}()
```

Pass ctx vào goroutine, check ctx.Done().

## Bài tập

- Endpoint `/proxy?url=...` — gọi HTTP downstream với timeout 5s, respect client cancel.
- Middleware `timeoutMW(5*time.Second)` áp dụng timeout chung cho mọi request.
- Pattern: store logger có request ID trong context, các function lấy ra dùng.
