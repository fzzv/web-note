package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// --- 第一种方法: 使用 os.Args ---
	// os.Args 是一个字符串切片，包含所有的命令行参数
	fmt.Println("--- 通过 os.Args 获取参数 ---")
	fmt.Printf("程序名称: %s\n", os.Args[0])
	// 从索引 1 开始遍历，获取用户输入的参数
	if len(os.Args) > 1 {
		fmt.Printf("用户输入的参数: %v\n", os.Args[1:])
	} else {
		fmt.Println("没有用户输入的参数。")
	}

	fmt.Println("\n" + "--- 通过 flag 包获取参数 ---")
	// --- 第二种方法: 使用 flag 包 ---
	// 1. 定义标志（Flags）
	// flag.String 用于定义一个 string 类型的标志
	// 参数依次是：标志名称, 默认值, 帮助信息
	name := flag.String("name", "world", "要问候的人的姓名")

	// flag.Int 用于定义一个 int 类型的标志
	age := flag.Int("age", 20, "用户的年龄")

	// flag.BoolVar 用于定义一个 bool 类型的标志
	verbose := flag.Bool("v", false, "是否启用详细模式")

	// 2. 解析标志
	// 这行代码必须在所有标志定义之后调用，它会解析命令行参数并为标志赋值
	flag.Parse()

	// 3. 使用解析后的标志值
	// 你需要使用 * 来解引用标志指针以获取实际值
	fmt.Printf("Hello, %s!\n", *name)
	fmt.Printf("你的年龄是: %d\n", *age)
	if *verbose {
		fmt.Println("已启用详细模式。")
	}

	// flag.Args() 会返回所有未被 flag 包解析的剩余参数
	if len(flag.Args()) > 0 {
		fmt.Printf("未解析的参数: %v\n", flag.Args())
	}
} 
