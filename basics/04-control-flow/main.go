// 04 - Control Flow: if / for / switch
package main

import "fmt"

func main() {
	// if — có thể khai báo biến ngay trong điều kiện
	if n := 10; n%2 == 0 {
		fmt.Println(n, "là số chẵn")
	} else {
		fmt.Println(n, "là số lẻ")
	}

	// for — Go chỉ có for (không có while, do-while)
	// 1) classic for
	sum := 0
	for i := 1; i <= 5; i++ {
		sum += i
	}
	fmt.Println("sum 1..5 =", sum)

	// 2) while-style
	count := 0
	for count < 3 {
		count++
	}
	fmt.Println("count =", count)

	// 3) infinite loop với break
	k := 0
	for {
		if k == 2 {
			break
		}
		k++
	}

	// 4) range
	nums := []int{10, 20, 30}
	for idx, val := range nums {
		fmt.Printf("nums[%d] = %d\n", idx, val)
	}

	// switch — không cần break, không fallthrough mặc định
	day := "Mon"
	switch day {
	case "Sat", "Sun":
		fmt.Println("Cuối tuần")
	case "Mon":
		fmt.Println("Đầu tuần")
	default:
		fmt.Println("Ngày thường")
	}

	// switch không có điều kiện = if/else if chuỗi
	score := 85
	switch {
	case score >= 90:
		fmt.Println("A")
	case score >= 80:
		fmt.Println("B")
	default:
		fmt.Println("C")
	}
}
