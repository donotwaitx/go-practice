// 10 - Graph (Adjacency List) + BFS + DFS
// Đồ thị: tập đỉnh (vertex) + tập cạnh (edge).
// Adjacency list: với mỗi đỉnh, lưu danh sách hàng xóm.
// Memory O(V + E), tốt cho đồ thị sparse (ít cạnh).
package main

import "fmt"

type Graph struct {
	adj      map[string][]string
	directed bool
}

func NewGraph(directed bool) *Graph {
	return &Graph{adj: map[string][]string{}, directed: directed}
}

func (g *Graph) AddVertex(v string) {
	if _, ok := g.adj[v]; !ok {
		g.adj[v] = nil
	}
}

func (g *Graph) AddEdge(u, v string) {
	g.AddVertex(u)
	g.AddVertex(v)
	g.adj[u] = append(g.adj[u], v)
	if !g.directed {
		g.adj[v] = append(g.adj[v], u)
	}
}

// BFS — duyệt theo tầng, dùng queue. Trả về thứ tự thăm.
func (g *Graph) BFS(start string) []string {
	if _, ok := g.adj[start]; !ok {
		return nil
	}
	visited := map[string]bool{start: true}
	queue := []string{start}
	order := []string{}

	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		order = append(order, v)
		for _, nb := range g.adj[v] {
			if !visited[nb] {
				visited[nb] = true
				queue = append(queue, nb)
			}
		}
	}
	return order
}

// DFS — đệ quy
func (g *Graph) DFS(start string) []string {
	visited := map[string]bool{}
	order := []string{}
	var dfs func(v string)
	dfs = func(v string) {
		if visited[v] {
			return
		}
		visited[v] = true
		order = append(order, v)
		for _, nb := range g.adj[v] {
			dfs(nb)
		}
	}
	dfs(start)
	return order
}

// Đường đi ngắn nhất (số cạnh) từ src → dst — dùng BFS
func (g *Graph) ShortestPath(src, dst string) []string {
	if src == dst {
		return []string{src}
	}
	parent := map[string]string{src: ""}
	queue := []string{src}
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		for _, nb := range g.adj[v] {
			if _, seen := parent[nb]; seen {
				continue
			}
			parent[nb] = v
			if nb == dst {
				// Truy ngược đường đi
				path := []string{dst}
				for cur := v; cur != ""; cur = parent[cur] {
					path = append([]string{cur}, path...)
				}
				return path
			}
			queue = append(queue, nb)
		}
	}
	return nil // không thông
}

func main() {
	// Đồ thị vô hướng:
	//   A - B - C
	//   |   |
	//   D - E - F
	g := NewGraph(false)
	g.AddEdge("A", "B")
	g.AddEdge("A", "D")
	g.AddEdge("B", "C")
	g.AddEdge("B", "E")
	g.AddEdge("D", "E")
	g.AddEdge("E", "F")

	fmt.Println("BFS từ A:", g.BFS("A"))
	fmt.Println("DFS từ A:", g.DFS("A"))
	fmt.Println("đường ngắn nhất A→F:", g.ShortestPath("A", "F"))
}
