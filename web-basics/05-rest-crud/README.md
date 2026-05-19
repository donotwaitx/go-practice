# 05 – REST CRUD API

Build một API CRUD đầy đủ. Tổng hợp những gì học từ bài 01-04 + concurrent-safe store.

## Chạy

```bash
go run ./web-basics/05-rest-crud
```

Test:
```bash
# LIST
curl http://localhost:8080/users

# CREATE
curl -X POST http://localhost:8080/users \
  -H 'Content-Type: application/json' \
  -d '{"name":"Charlie","email":"c@x.com"}'

# READ
curl http://localhost:8080/users/1

# UPDATE
curl -X PUT http://localhost:8080/users/1 \
  -H 'Content-Type: application/json' \
  -d '{"name":"Alice Smith","email":"alice@x.com"}'

# DELETE
curl -X DELETE http://localhost:8080/users/1
```

## REST conventions

| Method | Path | Action | Success status |
|--------|------|--------|----------------|
| GET | `/users` | List | 200 OK |
| GET | `/users/{id}` | Read one | 200 OK |
| POST | `/users` | Create | **201 Created** |
| PUT | `/users/{id}` | Replace toàn bộ | 200 OK |
| PATCH | `/users/{id}` | Update một phần | 200 OK |
| DELETE | `/users/{id}` | Delete | **204 No Content** |

### PUT vs PATCH

- **PUT**: gửi resource HOÀN CHỈNH. Bỏ field nào → field đó bị xóa.
- **PATCH**: gửi delta. Chỉ update field có trong payload.

Bài này dùng PUT cho đơn giản. PATCH cần handle "field không gửi" vs "field gửi null".

### Status codes phổ biến

| Code | Khi nào |
|------|---------|
| 200 OK | GET / PUT / PATCH thành công |
| 201 Created | POST tạo mới — kèm header `Location: /users/{id}` |
| 204 No Content | DELETE thành công, không body |
| 400 Bad Request | Client gửi data sai |
| 404 Not Found | Resource không tồn tại |
| 409 Conflict | Vi phạm constraint (vd email trùng) |
| 422 Unprocessable Entity | JSON đúng cú pháp nhưng business invalid |

## In-memory store với mutex

```go
type Store struct {
    mu     sync.RWMutex
    data   map[int]User
    nextID int
}
```

- **`sync.RWMutex`** — nhiều reader đồng thời (`RLock`), một writer (`Lock`).
- Read methods → `RLock` + `RUnlock`.
- Write methods → `Lock` + `Unlock`.

⚠️ Vì sao cần mutex? Mỗi HTTP request là 1 goroutine. Multiple goroutine cùng read+write map → **fatal error: concurrent map writes**.

## Anatomy của handler tốt

```go
func createHandler(w http.ResponseWriter, r *http.Request) {
    // 1. Parse input
    var u User
    if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
        writeErr(w, 400, "invalid JSON: "+err.Error())
        return
    }

    // 2. Validate
    if u.Name == "" {
        writeErr(w, 400, "name is required")
        return
    }

    // 3. Business logic
    created := store.Create(u)

    // 4. Response
    w.Header().Set("Location", fmt.Sprintf("/users/%d", created.ID))
    writeJSON(w, 201, created)
}
```

### Helper functions

Repeated logic → helper:
```go
func writeJSON(w, status, v) { ... }
func writeErr(w, status, msg) { ... }
```

Refactor sớm khi thấy lặp.

## Pitfall

### 1. Race condition trên store

```go
// Sai
if _, ok := store.Get(id); ok {
    store.Delete(id)  // ⚠️ giữa Get và Delete, ai đó có thể đã Delete
}

// Đúng — atomic check + delete trong 1 lock
ok := store.Delete(id)
if !ok { ... }
```

→ Đẩy logic phức tạp vào method của store dưới 1 lock.

### 2. PUT mà không clear field

Code bài này: `u.ID = id; data[id] = u` — replace hoàn toàn. Nếu client không gửi field nào → field đó về zero value.

→ Đúng với REST semantics của PUT. Nhưng cần tài liệu rõ cho client.

### 3. POST trả `id` ở body — quên header `Location`

```http
HTTP/1.1 201 Created
Location: /users/42         ← convention REST

{"id": 42, "name": "..."}
```

Client biết URL của resource mới — đỡ phải parse body.

### 4. Filter / pagination cho list

`GET /users` trả toàn bộ → OOM với 1 triệu user. Pattern:

```
GET /users?limit=20&offset=40         # offset pagination
GET /users?cursor=eyJpZCI6MTIzfQ      # cursor pagination
```

Bài này không có để giữ ngắn — production luôn cần.

### 5. ID auto-increment trong concurrent

`s.nextID++` — không atomic với map write. Trong lock thì OK. Nếu tách:
```go
id := atomic.AddInt32(&s.nextID, 1)  // hoặc dùng mutex
```

## Bài tập

- Implement PATCH với struct pointer fields.
- Thêm pagination `?limit&offset` cho list.
- Thêm filter `?q=alice` search theo name.
- Validate email format (sẽ làm trong bài 08).
- Thay in-memory bằng SQLite (sẽ làm trong bài 09).
- Thêm test cho mỗi handler (sẽ làm trong bài 10).
