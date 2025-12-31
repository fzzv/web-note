package main

import (
	"fmt"
	"strings"
)

func wordCount(s string) map[string]int {
	m := make(map[string]int)
	for _, w := range strings.Fields(s) {
		m[strings.ToLower(w)]++
	}
	return m
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
		"Alice":   85,
		"Bob":     92,
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

func main() {
	// 切片与底层数组共享
	base := make([]int, 0, 5)
	base = append(base, 1, 2, 3)
	alias := base
	alias[0] = 99 // 修改 alias 会影响 base（共享同一底层数组）
	fmt.Printf("base=%v len=%d cap=%d\n", base, len(base), cap(base))
	fmt.Printf("alias=%v len=%d cap=%d\n", alias, len(alias), cap(alias))

	// 深拷贝避免联动
	safe := append([]int(nil), base...)
	safe[1] = 77
	fmt.Printf("safe=%v (modify safe won't affect base)\n", safe)
	fmt.Printf("base(after safe)=%v\n", base)

	// map 基本操作与存在性检查
	m := map[string]int{"go": 2009}
	v, ok := m["go"]
	fmt.Printf("m['go']=%d ok=%v\n", v, ok)
	if _, exists := m["rust"]; !exists {
		m["rust"] = 2015
	}
	fmt.Println("map:", m)

	// 单词计数
	text := "Go go is expressive, concise, clean, and efficient"
	fmt.Println("wordCount:", wordCount(text))

	// 数组和切片
	arrayAndSliceExample()

	// Map
	mapExample()

	s1 := []int{1, 2, 3}
	s2 := s1
	s2 = append(s2, 4)
	s2[1] = 99
	fmt.Printf("s1=%v s2=%v\n", s1, s2)

	var data []int
	for i := range 2000 {
		data = append(data, i)
		fmt.Printf("len: %d, cap: %d\r\n", len(data), cap(data))
	}
}
