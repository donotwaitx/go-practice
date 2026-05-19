# Web Basics với Go

Học web development với Go qua 10 bài. Dùng chủ yếu `net/http` của standard library — đây là điểm mạnh của Go, không cần framework cồng kềnh.

## Lộ trình

| # | Chủ đề | Bao gồm |
|---|--------|---------|
| 01 | http-server | `ListenAndServe`, handler, ResponseWriter, status code |
| 02 | routing | `ServeMux` (Go 1.22+) với method + path variables |
| 03 | query-form | URL query, form-urlencoded, header |
| 04 | json | `encoding/json` encode/decode, struct tags |
| 05 | rest-crud | Full CRUD API với in-memory store + mutex |
| 06 | middleware | Logger, recover, request ID, auth — chain pattern |
| 07 | context | Cancellation, timeout, request-scoped values |
| 08 | validation | Input validation, structured error response |
| 09 | database | `database/sql` + SQLite (pure Go, no cgo) |
| 10 | testing | `httptest` — table-driven test, benchmark |

## Yêu cầu

- Đã học xong `basics/` (đặc biệt: structs, interfaces, errors, goroutines, context).
- Go 1.22+ để dùng enhanced ServeMux.

## Cách chạy

Mỗi bài là server HTTP listen trên `:8080`. Mở 2 terminal:

```bash
# Terminal 1: chạy server
go run ./web-basics/01-http-server

# Terminal 2: test bằng curl
curl http://localhost:8080/
```

Stop server bằng `Ctrl+C`.

## Cách test với curl (cheatsheet)

```bash
# GET đơn giản
curl http://localhost:8080/path

# GET có verbose (xem header + status)
curl -i http://localhost:8080/path

# POST với JSON
curl -X POST http://localhost:8080/users \
  -H 'Content-Type: application/json' \
  -d '{"name":"Alice","email":"a@x.com"}'

# Form
curl -X POST http://localhost:8080/login \
  -d 'username=alice&password=secret'

# Custom header
curl -H 'X-API-Key: secret' http://localhost:8080/secret

# Pretty-print JSON response
curl -s http://localhost:8080/users | jq
```

## Triết lý Go web development

Khác Node/Python, Go community không khuyến khích framework lớn:

1. **`net/http` đủ mạnh** cho 90% use case sau Go 1.22.
2. **Composable** — middleware là function, dễ chain.
3. **Concurrency built-in** — mỗi request 1 goroutine, không cần async/await.
4. **No magic** — code rõ ràng, không reflection trickery.

Khi nào cần framework? Gin, Echo, Fiber, Chi nhanh hơn về DX cho:
- Validation phức tạp
- Binding tự động
- Subrouter, group middleware
- Templating tích hợp

→ Học `net/http` trước, framework chỉ là "đường tắt" — phải hiểu cái nền tảng.

## Sau module này

- **Cache & rate limiting** — `golang.org/x/time/rate`, Redis.
- **Observability** — logging structured (`slog`), tracing (OpenTelemetry), metrics (Prometheus).
- **Production deploy** — graceful shutdown, signal handling, config (viper/env).
- **Real DB** — PostgreSQL với `pgx`, migrations với `goose`/`golang-migrate`.
- **Framework comparison** — thử Gin/Echo/Chi để so với raw `net/http`.
