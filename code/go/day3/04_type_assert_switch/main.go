package main

import "fmt"

func describe(i any) {
	// 断言
	if s, ok := i.(string); ok {
		fmt.Println("string:", s)
	}

	// type switch
	switch v := i.(type) {
	case int:
		fmt.Println("int:", v)
	case string:
		fmt.Println("string via switch:", v)
	default:
		fmt.Printf("unknown %T\n", v)
	}
}

func main() {
	describe(42)
	describe("go")
	describe(3.14)
} 
