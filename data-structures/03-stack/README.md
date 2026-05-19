# 03 – Stack (LIFO)

Last In, First Out — phần tử đẩy vào cuối là phần tử lấy ra đầu tiên.

## Chạy

```bash
go run ./data-structures/03-stack
```

## Khái niệm

```
push 1, 2, 3:
    │ 3 │  ← top
    │ 2 │
    │ 1 │
    └───┘

pop: lấy 3 ra
```

3 operation chính:
- **Push** — thêm phần tử lên top.
- **Pop** — lấy & xóa phần tử top.
- **Peek** — xem top mà không xóa.

Trong Go, dễ nhất là dùng slice với `append` (push) và slice down (pop).

## Generics

Bài này dùng **generics** (Go 1.18+) để stack chứa được bất kỳ kiểu nào:

```go
type Stack[T any] struct {
    data []T
}

func (s *Stack[T]) Push(v T) { ... }
func (s *Stack[T]) Pop() (T, bool) { ... }
```

→ `Stack[int]`, `Stack[string]`, `Stack[*User]` đều dùng chung code.

## Độ phức tạp

| Operation | Time | Note |
|-----------|------|------|
| Push | O(1) amortized | `append` đôi khi grow slice (O(n)) nhưng amortized O(1) |
| Pop | O(1) | Cắt phần tử cuối |
| Peek | O(1) | Index cuối |
| Len | O(1) | `len(slice)` |

## Trick: zero out khi Pop

```go
n := len(s.data) - 1
v := s.data[n]
s.data[n] = zero // ← giải phóng reference
s.data = s.data[:n]
```

Nếu T là pointer hoặc chứa pointer (slice, map, struct lớn), KHÔNG set zero → slice header thu nhỏ nhưng underlying array vẫn giữ pointer → memory leak.

## Ứng dụng

### 1. Balanced parentheses

```go
func balanced(expr string) bool {
    s := &Stack[rune]{}
    pairs := map[rune]rune{')': '(', ']': '[', '}': '{'}
    for _, c := range expr {
        switch c {
        case '(', '[', '{':
            s.Push(c)
        case ')', ']', '}':
            top, ok := s.Pop()
            if !ok || top != pairs[c] { return false }
        }
    }
    return s.IsEmpty()
}
```

### 2. Function call stack

Mỗi function call → push stack frame. Return → pop. Đệ quy quá sâu → **stack overflow**.

### 3. Undo history

Mỗi action push lên stack. Undo = pop.

### 4. Expression evaluation

Postfix evaluation, infix → postfix conversion (shunting yard).

### 5. DFS

DFS đệ quy = dùng call stack. DFS iterative = dùng stack tường minh.

## Pitfall

### Pop trên stack rỗng

```go
v := stack.Pop()   // ❌ panic nếu trả zero không safe
```

Nên trả về `(T, bool)` để caller check, hoặc trả về error.

### Concurrency

Stack này KHÔNG thread-safe. Multi goroutine cần mutex hoặc dùng channel.

## Bài tập

- Min Stack — push/pop O(1), getMin() cũng O(1).
- Daily Temperatures (LeetCode 739).
- Largest Rectangle in Histogram (stack monotone).
