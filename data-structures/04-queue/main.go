// 04 - Queue (FIFO)
// First In, First Out. Có 2 cách implement:
//  1. Slice đơn giản: Enqueue O(1), Dequeue O(n) (phải shift)
//  2. Slice với head index: O(1) amortized cả hai
//
// Bài này dùng cách 2 cho hiệu năng tốt hơn.
package main

import "fmt"

type Queue[T any] struct {
	data []T
	head int // index của phần tử đầu queue trong data
}

func (q *Queue[T]) Enqueue(v T) {
	q.data = append(q.data, v)
}

func (q *Queue[T]) Dequeue() (T, bool) {
	var zero T
	if q.head >= len(q.data) {
		return zero, false
	}
	v := q.data[q.head]
	q.data[q.head] = zero // giúp GC giải phóng
	q.head++
	// Nén lại nếu head quá lớn để tránh slice phình mãi
	if q.head > 64 && q.head*2 > len(q.data) {
		q.data = append([]T{}, q.data[q.head:]...)
		q.head = 0
	}
	return v, true
}

func (q *Queue[T]) Peek() (T, bool) {
	var zero T
	if q.head >= len(q.data) {
		return zero, false
	}
	return q.data[q.head], true
}

func (q *Queue[T]) Len() int      { return len(q.data) - q.head }
func (q *Queue[T]) IsEmpty() bool { return q.head >= len(q.data) }

func main() {
	q := &Queue[int]{}
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	fmt.Println("len:", q.Len()) // 3

	top, _ := q.Peek()
	fmt.Println("peek:", top) // 1

	for !q.IsEmpty() {
		v, _ := q.Dequeue()
		fmt.Println("dequeue:", v) // 1, 2, 3
	}

	// Queue of strings
	qs := &Queue[string]{}
	qs.Enqueue("first")
	qs.Enqueue("second")
	v, _ := qs.Dequeue()
	fmt.Println("first out:", v) // first
}
