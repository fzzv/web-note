package main

import (
	"fmt"
	"web-note/day2/internal/calc"
	strutil "web-note/day2/pkg/stringsutil"
)

func main() {
	fmt.Println("upper:", strutil.Upper("hello"))
	fmt.Println("sum:", calc.Sum([]int{1, 2, 3, 4}))
} 
