# go-practice

Repo học Go từ foundation. Code mẫu nằm trong `basics/`, mỗi chủ đề một file `main.go` chạy được độc lập.

## Mục lục

- [Cấu trúc dự án](#cấu-trúc-dự-án)
- [Giải thích `go.mod`](#giải-thích-gomod)
- [Cách chạy](#cách-chạy)
- [Nội dung từng bài](#nội-dung-từng-bài)

## Cấu trúc dự án

```
go-practice/
├── go.mod                      # khai báo module
├── README.md                   # file này
└── basics/
    ├── README.md               # roadmap ngắn gọn
    ├── 01-hello/main.go
    ├── 02-variables/main.go
    ├── 03-types/main.go
    ├── 04-control-flow/main.go
    ├── 05-functions/main.go
    ├── 06-arrays-slices/main.go
    ├── 07-maps/main.go
    ├── 08-structs/main.go
    ├── 09-pointers/main.go
    ├── 10-interfaces/main.go
    ├── 11-errors/main.go
    └── 12-goroutines/main.go
```

## Giải thích `go.mod`

`go.mod` là file mô tả **module** trong Go — đơn vị quản lý phụ thuộc (dependency). Bất kỳ project Go nào ngoài `$GOPATH` đều cần file này.

Nội dung file của repo:

```go
module github.com/donotwaitx/go-practice

go 1.26.3
```

Ý nghĩa từng dòng:

| Dòng | Ý nghĩa |
|------|---------|
| `module github.com/donotwaitx/go-practice` | **Module path** — định danh module. Đây cũng là prefix khi import nội bộ. Ví dụ nếu có package `basics/utils`, sẽ import bằng `github.com/donotwaitx/go-practice/basics/utils`. Theo convention, dùng URL Git repo. |
| `go 1.26.3` | **Phiên bản Go tối thiểu** để build module. Go toolchain dựa vào dòng này để bật/tắt feature mới (ví dụ generics có từ 1.18, range-over-func có từ 1.23). |

### Các directive khác (chưa dùng trong repo này)

```go
require (
    github.com/some/lib v1.2.3       // dependency và phiên bản
    golang.org/x/sync v0.7.0
)

replace github.com/old/pkg => github.com/new/pkg v1.0.0   // thay thế dependency

exclude github.com/buggy/lib v0.1.0                       // loại trừ version

retract v0.0.5                                            // đánh dấu version cần tránh
```

### Các lệnh thường dùng với module

```bash
go mod init <path>      # khởi tạo go.mod
go mod tidy             # thêm thiếu / xóa thừa dependency, đồng bộ go.sum
go mod download         # tải dependency về cache
go get <pkg>@<ver>      # thêm/nâng cấp dependency
go mod why <pkg>        # giải thích vì sao có dependency này
go list -m all          # liệt kê toàn bộ module trong dependency tree
```

### `go.sum` là gì?

Khi có dependency, Go sẽ tự tạo thêm `go.sum` chứa **checksum** của từng version — đảm bảo build reproducible và phát hiện dependency bị giả mạo. Repo này hiện không có `go.sum` vì chưa import package bên ngoài.

## Cách chạy

Từ thư mục gốc repo:

```bash
go run ./basics/01-hello
go run ./basics/02-variables
# ... tương tự cho các bài còn lại
```

Hoặc `cd` vào thư mục từng bài rồi `go run .`

Các lệnh hữu ích khác:

```bash
go build ./basics/...     # biên dịch tất cả package
go vet ./basics/...       # static analysis tìm bug tiềm ẩn
go fmt ./...              # auto-format toàn bộ
go test ./...             # chạy test (chưa có test trong repo)
```

## Nội dung từng bài

### 01 – Hello World

**File:** `basics/01-hello/main.go`

Cấu trúc chương trình Go tối thiểu.

- `package main` — package đặc biệt báo cho Go biết đây là executable.
- `import "fmt"` — import package chuẩn để in ra console.
- `func main()` — entrypoint, mọi chương trình chạy được đều bắt đầu từ đây.
- `fmt.Println` vs `fmt.Printf` — `Println` in kèm xuống dòng, `Printf` dùng verb (`%s`, `%d`, `%v`...) như C.

### 02 – Variables & Constants

**File:** `basics/02-variables/main.go`

Cách khai báo biến và hằng số.

- `var name type = value` — khai báo đầy đủ.
- `var name = value` — bỏ type, để Go suy luận.
- `name := value` — **short declaration**, chỉ dùng được trong hàm.
- `const Pi = 3.14` — hằng số, không thay đổi được.
- **Zero values:** mọi biến khai báo mà không gán đều có giá trị mặc định (`int=0`, `string=""`, `bool=false`, `pointer=nil`). Khác với C/Java, Go không có "biến chưa khởi tạo".
- Khai báo nhiều biến bằng `var (...)` group.

### 03 – Types & Conversion

**File:** `basics/03-types/main.go`

Kiểu dữ liệu cơ bản và cách chuyển đổi.

- **Số nguyên:** `int`, `int8/16/32/64`, `uint`... Trên 64-bit machine, `int` là 64-bit.
- **Số thực:** `float32`, `float64` (mặc định khi dùng `:=` với số thập phân).
- **String:** UTF-8 immutable. `len(s)` trả về số **byte**, không phải số ký tự.
- **Rune:** alias của `int32`, đại diện một code point Unicode (`'A'`).
- **Byte:** alias của `uint8`.
- **Bool:** `true` / `false`.
- **Conversion phải tường minh:** Go không tự ép kiểu như JS hay Python — `float64(i)`, `int(f)`. Ép sai có thể mất dữ liệu (truncate).
- **`strconv`:** package chuẩn chuyển đổi string ↔ number (`Atoi`, `Itoa`, `ParseFloat`).

### 04 – Control Flow

**File:** `basics/04-control-flow/main.go`

Cấu trúc điều khiển.

- **`if`** có thể khai báo biến ngay trong điều kiện: `if v := f(); v > 0 { ... }` — biến chỉ tồn tại trong scope đó.
- **`for`** là vòng lặp DUY NHẤT trong Go (không có `while`, `do-while`). 4 dạng:
  1. `for i := 0; i < n; i++` — classic.
  2. `for cond` — như `while`.
  3. `for {}` — vô hạn, thoát bằng `break`.
  4. `for i, v := range collection` — duyệt slice/map/string/channel.
- **`switch`** không cần `break` (mặc định không fallthrough — ngược C/Java). Mỗi `case` chạy độc lập. Nhiều giá trị: `case "Sat", "Sun":`.
- **`switch` không điều kiện** = chuỗi `if/else if` — gọn hơn nhiều khi check range.

### 05 – Functions

**File:** `basics/05-functions/main.go`

- **Multiple return values** — pattern phổ biến: `value, err := f()`. Đây là cách Go xử lý lỗi.
- **Named returns:** `func minmax(...) (min, max int)` — biến trả về có sẵn, có thể dùng `return` trần.
- **Variadic:** `func sum(nums ...int)` — nhận 0-N tham số, gói thành slice.
- **First-class function:** hàm có thể gán cho biến, truyền vào hàm khác, trả về từ hàm.
- **Closure:** hàm có thể capture biến từ scope bên ngoài và giữ trạng thái qua nhiều lần gọi.
- **`defer`** hoãn thực thi đến khi hàm cha kết thúc. Thứ tự LIFO. Thường dùng để: đóng file, unlock mutex, giải phóng resource.

### 06 – Arrays & Slices

**File:** `basics/06-arrays-slices/main.go`

- **Array `[N]T`:** kích thước cố định, là một phần của type. `[3]int` và `[5]int` là 2 kiểu KHÁC NHAU. Truyền array vào hàm = copy toàn bộ. Hiếm dùng trực tiếp.
- **Slice `[]T`:** view động trên một mảng. Đây là cấu trúc bạn dùng 99% thời gian.
  - `make([]T, len, cap)` — tạo slice với length/capacity định trước.
  - `append(s, x)` — thêm phần tử, có thể trả về slice mới nếu phải mở rộng underlying array.
  - `s[low:high]` — slicing, KHÔNG copy mà share underlying array → sửa slice con có thể đổi slice gốc.
  - `copy(dst, src)` — copy thật sự.
- **Pitfall:** giữ tham chiếu slice sau khi gọi `append` có thể lead to bug khó debug khi cap thay đổi.

### 07 – Maps

**File:** `basics/07-maps/main.go`

- `map[K]V` — hash table.
- **Khai báo:** `m := map[string]int{}` hoặc `m := make(map[string]int)`. **Lưu ý:** `var m map[string]int` tạo nil map — đọc OK nhưng GHI sẽ panic.
- **Comma-ok idiom:** `v, ok := m[key]` — phân biệt "không có key" với "có key value zero".
- **`delete(m, key)`** xóa key.
- **`range` map KHÔNG đảm bảo thứ tự** — đây là intentional của Go để code không phụ thuộc thứ tự ngầm định.

### 08 – Structs & Methods

**File:** `basics/08-structs/main.go`

- **Struct:** tập hợp các field có tên. Go không có class — struct + method = OOP-lite.
- **Method:** hàm gắn với receiver type: `func (u User) Greet() string`.
- **Value receiver vs pointer receiver:**
  - Value: làm việc trên BẢN SAO — không sửa được struct gốc.
  - Pointer (`*User`): làm việc trên gốc — sửa được. Cũng tránh copy nếu struct lớn.
  - Quy tắc: cần modify hoặc struct lớn → pointer receiver.
- **Embedded struct:** nhúng struct này vào struct khác → field và method được "promote". Đây là cách Go làm composition thay cho inheritance.
- **Anonymous struct:** struct không tên, dùng tạm cho data tạm thời.

### 09 – Pointers

**File:** `basics/09-pointers/main.go`

- `&x` — lấy địa chỉ của `x`.
- `*p` — dereference, lấy giá trị tại địa chỉ `p` trỏ tới.
- `new(T)` — cấp phát zero value và trả về pointer.
- Pointer cho phép: sửa giá trị từ hàm khác, tránh copy struct lớn, biểu diễn "không có giá trị" (nil).
- **Go KHÔNG có pointer arithmetic** (`p++`, `p+1`) như C → an toàn hơn nhưng vẫn linh hoạt.
- **Pitfall:** dereference nil pointer → panic.

### 10 – Interfaces

**File:** `basics/10-interfaces/main.go`

- **Interface = tập hợp method.** Type nào có đủ các method đó tự động "implement" interface — không cần khai báo `implements` như Java. Gọi là **structural typing** (duck typing tĩnh).
- **Empty interface `interface{}` / `any`** chứa được bất kỳ kiểu nào (từ Go 1.18 dùng `any` cho rõ).
- **Type assertion:** `x, ok := v.(Shape)` — thử cast value sang concrete type.
- **Type switch:** `switch x := v.(type) { case int: ...; case string: ... }` — pattern phổ biến khi xử lý `any`.
- **Best practice:** giữ interface NHỎ (1-3 methods). Đặt interface ở chỗ DÙNG, không phải chỗ implement (consumer-defined interfaces).

### 11 – Errors

**File:** `basics/11-errors/main.go`

- Go dùng **explicit error return** thay vì exception. Mọi hàm có thể fail đều trả về thêm `error`.
- Pattern chuẩn:
  ```go
  result, err := doSomething()
  if err != nil { return err }
  ```
- **Sentinel error:** `var ErrNotFound = errors.New("not found")` — error có thể so sánh bằng `==` hoặc `errors.Is`.
- **Custom error type:** struct implement `Error() string`.
- **Error wrapping:** `fmt.Errorf("context: %w", err)` — giữ chain, dùng `errors.Is` / `errors.As` để truy ngược.
- **`panic` / `recover`:** chỉ cho lỗi bất thường (programmer error, invariant violation). KHÔNG dùng thay cho error.

### 12 – Goroutines & Channels

**File:** `basics/12-goroutines/main.go`

Concurrency là điểm mạnh nhất của Go.

- **Goroutine:** lightweight thread, tạo bằng `go f()`. Một chương trình Go có thể chạy hàng nghìn/triệu goroutine.
- **`sync.WaitGroup`** — chờ một nhóm goroutine xong: `Add(n)` → `Done()` mỗi khi xong → `Wait()` block đến khi đủ.
- **Channel:** kênh giao tiếp giữa goroutines — `ch := make(chan int)`.
  - Unbuffered: gửi block đến khi có người nhận, và ngược lại.
  - Buffered (`make(chan T, n)`): gửi không block đến khi buffer đầy.
  - `close(ch)` báo "không còn data". `range ch` tự dừng khi channel close.
- **`select`** — chờ nhiều channel cùng lúc. Có thể kèm `default` (non-blocking) hoặc `time.After` (timeout).
- **Triết lý:** *"Don't communicate by sharing memory; share memory by communicating."* — dùng channel thay vì shared state + lock khi có thể.

## Tài liệu tham khảo

- [Tour of Go](https://go.dev/tour/) — interactive tutorial chính thức
- [Effective Go](https://go.dev/doc/effective_go) — best practice
- [Go by Example](https://gobyexample.com/) — code mẫu theo chủ đề
- [Standard library reference](https://pkg.go.dev/std)
