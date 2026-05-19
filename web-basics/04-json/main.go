// 04 - JSON encode/decode
// Hầu hết API hiện đại đều giao tiếp qua JSON. encoding/json là standard library.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email,omitempty"` // bỏ field nếu empty
	// Field không export (chữ thường) không xuất hiện trong JSON
	password string
}

// GET /user → trả JSON
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	u := User{ID: 1, Name: "Alice", Email: "alice@x.com"}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(u); err != nil {
		http.Error(w, "encode: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// POST /user → đọc JSON từ body
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var u User

	// Decoder đọc trực tiếp từ stream — không cần đọc hết vào memory
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() // strict — fail nếu có field lạ

	if err := dec.Decode(&u); err != nil {
		http.Error(w, "decode: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Echo lại để client xác nhận
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "created",
		"user":    u,
	})
}

// GET /users → trả slice
func listUsersHandler(w http.ResponseWriter, r *http.Request) {
	users := []User{
		{ID: 1, Name: "Alice", Email: "a@x.com"},
		{ID: 2, Name: "Bob"}, // Email rỗng → omitempty bỏ
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// Demo Marshal vs MarshalIndent (pretty print)
func init() {
	u := User{ID: 1, Name: "Demo"}
	b, _ := json.Marshal(u)
	fmt.Println("marshal       :", string(b))
	b, _ = json.MarshalIndent(u, "", "  ")
	fmt.Println("marshal indent:")
	fmt.Println(string(b))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /user", getUserHandler)
	mux.HandleFunc("POST /user", createUserHandler)
	mux.HandleFunc("GET /users", listUsersHandler)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
