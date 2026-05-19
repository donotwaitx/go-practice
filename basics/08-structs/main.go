// 08 - Structs & Methods
package main

import "fmt"

// Struct — tập hợp các field có tên
type User struct {
	ID    int
	Name  string
	Email string
}

// Struct lồng (embedded)
type Employee struct {
	User     // embedded — Employee có thể truy cập trực tiếp Name, Email
	Salary   float64
	Position string
}

// Method với value receiver — KHÔNG thay đổi original
func (u User) Greet() string {
	return "Hi, I'm " + u.Name
}

// Method với pointer receiver — CÓ thay đổi original
func (u *User) Rename(newName string) {
	u.Name = newName
}

func main() {
	// Khởi tạo
	u1 := User{ID: 1, Name: "Alice", Email: "a@x.com"}
	u2 := User{2, "Bob", "b@x.com"} // theo thứ tự field — không khuyến khích

	fmt.Println(u1.Greet())
	fmt.Println(u2.Greet())

	// Pointer receiver thay đổi struct
	u1.Rename("Alice Smith") // Go tự lấy &u1
	fmt.Println("renamed:", u1.Name)

	// Embedded struct
	emp := Employee{
		User:     User{ID: 10, Name: "Charlie", Email: "c@x.com"},
		Salary:   1000,
		Position: "Engineer",
	}
	fmt.Println(emp.Greet()) // promoted method
	fmt.Println(emp.Name, emp.Position)

	// Anonymous struct
	point := struct {
		X, Y int
	}{X: 3, Y: 4}
	fmt.Println("point:", point)
}
