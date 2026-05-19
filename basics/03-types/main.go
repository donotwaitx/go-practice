// 03 - Basic Types & Conversion
package main

import (
	"fmt"
	"strconv"
)

func main() {
	// Số nguyên: int, int8/16/32/64, uint...
	var a int = 42
	var b int64 = 9_000_000_000

	var c = -9223372036854775808 // type inference, Go tự suy luận kiểu từ giá trị

	// Số thực
	var f float64 = 3.14

	// Chuỗi và ký tự (rune = int32, byte = uint8)
	s := "Xin chào"
	r := 'A'
	by := byte('B')

	// Boolean
	t, fls := true, false

	// Conversion — Go KHÔNG tự ép kiểu, phải convert tường minh
	var i int = 100
	var fl float64 = float64(i)
	var u uint = uint(fl)

	// String <-> Number
	n, _ := strconv.Atoi("123") // string -> int
	str := strconv.Itoa(456)    // int -> string
	fl2, _ := strconv.ParseFloat("2.71", 64)

	fmt.Println(a, b, c, f)
	fmt.Println(s, "len bytes =", len(s), "first rune =", r, "byte =", by)
	fmt.Println(t, fls)
	fmt.Println("convert:", i, fl, u)
	fmt.Println("parsed:", n, str, fl2)
}
