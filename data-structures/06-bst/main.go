// 06 - Binary Search Tree (BST)
// Cây nhị phân tìm kiếm: với mỗi node, mọi giá trị bên LEFT < node < mọi giá trị bên RIGHT.
// Insert/Search/Delete trung bình O(log n), worst case O(n) nếu cây bị skewed.
package main

import "fmt"

type BSTNode struct {
	Val         int
	Left, Right *BSTNode
}

type BST struct {
	Root *BSTNode
}

func (t *BST) Insert(v int) {
	t.Root = insert(t.Root, v)
}

func insert(n *BSTNode, v int) *BSTNode {
	if n == nil {
		return &BSTNode{Val: v}
	}
	if v < n.Val {
		n.Left = insert(n.Left, v)
	} else if v > n.Val {
		n.Right = insert(n.Right, v)
	}
	// trùng giá trị: bỏ qua (BST này không cho duplicate)
	return n
}

func (t *BST) Search(v int) bool {
	cur := t.Root
	for cur != nil {
		if v == cur.Val {
			return true
		}
		if v < cur.Val {
			cur = cur.Left
		} else {
			cur = cur.Right
		}
	}
	return false
}

// Delete — 3 case:
//  1. Node lá: xóa thẳng
//  2. Node có 1 con: thay bằng con đó
//  3. Node có 2 con: thay bằng node nhỏ nhất bên phải (in-order successor)
func (t *BST) Delete(v int) {
	t.Root = del(t.Root, v)
}

func del(n *BSTNode, v int) *BSTNode {
	if n == nil {
		return nil
	}
	if v < n.Val {
		n.Left = del(n.Left, v)
	} else if v > n.Val {
		n.Right = del(n.Right, v)
	} else {
		// Tìm thấy
		if n.Left == nil {
			return n.Right
		}
		if n.Right == nil {
			return n.Left
		}
		// 2 con: lấy in-order successor (min của Right subtree)
		succ := minNode(n.Right)
		n.Val = succ.Val
		n.Right = del(n.Right, succ.Val)
	}
	return n
}

func minNode(n *BSTNode) *BSTNode {
	for n.Left != nil {
		n = n.Left
	}
	return n
}

// In-order = thứ tự tăng dần (đặc tính của BST)
func (t *BST) Sorted() []int {
	var out []int
	var walk func(*BSTNode)
	walk = func(n *BSTNode) {
		if n == nil {
			return
		}
		walk(n.Left)
		out = append(out, n.Val)
		walk(n.Right)
	}
	walk(t.Root)
	return out
}

func main() {
	t := &BST{}
	for _, v := range []int{50, 30, 70, 20, 40, 60, 80} {
		t.Insert(v)
	}
	fmt.Println("sorted:", t.Sorted()) // [20 30 40 50 60 70 80]

	fmt.Println("search 40:", t.Search(40)) // true
	fmt.Println("search 99:", t.Search(99)) // false

	t.Delete(30)
	fmt.Println("xóa 30 :", t.Sorted())

	t.Delete(50) // xóa root (2 con)
	fmt.Println("xóa 50 :", t.Sorted())
}
