# 09 – Database với `database/sql` + SQLite

`database/sql` là interface chuẩn cho mọi DB trong Go. Mỗi DB có driver riêng implement interface này. Bài này dùng SQLite qua pure-Go driver `modernc.org/sqlite`.

## Chạy

```bash
go run ./web-basics/09-database
```

Test:
```bash
curl http://localhost:8080/notes
curl -X POST http://localhost:8080/notes \
  -H 'Content-Type: application/json' \
  -d '{"title":"My Note","body":"Hello DB"}'
curl http://localhost:8080/notes/1
curl -X DELETE http://localhost:8080/notes/1
```

## Vì sao SQLite ở đây?

- **Đơn giản:** file đơn (hoặc `:memory:`), không cần cài server.
- **Production-ready:** Stripe, GitHub, Mozilla dùng cho 1 số service.
- **Pure Go driver** — `modernc.org/sqlite` không cần CGo → build static, deploy dễ.

Alternative: PostgreSQL (`pgx`), MySQL (`go-sql-driver/mysql`).

## Architecture

```
HTTP handler
    ↓ (gọi)
Repository  ← chỉ chỗ này biết SQL
    ↓ (gọi)
*sql.DB     ← connection pool
    ↓
SQLite file
```

Tách handler ↔ repo giúp:
- Test handler không cần DB thật (mock repo).
- Đổi DB ít chỗ phải sửa.
- SQL gom 1 chỗ, dễ review.

## API chính

### Mở DB

```go
db, err := sql.Open("sqlite", ":memory:")  // hoặc "file:notes.db"
defer db.Close()
```

`sql.Open` KHÔNG kết nối ngay — lazy. Gọi `db.Ping(ctx)` để verify.

### Connection pool

```go
db.SetMaxOpenConns(25)         // max conn đang dùng
db.SetMaxIdleConns(10)         // max idle
db.SetConnMaxLifetime(time.Hour)
db.SetConnMaxIdleTime(10*time.Minute)
```

**SQLite đặc biệt:** chỉ 1 writer cùng lúc. Set `MaxOpenConns(1)` cho writer hoặc dùng WAL mode.

### Query: nhiều rows

```go
rows, err := db.QueryContext(ctx, "SELECT id, name FROM users WHERE active = ?", true)
if err != nil { return err }
defer rows.Close()  // ⚠️ LUÔN close

for rows.Next() {
    var u User
    if err := rows.Scan(&u.ID, &u.Name); err != nil { return err }
    users = append(users, u)
}
if err := rows.Err(); err != nil { return err }  // check error sau loop
```

### QueryRow: 1 row

```go
var u User
err := db.QueryRowContext(ctx, "SELECT id, name FROM users WHERE id = ?", id).
    Scan(&u.ID, &u.Name)

if errors.Is(err, sql.ErrNoRows) {
    return ErrNotFound
}
```

### Exec: INSERT/UPDATE/DELETE

```go
res, err := db.ExecContext(ctx, "INSERT INTO users (name) VALUES (?)", name)
id, _ := res.LastInsertId()
n, _  := res.RowsAffected()
```

### Transaction

```go
tx, err := db.BeginTx(ctx, nil)
if err != nil { return err }
defer tx.Rollback() // no-op nếu đã commit

tx.ExecContext(ctx, "UPDATE balance SET amount = amount - ? WHERE id = ?", 100, from)
tx.ExecContext(ctx, "UPDATE balance SET amount = amount + ? WHERE id = ?", 100, to)

return tx.Commit()
```

## SQL injection — placeholder `?`

```go
// ❌ NGUY HIỂM — string concat
db.Query("SELECT * FROM users WHERE name = '" + userInput + "'")

// ✅ SAFE — placeholder, driver tự escape
db.Query("SELECT * FROM users WHERE name = ?", userInput)
```

Placeholder syntax tùy driver:
- SQLite, MySQL: `?`
- PostgreSQL: `$1`, `$2`, ...
- Oracle: `:name`

## Migrations

Bài này hardcode `CREATE TABLE IF NOT EXISTS` trong code. Production dùng tool:

- **`golang-migrate/migrate`** — file SQL, version up/down.
- **`pressly/goose`** — Go function migration cũng được.
- **`atlas`** — declarative schema (HCL/Terraform-like).

## Pitfall

### 1. Quên `rows.Close()`

```go
rows, _ := db.Query(...)
// ⚠️ thiếu defer rows.Close() → connection leak
```

Connection leak → pool cạn → "too many connections".

### 2. Quên kiểm tra `rows.Err()` sau loop

```go
for rows.Next() { ... }
// ⚠️ phải check rows.Err() — loop có thể dừng do error, không phải hết data
return rows.Err()
```

### 3. Scan vào sai type

```go
var id int
rows.Scan(&id)  // ❌ panic nếu cột là string
```

Scanner tự convert basic types, nhưng KHÔNG NULL → int. Dùng `sql.NullInt64`, `sql.NullString` cho cột nullable:

```go
var email sql.NullString
rows.Scan(&email)
if email.Valid { fmt.Println(email.String) }
```

Hoặc dùng pointer:
```go
var email *string
rows.Scan(&email)
```

### 4. `QueryRow` không phân biệt "no rows" vs "error"

```go
err := db.QueryRow(...).Scan(...)
if err == sql.ErrNoRows { /* not found */ }
if err != nil { /* real error */ }
```

Phải check `sql.ErrNoRows` riêng.

### 5. `LastInsertId` không đáng tin trên một số DB

PostgreSQL **không** support `LastInsertId`. Phải dùng:
```sql
INSERT INTO ... RETURNING id
```

### 6. Lock table với SQLite

SQLite single-writer. Multi-goroutine write → `database is locked`. Solution:
- `MaxOpenConns(1)` cho ghi tuần tự.
- WAL mode: `PRAGMA journal_mode = WAL`.
- Đẩy write vào single goroutine qua channel.

### 7. Time zone

SQLite không có timezone-aware datetime. Lưu dưới dạng UTC, convert ở app layer.

## ORM hay không?

Go community chia 2 trường phái:

**Anti-ORM** (phổ biến hơn):
- Raw SQL = mạnh, kiểm soát.
- ORM = magic, performance issue, học cost cao.

**ORM** (`gorm`, `ent`, `bun`):
- Productivity cho CRUD đơn giản.
- Code generation (ent) cho type-safe.

Middle ground:
- **`sqlc`** — viết SQL, generate Go code typed. Best of both worlds.

Khuyên: học `database/sql` trước. Nếu cần thì `sqlc`.

## Bài tập

- Thêm field `tags []string` cho Note — JSON column hoặc bảng riêng?
- Search: `GET /notes?q=foo` — FTS5 trong SQLite hoặc LIKE.
- Pagination: `?limit=10&cursor=<id>` với cursor-based.
- Transaction: transfer balance giữa 2 user, atomicity.
- Thay SQLite bằng PostgreSQL — chỉ đổi driver + placeholder.
