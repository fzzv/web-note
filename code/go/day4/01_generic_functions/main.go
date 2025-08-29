package main

import "fmt"

// Min 返回可比较类型中的较小值
func Min[T ~int | ~int64 | ~float64 | ~string](a, b T) T {
	if a < b { return a }
	return b
}

// Map 对切片进行映射
func Map[T any, R any](in []T, f func(T) R) []R {
	out := make([]R, 0, len(in))
	for _, v := range in { out = append(out, f(v)) }
	return out
}

func main() {
	fmt.Println("Min:", Min(3, 5), Min(3.2, -1.1), Min("b", "a"))
	fmt.Println("Map:", Map([]int{1,2,3}, func(x int) int { return x*x }))
} 
