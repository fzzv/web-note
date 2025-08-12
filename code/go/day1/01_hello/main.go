package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	fmt.Println("Hello, Go!")
	fmt.Printf("runtime=%s os=%s arch=%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Printf("now=%s\n", time.Now().Format(time.RFC3339))
} 
