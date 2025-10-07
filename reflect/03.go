package main

import (
	"fmt"
	"reflect"
)

type Calculator struct{}

func (c Calculator) Add(a, b int) int {
	return a + b
}
func main() {
	calculator := Calculator{}

	// 使用 reflect.ValueOf 取得計算器物件的reflect value
	val := reflect.ValueOf(calculator)

	// 使用 MethodByName 方法取得名為 "Add" 的方法的reflect value
	method := val.MethodByName("Add")

	// 準備方法的參數
	args := []reflect.Value{reflect.ValueOf(5), reflect.ValueOf(3)}

	// 調用方法並取得結果，然後轉換為整數型別
	result := method.Call(args)[0].Interface().(int)

	fmt.Printf("Result: %d\n", result) // Result: 8
}
