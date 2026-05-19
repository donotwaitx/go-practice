# Go Basics

Lộ trình học foundation của Go. Mỗi thư mục là một chủ đề, có file `main.go` chạy được độc lập.

## Lộ trình

| # | Chủ đề | Nội dung |
|---|--------|----------|
| 01 | hello | Cấu trúc chương trình, `package main`, `fmt` |
| 02 | variables | Biến, hằng, zero values, short declaration |
| 03 | types | Kiểu cơ bản, conversion, `strconv` |
| 04 | control-flow | `if`, `for`, `switch`, `range` |
| 05 | functions | Multiple return, variadic, closure, `defer` |
| 06 | arrays-slices | Array vs slice, `append`, `make`, `copy`, slicing |
| 07 | maps | Khởi tạo, comma-ok, `delete`, iterate |
| 08 | structs | Field, method, embedded struct, pointer receiver |
| 09 | pointers | `&`, `*`, `new`, pointer vs value semantics |
| 10 | interfaces | Polymorphism, empty interface (`any`), type assertion, type switch |
| 11 | errors | `errors.New`, custom error, `errors.Is/As`, wrapping, `panic`/`recover` |
| 12 | goroutines | `go`, `sync.WaitGroup`, channels, buffered channel, `select` |

## Cách chạy

Từ thư mục gốc của repo:

```bash
go run ./basics/01-hello
go run ./basics/02-variables
# ...
```

Hoặc `cd` vào từng thư mục rồi `go run .`

## Lệnh hữu ích

```bash
go run .          # chạy package hiện tại
go build .        # biên dịch ra binary
go fmt ./...      # format toàn bộ
go vet ./...      # static analysis
go test ./...     # chạy test
```
