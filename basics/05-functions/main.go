// 05 - Functions
package main

import "fmt"

// Hàm cơ bản
func add(a, b int) int {
	return a + b
}

// Nhiều giá trị trả về (rất phổ biến trong Go cho error handling)
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("không thể chia cho 0")
	}
	return a / b, nil // nil = không lỗi
}

// Named return values
func minmax(nums []int) (min, max int) {
	min, max = nums[0], nums[0]
	for _, n := range nums {
		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
	}
	return // naked return — trả về min, max
}

// naked return chỉ nên dùng cho hàm ngắn, đơn giản. Với hàm phức tạp, tốt hơn là trả về giá trị rõ ràng để dễ đọc.

// Variadic — số tham số biến đổi
// ...int nghĩa là sum có thể nhận 0 hoặc nhiều int, và chúng sẽ được gói thành slice nums []int
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// Function as value & closure
func makeCounter() func() int {
	count := 0
	// Closure giữ trạng thái: hàm trả về vẫn có thể truy cập và sửa đổi biến count, dù makeCounter đã kết thúc
	return func() int {
		count++
		return count
	}
}

func main() {
	fmt.Println("add:", add(2, 3))

	if q, err := divide(10, 3); err != nil {
		fmt.Println("err:", err)
	} else {
		fmt.Printf("10/3 = %.2f\n", q)
	}

	mn, mx := minmax([]int{4, 1, 9, 2, 7})
	fmt.Println("min, max:", mn, mx)

	fmt.Println("sum:", sum(1, 2, 3, 4, 5))

	// Closure giữ trạng thái
	counter := makeCounter() // counter là một closure, nó giữ tham chiếu đến biến count trong makeCounter 
	// Mỗi lần gọi counter(), nó sẽ tăng count và trả về giá trị mới, dù makeCounter đã kết thúc
	fmt.Println(counter(), counter(), counter()) // 1 2 3

	// defer — chạy khi hàm kết thúc, LIFO
	// Dùng defer để đảm bảo tài nguyên được giải phóng, file được đóng, mutex được unlock, v.v.
	// Các defer sẽ chạy theo thứ tự ngược lại với thứ tự chúng được gọi (Last In, First Out)
	defer fmt.Println("defer 1 (chạy cuối)")
	defer fmt.Println("defer 2")
	fmt.Println("main kết thúc")
}
