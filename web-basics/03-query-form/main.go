// 03 - Query string & Form data
// Cách đọc dữ liệu client gửi qua URL query và HTML form.
package main

import (
	"fmt"
	"log"
	"net/http"
)

// GET /search?q=golang&page=2&tags=web&tags=server
func searchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	keyword := q.Get("q")        // get một giá trị
	page := q.Get("page")        // empty nếu không có
	tags := q["tags"]            // get tất cả values của key (slice)
	hasFilter := q.Has("filter") // chỉ check tồn tại

	fmt.Fprintf(w, "keyword: %q\n", keyword)
	fmt.Fprintf(w, "page   : %q\n", page)
	fmt.Fprintf(w, "tags   : %v\n", tags)
	fmt.Fprintf(w, "filter?: %v\n", hasFilter)
}

// POST /login (form-urlencoded từ HTML form)
func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Bắt buộc gọi ParseForm trước khi r.PostForm dùng được
	if err := r.ParseForm(); err != nil {
		http.Error(w, "parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// PostFormValue lấy từ body POST/PUT (form-urlencoded hoặc multipart)
	user := r.PostFormValue("username")
	pass := r.PostFormValue("password")

	fmt.Fprintf(w, "username=%q\n", user)
	fmt.Fprintf(w, "password=%q (length=%d)\n", pass, len(pass))
}

// Header check
func headerHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	ct := r.Header.Get("Content-Type")
	fmt.Fprintf(w, "Authorization: %s\n", auth)
	fmt.Fprintf(w, "Content-Type : %s\n", ct)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /search", searchHandler)
	mux.HandleFunc("POST /login", loginHandler)
	mux.HandleFunc("GET /headers", headerHandler)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
