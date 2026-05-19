# 04 – Control Flow

`if`, `for`, `switch` — Go đơn giản hóa tối đa, không có `while`, `do-while`, `ternary`.

## Chạy

```bash
go run ./basics/04-control-flow
```

## `if`

```go
if x > 0 {
    // ...
} else if x == 0 {
    // ...
} else {
    // ...
}
```

**Đặc biệt:** có thể khai báo biến ngay trong điều kiện — biến chỉ sống trong `if/else` block:

```go
if v, err := doSomething(); err == nil {
    use(v)
} else {
    log(err)
}
// v, err không tồn tại ngoài này
```

Pattern này CỰC kỳ phổ biến trong Go.

⚠️ Không có ternary (`a ? b : c`). Phải dùng `if` đầy đủ. Đây là intentional — code rõ hơn.

## `for` — vòng lặp duy nhất

### 1. Classic (như C)

```go
for i := 0; i < 10; i++ {
    // ...
}
```

### 2. While-style

```go
for x < 100 {
    x *= 2
}
```

### 3. Infinite

```go
for {
    if done {
        break
    }
}
```

### 4. Range

Duyệt slice, array, string, map, channel:

```go
for i, v := range slice {}     // i = index, v = value
for k, v := range myMap {}     // k = key, v = value
for i, r := range "abc" {}     // i = byte index, r = rune
for v := range channel {}      // chỉ value, dừng khi channel close

for _, v := range slice {}     // bỏ index
for range n {}                 // Go 1.22+: lặp n lần
```

⚠️ Khi range, `v` là **copy** — sửa `v` KHÔNG sửa phần tử gốc.

## `switch`

### Switch theo giá trị

```go
switch day {
case "Sat", "Sun":         // nhiều giá trị một case
    fmt.Println("weekend")
case "Mon":
    fmt.Println("Monday")
default:
    fmt.Println("weekday")
}
```

**Khác C/Java:**
- Không cần `break` — chạy xong case là thoát.
- Muốn rơi xuống case sau: `fallthrough` (hiếm dùng).

### Switch không điều kiện = if/else if chain

```go
switch {
case score >= 90:
    return "A"
case score >= 80:
    return "B"
default:
    return "F"
}
```

Gọn hơn nhiều so với `if/else if/else`.

### Switch theo type (type switch)

```go
switch v := x.(type) {
case int:
    // v là int
case string:
    // v là string
}
```

Xem bài 10 (interfaces).

## `break`, `continue`, `goto`

- `break` thoát loop/switch.
- `continue` sang iteration tiếp.
- `goto label` — tồn tại nhưng hầu như không dùng.
- **Labeled break:** thoát nhiều loop lồng nhau.

```go
outer:
for i := 0; i < 10; i++ {
    for j := 0; j < 10; j++ {
        if cond {
            break outer  // thoát cả 2 loop
        }
    }
}
```
