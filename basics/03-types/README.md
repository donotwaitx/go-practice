# 03 – Types & Conversion

Kiểu dữ liệu cơ bản và cách chuyển đổi giữa chúng.

## Chạy

```bash
go run ./basics/03-types
```

## Số nguyên

| Type | Range |
|------|-------|
| `int8` / `uint8` | -128..127 / 0..255 |
| `int16` / `uint16` | -32768..32767 / 0..65535 |
| `int32` / `uint32` | ±2.1 tỷ / 0..4.3 tỷ |
| `int64` / `uint64` | ±9.2 quintillion |
| `int` / `uint` | **32 hoặc 64 bit tùy platform** (thường 64-bit) |
| `byte` | alias của `uint8` |
| `rune` | alias của `int32`, đại diện Unicode code point |

Mặc định khi viết `42` → kiểu `int`. Dùng `_` để dễ đọc số lớn: `9_000_000_000`.

## Số thực

- `float32`, `float64` — IEEE 754.
- Mặc định khi viết `3.14` → `float64`.
- ⚠️ So sánh float bằng `==` không đáng tin (lỗi precision). Nên so sánh với `math.Abs(a-b) < epsilon`.

## String

- UTF-8 encoded, **immutable**.
- `len(s)` trả về số **byte**, KHÔNG phải số ký tự. `len("Xin chào")` = 9, không phải 7.
- Để đếm số rune: `utf8.RuneCountInString(s)`.
- Index: `s[i]` trả về **byte**, không phải ký tự.
- Duyệt ký tự đúng cách: `for i, r := range s` — `r` là rune.

## Rune & Byte

```go
r := 'A'        // rune, type là int32, value 65
b := byte('B')  // byte, type là uint8, value 66
```

Dấu nháy đơn `' '` → rune. Dấu nháy kép `" "` → string.

## Boolean

`true` / `false`. Không tồn tại "truthy/falsy" như JS — `if 1 {}` ❌ compile error.

## Conversion — PHẢI tường minh

Go KHÔNG tự ép kiểu, ngay cả giữa các kiểu số:

```go
var i int = 100
var f float64 = float64(i)  // ✅
var f float64 = i           // ❌ cannot use i (type int) as float64
```

Syntax: `T(value)`.

⚠️ Conversion có thể MẤT dữ liệu:
```go
var f float64 = 3.9
var i int = int(f)   // i = 3 (truncate, không round)

var big int = 300
var b byte = byte(big)  // b = 44 (overflow uint8)
```

## String ↔ Number

Dùng package `strconv`:

```go
n, err := strconv.Atoi("123")          // string → int
s := strconv.Itoa(456)                 // int → string
f, err := strconv.ParseFloat("2.71", 64)
s := strconv.FormatFloat(2.71, 'f', 2, 64)
```

⚠️ **Đừng dùng `string(65)` để chuyển int → string** — nó trả về rune ký tự (`"A"`), không phải `"65"`. Go vet sẽ warn về việc này.

## Kiểm tra type runtime

```go
fmt.Printf("%T\n", x)   // in tên kiểu
```
