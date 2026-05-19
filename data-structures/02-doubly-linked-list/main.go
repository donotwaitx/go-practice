// 02 - Doubly Linked List
// Mỗi node có 2 con trỏ: Prev và Next. Cho phép duyệt 2 chiều,
// xóa node tại chỗ O(1) nếu đã có pointer đến node đó.
package main

import "fmt"

type DNode struct {
	Val        int
	Prev, Next *DNode
}

type DoublyList struct {
	Head, Tail *DNode
	size       int
}

func (l *DoublyList) PushFront(v int) {
	n := &DNode{Val: v, Next: l.Head}
	if l.Head != nil {
		l.Head.Prev = n
	} else {
		l.Tail = n
	}
	l.Head = n
	l.size++
}

func (l *DoublyList) PushBack(v int) {
	n := &DNode{Val: v, Prev: l.Tail}
	if l.Tail != nil {
		l.Tail.Next = n
	} else {
		l.Head = n
	}
	l.Tail = n
	l.size++
}

// PopFront O(1)
func (l *DoublyList) PopFront() (int, bool) {
	if l.Head == nil {
		return 0, false
	}
	v := l.Head.Val
	l.Head = l.Head.Next
	if l.Head != nil {
		l.Head.Prev = nil
	} else {
		l.Tail = nil
	}
	l.size--
	return v, true
}

// PopBack O(1)
func (l *DoublyList) PopBack() (int, bool) {
	if l.Tail == nil {
		return 0, false
	}
	v := l.Tail.Val
	l.Tail = l.Tail.Prev
	if l.Tail != nil {
		l.Tail.Next = nil
	} else {
		l.Head = nil
	}
	l.size--
	return v, true
}

func (l *DoublyList) Len() int { return l.size }

func (l *DoublyList) Forward() string {
	s := "head"
	for c := l.Head; c != nil; c = c.Next {
		s += fmt.Sprintf(" <-> %d", c.Val)
	}
	return s
}

func (l *DoublyList) Backward() string {
	s := "tail"
	for c := l.Tail; c != nil; c = c.Prev {
		s += fmt.Sprintf(" <-> %d", c.Val)
	}
	return s
}

func main() {
	l := &DoublyList{}
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	l.PushFront(0)
	fmt.Println(l.Forward())  // head <-> 0 <-> 1 <-> 2 <-> 3
	fmt.Println(l.Backward()) // tail <-> 3 <-> 2 <-> 1 <-> 0

	v, _ := l.PopFront()
	fmt.Println("popped front:", v, "→", l.Forward())

	v, _ = l.PopBack()
	fmt.Println("popped back:", v, "→", l.Forward())
}
