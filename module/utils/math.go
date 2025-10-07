package utils

import "errors"

func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("除數不能為零")
	}
	return a / b, nil
}
