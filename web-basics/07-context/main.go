// 07 - Context: cancellation, timeout, request-scoped values
// context.Context truyền cancellation signal & deadline xuống call chain.
// Mỗi request HTTP đã có sẵn context: r.Context()
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Giả lập một query DB chậm — biết tôn trọng context
func slowQuery(ctx context.Context, ms int) (string, error) {
	select {
	case <-time.After(time.Duration(ms) * time.Millisecond):
		return fmt.Sprintf("query done in %dms", ms), nil
	case <-ctx.Done():
		// Client hủy request hoặc timeout
		return "", ctx.Err()
	}
}

// GET /slow?ms=2000 — chậm có chủ đích
func slowHandler(w http.ResponseWriter, r *http.Request) {
	ms := 500
	if v := r.URL.Query().Get("ms"); v != "" {
		fmt.Sscanf(v, "%d", &ms)
	}

	result, err := slowQuery(r.Context(), ms)
	if err != nil {
		// Client cancel hoặc timeout → log nhưng không cần ghi response
		log.Printf("aborted: %v", err)
		http.Error(w, "request cancelled", http.StatusRequestTimeout)
		return
	}
	fmt.Fprintln(w, result)
}

// GET /timeout — đặt timeout 1 giây cho call downstream
func timeoutHandler(w http.ResponseWriter, r *http.Request) {
	// Tạo context con với deadline ngắn hơn
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel() // LUÔN cancel để giải phóng resource

	result, err := slowQuery(ctx, 3000) // 3s — sẽ timeout
	if err != nil {
		http.Error(w, "downstream timed out: "+err.Error(), http.StatusGatewayTimeout)
		return
	}
	fmt.Fprintln(w, result)
}

// Truyền giá trị qua context (cẩn thận — chỉ dùng cho request-scoped data)
type ctxKey string

const userKey ctxKey = "user"

func withUserMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Header.Get("X-User")
		if user == "" {
			user = "anonymous"
		}
		ctx := context.WithValue(r.Context(), userKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func whoamiHandler(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(userKey).(string)
	fmt.Fprintf(w, "you are: %s\n", user)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /slow", slowHandler)
	mux.HandleFunc("GET /timeout", timeoutHandler)
	mux.HandleFunc("GET /whoami", whoamiHandler)

	handler := withUserMW(mux)

	log.Println("listening on :8080")
	log.Println("thử: curl 'http://localhost:8080/slow?ms=3000'  rồi Ctrl+C giữa chừng")
	log.Println("    curl http://localhost:8080/timeout")
	log.Println("    curl -H 'X-User: alice' http://localhost:8080/whoami")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
