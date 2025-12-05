package main

import (
	"fmt"
	"log"
	"reflect"
	"unsafe"
)

/*
 * unsafe.Pointer 可以被GC追蹤
 * uintptr 不會被GC追蹤
 */

func voidPointer() {
	var a int = 10
	p := unsafe.Pointer(&a)
	log.Println(*(*int)(p))
	log.Println(*(*float64)(p))
}

type User struct {
	A int8
	B int32
	C int8
}

func memoryAlign() {
	t := reflect.TypeOf(User{})
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		log.Printf("Field %s: offset=%d, size=%d, align=%d\n", f.Name, f.Offset, f.Type.Size(), f.Type.Align())
	}
	log.Println("Size of User struct:", t.Size())
}

func Dump(ptr unsafe.Pointer, size uintptr) {
	// unsafe.Slice 可從 pointer 建一個 byte slice
	b := unsafe.Slice((*byte)(ptr), size)

	fmt.Printf("Raw memory (%d bytes):\n", size)
	for i := 0; i < len(b); i++ {
		if i%8 == 0 {
			fmt.Printf("\n0x%02X: ", i)
		}
		fmt.Printf("%02X ", b[i])
	}
	fmt.Println()
}

func main() {
	u := User{
		A: 1,
		B: 0x11223344,
		C: 2,
	}
	Dump(unsafe.Pointer(&u), unsafe.Sizeof(u))
}

func pointerOffset() {
	var arr = [5]int{10, 20, 30, 40, 50}
	basePtr := unsafe.Pointer(&arr[0])
	size := unsafe.Sizeof(arr[0])
	log.Println(uintptr(basePtr))
	log.Printf("Array base address: %p, element size: %d\n", basePtr, size)

	for i := 0; i < len(arr); i++ {
		elemPtr := unsafe.Pointer(uintptr(basePtr) + uintptr(i)*size)
		log.Printf("Element %d: %d\n", i, *(*int)(elemPtr))
	}
}

func zeroCopyByteSliceToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func zeroCopyStringToByteSlice(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
