package main

import (
	"errors"
	"fmt"
)

type DivideByZeroError struct{}

func (e DivideByZeroError) Error() string { return "divide by zero" }

// divide 是一个会返回两个值的函数：结果和错误
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, DivideByZeroError{}
		// 或者
		// return 0, errors.New("divide by zero")
	}
	return a / b, nil
}

func main() {
	// 正常情况
	q, err := divide(10, 2)
	// 检查错误是否为 nil
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println("10 / 2 =", q)

	// 错误情况
	_, err = divide(1, 0)
	if err != nil {
		var dz DivideByZeroError
		if errors.As(err, &dz) {
			fmt.Println("caught divide-by-zero:", dz)
		} else {
			fmt.Println("other error:", err)
		}
	}
} 
