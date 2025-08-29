package main

import "fmt"

type Pair[T any, U any] struct { First T; Second U }

func (p Pair[T, U]) String() string { return fmt.Sprintf("(%v,%v)", p.First, p.Second) }

func main() {
	p1 := Pair[int,string]{First:1, Second:"one"}
	p2 := Pair[string,float64]{First:"pi", Second:3.14}
	fmt.Println(p1.String())
	fmt.Println(p2.String())
} 
