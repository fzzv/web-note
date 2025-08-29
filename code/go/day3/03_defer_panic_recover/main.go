package main

import (
	"fmt"
)

func mayPanic(n int) {
	defer fmt.Println("defer in mayPanic runs (LIFO)")
	if n < 0 {
		panic("bad input")
	}
	fmt.Println("ok:", n)
}

func safeCall() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered:", r)
		}
	}()
	mayPanic(-1)
}

func main() {
	defer fmt.Println("defer in main")
	safeCall()
} 
