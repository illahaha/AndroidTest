package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"work"
)

func usage() {
	fmt.Println("[device config] [target config]")
}
func main() {

	lenOfArgs := len(os.Args)
	if lenOfArgs != 3 {
		fmt.Println("args error")
		usage()
		return
	}

	// use all the machine's cores
	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Println("main thread")

	var w sync.WaitGroup

	w.Add(1)

	go func() {
		work.Work()
		w.Done()
	}()

	w.Wait()
}
