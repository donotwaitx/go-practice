# 02 – Variables & Constants

Cách khai báo biến, hằng số, và khái niệm zero value.

## Chạy

```bash
go run ./basics/02-variables
```

## 4 cách khai báo biến

```go
var age int = 25      // 1. Đầy đủ: type + value
var name = "Trong"    // 2. Suy luận type từ value
var city string       // 3. Chỉ type, value = zero value
city := "Hanoi"       // 4. Short declaration — TỰ ĐỘNG suy luận
```

**Lưu ý:** `:=` **CHỈ dùng trong hàm**. Ngoài hàm (package level) phải dùng `var`.

## Khai báo nhiều biến

```go
var (
    x, y int    = 1, 2
    ok   bool   = true
    msg  string = "hello"
)
```

Hoặc inline: `var a, b = 1, "two"`

## Zero values

Biến khai báo mà không gán → có **zero value** mặc định:

| Type | Zero value |
|------|------------|
| `int`, `float64`, các số | `0` |
| `string` | `""` (chuỗi rỗng) |
| `bool` | `false` |
| pointer, slice, map, function, channel, interface | `nil` |
| struct | struct với tất cả field là zero value |

Khác C/C++ — Go KHÔNG có "biến chưa khởi tạo có giá trị rác".

## Hằng số (`const`)

```go
const Pi = 3.14159
const (
    StatusOK    = 200
    StatusError = 500
)
```

- Giá trị phải biết được tại **compile time**.
- Không có địa chỉ → không thể `&Pi`.
- Có thể khai báo cấp package hoặc trong hàm.

## Scope

- Khai báo **ngoài hàm** → cấp package (toàn bộ file/package dùng được).
- Khai báo **trong hàm** → chỉ tồn tại trong block đó.

## Pitfall thường gặp

```go
x := 5
x := 10   // ❌ no new variables on left side of :=
x = 10    // ✅ gán lại bằng =
```

`:=` yêu cầu có ÍT NHẤT MỘT biến mới. Khi đã khai báo rồi, dùng `=`.

```go
a, b := 1, 2
a, c := 3, 4   // ✅ OK — c là biến mới, a được gán lại
```

## Biến không dùng = compile error

```go
func main() {
    x := 10   // ❌ x declared and not used
}
```

Go ép bạn xóa code dead. Dùng `_` để discard nếu cần.
