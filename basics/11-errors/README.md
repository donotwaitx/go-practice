# 11 – Errors

Go dùng **explicit error returns** thay vì exceptions. Đây là điểm gây tranh cãi nhất của Go, nhưng có lý do.

## Chạy

```bash
go run ./basics/11-errors
```

## Triết lý

- Không có `try/catch`.
- Mọi hàm có thể fail đều trả thêm giá trị `error`.
- Caller PHẢI check và xử lý.

→ Errors là **giá trị**, không phải control flow. Code error path nằm cạnh code happy path, dễ thấy, dễ reason.

## Pattern chuẩn

```go
result, err := doSomething()
if err != nil {
    return nil, err
    // hoặc: return fmt.Errorf("doing X: %w", err)
}
// dùng result...
```

Nhìn quen với mọi developer Go.

## `error` là gì?

Interface built-in:

```go
type error interface {
    Error() string
}
```

→ Bất kỳ type nào có method `Error() string` đều là `error`.

## Tạo error

### `errors.New`

```go
err := errors.New("something went wrong")
```

### `fmt.Errorf` (có format)

```go
err := fmt.Errorf("user %d not found", id)
```

### Sentinel error — error có thể so sánh

```go
var ErrNotFound = errors.New("not found")

func find(id int) (*User, error) {
    if !exists(id) {
        return nil, ErrNotFound
    }
    // ...
}

// Caller
u, err := find(1)
if errors.Is(err, ErrNotFound) {
    // xử lý riêng
}
```

### Custom error type

```go
type ValidationError struct {
    Field string
    Msg   string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation: %s - %s", e.Field, e.Msg)
}
```

→ Chứa thêm context (field nào, vì sao). Caller extract bằng `errors.As`.

## Error wrapping

```go
if err != nil {
    return fmt.Errorf("loading config: %w", err)
}
```

`%w` (Go 1.13+) **giữ** error gốc trong chain. `%v` chỉ in string, mất chain.

→ Caller có thể truy ngược chain bằng `errors.Is` / `errors.As`.

## `errors.Is` vs `errors.As`

### `Is` — so sánh với sentinel error

```go
if errors.Is(err, ErrNotFound) {
    // err HOẶC một error trong chain == ErrNotFound
}
```

### `As` — extract error theo type

```go
var vErr *ValidationError
if errors.As(err, &vErr) {
    fmt.Println(vErr.Field)
}
```

Tìm trong chain một error có type khớp, gán vào `vErr`.

## Best practices

### 1. Thêm context khi wrap

```go
// ❌
return err

// ✅
return fmt.Errorf("processing user %d: %w", id, err)
```

Khi error log ra, đọc được "processing user 42: db connection refused" → biết ngay where & why.

### 2. Đừng vừa log vừa return

```go
// ❌ double reporting
log.Println(err)
return err

// ✅ chọn 1: hoặc log ở top-level, hoặc wrap rồi return
return fmt.Errorf("context: %w", err)
```

### 3. Kiểm tra error TRƯỚC khi dùng result

```go
result, err := f()
if err != nil { return err }
// dùng result an toàn
```

KHÔNG đảo: `if err != nil { ... }` rồi dùng result — Go cho phép nhưng nguy hiểm.

### 4. Đừng compare error bằng `==` với wrapped error

```go
if err == ErrNotFound { ... }       // ❌ fail nếu err đã wrap
if errors.Is(err, ErrNotFound) {}   // ✅
```

## `panic` / `recover`

Cho lỗi BẤT THƯỜNG, programmer error, không nên xảy ra:

```go
if i < 0 {
    panic("negative index — programmer error")
}
```

`recover` bắt panic trong defer:

```go
func safeCall() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("recovered:", r)
        }
    }()
    riskyOperation()
}
```

⚠️ **Đừng dùng panic thay error** cho business logic. Quy tắc:
- File không mở được → return error.
- Map bị nil (programmer bug) → panic OK.
- Invariant violation trong code → panic.

Một số HTTP framework dùng recover ở top-level để 1 request panic không kill server.

## So với try/catch

|  | Try/catch | Go errors |
|--|-----------|-----------|
| Khả thi miss error? | Có (silent catch, rethrow) | Khó miss — compiler buộc check |
| Boilerplate | Ít | Nhiều |
| Stack trace | Tự có | Phải tự thêm context khi wrap |
| Control flow | Implicit (throw nhảy xa) | Explicit |

Đánh đổi: code dài hơn, đổi lại rõ ràng và predictable.

## Ghi nhớ

- `error` luôn là giá trị return CUỐI cùng.
- Trả về `nil` cho error nếu thành công.
- `errors.Is`, `errors.As`, `%w` từ Go 1.13 — nên dùng thay vì compare bằng `==`.
