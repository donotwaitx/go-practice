# Data Structures với Go

Học các cấu trúc dữ liệu cốt lõi qua Go. Sau khi xong `basics/`, đây là bước tiếp theo.

Mỗi thư mục là một chủ đề độc lập, có `main.go` chạy được + `README.md` giải thích nguyên lý, độ phức tạp, và pitfall.

## Lộ trình

| # | Chủ đề | Operations chính | Độ phức tạp |
|---|--------|------------------|-------------|
| 01 | linked-list | Singly: Prepend, Append, Remove, Reverse | O(1) prepend, O(n) còn lại |
| 02 | doubly-linked-list | Push/Pop cả 2 đầu | O(1) cho push/pop |
| 03 | stack | Push, Pop, Peek (LIFO) — generic | O(1) amortized |
| 04 | queue | Enqueue, Dequeue (FIFO) — generic | O(1) amortized |
| 05 | binary-tree | 4 cách duyệt (in/pre/post/level-order) | O(n) |
| 06 | bst | Insert, Search, Delete | O(log n) trung bình |
| 07 | heap | Min-heap + Priority Queue (`container/heap`) | O(log n) push/pop |
| 08 | trie | Insert, Search, StartsWith, Autocomplete | O(L) với L = độ dài key |
| 09 | hash-table | Custom hash + separate chaining + resize | O(1) trung bình |
| 10 | graph | Adjacency list + BFS + DFS + shortest path | BFS/DFS O(V+E) |

## Cách chạy

Từ thư mục gốc của repo:

```bash
go run ./data-structures/01-linked-list
go run ./data-structures/02-doubly-linked-list
# ...
```

## Tại sao học data structures qua Go?

Go có 2 đặc tính giúp việc học data structures rất "sạch":

1. **Garbage collected** — không phải lo malloc/free, focus vào logic.
2. **Pointer tường minh nhưng an toàn** — biểu diễn được node-based DS như C, nhưng không có pointer arithmetic và nil-deref runtime nhanh chóng catch lỗi.
3. **Generics (1.18+)** — implement DS type-safe một lần, dùng cho mọi kiểu.

## Built-in của Go vs implement tay

Go std library đã có sẵn nhiều DS:

| Built-in | Use case | Bài tương ứng |
|----------|----------|---------------|
| `map[K]V` | Hash table | 09 |
| `container/list` | Doubly linked list | 02 |
| `container/heap` | Heap interface | 07 |
| `slice` | Dynamic array, stack base | 03, 04 |

→ Trong production luôn ưu tiên built-in. Implement tay là để **hiểu nguyên lý**, không phải để dùng thật.

## Kiến thức trước khi bắt đầu

Đã học xong `basics/`, đặc biệt:
- 08 (structs) — định nghĩa node, link
- 09 (pointers) — node trỏ đến node
- 10 (interfaces) — generic & interface satisfaction
- 11 (errors) — return value với ok pattern

## Sau khi xong

Lộ trình tiếp theo gợi ý:
- **Algorithms:** sorting (quicksort, mergesort), searching (binary search), dynamic programming.
- **Concurrent data structures:** sync.Map, channel-based queue, lock-free patterns.
- **Real-world Go:** HTTP server, database, JSON API.
