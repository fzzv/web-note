# Day2 学习笔记（函数、方法、指针/结构体、组合与可见性）

> 代码参考：`code/go/day2`
> 建议：先快速浏览要点，再对照示例运行与改动。

## 1. 核心要点速览
- 函数：多返回值很常见；错误用 `error` 返回值承载；支持可变参数 `...T`；命名返回值谨慎使用。
- 方法：值接收者 = 拷贝；指针接收者 = 可修改原对象并避免大对象拷贝；方法集与接口匹配要注意。
- 指针/结构体：零值可用；构造函数返回指针是常见范式；值拷贝与指针别名区别要清晰。
- 组合（embedding）：方法提升；优先组合而非继承；命名冲突时需显式选择。
- 可见性与包：首字母大写导出；`internal/` 仅父模块内可见；公共包可放 `pkg/`。

## 2. 函数（Functions）
- 多返回值与错误处理：
```go
b, err := os.ReadFile(path)
if err != nil { return nil, err }
```
- 命名返回值：可读性好但避免滥用；适合短函数。
- 可变参数：`func sum(xs ...int) int {}`；调用方可传切片 `sum(nums...)`。

示例：`code/go/day2/01_functions`

## 3. 方法与接收者（Methods）
- 值接收者：适合只读方法、小结构体；不会修改原对象。
- 指针接收者：可修改原对象；避免大对象拷贝；与接口实现一致性常选指针接收者。

示例：`code/go/day2/02_methods`

经验法则：
- 需要修改接收者、包含 `sync.Mutex` 等不可复制字段、或结构体较大 → 指针接收者
- 否则默认值接收者，保持简单

## 4. 指针与结构体（Pointers & Structs）
- 零值可用（尤其是 map/chan 需 `make`）。
- 值拷贝 vs 指针别名：`copy := *u` 与 `alias := u` 意义不同。

示例：`code/go/day2/03_pointers_structs`

## 5. 组合（Embedding）
- 结构体内匿名字段可“提升”其方法到外层。
- 命名冲突需显式选择：`s.Logger.Log("...")`。

示例：`code/go/day2/04_embedding`

## 6. 包与可见性（Packages & Visibility）
- 导出规则：标识符首字母大写 → 对包外可见。
- `internal/`：仅父模块及子包可导入，外部模块不可见。
- 公共复用包可放 `pkg/`（约定俗成）。

示例：
- 内部包：`code/go/day2/internal/calc`（`Sum`、`Average`）
- 公共包：`code/go/day2/pkg/stringsutil`（`Upper`、`Reverse`、`IsPalindrome`）
- 引用演示：`code/go/day2/05_packages_visibility`

## 7. 常见坑位
- 循环变量捕获：goroutine/闭包中先绑定 `v := v`。
- 方法集差异：`T` 与 `*T` 的方法集不同；接口实现取决于方法集。
- `nil` 与空切片：打印一样但行为不同，序列化/反射时需留意。
- `map` 并发写：非线程安全；需要锁或 `sync.Map`。

参考：
- Tour of Go（Functions/Methods）
- Effective Go（Functions、Methods、Embedding、Packages） 
