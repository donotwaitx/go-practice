// 02 - Variables & Constants
package main

import "fmt"

// Hằng số khai báo ngoài hàm
const Pi = 3.14159

// Biến cấp package
var appName = "go-practice"

func main() {
	// Khai báo đầy đủ
	var age int = 25
	// Suy luận kiểu
	var name = "Trong"
	// Short declaration — chỉ dùng trong hàm
	city := "Hanoi"

	// Khai báo nhiều biến
	var (
		x, y int    = 1, 2
		ok   bool   = true
		msg  string = "hello"
	)

	// Zero values: int=0, float=0.0, string="", bool=false, pointer=nil
	var z int
	var s string

	fmt.Println(appName, name, age, city)
	fmt.Println(x, y, ok, msg)
	fmt.Println("zero values:", z, fmt.Sprintf("%q", s))
	fmt.Println("Pi =", Pi)
}
