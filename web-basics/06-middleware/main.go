// 06 - Middleware
// Middleware là http.Handler bọc quanh handler khác, thêm logic chung
// (log, recover, auth, CORS, request ID...).
// Pattern: func(next http.Handler) http.Handler
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"crypto/rand"
	"encoding/hex"
)

// === Middleware 1: Logger — log mỗi request ===
func loggerMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Wrap ResponseWriter để capture status code
		sw := &statusWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(sw, r)
		log.Printf("%s %s -> %d (%s)", r.Method, r.URL.Path, sw.status, time.Since(start))
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (sw *statusWriter) WriteHeader(code int) {
	sw.status = code
	sw.ResponseWriter.WriteHeader(code)
}

// === Middleware 2: Recoverer — bắt panic, trả 500 ===
func recoverMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic: %v\n%s", rec, debug.Stack())
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// === Middleware 3: Request ID — gán ID duy nhất cho mỗi request ===
type ctxKey string

const requestIDKey ctxKey = "requestID"

func requestIDMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			b := make([]byte, 8)
			rand.Read(b)
			id = hex.EncodeToString(b)
		}
		w.Header().Set("X-Request-ID", id)
		ctx := context.WithValue(r.Context(), requestIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// === Middleware 4: Auth — check API key ===
func apiKeyMW(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-API-Key") != key {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// Chain helper: gộp nhiều middleware, áp dụng từ phải qua trái
//
//	chain(h, A, B, C) = A(B(C(h)))  — A chạy outer nhất
func chain(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := r.Context().Value(requestIDKey).(string)
	fmt.Fprintf(w, "hello (request_id=%s)\n", id)
}

func panicHandler(w http.ResponseWriter, r *http.Request) {
	panic("intentional panic for demo")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello", helloHandler)
	mux.HandleFunc("GET /panic", panicHandler)
	mux.HandleFunc("GET /secret", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "🔒 secret content")
	})

	// Chain global middleware cho mọi route
	handler := chain(mux,
		loggerMW,
		recoverMW,
		requestIDMW,
	)

	log.Println("listening on :8080")
	log.Println("test: curl -H 'X-API-Key: secret' http://localhost:8080/secret")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
