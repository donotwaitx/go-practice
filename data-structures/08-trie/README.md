# 08 – Trie (Prefix Tree)

Cây dùng cho chuỗi: mỗi cạnh là 1 ký tự, đường đi từ root đến node là 1 prefix. Cực mạnh cho autocomplete, spell-check, IP routing.

## Chạy

```bash
go run ./data-structures/08-trie
```

## Khái niệm

Insert "go", "golang", "google", "good":

```
         (root)
           │
           g
           │
           o   ← isEnd (kết thúc "go")
          /|\
         d o l
         │ │ │
        (e)g a
         ⋮ │ │
           l n
           │ │
           e g
```

Mỗi node:
```go
type TrieNode struct {
    children map[rune]*TrieNode
    isEnd    bool  // đánh dấu node này là kết thúc 1 từ
}
```

`isEnd` quan trọng vì `go` lẫn `golang` cùng tồn tại, mà `go` không phải prefix-only.

## Operations

| Operation | Time | Space |
|-----------|------|-------|
| Insert(word) | O(L) | O(L) worst |
| Search(word) | O(L) | O(1) |
| StartsWith(prefix) | O(L) | O(1) |
| Autocomplete(prefix) | O(L + n) | O(n) |

L = độ dài word/prefix, n = tổng ký tự trong kết quả.

→ **Không phụ thuộc số từ đã insert** trong Search/StartsWith. Đây là điểm sức mạnh của trie.

## Search vs StartsWith

- **Search** trả `true` chỉ khi từ ĐÚNG đã được insert (kiểm tra `isEnd`).
- **StartsWith** trả `true` nếu prefix có trong trie, bất kể có là từ hoàn chỉnh hay không.

```go
t.Insert("apple")
t.Search("app")      // false — không insert "app"
t.StartsWith("app")  // true
t.Search("apple")    // true
```

## Autocomplete

Tìm node ứng với prefix, rồi DFS toàn bộ subtree để liệt kê từ:

```go
func (t *Trie) Autocomplete(prefix string) []string {
    node := t.find(prefix)
    if node == nil { return nil }
    var results []string
    var dfs func(n *TrieNode, path string)
    dfs = func(n *TrieNode, path string) {
        if n.isEnd { results = append(results, path) }
        for c, child := range n.children {
            dfs(child, path+string(c))
        }
    }
    dfs(node, prefix)
    return results
}
```

## So với hash set

| | Hash set | Trie |
|---|---------|------|
| Search 1 từ | O(L), good hash | O(L), no hash |
| Prefix search | ❌ O(n) scan all | ✅ O(L) |
| Sorted enumerate | ❌ random order | ✅ DFS in alphabet order |
| Memory | Compact | Tốn hơn (mỗi node có map) |
| Best for | Membership test | Prefix queries |

## Optimization

### 1. Array thay vì map

Nếu alphabet nhỏ (vd ASCII a-z), dùng `[26]*TrieNode` thay map → nhanh hơn nhiều, ít allocation:

```go
type TrieNode struct {
    children [26]*TrieNode
    isEnd    bool
}
// idx = c - 'a'
```

### 2. Radix tree / Compressed trie

Gộp các node có 1 con duy nhất → ít node hơn, ít cache miss. Linux kernel routing table dùng dạng này (Patricia trie).

### 3. DAWG / Suffix automaton

Cho compression cao hơn nhưng phức tạp hơn nhiều.

## Use cases thực tế

- **Autocomplete** (search bar, IDE).
- **Spell check** — kết hợp Levenshtein distance.
- **IP routing table** (Patricia trie).
- **Word games** (Boggle, Scrabble) — prune branch khi prefix không hợp lệ.
- **DNA/string matching** với fixed alphabet.

## Pitfall

### Memory blowup

Trie với alphabet lớn (Unicode) + ít từ → tốn memory vô lý. Cân nhắc compressed trie hoặc hash trie.

### Quên `isEnd`

Insert chỉ tạo node theo path → KHÔNG đánh dấu kết thúc → Search sẽ luôn trả false. Bug phổ biến khi mới học trie.

### Concurrent insert

Map của Go không thread-safe → concurrent insert vào trie cần lock.

## Bài tập

- Implement Trie với chỉ Insert / Search / StartsWith (LeetCode 208).
- Word Search II (LeetCode 212) — trie + DFS trên board.
- Longest Word in Dictionary.
- Replace Words với root mapping (LeetCode 648).
