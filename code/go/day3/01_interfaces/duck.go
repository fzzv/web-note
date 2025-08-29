package main

import "fmt"

// 定义一个 "Duck" 接口
// 接口定义了行为，但没有实现
type Duck interface {
    Quack()
    Walk()
}

// 定义一个 Bird 结构体，它不是一只真正的鸭子
type Bird struct{}

// Bird 实现了 Quack() 方法
func (b Bird) Quack() {
    fmt.Println("嘎嘎叫...")
}

// Bird 实现了 Walk() 方法
func (b Bird) Walk() {
    fmt.Println("走路...")
}

// 定义一个 Person 结构体
type Person struct{}

// Person 也实现了 Quack() 方法
func (p Person) Quack() {
    fmt.Println("我也能模仿鸭子叫...")
}

// Person 也实现了 Walk() 方法
func (p Person) Walk() {
    fmt.Println("用两条腿走路...")
}

// 定义一个函数，它只接受一个 "Duck" 接口类型的参数
// 这个函数不关心传入的是什么具体类型，只关心它有没有 Quack() 和 Walk() 方法
func SqueezeDuck(d Duck) {
    d.Quack()
    d.Walk()
}

func main() {
    // 尽管 Bird 和 Person 结构体没有显式声明实现了 Duck 接口
    // 但因为它们都拥有 Quack() 和 Walk() 方法
    // 所以它们都可以被当作 Duck 类型来使用

    // 将 Bird 实例传入 SqueezeDuck
    var bird Bird
    SqueezeDuck(bird)

    // 将 Person 实例传入 SqueezeDuck
    var person Person
    SqueezeDuck(person)
}
