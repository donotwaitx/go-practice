// 11 - Errors
// Go dùng explicit error return values thay vì exceptions.
package main

import (
	"errors"
	"fmt"
)

// Sentinel error — error có thể so sánh
var ErrNotFound = errors.New("không tìm thấy")

// Custom error type
type ValidationError struct {
	Field string
	Msg   string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation: %s - %s", e.Field, e.Msg)
}

func findUser(id int) (string, error) {
	users := map[int]string{1: "Alice", 2: "Bob"}
	name, ok := users[id]
	if !ok {
		return "", ErrNotFound
	}
	return name, nil
}

func validateAge(age int) error {
	if age < 0 {
		return &ValidationError{Field: "age", Msg: "phải >= 0"}
	}
	if age > 150 {
		// Wrap error để giữ context
		return fmt.Errorf("validate age %d: %w", age, &ValidationError{Field: "age", Msg: "quá lớn"})
	}
	return nil
}

func main() {
	// Pattern chuẩn của Go: check error ngay sau khi gọi
	name, err := findUser(3)
	if err != nil {
		// errors.Is — so sánh với sentinel error (kể cả khi đã wrap)
		if errors.Is(err, ErrNotFound) {
			fmt.Println("user 3: không tồn tại")
		} else {
			fmt.Println("error:", err)
		}
	} else {
		fmt.Println("user 3:", name)
	}

	// Happy path
	if name, err := findUser(1); err == nil {
		fmt.Println("user 1:", name)
	}

	// errors.As — extract error theo type cụ thể
	err = validateAge(200)
	var vErr *ValidationError
	if errors.As(err, &vErr) {
		fmt.Println("got validation error on field:", vErr.Field)
	}
	fmt.Println("full:", err)

	// panic chỉ dùng cho lỗi bất thường không recover được
	// defer + recover có thể bắt panic — nhưng không nên lạm dụng
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered:", r)
		}
	}()
	// panic("oops") // bỏ comment để thử
}
