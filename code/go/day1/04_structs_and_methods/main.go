package main

import "fmt"

type User struct {
	ID   int
	Name string
	Age  int
}

func NewUser(id int, name string, age int) *User {
	return &User{ID: id, Name: name, Age: age}
}

// 值接收者：不修改原对象
func (u User) Greeting() string {
	return fmt.Sprintf("Hi, I'm %s (%d)", u.Name, u.Age)
}

// 指针接收者：可修改原对象
func (u *User) Rename(name string) {
	u.Name = name
}

func (u User) IsAdult() bool { return u.Age >= 18 }

// 结构体
// 定义一个名为 Person 的结构体，它有 Name 和 Age 两个字段
// 相当于 JavaScript 中的 const person = { name: "...", age: ... };
type Person struct {
	Name string
	Age  int
}

// 这是一个使用值接收者的方法。它不会修改原始结构体。
// 我们可以称之为“只读”方法。
func (p Person) SayHello() {
	fmt.Printf("Hello, my name is %s and I am %d years old.\n", p.Name, p.Age)
}

// 这是一个使用指针接收者的方法。它会修改原始结构体。
// 这种方式是 Go 语言中修改结构体字段的惯用做法。
func (p *Person) GrowUp() {
	p.Age++
	fmt.Printf("%s has grown up and is now %d years old.\n", p.Name, p.Age)
}

// 这是一个普通函数，它接收一个 Person 结构体的副本作为参数
func changeNameByValue(p Person, newName string) {
	fmt.Println("--- 在 changeNameByValue 函数内部 ---")
	p.Name = newName
	fmt.Printf("函数内部: 姓名被修改为 %s\n", p.Name)
	fmt.Println("---------------------------------")
}

// 这是一个普通函数，它接收一个指向 Person 结构体的指针作为参数
func changeNameByPointer(p *Person, newName string) {
	fmt.Println("--- 在 changeNameByPointer 函数内部 ---")
	// 通过指针 p，我们可以访问并修改原始的结构体
	p.Name = newName
	fmt.Printf("函数内部: 姓名被修改为 %s\n", p.Name)
	fmt.Println("-----------------------------------")
}

func structExample() {
	// 1. 创建一个 Person 结构体实例
	person := Person{Name: "Alice", Age: 25}
	fmt.Println("原始结构体:", person)

	// 2. 使用值接收者的方法，不会修改原始结构体
	person.SayHello()

	// 3. 将 person 结构体传给一个普通函数
	// 注意，这里传递的是一个副本
	changeNameByValue(person, "Bob")
	fmt.Println("函数调用后, 原始结构体:", person) // 原始结构体没有被修改

	fmt.Println("\n--- 开始使用指针 ---")

	// 4. 使用 & 运算符获取 person 结构体的指针
	personPointer := &person
	fmt.Println("person 结构体的指针:", personPointer)

	// 5. 将 personPointer 传给一个函数
	// 函数通过指针可以修改原始结构体
	changeNameByPointer(personPointer, "Charlie")
	fmt.Println("函数调用后, 原始结构体:", person) // 原始结构体已被修改

	// 6. 使用指针接收者的方法。
	// Go 编译器会自动将 person 转换为指针 &person 来调用该方法
	person.GrowUp()
	fmt.Println("方法调用后, 原始结构体:", person) // 原始结构体已被修改
}

func main() {
	u := NewUser(1, "Alice", 20)
	fmt.Println(u.Greeting(), "adult=", u.IsAdult())
	u.Rename("Alicia")
	fmt.Println(u.Greeting())
	
	// 结构体
	structExample()
} 
