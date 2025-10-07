package main

import (
	"fmt"
	"os"
	"runtime/trace"
)

func main() {
	f, _ := os.Create("trace.out")
	trace.Start(f)
	defer trace.Stop()
	for range 1000 {
		fmt.Println("Hello world")
	}
}

/*
Profile : application resource situation. cpu, memory, etc. Identify hotspots and inefficiencies.
Trace : application executaion time and events-detect race conditions and multithreaded issues.
*/
