# 07 – Maps

Hash table built-in. Key có thể là bất kỳ kiểu so sánh được (comparable).

## Chạy

```bash
go run ./basics/07-maps
```

## Khai báo

```go
// Literal
ages := map[string]int{
    "Alice": 30,
    "Bob":   25,
}

// make
m := make(map[string]int)
m := make(map[string]int, 100)  // hint capacity

// nil map — chỉ đọc được, GHI sẽ PANIC
var m map[string]int
m["x"] = 1  // ❌ panic: assignment to entry in nil map
```

⚠️ **Luôn `make` hoặc literal trước khi ghi.**

## Đọc / Ghi / Xóa

```go
ages["Charlie"] = 28      // thêm hoặc update
v := ages["Alice"]        // đọc — nếu không có, trả về zero value
delete(ages, "Bob")       // xóa key (an toàn nếu key không tồn tại)
n := len(ages)            // số phần tử
```

## Comma-ok idiom

Vì map trả zero value khi key không tồn tại, không thể phân biệt:
- Key không có
- Key có với value zero

Dùng dạng 2 giá trị trả về:

```go
if age, ok := ages["David"]; ok {
    fmt.Println("Có David:", age)
} else {
    fmt.Println("Không có")
}
```

`ok` là `bool` — `true` nếu key tồn tại.

## Iterate

```go
for k, v := range m {
    fmt.Println(k, v)
}

for k := range m { ... }    // chỉ key
```

⚠️ **Thứ tự duyệt KHÔNG xác định** — Go intentional randomize để code không phụ thuộc thứ tự.

Cần thứ tự ổn định → extract keys, sort:
```go
keys := make([]string, 0, len(m))
for k := range m {
    keys = append(keys, k)
}
sort.Strings(keys)
for _, k := range keys { /* ... */ }
```

## Kiểu key

Bất kỳ kiểu nào **comparable** (so sánh bằng `==`):
- ✅ `string`, số, `bool`, pointer, interface, struct (nếu mọi field comparable), array.
- ❌ `slice`, `map`, `function` — không comparable.

```go
type Point struct{ X, Y int }
m := map[Point]string{}  // OK
```

## Map với value phức tạp

```go
groups := map[string][]string{
    "frontend": {"Alice", "Bob"},
    "backend":  {"Charlie"},
}
groups["backend"] = append(groups["backend"], "Dave")
```

⚠️ Không thể sửa trực tiếp field của struct lưu trong map:
```go
people := map[string]Person{"a": {Name: "Alice"}}
people["a"].Name = "Alicia"  // ❌ cannot assign to struct field
```

Phải lấy ra, sửa, gán lại — hoặc lưu pointer (`map[string]*Person`).

## Concurrency

Map **KHÔNG safe** cho concurrent read+write. Race detector sẽ panic.
- Read-only đa luồng: OK.
- Có write → cần `sync.Mutex` hoặc dùng `sync.Map`.

## Ghi nhớ

- Pass map vào hàm = pass reference (cheap). Sửa trong hàm = sửa caller's map.
- Map không có capacity API như slice — không có `cap()`.
- Để clear map: `clear(m)` (Go 1.21+) hoặc tạo map mới.
