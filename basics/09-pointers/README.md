# 09 – Pointers

Pointer là biến chứa **địa chỉ memory** của biến khác. Go có pointer như C nhưng đơn giản và an toàn hơn.

## Chạy

```bash
go run ./basics/09-pointers
```

## Cú pháp

```go
x := 10
p := &x       // & lấy địa chỉ — p là *int trỏ đến x
fmt.Println(p)   // 0xc0000180a8 (địa chỉ)
fmt.Println(*p)  // 10 — dereference, lấy giá trị

*p = 20       // sửa giá trị x qua pointer
fmt.Println(x)  // 20
```

| Toán tử | Ý nghĩa |
|---------|---------|
| `&x` | Address-of — lấy địa chỉ của `x` |
| `*p` | Dereference — lấy giá trị tại địa chỉ `p` |
| `*int` | Kiểu: pointer trỏ đến int |

## Tại sao cần pointer?

### 1. Sửa giá trị từ hàm khác

```go
func double(n *int) {
    *n = *n * 2
}

x := 5
double(&x)
fmt.Println(x)  // 10
```

Không có pointer → hàm chỉ thấy bản copy, sửa không ảnh hưởng gốc.

### 2. Tránh copy struct lớn

```go
func process(u *User) { ... }  // truyền 8 byte (pointer)
func process(u User) { ... }   // copy toàn bộ struct
```

### 3. Biểu diễn "không có giá trị" với nil

```go
var p *User       // nil
if p == nil { ... }
```

## `new()` vs `&T{}`

```go
p1 := new(int)       // pointer đến int = 0
p2 := &User{}        // pointer đến User zero value
p3 := &User{Name: "A"}  // pointer đến User đã init
```

`new(T)` cấp phát T zero value, trả về `*T`. Tương đương `var t T; p := &t`.

Trong thực tế `&T{}` phổ biến hơn vì cho phép init field.

## Nil pointer

```go
var p *int
fmt.Println(p)    // <nil>
fmt.Println(*p)   // ❌ panic: runtime error: invalid memory address
```

Dereference nil → **panic**. Luôn check nil khi không chắc.

## Pointer to struct — Go tự syntactic sugar

```go
type Point struct { X, Y int }
p := &Point{X: 1, Y: 2}

// Cả hai tương đương:
fmt.Println((*p).X)
fmt.Println(p.X)       // Go tự dereference
```

Tương tự với method — Go tự `&` hoặc `*` để khớp receiver:

```go
type Counter struct{ value int }
func (c *Counter) Inc() { c.value++ }

c := Counter{}
c.Inc()    // Go tự chuyển thành (&c).Inc()
```

## ❌ Không có pointer arithmetic

Khác C:
```c
int* p = arr;
p++;  // C: trỏ đến phần tử tiếp
```

Go không cho phép `p + 1`, `p++` với pointer. → An toàn hơn, không có buffer overflow kiểu C.

## Khi nào dùng / không dùng pointer?

**Dùng pointer khi:**
- Cần modify từ caller.
- Struct lớn (tránh copy).
- Method cần modify receiver.
- Optional value (`*int` nilable thay vì `int`).

**Tránh pointer khi:**
- Type cơ bản nhỏ (`int`, `bool`, `string`...) — copy rẻ hơn pointer indirection.
- Bất biến (immutable) — value semantics rõ ràng hơn.

## Ghi nhớ

- Pointer không phải "tham chiếu" của Java — nó thật sự là địa chỉ memory.
- Slice, map, channel, function đã chứa pointer bên trong — truyền chúng KHÔNG cần `&`.
- `&` chỉ dùng được với biến có địa chỉ (addressable). `&User{}` được, `&5` thì không. `&someFunc()` cũng không.
