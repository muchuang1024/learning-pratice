package main

import (
	"os"
	"fmt"
	"runtime/trace"
	"unsafe"
)

func main() {

	args := 1.0
	v := unsafe.Sizeof(args)
	fmt.Println(v)

	//创建trace文件
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	//启动trace goroutine
	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	defer trace.Stop()

	//main
	fmt.Println("Hello World")
}