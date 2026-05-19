// 03 - Stack (LIFO)
// Last In, First Out. Implement bằng slice — Push/Pop O(1) amortized.
// Dùng generics để stack chứa được bất kỳ kiểu nào.
package main

import "fmt"

type Stack[T any] struct {
	data []T
}

func (s *Stack[T]) Push(v T) {
	s.data = append(s.data, v)
}

func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.data) == 0 {
		return zero, false
	}
	n := len(s.data) - 1
	v := s.data[n]
	s.data[n] = zero // tránh giữ reference, giúp GC
	s.data = s.data[:n]
	return v, true
}

func (s *Stack[T]) Peek() (T, bool) {
	var zero T
	if len(s.data) == 0 {
		return zero, false
	}
	return s.data[len(s.data)-1], true
}

func (s *Stack[T]) Len() int      { return len(s.data) }
func (s *Stack[T]) IsEmpty() bool { return len(s.data) == 0 }

// Ứng dụng: kiểm tra dấu ngoặc cân bằng
func balanced(expr string) bool {
	s := &Stack[rune]{}
	pairs := map[rune]rune{')': '(', ']': '[', '}': '{'}
	for _, c := range expr {
		switch c {
		case '(', '[', '{':
			s.Push(c)
		case ')', ']', '}':
			top, ok := s.Pop()
			if !ok || top != pairs[c] {
				return false
			}
		}
	}
	return s.IsEmpty()
}

func main() {
	s := &Stack[int]{}
	s.Push(1)
	s.Push(2)
	s.Push(3)
	fmt.Println("len:", s.Len()) // 3

	top, _ := s.Peek()
	fmt.Println("peek:", top) // 3

	for !s.IsEmpty() {
		v, _ := s.Pop()
		fmt.Println("pop:", v) // 3, 2, 1
	}

	// Stack of strings
	ss := &Stack[string]{}
	ss.Push("hello")
	ss.Push("world")
	v, _ := ss.Pop()
	fmt.Println("string pop:", v)

	// Ứng dụng
	fmt.Println("balanced '({[]})':", balanced("({[]})")) // true
	fmt.Println("balanced '([)]':", balanced("([)]"))     // false
}
