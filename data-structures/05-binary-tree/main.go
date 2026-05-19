// 05 - Binary Tree & Traversals
// Cây nhị phân: mỗi node có tối đa 2 con (Left, Right).
// Bài này tập trung vào 4 cách duyệt cơ bản.
package main

import "fmt"

type TreeNode struct {
	Val         int
	Left, Right *TreeNode
}

// Xây cây mẫu:
//
//	    1
//	   / \
//	  2   3
//	 / \   \
//	4   5   6
func sampleTree() *TreeNode {
	return &TreeNode{
		Val: 1,
		Left: &TreeNode{
			Val:   2,
			Left:  &TreeNode{Val: 4},
			Right: &TreeNode{Val: 5},
		},
		Right: &TreeNode{
			Val:   3,
			Right: &TreeNode{Val: 6},
		},
	}
}

// 1. In-order: Left -> Root -> Right  (cho BST sẽ ra thứ tự tăng dần)
func inorder(n *TreeNode, out *[]int) {
	if n == nil {
		return
	}
	inorder(n.Left, out)
	*out = append(*out, n.Val)
	inorder(n.Right, out)
}

// 2. Pre-order: Root -> Left -> Right
func preorder(n *TreeNode, out *[]int) {
	if n == nil {
		return
	}
	*out = append(*out, n.Val)
	preorder(n.Left, out)
	preorder(n.Right, out)
}

// 3. Post-order: Left -> Right -> Root
func postorder(n *TreeNode, out *[]int) {
	if n == nil {
		return
	}
	postorder(n.Left, out)
	postorder(n.Right, out)
	*out = append(*out, n.Val)
}

// 4. Level-order (BFS): theo từng tầng — dùng queue
func levelorder(root *TreeNode) []int {
	if root == nil {
		return nil
	}
	result := []int{}
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		result = append(result, n.Val)
		if n.Left != nil {
			queue = append(queue, n.Left)
		}
		if n.Right != nil {
			queue = append(queue, n.Right)
		}
	}
	return result
}

// Chiều cao cây — đệ quy
func height(n *TreeNode) int {
	if n == nil {
		return 0
	}
	l, r := height(n.Left), height(n.Right)
	if l > r {
		return l + 1
	}
	return r + 1
}

func main() {
	root := sampleTree()

	var in, pre, post []int
	inorder(root, &in)
	preorder(root, &pre)
	postorder(root, &post)

	fmt.Println("inorder   :", in)               // 4 2 5 1 3 6
	fmt.Println("preorder  :", pre)              // 1 2 4 5 3 6
	fmt.Println("postorder :", post)             // 4 5 2 6 3 1
	fmt.Println("levelorder:", levelorder(root)) // 1 2 3 4 5 6
	fmt.Println("height    :", height(root))     // 3
}
