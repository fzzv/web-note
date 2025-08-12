package main

import (
	"fmt"
	"math"
)

// 常量
const (
	pi       = 3.1415926535
	_  = iota // 0
	KB = 1 << (10 * iota)
	MB
	GB
)

// 分支与 switch
func branchAndSwitchExample() {
	x := 10
	y := 20

	if x%2 == 0 {
		fmt.Println("x is even")
	} else {
		fmt.Println("x is odd")
	}

	switch y {
	case 1, 2, 3:
		fmt.Println("y in [1..3]")
	case 7:
		fmt.Println("y is lucky 7")
	default:
		fmt.Println("y is other")
	}	
}

// 函数
func addAndDiff(a, b int) (sum int, diff int) {
	sum = a + b
	diff = a - b
	return
}

func sumLoop(n int) int {
	s := 0
	for i := 1; i <= n; i++ {
		s += i
	}
	return s
}

func sumFormula(n int) int {
	return n * (n + 1) / 2
}

func functionExample() {
	x := 10
	y := 20

	s, d := addAndDiff(x, y)
	fmt.Printf("sum=%d diff=%d\n", s, d)

	fmt.Printf("KB=%d MB=%d GB=%d\n", KB, MB, GB)

	fmt.Printf("sumLoop(100)=%d sumFormula(100)=%d\n", sumLoop(100), sumFormula(100))

	s1 := "你好, Go"
	for idx, r := range s1 { // 按 Unicode 码点遍历
			fmt.Printf("idx=%d rune=%c\n", idx, r)
	}
}

// 主函数
func main() {
	var x int = 42
	y := 7
	f := 3.5
	fmt.Printf("x=%d y=%d f=%.2f pi=%.3f\n", x, y, f, pi)

	// 类型转换
	fx := float64(x) + f
	fmt.Printf("fx=%.2f sqrt(x)=%.2f\n", fx, math.Sqrt(float64(x)))

	// 分支与 switch
	branchAndSwitchExample()

	// 函数
	functionExample()
} 
