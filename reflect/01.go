package main

import (
	"fmt"
	"reflect"
)

func main() {
	var num int
	typ := reflect.TypeOf(num)   // 取得變數的型態
	val := reflect.ValueOf(num)  // 取得變數的值
	zeroVal := reflect.Zero(typ) // 設定該變數為 zero value

	fmt.Printf("Type: %v\n", typ)
	fmt.Printf("Value: %v\n", val)
	fmt.Printf("Zero Value: %v\n", zeroVal)

}
