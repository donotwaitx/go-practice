# 06 – Arrays & Slices

Khác biệt giữa array (kích thước cố định) và slice (kích thước động) — bài quan trọng nhất phải hiểu kỹ.

## Chạy

```bash
go run ./basics/06-arrays-slices
```

## Array — `[N]T`

```go
var arr [3]int          // [0 0 0]
arr2 := [3]string{"a", "b", "c"}
arr3 := [...]int{1, 2, 3}  // compiler tự đếm size
```

Đặc điểm:
- **Size là một phần của type:** `[3]int` và `[5]int` là 2 kiểu KHÁC NHAU. Không cast được.
- **Value semantics:** gán array hoặc truyền vào hàm = **copy toàn bộ**. Tốn memory với array lớn.

→ Trong code Go thực tế, **hiếm khi dùng array trực tiếp**. Hầu như luôn dùng slice.

## Slice — `[]T`

Đây là cấu trúc bạn sẽ dùng 99% thời gian.

### Khái niệm

Slice là một **view động** trên một underlying array. Mỗi slice gồm 3 thành phần:
- **pointer** → trỏ đến phần tử đầu trong array gốc
- **length** → số phần tử slice "thấy"
- **capacity** → số phần tử array gốc còn lại từ pointer

```
underlying array: [10, 20, 30, 40, 50]
slice s = arr[1:3]:
   ptr ──┐
         ▼
        [20, 30]    len=2, cap=4 (từ index 1 đến hết)
```

### Khởi tạo

```go
s := []int{1, 2, 3}              // literal
s := make([]int, 5)              // length=5, cap=5, all zero
s := make([]int, 0, 10)          // length=0, cap=10
var s []int                      // nil slice — len=0, cap=0, nhưng append được
```

### `append`

```go
s = append(s, 4)        // thêm 1
s = append(s, 5, 6, 7)  // thêm nhiều
s = append(s, other...) // append cả slice khác
```

⚠️ **`append` có thể trả về slice mới** nếu cap không đủ:
- Cap đủ → ghi tại chỗ, trả về cùng slice.
- Cap không đủ → cấp phát array mới (thường gấp đôi), copy data, trả về slice mới.

Luôn gán lại: `s = append(s, x)`.

### Slicing — `s[low:high]`

`high` **không bao gồm**. Mặc định `low=0`, `high=len`:

```go
a := []int{10, 20, 30, 40, 50}
a[1:4]  // [20 30 40]
a[:2]   // [10 20]
a[3:]   // [40 50]
a[:]    // toàn bộ
```

### ⚠️ Pitfall: Slice share underlying array

```go
a := []int{1, 2, 3, 4, 5}
b := a[1:3]    // [2 3]
b[0] = 999
fmt.Println(a) // [1 999 3 4 5]  ← a cũng bị đổi!
```

Sửa slice con sửa luôn slice gốc vì chúng share array.

**Giải pháp:** `copy()` để clone:

```go
cp := make([]int, len(a))
copy(cp, a)
cp[0] = -1     // không ảnh hưởng a
```

### `len` vs `cap`

```go
s := make([]int, 3, 10)
len(s)  // 3 — số phần tử có thể truy cập (s[0]..s[2])
cap(s)  // 10 — kích thước underlying array
```

`append` chỉ realloc khi `len == cap`.

### Iterate

```go
for i, v := range a { ... }
for _, v := range a { ... }
for i := range a { ... }
```

## Khi nào dùng array thật?

- Cần kích thước fixed compile-time (cryptography, hash).
- Performance-critical, muốn tránh indirection của slice header.
- Map key (slice không hash được, array hash được).

99% còn lại → slice.

## Ghi nhớ

- Truyền slice vào hàm = truyền slice header (pointer + len + cap), KHÔNG copy data.
- Hàm có thể sửa phần tử của slice gốc thông qua slice tham số.
- Nhưng `append` trong hàm có thể tạo slice mới → caller không thấy.
