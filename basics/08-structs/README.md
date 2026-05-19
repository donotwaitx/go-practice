# 08 – Structs & Methods

Go không có class — struct + method = mô hình OOP của Go.

## Chạy

```bash
go run ./basics/08-structs
```

## Struct

Tập hợp các field có tên và type:

```go
type User struct {
    ID    int
    Name  string
    Email string
}
```

Truy cập field bằng `.`:

```go
u.Name = "Alice"
fmt.Println(u.ID)
```

## Khởi tạo

```go
// Đặt theo tên field — KHUYẾN KHÍCH
u := User{ID: 1, Name: "Alice", Email: "a@x.com"}

// Theo thứ tự — không khuyến khích, fragile khi thêm field
u := User{1, "Alice", "a@x.com"}

// Zero value
var u User       // {0 "" ""}
u := User{}      // {0 "" ""}

// Pointer
p := &User{Name: "Bob"}
```

## Methods

Hàm gắn với receiver type. Cú pháp:

```go
func (u User) Greet() string {
    return "Hi, " + u.Name
}
```

`(u User)` là **receiver** — giống `this` trong các ngôn ngữ khác, nhưng tường minh.

Gọi: `u.Greet()`.

## Value receiver vs Pointer receiver

```go
func (u User) ChangeName(n string)  { u.Name = n }   // value
func (u *User) Rename(n string)     { u.Name = n }   // pointer
```

|  | Value receiver | Pointer receiver |
|--|----------------|------------------|
| Nhận | Bản copy | Pointer trỏ đến gốc |
| Sửa gốc được? | ❌ Không | ✅ Có |
| Copy struct? | ✅ (tốn với struct lớn) | ❌ |
| Gọi trên value và pointer? | Cả 2 | Cả 2 (Go tự `&`) |

**Quy tắc chọn:**
1. Cần modify struct → **pointer receiver**.
2. Struct lớn (nhiều field) → **pointer receiver** để tránh copy.
3. **Nhất quán:** nếu một method dùng pointer, các method khác cũng nên pointer.

## Embedded structs (composition)

Go không có inheritance, dùng **embedding**:

```go
type Employee struct {
    User              // embedded — không có tên field
    Salary   float64
    Position string
}

emp := Employee{
    User:     User{Name: "Charlie"},
    Salary:   1000,
    Position: "Engineer",
}

fmt.Println(emp.Name)      // truy cập trực tiếp field của User (promoted)
fmt.Println(emp.Greet())   // gọi method của User (promoted)
fmt.Println(emp.User.Name) // cũng OK
```

Đây là cách Go làm code reuse — composition over inheritance.

## Anonymous struct

Struct không tên, dùng khi cần một type tạm:

```go
point := struct {
    X, Y int
}{X: 3, Y: 4}

// Phổ biến trong test, config
config := struct {
    Host string
    Port int
}{"localhost", 8080}
```

## Tags

Metadata cho field, đọc bằng reflection — thường dùng cho JSON, DB ORM:

```go
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email,omitempty"`
}
```

`encoding/json` sẽ dùng tag để map field → key JSON.

## So sánh struct

```go
u1 := User{1, "A", "a"}
u2 := User{1, "A", "a"}
u1 == u2  // true nếu mọi field comparable
```

Struct chứa slice/map/func → không comparable.

## Ghi nhớ

- Không có constructor — viết hàm `NewUser(...)` trả về `*User` hoặc `User`.
- Không có visibility modifier — **field/method viết hoa chữ đầu = exported** (public), viết thường = unexported (private trong package).
- Không có inheritance, có composition qua embedding.
