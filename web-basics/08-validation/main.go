// 08 - Validation & Error Responses
// Validate input, trả error JSON có cấu trúc, dùng custom error types.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"strings"
	"unicode/utf8"
)

// === Domain types ===

type CreateUserReq struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Age      int    `json:"age"`
}

// === Validation ===

// FieldError = lỗi của 1 field
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors gom nhiều field error — implement error
type ValidationErrors []FieldError

func (v ValidationErrors) Error() string {
	parts := make([]string, len(v))
	for i, e := range v {
		parts[i] = fmt.Sprintf("%s: %s", e.Field, e.Message)
	}
	return strings.Join(parts, "; ")
}

func (r CreateUserReq) Validate() error {
	var errs ValidationErrors

	if strings.TrimSpace(r.Name) == "" {
		errs = append(errs, FieldError{"name", "required"})
	} else if utf8.RuneCountInString(r.Name) > 50 {
		errs = append(errs, FieldError{"name", "max 50 chars"})
	}

	if _, err := mail.ParseAddress(r.Email); err != nil {
		errs = append(errs, FieldError{"email", "invalid email"})
	}

	if len(r.Password) < 8 {
		errs = append(errs, FieldError{"password", "min 8 chars"})
	}

	if r.Age < 0 || r.Age > 150 {
		errs = append(errs, FieldError{"age", "must be 0..150"})
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

// === Error response ===

type APIError struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Fields  any    `json:"fields,omitempty"`
}

func (e *APIError) Error() string { return e.Message }

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, err error) {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		writeJSON(w, apiErr.Status, apiErr)
		return
	}
	var vErr ValidationErrors
	if errors.As(err, &vErr) {
		writeJSON(w, http.StatusBadRequest, &APIError{
			Code:    "validation_failed",
			Message: "request validation failed",
			Fields:  vErr,
		})
		return
	}
	// Fallback
	writeJSON(w, http.StatusInternalServerError, &APIError{
		Code:    "internal_error",
		Message: err.Error(),
	})
}

// === Handler ===

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateUserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, &APIError{Status: http.StatusBadRequest, Code: "invalid_json", Message: err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		writeErr(w, err)
		return
	}

	// (giả lập) lưu vào DB...
	writeJSON(w, http.StatusCreated, map[string]any{
		"message": "user created",
		"name":    req.Name,
		"email":   req.Email,
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", createUserHandler)

	log.Println("listening on :8080")
	log.Println("thử: curl -X POST localhost:8080/users -d '{\"name\":\"\",\"email\":\"bad\",\"password\":\"123\",\"age\":-1}'")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
