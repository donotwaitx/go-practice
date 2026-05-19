// 01 - Singly Linked List
// Cấu trúc dữ liệu tuyến tính, mỗi node trỏ đến node kế tiếp.
package main

import "fmt"

type Node struct {
	Val  int
	Next *Node
}

type LinkedList struct {
	Head *Node
	size int
}

// Prepend O(1) — thêm đầu danh sách
func (l *LinkedList) Prepend(v int) {
	l.Head = &Node{Val: v, Next: l.Head}
	l.size++
}

// Append O(n) — thêm cuối danh sách
func (l *LinkedList) Append(v int) {
	n := &Node{Val: v}
	l.size++
	if l.Head == nil {
		l.Head = n
		return
	}
	cur := l.Head
	for cur.Next != nil {
		cur = cur.Next
	}
	cur.Next = n
}

// Remove O(n) — xóa node đầu tiên có giá trị v
func (l *LinkedList) Remove(v int) bool {
	if l.Head == nil {
		return false
	}
	if l.Head.Val == v {
		l.Head = l.Head.Next
		l.size--
		return true
	}
	cur := l.Head
	for cur.Next != nil && cur.Next.Val != v {
		cur = cur.Next
	}
	if cur.Next == nil {
		return false
	}
	cur.Next = cur.Next.Next
	l.size--
	return true
}

// Find O(n)
func (l *LinkedList) Find(v int) bool {
	for cur := l.Head; cur != nil; cur = cur.Next {
		if cur.Val == v {
			return true
		}
	}
	return false
}

// Reverse O(n) - đảo ngược danh sách tại chỗ
func (l *LinkedList) Reverse() {
	var prev *Node
	cur := l.Head
	for cur != nil {
		next := cur.Next
		cur.Next = prev
		prev = cur
		cur = next
	}
	l.Head = prev
}

func (l *LinkedList) Len() int { return l.size }

func (l *LinkedList) String() string {
	s := "["
	for cur := l.Head; cur != nil; cur = cur.Next {
		if cur != l.Head {
			s += " -> "
		}
		s += fmt.Sprintf("%d", cur.Val)
	}
	return s + "]"
}

func main() {
	l := &LinkedList{}
	l.Append(1)
	l.Append(2)
	l.Append(3)
	l.Prepend(0)
	fmt.Println("list:", l, "len:", l.Len()) // [0 -> 1 -> 2 -> 3]

	fmt.Println("find 2:", l.Find(2))
	fmt.Println("find 9:", l.Find(9))

	l.Remove(2)
	fmt.Println("sau remove 2:", l)

	l.Reverse()
	fmt.Println("reversed:", l)
}
