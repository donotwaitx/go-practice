// 08 - Trie (Prefix Tree)
// Cấu trúc dữ liệu cây cho chuỗi. Mỗi cạnh = 1 ký tự, path từ root → node = prefix.
// Cực kỳ hiệu quả cho autocomplete, spell-check, IP routing.
// Insert/Search/StartsWith đều O(L) với L = độ dài chuỗi.
package main

import "fmt"

type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool // đánh dấu kết thúc 1 từ
}

type Trie struct {
	root *TrieNode
}

func NewTrie() *Trie {
	return &Trie{root: &TrieNode{children: map[rune]*TrieNode{}}}
}

func (t *Trie) Insert(word string) {
	cur := t.root
	for _, c := range word {
		nxt, ok := cur.children[c]
		if !ok {
			nxt = &TrieNode{children: map[rune]*TrieNode{}}
			cur.children[c] = nxt
		}
		cur = nxt
	}
	cur.isEnd = true
}

// Search: từ có tồn tại CHÍNH XÁC trong trie?
func (t *Trie) Search(word string) bool {
	n := t.find(word)
	return n != nil && n.isEnd
}

// StartsWith: có từ nào bắt đầu bằng prefix không?
func (t *Trie) StartsWith(prefix string) bool {
	return t.find(prefix) != nil
}

func (t *Trie) find(s string) *TrieNode {
	cur := t.root
	for _, c := range s {
		nxt, ok := cur.children[c]
		if !ok {
			return nil
		}
		cur = nxt
	}
	return cur
}

// Autocomplete: liệt kê mọi từ có prefix cho trước
func (t *Trie) Autocomplete(prefix string) []string {
	node := t.find(prefix)
	if node == nil {
		return nil
	}
	var results []string
	var dfs func(n *TrieNode, path string)
	dfs = func(n *TrieNode, path string) {
		if n.isEnd {
			results = append(results, path)
		}
		for c, child := range n.children {
			dfs(child, path+string(c))
		}
	}
	dfs(node, prefix)
	return results
}

func main() {
	t := NewTrie()
	for _, w := range []string{"go", "golang", "google", "good", "bee", "beer"} {
		t.Insert(w)
	}

	fmt.Println("search 'go'    :", t.Search("go"))      // true
	fmt.Println("search 'goo'   :", t.Search("goo"))     // false (không phải từ đầy đủ)
	fmt.Println("starts 'goo'   :", t.StartsWith("goo")) // true

	fmt.Println("autocomplete 'go' :", t.Autocomplete("go"))
	fmt.Println("autocomplete 'be' :", t.Autocomplete("be"))
	fmt.Println("autocomplete 'xy' :", t.Autocomplete("xy")) // []
}
