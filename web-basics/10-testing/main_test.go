// Tests cho calculator API — minh họa các pattern phổ biến của Go testing.
package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// === Test 1: dùng httptest.NewRecorder — test handler trực tiếp, KHÔNG khởi server ===
func TestCalc_Add(t *testing.T) {
	body := `{"a":2,"b":3,"op":"add"}`
	req := httptest.NewRequest(http.MethodPost, "/calc", strings.NewReader(body))
	rec := httptest.NewRecorder()

	NewServer().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp CalcResp
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Result != 5 {
		t.Errorf("expected 5, got %d", resp.Result)
	}
}

// === Test 2: Table-driven test — pattern chuẩn của Go ===
func TestCalc_Operations(t *testing.T) {
	cases := []struct {
		name     string
		body     string
		wantCode int
		wantRes  int
		wantErr  string
	}{
		{"add", `{"a":10,"b":5,"op":"add"}`, 200, 15, ""},
		{"sub", `{"a":10,"b":5,"op":"sub"}`, 200, 5, ""},
		{"mul", `{"a":10,"b":5,"op":"mul"}`, 200, 50, ""},
		{"div ok", `{"a":10,"b":5,"op":"div"}`, 200, 2, ""},
		{"div zero", `{"a":10,"b":0,"op":"div"}`, 400, 0, "divide by zero"},
		{"unknown op", `{"a":1,"b":2,"op":"mod"}`, 400, 0, "unknown op: mod"},
		{"bad json", `{not json}`, 400, 0, "invalid json"},
	}

	srv := NewServer()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/calc", strings.NewReader(c.body))
			rec := httptest.NewRecorder()
			srv.ServeHTTP(rec, req)

			if rec.Code != c.wantCode {
				t.Errorf("code: want %d, got %d", c.wantCode, rec.Code)
			}

			var resp CalcResp
			json.NewDecoder(rec.Body).Decode(&resp)
			if resp.Result != c.wantRes {
				t.Errorf("result: want %d, got %d", c.wantRes, resp.Result)
			}
			if resp.Error != c.wantErr {
				t.Errorf("error: want %q, got %q", c.wantErr, resp.Error)
			}
		})
	}
}

// === Test 3: httptest.NewServer — khởi server thật, gọi qua HTTP client ===
func TestCalc_RealServer(t *testing.T) {
	server := httptest.NewServer(NewServer())
	defer server.Close()

	resp, err := http.Post(server.URL+"/calc", "application/json",
		strings.NewReader(`{"a":7,"b":8,"op":"mul"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var result CalcResp
	json.NewDecoder(resp.Body).Decode(&result)
	if result.Result != 56 {
		t.Errorf("want 56, got %d", result.Result)
	}
}

// === Benchmark: đo hiệu năng handler ===
func BenchmarkCalc_Add(b *testing.B) {
	body := `{"a":2,"b":3,"op":"add"}`
	srv := NewServer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/calc", strings.NewReader(body))
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, req)
	}
}
