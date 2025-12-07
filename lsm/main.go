package main

import (
	"fmt"
	"lsm/lsm"
	"os"
)

func main() {
	dir := "./data_lsm_example"
	os.RemoveAll(dir)
	os.MkdirAll(dir, 0o755)

	lsm, err := lsm.OpenLSM(dir)
	if err != nil {
		panic(err)
	}

	// Put some values
	lsm.Put("apple", []byte("red"))
	lsm.Put("banana", []byte("yellow"))
	lsm.Put("cherry", []byte("darkred"))

	// flush mem -> sst
	if err := lsm.FlushMemToSST(); err != nil {
		panic(err)
	}

	lsm.Put("banana", []byte("green")) // override
	lsm.Put("date", []byte("brown"))

	// flush again
	if err := lsm.FlushMemToSST(); err != nil {
		panic(err)
	}

	// get values
	keys := []string{"apple", "banana", "cherry", "date", "fig"}
	for _, k := range keys {
		v, ok, err := lsm.Get(k)
		if err != nil {
			panic(err)
		}
		if ok {
			fmt.Printf("%s => %s\n", k, string(v))
		} else {
			fmt.Printf("%s => <nil>\n", k)
		}
	}

	// compact the two newest SST
	if err := lsm.CompactTwoNewest(); err != nil {
		panic(err)
	}

	fmt.Println("After compaction:")
	for _, k := range keys {
		v, ok, err := lsm.Get(k)
		if err != nil {
			panic(err)
		}
		if ok {
			fmt.Printf("%s => %s\n", k, string(v))
		} else {
			fmt.Printf("%s => <nil>\n", k)
		}
	}
}
