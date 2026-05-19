// 05 - REST CRUD API
// Full Create/Read/Update/Delete cho resource User, dùng in-memory store có mutex.
//
//	GET    /users        - list tất cả
//	POST   /users        - tạo mới
//	GET    /users/{id}   - lấy 1
//	PUT    /users/{id}   - update toàn bộ
//	DELETE /users/{id}   - xóa
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Store — trong production dùng database, ở đây dùng map + mutex
type Store struct {
	mu     sync.RWMutex
	data   map[int]User
	nextID int
}

func NewStore() *Store {
	return &Store{data: map[int]User{}, nextID: 1}
}

func (s *Store) List() []User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]User, 0, len(s.data))
	for _, u := range s.data {
		out = append(out, u)
	}
	return out
}

func (s *Store) Get(id int) (User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.data[id]
	return u, ok
}

func (s *Store) Create(u User) User {
	s.mu.Lock()
	defer s.mu.Unlock()
	u.ID = s.nextID
	s.nextID++
	s.data[u.ID] = u
	return u
}

func (s *Store) Update(id int, u User) (User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[id]; !ok {
		return User{}, false
	}
	u.ID = id
	s.data[id] = u
	return u, true
}

func (s *Store) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[id]; !ok {
		return false
	}
	delete(s.data, id)
	return true
}

// Helper: trả JSON với status code
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func main() {
	store := NewStore()
	// Seed
	store.Create(User{Name: "Alice", Email: "a@x.com"})
	store.Create(User{Name: "Bob", Email: "b@x.com"})

	mux := http.NewServeMux()

	mux.HandleFunc("GET /users", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, store.List())
	})

	mux.HandleFunc("POST /users", func(w http.ResponseWriter, r *http.Request) {
		var u User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			writeErr(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
			return
		}
		if u.Name == "" {
			writeErr(w, http.StatusBadRequest, "name is required")
			return
		}
		created := store.Create(u)
		w.Header().Set("Location", fmt.Sprintf("/users/%d", created.ID))
		writeJSON(w, http.StatusCreated, created)
	})

	mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			writeErr(w, http.StatusBadRequest, "id must be integer")
			return
		}
		u, ok := store.Get(id)
		if !ok {
			writeErr(w, http.StatusNotFound, "user not found")
			return
		}
		writeJSON(w, http.StatusOK, u)
	})

	mux.HandleFunc("PUT /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			writeErr(w, http.StatusBadRequest, "id must be integer")
			return
		}
		var u User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			writeErr(w, http.StatusBadRequest, "invalid JSON")
			return
		}
		updated, ok := store.Update(id, u)
		if !ok {
			writeErr(w, http.StatusNotFound, "user not found")
			return
		}
		writeJSON(w, http.StatusOK, updated)
	})

	mux.HandleFunc("DELETE /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			writeErr(w, http.StatusBadRequest, "id must be integer")
			return
		}
		if !store.Delete(id) {
			writeErr(w, http.StatusNotFound, "user not found")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
