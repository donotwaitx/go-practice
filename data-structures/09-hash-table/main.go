// 09 - Hash Table (custom implementation)
// Go có sẵn map built-in, nhưng implement tay để hiểu nguyên lý.
// Xử lý collision bằng "separate chaining" (mỗi bucket là một linked list).
// Average O(1) cho Get/Set/Delete, worst case O(n) khi nhiều collision.
package main

import "fmt"

type entry struct {
	key   string
	value int
	next  *entry
}

type HashTable struct {
	buckets []*entry
	size    int
}

func NewHashTable(cap int) *HashTable {
	if cap < 8 {
		cap = 8
	}
	return &HashTable{buckets: make([]*entry, cap)}
}

// Hàm hash đơn giản — djb2
func (h *HashTable) hash(key string) int {
	var hash uint64 = 5381
	for _, c := range key {
		hash = ((hash << 5) + hash) + uint64(c) // hash * 33 + c
	}
	return int(hash % uint64(len(h.buckets)))
}

func (h *HashTable) Set(key string, value int) {
	idx := h.hash(key)
	// Update nếu key đã tồn tại
	for e := h.buckets[idx]; e != nil; e = e.next {
		if e.key == key {
			e.value = value
			return
		}
	}
	// Insert mới — chèn đầu chuỗi
	h.buckets[idx] = &entry{key: key, value: value, next: h.buckets[idx]}
	h.size++

	// Resize khi load factor > 0.75
	if float64(h.size)/float64(len(h.buckets)) > 0.75 {
		h.resize(len(h.buckets) * 2)
	}
}

func (h *HashTable) Get(key string) (int, bool) {
	for e := h.buckets[h.hash(key)]; e != nil; e = e.next {
		if e.key == key {
			return e.value, true
		}
	}
	return 0, false
}

func (h *HashTable) Delete(key string) bool {
	idx := h.hash(key)
	var prev *entry
	for e := h.buckets[idx]; e != nil; e = e.next {
		if e.key == key {
			if prev == nil {
				h.buckets[idx] = e.next
			} else {
				prev.next = e.next
			}
			h.size--
			return true
		}
		prev = e
	}
	return false
}

func (h *HashTable) Len() int { return h.size }

// Resize: cấp phát bucket mới to hơn, rehash toàn bộ entry
func (h *HashTable) resize(newCap int) {
	old := h.buckets
	h.buckets = make([]*entry, newCap)
	h.size = 0
	for _, head := range old {
		for e := head; e != nil; e = e.next {
			h.Set(e.key, e.value)
		}
	}
}

func main() {
	h := NewHashTable(4)
	h.Set("apple", 1)
	h.Set("banana", 2)
	h.Set("cherry", 3)
	h.Set("apple", 100) // update

	v, _ := h.Get("apple")
	fmt.Println("apple :", v) // 100

	_, ok := h.Get("orange")
	fmt.Println("orange:", ok) // false

	h.Delete("banana")
	_, ok = h.Get("banana")
	fmt.Println("banana sau delete:", ok) // false

	// Stress test — đẩy nhiều key để thấy resize hoạt động
	for i := 0; i < 20; i++ {
		h.Set(fmt.Sprintf("k%d", i), i)
	}
	fmt.Println("len :", h.Len())                                         // 22
	fmt.Println("cap :", len(h.buckets))                                  // tự grow
	fmt.Println("k15 :", func() int { v, _ := h.Get("k15"); return v }()) // 15
}
