// 10 - Testing HTTP với httptest
// File main demo. Test thật nằm trong main_test.go.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type Calculator struct{}

type CalcReq struct {
	A, B int
	Op   string
}

type CalcResp struct {
	Result int    `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func (c *Calculator) Handle(w http.ResponseWriter, r *http.Request) {
	var req CalcReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, CalcResp{Error: "invalid json"})
		return
	}

	var result int
	switch strings.ToLower(req.Op) {
	case "add":
		result = req.A + req.B
	case "sub":
		result = req.A - req.B
	case "mul":
		result = req.A * req.B
	case "div":
		if req.B == 0 {
			respondJSON(w, http.StatusBadRequest, CalcResp{Error: "divide by zero"})
			return
		}
		result = req.A / req.B
	default:
		respondJSON(w, http.StatusBadRequest, CalcResp{Error: "unknown op: " + req.Op})
		return
	}

	respondJSON(w, http.StatusOK, CalcResp{Result: result})
}

func respondJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func NewServer() http.Handler {
	mux := http.NewServeMux()
	c := &Calculator{}
	mux.HandleFunc("POST /calc", c.Handle)
	return mux
}

func main() {
	log.Println("listening on :8080 — chạy `go test` để xem test")
	log.Fatal(http.ListenAndServe(":8080", NewServer()))
}
