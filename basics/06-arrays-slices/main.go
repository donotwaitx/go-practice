// 06 - Arrays & Slices
package main

import "fmt"

func main() {
	// ARRAY — kích thước cố định, là một phần của kiểu (length is type)
	var arr [3]int
	arr[0] = 10
	arr[1] = 20
	arr[2] = 30
	fmt.Println("array:", arr, "len:", len(arr))

	arr2 := [3]string{"a", "b", "c"}
	fmt.Println("array2:", arr2)

	// SLICE — kích thước động, đây mới là cấu trúc thường dùng
	s := []int{1, 2, 3}
	s = append(s, 4, 5)
	fmt.Println("slice:", s, "len:", len(s), "cap:", cap(s))

	// make([]T, len, cap)
	buf := make([]byte, 0, 10)
	buf = append(buf, 'H', 'i')
	fmt.Println("buf:", string(buf), "cap:", cap(buf))

	// Slicing: s[low:high] — high không bao gồm
	a := []int{10, 20, 30, 40, 50}
	fmt.Println(a[1:4]) // [20 30 40]
	fmt.Println(a[:2])  // [10 20]
	fmt.Println(a[3:])  // [40 50]

	// CHÚ Ý: slice share underlying array
	b := a[1:3]
	b[0] = 999
	fmt.Println("a sau khi đổi b:", a) // a[1] cũng bị đổi

	// copy để tránh chia sẻ
	cp := make([]int, len(a))
	copy(cp, a)
	cp[0] = -1
	fmt.Println("a:", a, "cp:", cp)

	// Iterate
	for i, v := range a {
		fmt.Printf("a[%d]=%d ", i, v)
	}
	fmt.Println()
}
