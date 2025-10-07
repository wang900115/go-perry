package main

import (
	"fmt"
	"reflect"
)

func main() {
	// 取得整數類型的 reflect.Type
	intType := reflect.TypeOf(0)

	// 建立一個新的整數變數
	newInt := reflect.New(intType).Elem()

	// 設定變數的value
	newInt.SetInt(42)

	// 取得變數的value
	fmt.Println("New Integer Value:", newInt.Int()) // 42
}
