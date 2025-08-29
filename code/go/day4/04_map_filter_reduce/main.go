package main

import "fmt"

func Map[T any, R any](in []T, f func(T) R) []R {
	out := make([]R, 0, len(in))
	for _, v := range in { out = append(out, f(v)) }
	return out
}

func Filter[T any](in []T, pred func(T) bool) []T {
	out := make([]T, 0, len(in))
	for _, v := range in { if pred(v) { out = append(out, v) } }
	return out
}

func Reduce[T any, R any](in []T, init R, f func(R, T) R) R {
	agg := init
	for _, v := range in { agg = f(agg, v) }
	return agg
}

func main() {
	xs := []int{1,2,3,4,5}
	squares := Map(xs, func(x int) int { return x*x })
	evens := Filter(xs, func(x int) bool { return x%2==0 })
	sum := Reduce(xs, 0, func(acc, x int) int { return acc + x })
	fmt.Println("squares:", squares)
	fmt.Println("evens:", evens)
	fmt.Println("sum:", sum)
} 
