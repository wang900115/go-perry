package main

import (
	"errors"
	"fmt"
)

func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("除數不能為0")
	}
	return a / b, nil
}

func main() {
	result, err := divide(10, 2)
	if err != nil {
		fmt.Println("錯誤:", err)
	} else {
		fmt.Println("結果:", result)
	}

	result, err = divide(10, 0)
	if err != nil {
		fmt.Println("錯誤:", err)
	} else {
		fmt.Println("結果:", result)
	}
}
