# Day4：泛型（Generics）与常用范式

## 今日目标
- 掌握泛型函数与泛型类型的声明与使用
- 理解约束（constraints）与类型集（type set）
- 以泛型实现常用集合算法：Map/Filter/Reduce、Set
- 避免常见坑位：比较、约束过度、装箱/逃逸导致的性能问题

## 学习步骤
1. 阅读并运行 `code/go/day4` 5 个示例
2. 完成本页“练习与思考”
3. 回顾官方博客：`https://go.dev/blog/intro-generics`、`https://go.dev/blog/constraints`

## 练习与思考
- 在 `04_map_filter_reduce` 中增加 `Find` 与 `GroupBy`
- 在 `05_set` 中增加 `Intersect`、`Diff` 与 `IsSubset`
- 在 `03_constraints` 中扩展数值约束，支持 `~int8` 等底层类型

## 运行指引
```bash
cd code/go/day4
(go mod init web-note/day4) 2> NUL || echo ok
go run ./01_generic_functions
go run ./02_generic_types
go run ./03_constraints
go run ./04_map_filter_reduce
go run ./05_set
``` 
