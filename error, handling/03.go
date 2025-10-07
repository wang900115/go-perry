package main

import "fmt"

func riskyFunction() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("恢復自:", r)
		}
	}()
	panic("發生嚴重錯誤")
}

func main() {
	fmt.Println("開始")
	riskyFunction()
	fmt.Println("結束")

}
