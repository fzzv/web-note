# Day2：函数、方法、指针/结构体与包可见性

## 今日目标
- 深入函数与多返回值、命名返回值、可变参数
- 掌握方法接收者（值/指针）差异与选择
- 理解指针与结构体的内存语义与初始化方式
- 了解组合（embedding）与可见性（导出/未导出、`internal/` 机制）

## 学习步骤
1. 阅读 `code/go/day2` 中五个示例，逐个 `go run` 并修改练习
2. 参考《Effective Go》：Functions、Methods、Initialization、Embedding、Packages章节
3. 完成本页“练习与思考”

## 关键知识点速记
- 多返回值搭配错误处理是 Go 常态
- 值接收者：方法接收一份拷贝；指针接收者：可修改原对象，避免拷贝大结构
- 组合（embedding）是复用与扩展手段，优先于继承
- 可见性：首字母大写导出；`internal/` 只能在父模块及其子包内使用

## 练习与思考
- 将方法 `Rename` 改为值接收者，观察调用者变量是否变化，说明原因
- 为 `stringsutil` 包添加 `Reverse` 与 `IsPalindrome` 并写表驱动测试
- 在 `internal/calc` 中实现 `Average([]int) (float64, error)`，空切片返回错误

## 运行指引
```bash
cd code/go/day2
# 初始化或更新依赖（如果需要第三方包）
go mod init web-note/day2 2> NUL || echo ok

go run ./01_functions
go run ./02_methods
go run ./03_pointers_structs
go run ./04_embedding
go run ./05_packages_visibility
``` 
