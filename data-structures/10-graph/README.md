# 10 – Graph (Adjacency List) + BFS + DFS

Đồ thị = tập **đỉnh (vertex)** + tập **cạnh (edge)** nối các đỉnh. Cơ sở cho social network, routing, dependency, AI search...

## Chạy

```bash
go run ./data-structures/10-graph
```

## Phân loại đồ thị

| Tiêu chí | Loại |
|----------|------|
| Hướng | **Directed** (cạnh có chiều) / **Undirected** |
| Trọng số | **Weighted** / **Unweighted** |
| Chu trình | **Cyclic** / **Acyclic** (DAG) |
| Kết nối | **Connected** / **Disconnected** |

## Cách biểu diễn

### 1. Adjacency list (bài này dùng)

Mỗi đỉnh giữ danh sách hàng xóm.

```go
type Graph struct {
    adj map[string][]string
}
```

- Memory: O(V + E)
- Tốt cho **sparse graph** (E << V²).
- Iterate hàng xóm nhanh.

### 2. Adjacency matrix

Mảng 2 chiều `matrix[u][v] = 1` nếu có cạnh.

- Memory: O(V²) — luôn
- Check cạnh `u-v`: O(1)
- Tốt cho **dense graph** (E ~ V²) hoặc cần nhiều check cạnh.

### 3. Edge list

List các tuple `(u, v, weight)`. Đơn giản nhất, dùng cho Kruskal MST.

## BFS — Breadth-First Search

Duyệt theo TẦNG. Dùng queue.

```go
visited := map[string]bool{start: true}
queue := []string{start}
for len(queue) > 0 {
    v := queue[0]
    queue = queue[1:]
    // process v
    for _, nb := range adj[v] {
        if !visited[nb] {
            visited[nb] = true
            queue = append(queue, nb)
        }
    }
}
```

**Ứng dụng:**
- Đường ngắn nhất trong unweighted graph.
- Tìm friend cấp 1, 2 trong mạng xã hội.
- Web crawler theo độ sâu.

## DFS — Depth-First Search

Đi sâu trước, rồi backtrack. Dễ implement đệ quy.

```go
var dfs func(v string)
dfs = func(v string) {
    if visited[v] { return }
    visited[v] = true
    // process v
    for _, nb := range adj[v] { dfs(nb) }
}
```

**Ứng dụng:**
- Topological sort (DAG).
- Detect cycle.
- Connected components.
- Maze solving, backtracking.

## So sánh BFS vs DFS

| | BFS | DFS |
|---|-----|-----|
| Data structure | Queue | Stack (hoặc đệ quy) |
| Memory | O(W) — bề rộng | O(H) — chiều sâu |
| Tìm đường ngắn nhất | ✅ unweighted | ❌ |
| Detect cycle | Khó hơn | ✅ dễ |
| Topo sort | ❌ | ✅ (Kahn's algorithm dùng BFS cũng OK) |

## Shortest path

| Loại đồ thị | Thuật toán |
|-------------|------------|
| Unweighted | BFS — O(V + E) |
| Weighted, non-negative | Dijkstra — O((V+E) log V) với heap |
| Có cạnh âm | Bellman-Ford — O(V·E) |
| All pairs | Floyd-Warshall — O(V³) |

Bài này có shortest path bằng BFS (unweighted).

## Độ phức tạp

| Operation | Adjacency list |
|-----------|----------------|
| Add vertex | O(1) |
| Add edge | O(1) |
| Check edge u-v | O(deg(u)) |
| BFS / DFS | O(V + E) |
| Memory | O(V + E) |

## Pitfall

### 1. Quên đánh dấu visited

BFS/DFS trên graph có cycle → infinite loop nếu không track visited.

→ Khác với **tree** (không có cycle), graph BẮT BUỘC visited set.

### 2. Directed vs undirected

```go
g.adj[u] = append(g.adj[u], v)
if !g.directed {
    g.adj[v] = append(g.adj[v], u)  // cạnh 2 chiều
}
```

Quên dòng thứ 2 → tưởng undirected nhưng thật ra là directed.

### 3. Recursion depth

DFS đệ quy trên graph triệu node có thể overflow. → Iterative DFS với explicit stack.

### 4. Slice modification trong loop

```go
for _, nb := range g.adj[v] {
    g.AddEdge(v, "x")  // ⚠️ undefined behavior — modify slice đang range
}
```

### 5. Self-loop & multi-edge

Đồ thị "thuần" thường không cho self-loop (`v→v`) hay multi-edge (2 cạnh `u→v`). Code bài này không filter — bạn cần check nếu cần.

## Use cases thực tế

- **Social network** — friend recommendation, shortest connection.
- **Maps / navigation** — Dijkstra trên road network.
- **Package manager** — dependency graph + topo sort.
- **Compiler** — call graph, control flow graph.
- **Crawler / spider** — BFS over web.
- **Game AI** — A* pathfinding (BFS + heuristic).

## Bài tập

- Number of Islands (LeetCode 200) — DFS/BFS trên grid.
- Course Schedule (LeetCode 207) — detect cycle bằng DFS.
- Clone Graph (LeetCode 133) — DFS với map.
- Word Ladder (LeetCode 127) — BFS với word transform.
- Implement Dijkstra với `container/heap`.
- Topological Sort (Kahn's algorithm).
