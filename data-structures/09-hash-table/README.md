# 09 – Hash Table (custom)

Hiểu cách `map` của Go hoạt động bằng cách tự implement. Bài này dùng **separate chaining** để xử lý collision.

## Chạy

```bash
go run ./data-structures/09-hash-table
```

## Khái niệm

```
buckets:
  [0] → ("apple", 1) → nil
  [1] → nil
  [2] → ("banana", 2) → ("cherry", 3) → nil   ← collision: 2 key cùng index
  [3] → nil
```

3 thành phần cốt lõi:

1. **Hash function** — biến key thành integer.
2. **Bucket array** — chỉ số = `hash(key) % len(buckets)`.
3. **Collision resolution** — khi 2 key về cùng bucket.

## Hash function

Bài này dùng **djb2**:

```go
func hash(key string) uint64 {
    var h uint64 = 5381
    for _, c := range key {
        h = ((h << 5) + h) + uint64(c)  // h * 33 + c
    }
    return h
}
```

Yêu cầu hash tốt:
- **Deterministic** — cùng key → cùng hash.
- **Uniform** — phân bố đều, giảm collision.
- **Fast**.

Go std dùng **AES-NI hardware hash** với random seed (chống hash flooding attack).

## Collision resolution

### Cách 1: Separate chaining (bài này)

Mỗi bucket là 1 linked list các entry. Collision = chèn vào head list.

✅ Đơn giản, không phụ thuộc load factor.
❌ Tốn memory (pointer extra), cache locality kém.

### Cách 2: Open addressing

Nếu bucket bị chiếm, dò sang bucket khác theo quy tắc (linear probing, quadratic, double hashing).

✅ Cache friendly hơn.
❌ Phức tạp khi delete (cần tombstone).

Go std map dùng dạng **open addressing với bucket có 8 slot** — hybrid của cả 2.

## Load factor & resize

```
load_factor = số entry / số bucket
```

Khi load factor > threshold (vd 0.75), bảng cần **resize** — gấp đôi bucket, rehash toàn bộ entry:

```go
if float64(h.size)/float64(len(h.buckets)) > 0.75 {
    h.resize(len(h.buckets) * 2)
}
```

→ Insert thường O(1), thỉnh thoảng O(n) khi resize. **Amortized O(1)**.

## Độ phức tạp

| Operation | Trung bình | Worst case |
|-----------|-----------|------------|
| Set | O(1) amortized | O(n) (collision dồn 1 bucket) |
| Get | O(1) | O(n) |
| Delete | O(1) | O(n) |

Worst case chỉ xảy ra với hash function tệ hoặc adversarial input (hash flooding attack).

## So với Go `map` built-in

| | Custom (bài này) | Go map |
|---|------------------|--------|
| Hash | djb2 | AES-NI + random seed |
| Collision | Chaining | 8-slot buckets + chaining |
| Iterate order | Theo bucket | Random (intentional) |
| Concurrent | ❌ panic | ❌ panic (cần `sync.Map`) |
| Hash flooding | Vulnerable | Resistant |

→ Trong production luôn dùng Go map.

## Pitfall

### 1. Map literal trong test = không deterministic

Order iterate khác nhau mỗi lần run → test phụ thuộc order sẽ flaky.

### 2. Modify trong khi iterate

Go map panic nếu thêm key trong loop range. Custom hash table cũng có thể bug.

### 3. Key chứa nil/NaN

NaN không bằng chính nó (`NaN != NaN`) → insert NaN không lookup được.
Pointer key: 2 pointer trỏ cùng object so sánh bằng được, nhưng nội dung khác địa chỉ → khác key.

### 4. Hash function phải stable

Nếu key là struct với pointer → hash từ pointer value sẽ đổi giữa các run (random heap layout). Phải hash từ field nội dung.

## Use cases (vs alternatives)

| Cần | Dùng |
|-----|------|
| Lookup O(1), không cần order | Hash table / map |
| Lookup O(log n), cần sorted/range | BST / red-black tree |
| Lookup O(L) cho prefix | Trie |
| Set nhỏ (< 10 phần tử) | Linear scan slice |

## Bài tập

- Implement LRU Cache (kết hợp hash + doubly linked list).
- Implement RandomSet với O(1) insert/remove/getRandom.
- 2-Sum với hash.
- Group Anagrams.
