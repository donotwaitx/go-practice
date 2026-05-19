// 02 - Routing với ServeMux (Go 1.22+)
// Từ Go 1.22, http.ServeMux hỗ trợ HTTP method và path variables:
//
//	"GET /users/{id}"   - match đúng method GET và bắt id
//	"POST /users"       - match POST
//	"GET /files/{path...}"  - wildcard path còn lại
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// Method + path cố định
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Home page")
	})

	// Path variable — lấy bằng r.PathValue
	mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Fprintf(w, "Get user %s\n", id)
	})

	// Khác method cùng path — handler khác nhau
	mux.HandleFunc("POST /users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Create user (POST)")
	})

	mux.HandleFunc("DELETE /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Delete user %s\n", r.PathValue("id"))
	})

	// Nested resource
	mux.HandleFunc("GET /users/{id}/posts/{postID}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "user=%s post=%s\n", r.PathValue("id"), r.PathValue("postID"))
	})

	// Wildcard: bắt toàn bộ path còn lại
	mux.HandleFunc("GET /files/{path...}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Serve file: %s\n", r.PathValue("path"))
	})

	// Host-specific (match chỉ khi Host header khớp)
	mux.HandleFunc("api.localhost/v1/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong from api host")
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
