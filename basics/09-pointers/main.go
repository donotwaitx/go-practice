// 09 - Pointers
package main

import "fmt"

type Counter struct {
	value int
}

// Pointer receiver — modify struct
func (c *Counter) Inc() {
	c.value++
}

// Hàm nhận pointer — modify caller's variable
func double(n *int) {
	*n = *n * 2
}

func main() {
	x := 10
	p := &x // p là pointer trỏ đến x
	fmt.Println("x:", x, "p:", p, "*p:", *p)

	*p = 20 // sửa qua pointer
	fmt.Println("x sau:", x)

	double(&x)
	fmt.Println("after double:", x)

	// nil pointer
	var pp *int
	fmt.Println("nil pointer:", pp)
	// fmt.Println(*pp) // PANIC: nil pointer dereference

	// new() trả về pointer đến zero value
	n := new(int)
	*n = 42
	fmt.Println("new int:", *n)

	// Pointer to struct
	c := &Counter{}
	c.Inc()
	c.Inc()
	c.Inc()
	fmt.Println("counter:", c.value) // 3

	// Go KHÔNG có pointer arithmetic như C
	// p++ không hợp lệ với pointer
}
