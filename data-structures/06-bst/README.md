# 06 – Binary Search Tree (BST)

Cây nhị phân với **invariant**: với mọi node, mọi giá trị trong Left subtree < node < mọi giá trị trong Right subtree.

## Chạy

```bash
go run ./data-structures/06-bst
```

## Khái niệm

```
       50
      /  \
    30    70
   / \    / \
  20  40 60  80
```

Insert/Search/Delete chạy theo branch — trung bình O(log n).

## Operations

### Insert

So sánh với node hiện tại, đi trái nếu nhỏ hơn, đi phải nếu lớn hơn, đến nil thì gắn vào:

```go
func insert(n *BSTNode, v int) *BSTNode {
    if n == nil { return &BSTNode{Val: v} }
    if v < n.Val      { n.Left  = insert(n.Left, v) }
    else if v > n.Val { n.Right = insert(n.Right, v) }
    return n
}
```

### Search

```go
for cur != nil {
    if v == cur.Val { return true }
    if v < cur.Val  { cur = cur.Left }
    else            { cur = cur.Right }
}
return false
```

### Delete — case khó nhất

3 trường hợp:

1. **Node là lá** → set parent's pointer về nil.
2. **Node có 1 con** → thay thế node bằng con của nó.
3. **Node có 2 con** → tìm **in-order successor** (node nhỏ nhất bên phải), copy giá trị, rồi delete successor (chắc chắn rơi vào case 1 hoặc 2).

```go
func del(n *BSTNode, v int) *BSTNode {
    if n == nil { return nil }
    if v < n.Val      { n.Left  = del(n.Left, v) }
    else if v > n.Val { n.Right = del(n.Right, v) }
    else {
        if n.Left == nil  { return n.Right }
        if n.Right == nil { return n.Left }
        succ := minNode(n.Right)
        n.Val = succ.Val
        n.Right = del(n.Right, succ.Val)
    }
    return n
}
```

## Đặc tính then chốt: In-order = Sorted

Duyệt in-order BST cho dãy **tăng dần**. Đây là kết quả trực tiếp của BST invariant.

→ BST là một cách "sort dần dần khi insert".

## Độ phức tạp

| Operation | Trung bình | Worst case |
|-----------|-----------|------------|
| Insert | O(log n) | O(n) |
| Search | O(log n) | O(n) |
| Delete | O(log n) | O(n) |
| In-order | O(n) | O(n) |

Worst case xảy ra khi cây **bị skewed** thành linked list (insert dãy đã sort).

## Vấn đề: BST không tự cân bằng

Insert 1, 2, 3, 4, 5 → cây nghiêng:

```
1
 \
  2
   \
    3
     \
      4
       \
        5
```

→ Search O(n).

## Self-balancing BST

Để giữ O(log n) worst case, cần cây tự cân bằng:
- **AVL Tree** — rotate sau insert/delete để giữ height diff ≤ 1.
- **Red-Black Tree** — màu đỏ/đen + invariant, balance lỏng hơn AVL nhưng nhanh hơn cho insert/delete.
- **B-Tree, B+ Tree** — fan-out lớn, dùng trong DB và filesystem.

Go không có RBT/AVL std built-in. Map của Go là hash table, không phải BST.

## Use cases

- **Sorted set/map** cần range query, predecessor/successor.
- **Database index** (B-tree là dạng generalized BST).
- **Auto-complete** với prefix sort.

## Pitfall

### Duplicate

BST cơ bản không hỗ trợ duplicate. Tùy convention:
- Bỏ qua (như code bài này).
- Cho phép ở phía right (hoặc left).
- Mỗi node giữ `count` thay vì duplicate node.

### So sánh non-comparable type

Code bài này dùng `int`. Với generics + custom type, cần truyền comparator (kiểu như Java `Comparator<T>`).

### Recursion depth

Cây skewed → đệ quy sâu → stack overflow. Iterative version an toàn hơn.

## Bài tập

- Validate BST (check invariant từ bộ test).
- Find k-th smallest element.
- Convert sorted array → balanced BST.
- Find Lowest Common Ancestor trong BST (O(log n)).
- Recover BST sau khi 2 node bị swap.
