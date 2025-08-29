# Day1：Go 基础起步

## 今日目标
- 完成环境与工具验证。
- 掌握变量/常量、基础类型、控制流、切片与 map 的基本用法。
- 能运行与修改 6 个示例程序。

## 先决条件
- 已安装 Go：在终端执行 `go version` 应能看到版本信息。
- 推荐编辑器：VS Code + Go 扩展。

## 学习步骤（约 2–3 小时）
1. 浏览官方教程前半部分：`https://go.dev/tour`（基础语法、流程控制、方法与接口前）
2. 阅读《Effective Go》相关章节：命名、注释、包、错误处理（前半）。
3. 按顺序运行本目录中的示例：
   - `go run ./01_hello`
   - `go run ./02_basics`
   - `go run ./03_slices_and_maps`
   - `go run ./04_structs_and_methods`
   - `go run ./05_errors`
   - `go run ./06_cli_args -name Alice -times 2`
4. 修改代码：在每个示例中完成下列小练习（如下）。

## 小练习
- 01_hello：
  - 输出运行时版本与当前时间（`time.Now()`）。
- 02_basics：
  - 写一个 `sum(n int) int` 计算 1..n 的和，并用循环与公式两种实现对比。
- 03_slices_and_maps：
  - 演示切片共享底层数组导致的联动修改，再用 `append([]T(nil), s...)` 做深拷贝避免。
  - 写一个单词计数器：给定字符串，统计每个单词出现次数（map）。
- 04_structs_and_methods：
  - 为 `User` 增加 `Age` 字段与 `IsAdult()` 方法（≥18）。
- 05_errors：
  - 实现一个 `readJSON(path string) (map[string]any, error)`，文件不存在时返回自定义错误类型并在 `main` 中用 `errors.Is/As` 判定。
- 06_cli_args：
  - 增加 `-uppercase` 布尔参数，控制是否大写输出。

## 验收清单
- 能解释切片的 `len` / `cap`，以及 `append` 导致的扩容与拷贝行为。
- 能编写函数返回值携带 `error` 并在调用端分支处理。
- 能使用 `flag` 包解析命令行参数并提供默认值。

## 延伸阅读
- Go by Example：`https://gobyexample.com/`
- 并发流水线：`https://go.dev/blog/pipelines`

> 下一步：完成 Day1 后，继续 Day2 函数/方法与包的深入与工程化实践。 
