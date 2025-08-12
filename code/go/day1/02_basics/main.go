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

// 数组和切片
// changeArrayByValue 接收一个数组的副本
func changeArrayByValue(arr [3]int) {
	fmt.Println("--- 在 changeArrayByValue 函数内部 ---")
	arr[0] = 99 // 修改副本的第一个元素
	fmt.Println("函数内数组:", arr)
	fmt.Println("--------------------------------")
}

// changeSliceByReference 接收一个切片的引用
func changeSliceByReference(s []int) {
	fmt.Println("--- 在 changeSliceByReference 函数内部 ---")
	s[0] = 99 // 修改底层数组的第一个元素
	fmt.Println("函数内切片:", s)
	fmt.Println("--------------------------------")
}

func arrayAndSliceExample() {
	fmt.Println("### 数组 (Array) 示例 ###")
	// 声明并初始化一个长度为 3 的数组
	var arr [3]int = [3]int{1, 2, 3}
	fmt.Println("原始数组:", arr)

	// 尝试通过函数修改数组
	changeArrayByValue(arr)
	fmt.Println("函数调用后, 原始数组:", arr) // 原始数组没有被修改

	// 数组是值类型，赋值时会创建副本
	arr2 := arr
	arr2[0] = 100
	fmt.Println("arr2 修改后, arr 原始数组:", arr) // 原始数组没有被修改

	fmt.Println("\n### 切片 (Slice) 示例 ###")
	// 声明并初始化一个切片
	var s2 []int = []int{10, 20, 30}
	fmt.Println("原始切片:", s2)
	fmt.Printf("切片长度: %d, 容量: %d\n", len(s2), cap(s2))

	// 尝试通过函数修改切片
	changeSliceByReference(s2)
	fmt.Println("函数调用后, 原始切片:", s2) // 原始切片被修改了

	// 切片是引用类型，赋值时只复制指针、长度和容量
	s3 := s2
	s3[0] = 100
	fmt.Println("s3 修改后, s2 原始切片:", s2) // s2 原始切片也跟着被修改了

	// 切片的扩容
	s4 := append(s2, 40)
	fmt.Println("追加元素后, 新切片:", s4)
	fmt.Printf("新切片长度: %d, 容量: %d\n", len(s4), cap(s4))	
}

// Map
func mapExample() {
	// 1. 使用字面量创建并初始化一个映射
	// key 类型是 string，value 类型是 int
	students := map[string]int{
		"Alice": 85,
		"Bob":   92,
		"Charlie": 78,
	}
	fmt.Println("\n--- Map 示例 ---")
	fmt.Println("初始化的映射:", students)

	// 2. 增加新的元素
	students["David"] = 95
	fmt.Println("增加元素后的映射:", students)

	// 3. 访问元素
	aliceScore := students["Alice"]
	fmt.Printf("Alice 的分数是: %d\n", aliceScore)

	// 4. 检查键是否存在 (非常重要的 Go 惯用写法)
	// 映射访问会返回两个值：值本身和表示键是否存在的布尔值
	if score, ok := students["Eve"]; ok {
		fmt.Printf("Eve 的分数是: %d\n", score)
	} else {
		fmt.Println("Eve 不在映射中。")
	}

	// 5. 修改元素
	students["Bob"] = 99
	fmt.Println("修改 Bob 的分数后:", students)

	// 6. 删除元素
	delete(students, "Charlie")
	fmt.Println("删除 Charlie 后:", students)

	// 7. 遍历映射 (注意：遍历顺序是无序的)
	fmt.Println("\n--- 遍历映射 ---")
	for name, score := range students {
		fmt.Printf("学生: %s, 分数: %d\n", name, score)
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

	// 结构体
	structExample()

	// 数组和切片
	arrayAndSliceExample()

	// Map
	mapExample()
} 
