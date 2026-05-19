// 07 - Maps
package main

import "fmt"

func main() {
	// Khai báo & khởi tạo
	ages := map[string]int{
		"Alice": 30,
		"Bob":   25,
	}

	// Thêm/cập nhật
	ages["Charlie"] = 28
	ages["Alice"] = 31

	// Đọc
	fmt.Println("Alice:", ages["Alice"])

	// Comma-ok idiom: kiểm tra key tồn tại
	if age, ok := ages["David"]; ok {
		fmt.Println("David:", age)
	} else {
		fmt.Println("Không có David")
	}

	// Xóa
	delete(ages, "Bob")

	// Duyệt — KHÔNG đảm bảo thứ tự
	for name, age := range ages {
		fmt.Printf("%s -> %d\n", name, age)
	}

	// Map cần được make trước khi dùng (nếu khai báo var)
	var m map[string]int     // nil map — KHÔNG ghi được
	m = make(map[string]int) // giờ ghi được
	m["x"] = 1
	fmt.Println("m:", m, "len:", len(m))

	// Map với value là struct/slice
	groups := map[string][]string{
		"frontend": {"Alice", "Bob"},
		"backend":  {"Charlie"},
	}
	groups["backend"] = append(groups["backend"], "Dave")
	fmt.Println(groups)
}
