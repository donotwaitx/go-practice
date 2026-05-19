# 01 – Singly Linked List

Mỗi node có 1 giá trị + 1 pointer trỏ đến node kế. Truy cập tuần tự, không random access.

## Chạy

```bash
go run ./data-structures/01-linked-list
```

## Khái niệm

```
Head → [1|*] → [2|*] → [3|*] → nil
```

Mỗi node:
```go
type Node struct {
    Val  int
    Next *Node
}
```

List chỉ giữ `Head`. Để duyệt đến cuối → đi từ head qua `Next` đến khi `nil`.

## Độ phức tạp

| Operation | Time | Note |
|-----------|------|------|
| Prepend (thêm đầu) | O(1) | Chỉ đổi `Head` |
| Append (thêm cuối) | O(n) | Phải đi đến tail. Nếu giữ pointer `tail` riêng → O(1) |
| Remove giá trị v | O(n) | Tìm rồi nối lại |
| Find | O(n) | Tuyến tính |
| Reverse | O(n), O(1) memory | Đảo `Next` từng node |
| Random access (n-th node) | O(n) | ❌ Không như array |

## Ưu / Nhược

✅ **Ưu:**
- Insert/delete giữa list O(1) nếu đã có pointer node đó.
- Không cần cấp phát liên tục như array → không tốn time cho resize.

❌ **Nhược:**
- Random access tệ (O(n)).
- Tốn memory: mỗi node ngoài data còn thêm pointer.
- Cache locality kém (node nằm rải rác trong memory) → chậm hơn slice trong thực tế cho hầu hết workload.

## Trong Go — khi nào dùng?

Hầu như **không bao giờ** dùng singly linked list trong production Go. Slice (`[]T`) tốt hơn 99% trường hợp do:
- Cache friendly.
- API rất gọn (`append`, `[i]`).

Linked list chỉ có ý nghĩa khi:
- Frequent insert/delete ở giữa, có sẵn pointer.
- Implement LRU cache (kết hợp doubly linked list + map).
- Học thuật / interview / xây DS phức tạp hơn.

Go std có `container/list` (doubly linked list) — vẫn hiếm dùng.

## Pitfall

### 1. Quên update `Head` khi list rỗng

```go
func (l *LinkedList) Append(v int) {
    n := &Node{Val: v}
    if l.Head == nil {
        l.Head = n   // ⚠️ phải check
        return
    }
    // ...
}
```

### 2. Memory leak với linked list lớn

Nếu vẫn giữ tham chiếu đến head sau khi pop, GC không thu hồi được phần đã pop. → set `n.Next = nil` khi remove.

### 3. Infinite loop nếu tạo cycle

```go
last.Next = l.Head   // ❌ list trở thành vòng tròn
```

→ Duyệt với `for cur != nil` sẽ chạy mãi mãi. Detect cycle bằng thuật toán Floyd (slow/fast pointer).

## Bài tập đề xuất

- Tìm phần tử thứ k từ cuối (one-pass).
- Detect cycle.
- Merge 2 sorted linked lists.
- Tính độ dài, kiểm tra palindrome.
