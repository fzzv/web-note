package main

import (
	"fmt"
	"os"
)

func addAndDiff(a, b int) (sum int, diff int) {
	sum = a + b
	diff = a - b
	return
}

func sum(values ...int) int {
	s := 0
	for _, v := range values {
		s += v
	}
	return s
}

func readFile(path string) ([]byte, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func main() {
	s, d := addAndDiff(10, 3)
	fmt.Println("sum, diff:", s, d)
	nums := []int{1, 2, 3, 4}
	fmt.Println("sum variadic:", sum(nums...))

	if _, err := readFile("not-exist.txt"); err != nil {
		fmt.Println("read error:", err)
	}
} 
