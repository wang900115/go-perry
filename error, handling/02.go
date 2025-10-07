package main

import "fmt"

type DivisionError struct {
	A, B int
}

func (e *DivisionError) Error() string {
	return fmt.Sprintf("無法除以 %d 除以 %d", e.A, e.B)
}

func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, &DivisionError{A: a, B: b}
	}
	return a / b, nil
}

func main() {
	if result, err := divide(10, 0); err != nil {
		fmt.Println("錯誤:", err)
	} else {
		fmt.Println("結果:", result)
	}
}
