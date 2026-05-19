# 01 – Hello World

Bài đầu tiên: cấu trúc tối thiểu của một chương trình Go.

## Chạy

```bash
go run ./basics/01-hello
```

Output:
```
Hello, Go!
Đây là Go, phiên bản 1
```

## Khái niệm

### `package main`

Mọi file Go đều thuộc về một **package**. Package có tên đặc biệt `main` báo cho Go biết đây là **chương trình thực thi** (executable), không phải thư viện.

- `package main` → tạo binary chạy được.
- `package <bất kỳ tên khác>` → thư viện, không chạy độc lập.

### `import "fmt"`

Import package từ thư viện chuẩn. `fmt` chứa các hàm format I/O (Print, Println, Printf, Sprintf...).

Nhiều import:
```go
import (
    "fmt"
    "os"
    "strings"
)
```

### `func main()`

**Entrypoint** — hàm Go chạy khi khởi động chương trình. Bắt buộc trong package `main`, không nhận tham số, không trả giá trị.

### `fmt.Println` vs `fmt.Printf`

| Hàm | Khi dùng |
|-----|----------|
| `Println(a, b, c)` | In các giá trị, ngăn cách bằng space, tự xuống dòng. |
| `Printf("...", args)` | Có verb format: `%s` (string), `%d` (int), `%v` (giá trị mặc định), `%T` (type), `%f` (float). KHÔNG tự xuống dòng — phải `\n`. |
| `Print(...)` | Như `Println` nhưng không xuống dòng. |
| `Sprintf(...)` | Trả về string thay vì in. |

## Ghi nhớ

- Tên file không cần trùng tên package.
- Một thư mục chỉ có MỘT package (trừ file `_test.go`).
- Go cực kỳ strict: import thừa hoặc biến thừa → **compile error**, không phải warning.
