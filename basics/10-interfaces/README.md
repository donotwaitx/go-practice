# 10 – Interfaces

Polymorphism kiểu Go: **structural typing**. Không cần `implements`, type tự "khớp" interface nếu có đủ methods.

## Chạy

```bash
go run ./basics/10-interfaces
```

## Định nghĩa interface

```go
type Shape interface {
    Area() float64
    Perimeter() float64
}
```

Interface = **tập hợp method signature**. Không có field, không có implementation.

## Implement — TỰ ĐỘNG

```go
type Rectangle struct{ W, H float64 }

func (r Rectangle) Area() float64      { return r.W * r.H }
func (r Rectangle) Perimeter() float64 { return 2 * (r.W + r.H) }
```

`Rectangle` có cả 2 method `Shape` yêu cầu → tự động implement `Shape`. **Không cần khai báo `implements Shape`** như Java.

Đây gọi là **structural typing** (hoặc duck typing tĩnh): "If it walks like a duck and quacks like a duck, it's a duck."

## Polymorphism

```go
func describe(s Shape) {
    fmt.Printf("area=%.2f\n", s.Area())
}

describe(Rectangle{3, 4})
describe(Circle{5})
```

Hàm nhận `Shape` — chấp nhận bất kỳ type nào implement đủ methods.

## Empty interface — `interface{}` / `any`

Không có method → mọi type đều implement → chứa được bất kỳ kiểu nào:

```go
var x any = 42
x = "hello"
x = Rectangle{}
```

Từ Go 1.18, `any` là alias chính thức của `interface{}`. Dùng `any` cho rõ.

⚠️ Lạm dụng `any` là **anti-pattern** — mất type safety. Chỉ dùng khi thật sự cần (vd: `fmt.Println`, JSON decode dynamic).

## Type assertion

Lấy concrete type từ interface:

```go
var s Shape = Rectangle{3, 4}

r := s.(Rectangle)        // ❌ panic nếu s không phải Rectangle
r, ok := s.(Rectangle)    // ✅ ok=false nếu không khớp, không panic
```

Luôn dùng dạng 2 giá trị nếu không chắc chắn.

## Type switch

Phân nhánh theo type khi xử lý interface (đặc biệt `any`):

```go
func inspect(v any) {
    switch x := v.(type) {
    case int:
        fmt.Println("int:", x*2)
    case string:
        fmt.Println("string:", x)
    case Shape:
        fmt.Println("shape area:", x.Area())
    default:
        fmt.Println("unknown")
    }
}
```

`x` được tự động cast sang type của case match.

## Interface composition

Gộp nhiều interface:

```go
type Reader interface { Read(p []byte) (n int, err error) }
type Writer interface { Write(p []byte) (n int, err error) }

type ReadWriter interface {
    Reader
    Writer
}
```

`io.ReadWriter` chuẩn được định nghĩa kiểu này.

## Best practices

### 1. Interface nhỏ

Idiomatic Go: interface 1-3 methods. Ví dụ `io.Reader` chỉ có `Read`. → Dễ implement, dễ test với mock.

> "The bigger the interface, the weaker the abstraction." — Rob Pike

### 2. Define ở chỗ DÙNG (consumer-defined)

```go
// ❌ Đừng export interface từ package implement
// ✅ Định nghĩa interface ở package CONSUMER, dùng concrete type ở producer

// package db
func GetUser(id int) (*User, error) { ... }

// package report
type UserGetter interface {
    GetUser(id int) (*User, error)
}
func Generate(g UserGetter) { ... }
```

→ Decoupling tốt hơn.

### 3. Accept interfaces, return structs

Hàm public:
- Tham số nên là interface (linh hoạt với caller).
- Trả về nên là struct/concrete type (rõ ràng với caller).

## Nil interface vs interface chứa nil

Pitfall kinh điển:

```go
var p *MyError = nil
var err error = p
fmt.Println(err == nil)  // false ❗
```

Interface có 2 phần: (type, value). Ở đây type=`*MyError`, value=nil → interface KHÔNG nil dù value nil.

**Tránh:** đừng wrap typed nil vào interface. Trả về `nil` literal cho error.

## Ghi nhớ

- Interface satisfaction là **structural**, không phải nominal.
- Empty interface `any` = "bất kỳ type nào", nhưng phải type assert mới dùng được.
- Method set: value type `T` có methods với receiver `T`. Pointer `*T` có methods với cả `T` và `*T`. → Pointer thường dễ "fit" interface hơn.
