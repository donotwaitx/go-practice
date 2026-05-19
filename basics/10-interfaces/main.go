// 10 - Interfaces
package main

import "fmt"

// Interface = tập hợp method
type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	W, H float64
}

func (r Rectangle) Area() float64      { return r.W * r.H }
func (r Rectangle) Perimeter() float64 { return 2 * (r.W + r.H) }

type Circle struct {
	R float64
}

func (c Circle) Area() float64      { return 3.14159 * c.R * c.R }
func (c Circle) Perimeter() float64 { return 2 * 3.14159 * c.R }

// Hàm nhận interface — polymorphism
func describe(s Shape) {
	fmt.Printf("area=%.2f, perimeter=%.2f\n", s.Area(), s.Perimeter())
}

// Empty interface — chứa bất kỳ kiểu nào (giống any trong Go 1.18+)
func printAny(v any) {
	fmt.Printf("type=%T value=%v\n", v, v)
}

// Type assertion & type switch
func inspect(v any) {
	switch x := v.(type) {
	case int:
		fmt.Println("int gấp đôi:", x*2)
	case string:
		fmt.Println("string upper:", x)
	case Shape:
		fmt.Println("shape area:", x.Area())
	default:
		fmt.Println("unknown type")
	}
}

func main() {
	shapes := []Shape{
		Rectangle{W: 3, H: 4},
		Circle{R: 5},
	}
	for _, s := range shapes {
		describe(s)
	}

	printAny(42)
	printAny("hello")
	printAny(Rectangle{W: 1, H: 2})

	inspect(10)
	inspect("hi")
	inspect(Circle{R: 2})
}
