package main

import "fmt"

type Counter struct{ n int }

func (c Counter) Value() int { // 值接收者：读取
	return c.n
}

func (c *Counter) Inc() { // 指针接收者：修改
	c.n++
}

func main() {
	c := Counter{n: 0}
	fmt.Println("value:", c.Value())
	c.Inc()
	fmt.Println("value after inc:", c.Value())
} 
