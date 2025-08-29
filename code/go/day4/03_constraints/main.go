package main

import "fmt"

type Number interface {
	~int | ~int64 | ~float64
}

func Sum[T Number](xs []T) T {
	var s T
	for _, v := range xs { s += v }
	return s
}

func main() {
	fmt.Println("Sum int:", Sum([]int{1,2,3}))
	fmt.Println("Sum float64:", Sum([]float64{1.5,2.5}))
} 
