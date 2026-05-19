# 05 – Functions

Functions là first-class citizen trong Go.

## Chạy

```bash
go run ./basics/05-functions
```

## Syntax cơ bản

```go
func add(a int, b int) int {
    return a + b
}

// Cùng type → gộp tham số
func add(a, b int) int { ... }
```

Cú pháp: `func tên(tham số) kiểu_trả_về { ... }`. Type đặt SAU tên biến (ngược với C/Java) — vì Go đọc trái-sang-phải tự nhiên.

## Multiple return values

Pattern then chốt của Go — đặc biệt cho error handling:

```go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("không thể chia cho 0")
    }
    return a / b, nil
}

result, err := divide(10, 3)
if err != nil { /* handle */ }
```

⚠️ Phải nhận đủ giá trị trả về. Bỏ qua bằng `_`:
```go
result, _ := divide(10, 3)
```

## Named return values

```go
func minmax(nums []int) (min, max int) {
    min, max = nums[0], nums[0]
    // ...
    return  // "naked return" — trả về min, max
}
```

`min`, `max` được khai báo sẵn, init về zero value. `return` trần trả về chúng.

**Khi nào dùng:** hàm ngắn, tên biến tự documented. Hàm dài → nên `return value` rõ ràng để dễ đọc.

## Variadic functions

```go
func sum(nums ...int) int {
    total := 0
    for _, n := range nums {  // nums là []int
        total += n
    }
    return total
}

sum(1, 2, 3)        // 6
sum()               // 0
sum([]int{1,2,3}...) // expand slice
```

`...int` ở tham số = nhận 0 đến N int, gói thành `[]int`.

## Function as value

Hàm là một giá trị → gán cho biến, truyền vào hàm khác, trả về:

```go
add := func(a, b int) int { return a + b }
result := add(2, 3)

// Function type
var op func(int, int) int = add
```

## Closure

Hàm có thể "capture" biến từ scope bên ngoài và giữ nó:

```go
func makeCounter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}

c := makeCounter()
c() // 1
c() // 2
c() // 3
```

Mỗi `makeCounter()` tạo `count` riêng. Closure giữ tham chiếu, không phải copy → mỗi counter độc lập.

## `defer`

Hoãn lệnh đến khi **hàm cha return**. Thường dùng cho cleanup:

```go
func readFile() {
    f, _ := os.Open("data.txt")
    defer f.Close()   // tự đóng khi hàm xong, kể cả khi panic
    // ...
}
```

**Quy tắc:**
1. `defer` chạy theo **LIFO** (Last In First Out).
2. Tham số `defer` được evaluate NGAY khi gặp `defer`, không phải lúc chạy:
   ```go
   x := 10
   defer fmt.Println(x)  // sẽ in 10
   x = 20                // không ảnh hưởng
   ```
3. Pattern phổ biến: `defer mu.Unlock()` sau `mu.Lock()`, `defer rows.Close()` sau query.

## Recursion

Bình thường — Go hỗ trợ đầy đủ. Không có tail-call optimization, nên đệ quy sâu sẽ stack overflow.

## Ghi nhớ

- Không có overloading (cùng tên, khác signature) — không có.
- Không có default parameter — phải dùng variadic hoặc struct.
- Không có generics method trên type (vẫn có generics function/type từ 1.18).
