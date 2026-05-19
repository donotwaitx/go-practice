# 12 – Goroutines & Channels

Concurrency là **siêu năng lực** của Go. Goroutine rẻ (vài KB stack), channel là kênh giao tiếp an toàn giữa chúng.

## Chạy

```bash
go run ./basics/12-goroutines
```

## Goroutine

Lightweight thread, quản lý bởi Go runtime, không phải OS thread.

```go
go someFunc()           // chạy someFunc bất đồng bộ
go func() { ... }()     // anonymous goroutine
```

So sánh:
- OS thread: ~1MB stack, vài ngàn là nhiều.
- Goroutine: ~2KB stack ban đầu (tự grow), **chạy hàng triệu** trên một máy bình thường.

⚠️ `main()` kết thúc → toàn bộ goroutines bị kill, không chờ. Phải dùng cơ chế đồng bộ.

## `sync.WaitGroup`

Chờ nhiều goroutine xong:

```go
var wg sync.WaitGroup

for i := 1; i <= 3; i++ {
    wg.Add(1)                // tăng counter
    go func(id int) {
        defer wg.Done()      // giảm counter khi xong
        work(id)
    }(i)
}

wg.Wait()  // block đến khi counter = 0
```

**Quy tắc:**
1. `Add` PHẢI gọi TRƯỚC `go` (không phải bên trong goroutine).
2. `Done` đặt trong `defer` để chắc chắn chạy.
3. Pass loop variable vào goroutine (`func(id int)`) — tránh closure bug.

## Channels

"Kênh" để goroutines trao đổi data. **Type-safe**.

```go
ch := make(chan int)       // unbuffered
ch := make(chan int, 5)    // buffered, capacity 5

ch <- 42      // gửi
v := <-ch     // nhận
close(ch)     // đóng
```

### Unbuffered channel

- Gửi **block** đến khi có người nhận.
- Nhận **block** đến khi có người gửi.
- → Đồng bộ, "rendezvous".

### Buffered channel

- Gửi không block nếu buffer chưa đầy.
- Nhận không block nếu buffer chưa rỗng.

### `close` và `range`

```go
go func() {
    for i := 1; i <= 5; i++ {
        ch <- i
    }
    close(ch)    // báo "không còn data"
}()

for v := range ch {   // tự dừng khi channel close
    fmt.Println(v)
}
```

⚠️ Quy tắc về `close`:
- **CHỈ producer (sender) close** channel.
- Close một channel đã close → **panic**.
- Gửi vào channel đã close → **panic**.
- Nhận từ channel đã close → trả zero value, `ok=false`.

```go
v, ok := <-ch
// ok=false nếu channel đã close và buffer rỗng
```

## Direction

Khi truyền channel làm tham số, giới hạn chiều để rõ intent:

```go
func producer(ch chan<- int) { ... }  // chỉ gửi
func consumer(ch <-chan int) { ... }  // chỉ nhận
```

## `select`

Chờ NHIỀU channel cùng lúc — case nào ready chạy case đó:

```go
select {
case v := <-c1:
    fmt.Println("c1:", v)
case v := <-c2:
    fmt.Println("c2:", v)
case c3 <- 10:
    fmt.Println("sent to c3")
case <-time.After(time.Second):
    fmt.Println("timeout")
default:
    fmt.Println("no channel ready")  // non-blocking
}
```

**Use case phổ biến:**

### Timeout

```go
select {
case res := <-resultCh:
    handle(res)
case <-time.After(2 * time.Second):
    return errors.New("timeout")
}
```

### Cancellation với context

```go
select {
case <-ctx.Done():
    return ctx.Err()
case work <- task:
    // sent
}
```

## Patterns phổ biến

### Worker pool

```go
jobs := make(chan int, 100)
results := make(chan int, 100)

for w := 1; w <= 3; w++ {
    go func() {
        for j := range jobs {
            results <- j * 2
        }
    }()
}

for j := 1; j <= 5; j++ { jobs <- j }
close(jobs)
// nhận kết quả từ results...
```

### Fan-in / Fan-out

- **Fan-out:** nhiều worker đọc từ cùng 1 channel.
- **Fan-in:** nhiều channel gửi vào 1 channel chung.

## Pitfalls

### 1. Deadlock

```go
ch := make(chan int)
ch <- 1   // ❌ deadlock — không ai nhận
```

Unbuffered channel gửi mà không có receiver → fatal deadlock.

### 2. Goroutine leak

```go
go func() {
    v := <-ch   // không bao giờ nhận được → goroutine bị treo mãi
}()
```

Goroutine không có cách thoát → leak. Luôn có exit path (close channel, context).

### 3. Loop variable capture (đã fix từ Go 1.22)

```go
// Go < 1.22 — BUG
for _, x := range items {
    go func() {
        process(x)   // tất cả goroutines thấy cùng x cuối
    }()
}

// Fix: pass làm parameter
for _, x := range items {
    go func(x Item) { process(x) }(x)
}
```

Go 1.22+ tự fix: mỗi iteration có biến x riêng.

### 4. Data race

Đọc/ghi cùng biến từ nhiều goroutine mà không sync → race condition. Check bằng:

```bash
go run -race main.go
```

Race detector là tool tuyệt vời, luôn chạy test với `-race`.

## Triết lý

> **Don't communicate by sharing memory; share memory by communicating.**

Thay vì shared state + mutex, ưu tiên message passing qua channel. Mutex vẫn cần (đặc biệt cho map shared) nhưng channel rõ intent hơn.

## Ghi nhớ

- `go` cheap nhưng KHÔNG free. Vẫn cần kiểm soát số lượng.
- Channel có buffer ≠ async queue — có giới hạn. Đầy = block.
- `sync.Mutex`, `sync.RWMutex`, `sync.Once`, `sync.Map` cho các kịch bản khác.
- `context.Context` là cách chuẩn để cancel/timeout cascade qua nhiều goroutine.
