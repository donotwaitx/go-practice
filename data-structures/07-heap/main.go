// 07 - Heap & Priority Queue
// Heap = cây nhị phân hoàn chỉnh, parent <= con (min-heap) hoặc >= con (max-heap).
// Go có sẵn package `container/heap` — bạn implement interface heap.Interface,
// Go lo phần còn lại. Push/Pop O(log n), Peek O(1).
package main

import (
	"container/heap"
	"fmt"
)

// === Min-heap đơn giản ===
type IntMinHeap []int

func (h IntMinHeap) Len() int           { return len(h) }
func (h IntMinHeap) Less(i, j int) bool { return h[i] < h[j] } // đổi `<` → `>` để có max-heap
func (h IntMinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntMinHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *IntMinHeap) Pop() any {
	old := *h
	n := len(old)
	v := old[n-1]
	*h = old[:n-1]
	return v
}

// === Priority Queue với struct value ===
type Task struct {
	Name     string
	Priority int // càng nhỏ càng ưu tiên
}

type TaskQueue []*Task

func (pq TaskQueue) Len() int           { return len(pq) }
func (pq TaskQueue) Less(i, j int) bool { return pq[i].Priority < pq[j].Priority }
func (pq TaskQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }
func (pq *TaskQueue) Push(x any)        { *pq = append(*pq, x.(*Task)) }
func (pq *TaskQueue) Pop() any {
	old := *pq
	n := len(old)
	v := old[n-1]
	*pq = old[:n-1]
	return v
}

func main() {
	// Min-heap
	h := &IntMinHeap{5, 2, 8, 1, 9, 3}
	heap.Init(h)
	fmt.Println("heap min:", (*h)[0]) // 1

	heap.Push(h, 0)
	for h.Len() > 0 {
		fmt.Print(heap.Pop(h), " ") // 0 1 2 3 5 8 9
	}
	fmt.Println()

	// Priority queue
	pq := &TaskQueue{}
	heap.Init(pq)
	heap.Push(pq, &Task{Name: "write report", Priority: 3})
	heap.Push(pq, &Task{Name: "fix bug", Priority: 1})
	heap.Push(pq, &Task{Name: "team meeting", Priority: 2})

	fmt.Println("--- xử lý task theo ưu tiên ---")
	for pq.Len() > 0 {
		t := heap.Pop(pq).(*Task)
		fmt.Printf("[P%d] %s\n", t.Priority, t.Name)
	}
}
