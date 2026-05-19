// 01 - HTTP Server
// Server HTTP đơn giản nhất với net/http (standard library, không dependency).
package main

import (
	"fmt"
	"log"
	"net/http"
)

// 1) Handler dùng HandlerFunc — đơn giản nhất
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}

// 2) Handler in thông tin request — xem những gì server nhận được
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Method: %s\n", r.Method)
	fmt.Fprintf(w, "Path  : %s\n", r.URL.Path)
	fmt.Fprintf(w, "Host  : %s\n", r.Host)
	fmt.Fprintf(w, "UA    : %s\n", r.UserAgent())
	fmt.Fprintf(w, "Addr  : %s\n", r.RemoteAddr)
}

// 3) Set status & header trước khi ghi body
func teapotHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusTeapot) // 418
	fmt.Fprintln(w, "I'm a teapot ☕")
}

func main() {
	// http.HandleFunc đăng ký vào DefaultServeMux
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/teapot", teapotHandler)

	addr := ":8080"
	log.Println("Server listening at http://localhost" + addr)
	// ListenAndServe block đến khi error. Nil = dùng DefaultServeMux.
	log.Fatal(http.ListenAndServe(addr, nil))
}
