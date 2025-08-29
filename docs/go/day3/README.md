# Day3：接口、错误、defer/panic/recover、类型断言与 IO 接口

## 今日目标
- 掌握接口的隐式实现、接口组合与依赖倒置
- 规范化错误处理与错误包装，清楚何时使用 panic/recover
- 熟练使用 `defer`，理解其执行时机与常见坑
- 理解类型断言与 type switch 的使用场景
- 掌握 `io.Reader`/`io.Writer` 抽象与常见用法

## 学习步骤
1. 阅读并运行 `code/go/day3` 下 5 个示例
2. 完成本页“练习与思考”
3. 回顾官方文档：`https://go.dev/doc/effective_go`（Errors、Panic、Defer、Interfaces）

## 练习与思考
- 为接口实现新增一个 mock，并以依赖注入方式在 main 中切换
- 用 `%w` 包装错误，并用 `errors.Is/As` 判断根因
- 修改 defer 顺序与闭包变量，观察输出差异并解释
- 实现一个通用的 `CopyN(r io.Reader, w io.Writer, n int64)`

## 运行指引
```bash
cd code/go/day3
# 初始化模块
(go mod init web-note/day3) 2> NUL || echo ok
# 运行示例
go run ./01_interfaces
go run ./02_errors
go run ./03_defer_panic_recover
go run ./04_type_assert_switch
go run ./05_io_interface
``` 
