# 05 – Binary Tree & Traversals

Cây nhị phân: mỗi node có tối đa 2 con (`Left`, `Right`). Đây là base cho BST, heap, segment tree...

## Chạy

```bash
go run ./data-structures/05-binary-tree
```

## Khái niệm

```
       1          ← root, height 3
      / \
     2   3
    / \   \
   4   5   6     ← lá
```

```go
type TreeNode struct {
    Val         int
    Left, Right *TreeNode
}
```

## 4 cách duyệt

### Depth-First (DFS) — đệ quy

| Tên | Thứ tự | Ứng dụng |
|-----|--------|----------|
| **In-order** | L → Root → R | BST in-order = sorted ascending |
| **Pre-order** | Root → L → R | Copy/clone cây, serialize |
| **Post-order** | L → R → Root | Xóa cây (free con trước), expression eval |

Implement chung pattern:

```go
func inorder(n *TreeNode, out *[]int) {
    if n == nil { return }
    inorder(n.Left, out)
    *out = append(*out, n.Val)   // pre: dòng này ở đầu, post: ở cuối
    inorder(n.Right, out)
}
```

### Breadth-First (BFS) — Level-order

Duyệt theo từng tầng, dùng queue:

```go
queue := []*TreeNode{root}
for len(queue) > 0 {
    n := queue[0]
    queue = queue[1:]
    // xử lý n
    if n.Left != nil  { queue = append(queue, n.Left) }
    if n.Right != nil { queue = append(queue, n.Right) }
}
```

## Output của cây mẫu

| Traversal | Kết quả |
|-----------|---------|
| In-order | 4 2 5 1 3 6 |
| Pre-order | 1 2 4 5 3 6 |
| Post-order | 4 5 2 6 3 1 |
| Level-order | 1 2 3 4 5 6 |

## Độ phức tạp

| Operation | Time | Space (worst) |
|-----------|------|---------------|
| Traversal | O(n) | O(h) = O(n) đệ quy, hoặc O(w) BFS với w = bề rộng |
| Height | O(n) | O(h) call stack |

Với cây cân bằng: `h = log n`. Cây skewed (lệch hẳn về 1 phía): `h = n`.

## Đệ quy vs Iterative

Đệ quy code rất sạch nhưng risk stack overflow với cây sâu (đặc biệt skewed):
- Cây cân bằng n = 1 triệu → h ≈ 20, OK.
- Cây skewed n = 1 triệu → h = 1 triệu → stack overflow.

Production với input untrusted → iterative + explicit stack:

```go
stack := []*TreeNode{}
cur := root
for cur != nil || len(stack) > 0 {
    for cur != nil {
        stack = append(stack, cur)
        cur = cur.Left
    }
    n := stack[len(stack)-1]
    stack = stack[:len(stack)-1]
    visit(n)
    cur = n.Right
}
```

## Khái niệm liên quan

- **Lá (leaf):** node không có con.
- **Internal node:** node có ít nhất 1 con.
- **Height:** số cạnh từ root đến lá xa nhất.
- **Depth của node:** số cạnh từ root đến nó.
- **Cây cân bằng:** chênh height giữa Left/Right subtree mọi node ≤ 1.
- **Cây hoàn chỉnh (complete):** mọi tầng đầy, trừ tầng cuối được đầy từ trái → phải.
- **Cây đầy (full):** mọi node có 0 hoặc 2 con.

## Pitfall

### Quên check nil

```go
visit(n.Left.Val)   // ❌ panic nếu Left == nil
```

### Đệ quy không tail-recursive

Go không có tail call optimization → đệ quy đến giới hạn stack OS (~1MB default).

### Mutate cây trong khi duyệt

Tránh sửa structure (Left/Right) trong duyệt — có thể tạo cycle hoặc skip node.

## Bài tập

- Đếm node, đếm lá.
- Kiểm tra hai cây có giống nhau không.
- Lowest Common Ancestor.
- Serialize / deserialize binary tree.
- Mirror / invert tree.
