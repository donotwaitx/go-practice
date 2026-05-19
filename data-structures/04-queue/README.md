# 04 – Queue (FIFO)

First In, First Out — phần tử vào trước ra trước. Như xếp hàng mua vé.

## Chạy

```bash
go run ./data-structures/04-queue
```

## Khái niệm

```
enqueue: 1, 2, 3
            head ↓        ↓ tail
            ┌───┬───┬───┐
            │ 1 │ 2 │ 3 │
            └───┴───┴───┘

dequeue: lấy 1 ra (đầu hàng)
```

2 operation chính:
- **Enqueue** — thêm vào cuối.
- **Dequeue** — lấy & xóa từ đầu.

## Implement: 3 cách

### Cách 1: slice naive (CHẬM)

```go
q := []int{}
q = append(q, v)         // enqueue O(1)
v := q[0]; q = q[1:]     // dequeue: O(1) nhưng underlying array không thu hồi, leak
```

Vấn đề: `q = q[1:]` không thực sự dequeue — chỉ dời pointer. Phần đã dequeue vẫn giữ trong underlying array → memory leak với queue lâu chạy.

### Cách 2: slice với head index (BÀI NÀY DÙNG)

Giữ thêm `head int` trỏ đến phần tử đầu logical. Periodically compact khi `head` quá to:

```go
type Queue[T any] struct {
    data []T
    head int
}
```

- Enqueue: `append`.
- Dequeue: tăng `head`, set `data[head] = zero`.
- Compact khi `head > 64 && head*2 > len(data)`.

→ O(1) amortized cả 2 chiều.

### Cách 3: Circular buffer (ring buffer)

Mảng cố định + 2 index head/tail wrap quanh. Tối ưu nếu biết max size trước.

### Cách 4: Doubly linked list

Push/pop 2 đầu = queue.

## Độ phức tạp

| Operation | Cách 2 (head index) |
|-----------|---------------------|
| Enqueue | O(1) amortized |
| Dequeue | O(1) amortized |
| Peek | O(1) |
| Len | O(1) |

## Ứng dụng

### 1. BFS (Breadth-First Search)

Duyệt cây/đồ thị theo tầng — đẩy node vào queue, lấy ra xử lý, đẩy hàng xóm vào.

### 2. Task scheduling

Producer push job, consumer pull job theo thứ tự.

### 3. Print spooler, request queue

Server xử lý request theo thứ tự đến.

### 4. Streaming / buffer

Sliding window data.

## Trong Go: channels

Channel chính là **thread-safe queue**:

```go
ch := make(chan int, 10) // buffered queue size 10
ch <- 1                  // enqueue
v := <-ch                // dequeue
```

→ Trong code Go thực tế, **dùng channel** thay vì tự implement queue khi cần concurrency.

## Pitfall

### Slice naive leak

```go
q = q[1:]   // ⚠️ data[0] vẫn không được GC
```

Production cần circular buffer hoặc head/compact pattern.

### Empty check

Phân biệt "queue rỗng" với "Dequeue trả zero":
```go
v, ok := q.Dequeue()
if !ok { /* empty */ }
```

### Channel khác slice queue

- Channel có capacity giới hạn, đầy thì block.
- Slice queue grow vô hạn (đến hết RAM).

Chọn theo use case.

## Bài tập

- Implement Queue dùng 2 stack.
- Implement Stack dùng 2 queue.
- Sliding Window Maximum (LeetCode 239) — dùng deque.
- BFS shortest path trong grid.
