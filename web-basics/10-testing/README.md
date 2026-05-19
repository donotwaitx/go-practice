# 10 – Testing HTTP với `httptest`

Test handler không cần khởi server thật. `httptest` của stdlib cung cấp mọi thứ cần.

## Chạy

```bash
go test ./web-basics/10-testing -v
go test ./web-basics/10-testing -v -run Add        # chỉ test có "Add"
go test ./web-basics/10-testing -bench=.            # chạy benchmark
go test ./web-basics/10-testing -cover              # coverage
```

## 3 cách test HTTP handler

### Cách 1: `httptest.NewRecorder` — KHÔNG khởi server

Nhanh nhất. Test handler trực tiếp, không qua network.

```go
func TestAdd(t *testing.T) {
    req := httptest.NewRequest("POST", "/calc", strings.NewReader(`{"a":2,"b":3,"op":"add"}`))
    rec := httptest.NewRecorder()
    
    handler.ServeHTTP(rec, req)
    
    if rec.Code != 200 { t.Errorf("want 200, got %d", rec.Code) }
    // rec.Body chứa response body
}
```

→ **99% test handler nên dùng cách này.** Đơn giản, nhanh, deterministic.

### Cách 2: `httptest.NewServer` — KHỞI server thật

Cho integration test có client thật (HTTP request đi qua network loopback):

```go
server := httptest.NewServer(handler)
defer server.Close()

resp, _ := http.Get(server.URL + "/path")
```

→ Dùng khi test client code, hoặc handler dùng `r.RemoteAddr`, `r.TLS`...

### Cách 3: Black box với `http.Client`

Run server bình thường, test request từ ngoài. Cho E2E test, không phải unit.

## Table-driven test — pattern chuẩn

Mỗi test case là một entry trong slice:

```go
cases := []struct {
    name     string
    body     string
    wantCode int
    wantRes  int
}{
    {"add", `{"a":1,"b":2,"op":"add"}`, 200, 3},
    {"sub", `{"a":5,"b":2,"op":"sub"}`, 200, 3},
    {"div zero", `{"a":1,"b":0,"op":"div"}`, 400, 0},
}

for _, c := range cases {
    t.Run(c.name, func(t *testing.T) {
        // ... run test, assert
    })
}
```

**Lợi ích:**
- Thêm case mới = thêm 1 dòng.
- Run riêng từng case với `-run TestX/case_name`.
- Failure rõ ràng — biết case nào sai.

## Test helpers

### Assert

Go stdlib không có `assert`. Pattern thường:

```go
if got != want {
    t.Errorf("got %v, want %v", got, want)
}
```

`t.Errorf` ghi error nhưng test tiếp. `t.Fatalf` ghi rồi DỪNG.

Helper function:
```go
func assertEqual(t *testing.T, got, want any) {
    t.Helper()  // báo line number lỗi là caller, không phải dòng này
    if got != want { t.Errorf("got %v, want %v", got, want) }
}
```

Library phổ biến: `testify` (`assert.Equal`, `require.NoError`...). Stdlib đủ cho project nhỏ.

### `t.Helper()`

```go
func mustDecode(t *testing.T, body io.Reader, v any) {
    t.Helper()  // ⚠️ quan trọng
    if err := json.NewDecoder(body).Decode(v); err != nil {
        t.Fatal(err)
    }
}
```

Khi assert fail trong helper, Go report **line của caller** thay vì line trong helper.

## Coverage

```bash
go test -cover ./...                          # in % coverage
go test -coverprofile=cover.out ./...         # save profile
go tool cover -html=cover.out                 # mở browser xem chi tiết
```

## Benchmark

```go
func BenchmarkCalc(b *testing.B) {
    body := `{"a":2,"b":3,"op":"add"}`
    h := NewServer()
    b.ResetTimer()  // bỏ thời gian setup
    for i := 0; i < b.N; i++ {
        req := httptest.NewRequest("POST", "/calc", strings.NewReader(body))
        rec := httptest.NewRecorder()
        h.ServeHTTP(rec, req)
    }
}
```

Chạy:
```bash
go test -bench=. -benchmem
# BenchmarkCalc-8   1000000   1234 ns/op   456 B/op   7 allocs/op
```

## Subtest & parallel

```go
for _, c := range cases {
    c := c  // ⚠️ capture loop var trước Go 1.22, từ 1.22 không cần
    t.Run(c.name, func(t *testing.T) {
        t.Parallel()  // chạy song song với subtest khác
        // ...
    })
}
```

Subtest parallel chạy song song trong cùng test func. Tăng tốc khi nhiều case I/O.

## Pitfall

### 1. Test không deterministic

```go
// ❌
time.Sleep(100 * time.Millisecond)  // flaky

// ❌
go doWork()
// rồi assert ngay → có thể chưa xong
```

→ Đồng bộ qua channel, mock thời gian, hoặc chấp nhận retry.

### 2. Global state giữa các test

```go
var counter int  // global

func TestA(t *testing.T) { counter++ }
func TestB(t *testing.T) { ... }   // counter đã = 1, có thể sai
```

→ Mỗi test setup fresh state. Hoặc reset trong `t.Cleanup()`.

### 3. Quên close

```go
server := httptest.NewServer(h)
// ⚠️ thiếu defer server.Close() → port leak
```

### 4. Compare struct với pointer

```go
got := &User{ID: 1}
want := &User{ID: 1}
if got != want { ... }  // ❌ luôn fail — compare pointer khác
```

→ Compare value: `*got != *want`. Hoặc dùng `reflect.DeepEqual` / `cmp.Diff`.

### 5. Test private package

Test cùng package (`package x`) đọc được private. Test "black box" (`package x_test`) chỉ thấy public. Chọn theo intent.

### 6. Đo coverage không đủ

100% coverage ≠ test tốt. Vẫn có thể miss edge case. Quan trọng là **test có failure mode rõ ràng**.

## Fuzz testing (Go 1.18+)

```go
func FuzzAdd(f *testing.F) {
    f.Add(2, 3)  // seed
    f.Fuzz(func(t *testing.T, a, b int) {
        result := a + b
        if result < a && b > 0 { t.Error("overflow not detected") }
    })
}
```

Chạy: `go test -fuzz=FuzzAdd -fuzztime=10s`. Tự sinh input, tìm crash.

## CI tích hợp

`.github/workflows/ci.yml` của repo đã chạy `go test -v -race ./...` cho mọi push/PR. → Mỗi commit đảm bảo test pass.

## Bài tập

- Test handler PUT user — cover cases: id không tồn tại, JSON sai, name empty, OK.
- Mock repository: tạo `interface NoteRepo`, mock cho test.
- Test middleware (Logger, Recover) bằng cách wrap handler giả.
- Fuzz JSON decoder để tìm panic input.
- Benchmark: so sánh `encoding/json` vs `easyjson`/`sonic`.
