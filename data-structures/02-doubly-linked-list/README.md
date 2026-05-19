# 02 – Doubly Linked List

Mỗi node có 2 pointer: `Prev` và `Next`. Duyệt được 2 chiều, xóa O(1) nếu có pointer.

## Chạy

```bash
go run ./data-structures/02-doubly-linked-list
```

## Khái niệm

```
nil ← [Prev|1|Next] ⇄ [Prev|2|Next] ⇄ [Prev|3|Next] → nil
       Head                                Tail
```

List giữ cả `Head` và `Tail` → push/pop ở cả 2 đầu đều O(1).

```go
type DNode struct {
    Val        int
    Prev, Next *DNode
}
```

## Độ phức tạp

| Operation | Time |
|-----------|------|
| PushFront / PushBack | O(1) |
| PopFront / PopBack | O(1) |
| Insert/Remove giữa list (đã có pointer) | O(1) |
| Find | O(n) |

## So với singly

| | Singly | Doubly |
|---|--------|--------|
| Memory mỗi node | 1 pointer extra | 2 pointer extra |
| Duyệt ngược | ❌ | ✅ |
| Xóa node khi có pointer | O(n) (cần prev) | O(1) |
| Insert trước node X | O(n) | O(1) |

## Khi nào dùng

### 1. LRU Cache

Pattern kinh điển: doubly linked list + hash map.
- Map: `key → node pointer` (O(1) lookup).
- List: thứ tự sử dụng (đầu = mới nhất, cuối = cũ nhất).
- Khi access key → di chuyển node về đầu (O(1) nhờ pointer).
- Khi đầy → xóa tail (O(1)).

### 2. Deque (Double-Ended Queue)

Push/pop cả 2 đầu — phục vụ sliding window, work-stealing scheduler...

### 3. Browser history, undo/redo

Cần đi ngược/xuôi → doubly fit hoàn hảo.

## Trong Go

`container/list` của std library chính là doubly linked list:

```go
import "container/list"

l := list.New()
e1 := l.PushBack("hello")
e2 := l.PushFront("world")
l.Remove(e1)
for e := l.Front(); e != nil; e = e.Next() {
    fmt.Println(e.Value)
}
```

→ Dùng cái này trong production thay vì tự implement.

## Pitfall

### Update đủ 4 pointer khi remove

Khi xóa node X khỏi giữa list, phải update:
- `X.Prev.Next = X.Next`
- `X.Next.Prev = X.Prev`

Quên một là list vỡ. Edge case `Prev == nil` hoặc `Next == nil` cần xử lý riêng.

### Sentinel / dummy node

Cách phổ biến để giảm edge case: thêm 2 sentinel `head`, `tail` không chứa data. List thật nằm giữa chúng. → Không bao giờ phải check `nil`.

## Bài tập

- Implement LRU Cache với capacity K (LeetCode 146).
- Implement deque có Push/Pop ở cả 2 đầu O(1).
- Reverse một doubly linked list tại chỗ.
