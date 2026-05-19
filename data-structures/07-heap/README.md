# 07 – Heap & Priority Queue

Heap là cây nhị phân hoàn chỉnh với invariant:
- **Min-heap:** parent ≤ con → root là min.
- **Max-heap:** parent ≥ con → root là max.

Heap là implementation chính của Priority Queue (PQ).

## Chạy

```bash
go run ./data-structures/07-heap
```

## Khái niệm

```
Min-heap:
       1
      / \
     2   3
    / \  /
   5  8 9
```

Heap được lưu trong **array**, không cần pointer:
- Node tại index `i`:
  - Parent: `(i-1) / 2`
  - Left:   `2*i + 1`
  - Right:  `2*i + 2`

## Go `container/heap`

Go std cung cấp interface, bạn implement → Go làm phần thuật toán (sift up/down):

```go
type Interface interface {
    sort.Interface       // Len, Less, Swap
    Push(x any)
    Pop() any
}
```

Implement 5 method là dùng được:

```go
type IntMinHeap []int

func (h IntMinHeap) Len() int           { return len(h) }
func (h IntMinHeap) Less(i, j int) bool { return h[i] < h[j] }  // < cho min-heap
func (h IntMinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntMinHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *IntMinHeap) Pop() any {
    old := *h
    n := len(old)
    v := old[n-1]
    *h = old[:n-1]
    return v
}
```

⚠️ **KHÔNG gọi trực tiếp** `h.Push(x)`. Phải gọi `heap.Push(h, x)`:
- `h.Push` chỉ append, không sắp xếp.
- `heap.Push` gọi `h.Push` rồi sift-up để giữ heap invariant.

Tương tự với Pop.

## API

| Function | Mô tả | Time |
|----------|-------|------|
| `heap.Init(h)` | Heapify một slice đã có | O(n) |
| `heap.Push(h, x)` | Thêm phần tử | O(log n) |
| `heap.Pop(h)` | Lấy & xóa min/max | O(log n) |
| `h[0]` | Peek min/max | O(1) |
| `heap.Fix(h, i)` | Sửa vị trí i sau khi đổi giá trị | O(log n) |
| `heap.Remove(h, i)` | Xóa phần tử tại i | O(log n) |

## Max-heap

Đơn giản: đổi `Less` từ `<` thành `>`:

```go
func (h IntMaxHeap) Less(i, j int) bool { return h[i] > h[j] }
```

## Priority Queue

PQ thường lưu struct, sắp xếp theo field `Priority`:

```go
type Task struct {
    Name     string
    Priority int
}

type TaskQueue []*Task
func (pq TaskQueue) Less(i, j int) bool { return pq[i].Priority < pq[j].Priority }
// ... các method khác
```

## Ứng dụng

### 1. Top K elements

Tìm K phần tử lớn nhất / nhỏ nhất trong O(n log k):
- Dùng min-heap size K.
- Mỗi phần tử: nếu lớn hơn root → push, pop nhỏ nhất.

### 2. Dijkstra shortest path

PQ chứa `(distance, vertex)`, luôn lấy vertex có distance min để relax.

### 3. Huffman coding

Build cây Huffman bằng cách merge 2 frequency thấp nhất → push lại.

### 4. Event scheduling

Mỗi event có timestamp. PQ trả event sớm nhất kế tiếp.

### 5. K-way merge

Merge K sorted list → heap chứa head của mỗi list.

## Pitfall

### Gọi sai API

```go
h.Push(x)        // ❌ không sift up, vỡ invariant
heap.Push(h, x)  // ✅
```

### Update priority

Sửa value trong heap **không** tự cân bằng. Phải:
```go
pq[i].Priority = newPri
heap.Fix(pq, i)
```

Cần biết index `i` — nếu không, phải lưu thêm map `key → index` (gọi là **indexed PQ**).

### Pop trên heap rỗng

`heap.Pop` panic nếu rỗng. Check `h.Len() > 0` trước.

### Less là `<=` (sai)

Phải dùng `<` strict cho min-heap, không phải `<=`. Equal phần tử có thể swap loạn xạ → vẫn đúng heap invariant nhưng output không stable.

## Bài tập

- K largest elements (LeetCode 215).
- Merge K sorted lists (LeetCode 23).
- Find median from data stream (LeetCode 295) — 2 heap.
- Task scheduler — re-arrange tasks với cooldown.
