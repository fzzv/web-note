# Day3 详细笔记：接口、错误、defer/panic/recover、类型断言与 IO

> 代码参考：`code/go/day3`
> 建议：先快速浏览要点，再对照示例运行与改动。

## 1) 接口与依赖倒置（01_interfaces）
- 隐式实现：类型只要实现了接口所需方法即可，无需显式 `implements`（鸭子类型）。
- 依赖倒置：高层逻辑依赖接口而非具体实现（`Alert(n Notifier, ...)`）。
- 方法集与接收者：
  - 若接口方法签名与类型的方法不匹配（值/指针接收者差异），将无法实现接口。
  - 一般规则：值接收者的方法属于 `T` 和 `*T`；指针接收者的方法仅属于 `*T`。
- 接口值结构：包含“动态类型 + 动态值”。注意接口为 `nil` 与“动态值为 `nil`”的区别。

> **什么是依赖倒置原则？**
>
> 它的核心思想可以总结为两点：
>
> 1. **高层模块不应该依赖于低层模块，两者都应该依赖于抽象。**
> 2. **抽象不应该依赖于细节，细节应该依赖于抽象。**

示例片段：
```go
type Notifier interface{ Notify(string) error }
func Alert(n Notifier, msg string) error { return n.Notify(msg) }
```

实践建议：
- 对外暴露接口，对内以构造函数注入具体实现，便于测试与替换。

---

## 2) 错误处理与包装（02_errors）
- 返回 `error` 而非抛异常：Go 将错误作为显式返回值链式传递。
- 包装错误：使用 `fmt.Errorf("...: %w", err)` 增加上下文，保持根因可追溯。
- 判断根因：`errors.Is(err, target)`；类型匹配：`errors.As(err, *targetType)`。
- 哨兵错误（sentinel）与自定义错误类型：
  - 小心哨兵错误的全局耦合；推荐返回语义更清晰的自定义错误或包装上下文。
- 不要滥用 `panic`：仅用于不可恢复的编程错误或初始化致命缺陷。

示例片段：
```go
var ErrNotFound = errors.New("not found")
return nil, fmt.Errorf("read config %q: %w", path, err)
if errors.Is(err, ErrNotFound) { /* ... */ }
```

实践建议：
- 始终在靠近产生错误处添加上下文；在边界（API 层）统一处理与映射错误码。

---

## 3) defer / panic / recover（03_defer_panic_recover）
- `defer`：按照 LIFO 顺序在函数返回前执行；参数在 `defer` 语句处即求值。
- `panic`：触发栈展开，依次执行已登记的 `defer`。
- `recover`：仅在同一 goroutine 的延迟函数中调用才有效，用于捕获 `panic`，避免程序崩溃。
- 使用原则：
  - 资源释放（文件/连接）、指标上报、日志收尾适合用 `defer`。
  - `recover` 只在明确边界（如 goroutine 顶层、进程入口）兜底，不做正常流程控制。

示例片段：
```go
defer func(){ if r := recover(); r != nil { log.Println("recovered", r) } }()
```

常见坑：
- `defer` 闭包捕获外部变量，易被后续修改影响结果；需要时将值拷贝到局部变量。

---

## 4) 类型断言与 type switch（04_type_assert_switch）
- 断言：`v, ok := i.(T)`；失败返回 `ok=false`，若省略 `, ok` 将直接 `panic`。
- type switch：对接口动态类型分支处理，避免多次断言。
- `any`：是 `interface{}` 的别名；更语义化，建议用于“任意类型”场景。

示例片段：
```go
switch v := i.(type) {
case int:    // ...
case string: // ...
default:     // unknown
}
```

实践建议：
- 尽量通过多态/接口消化分支；必要时用 type switch，避免到处 `.(type)` 分支。

---

## 5) IO 接口与流式处理（05_io_interface）
- 抽象：`io.Reader`/`io.Writer` 使组件以流式方式解耦（来源与去向可替换）。
- 常用实现：`strings.Reader`、`bytes.Buffer`、`os.File`、`net.Conn`。
- 拓展函数：`io.Copy`、`io.CopyN`、`io.TeeReader`、`io.LimitedReader`。
- 语义要点：
  - `CopyN` 当源不足 `n` 字节时返回错误（通常 `EOF`/`UnexpectedEOF`），但可能同时写入了部分数据。
  - 处理错误时注意“部分成功”的场景：先用返回的已传输字节数做幂等补偿或回滚。

示例片段：
```go
n, err := io.CopyN(dst, src, 5) // 可能 n>0 且 err!=nil
```

实践建议：
- 尽量以 Reader/Writer 处理大文件/网络数据，避免整块读入内存。

---

## 6) 实战模式与测试建议
- 依赖注入：构造函数/参数注入接口（如 `Notifier`），主流程只依赖抽象，便于替换为 mock。
- 错误分层：底层返回根因并包装，上层根据场景统一转译（HTTP 状态码/业务错误码）。
- `defer` 兜底：在 goroutine 顶层或 `main` 入口做 panic 兜底，打印堆栈与必要上下文。
- IO 单元测试：使用 `bytes.Buffer`/`strings.Reader` 构造可控的 Reader/Writer，无需依赖文件系统。

---

## 7) 检查清单（Checklist）
- 接口：是否避免了具体实现耦合？方法接收者选择是否合理？
- 错误：是否使用 `%w` 包装？是否在边界统一处理？`errors.Is/As` 是否覆盖到位？
- defer：是否存在闭包变量被后续修改的问题？是否有资源忘记释放？
- 断言：是否使用了 `, ok` 防护？type switch 是否覆盖默认分支？
- IO：是否处理了部分拷贝/短读？是否避免了不必要的大块内存？

---

## 8) 可直接复用的代码模板
- 错误包装与判定：
```go
if err != nil { return fmt.Errorf("do X: %w", err) }
if errors.Is(err, ErrNotFound) { /* handle */ }
```
- panic 兜底：
```go
defer func(){ if r := recover(); r != nil { log.Printf("panic: %v", r) } }()
```
- Reader/Writer 适配：
```go
var (
    r io.Reader = strings.NewReader(input)
    w io.Writer = new(bytes.Buffer)
)
```

---

## 9) 进一步练习
- 定义自定义错误类型（实现 `Error()`），并用 `errors.As` 断言具体类型。
- 为 `Notifier` 增加一个基于 Webhook 的实现；通过环境变量选择不同实现运行。
- 为 IO 示例加入限速/限流（`io.LimitedReader` + 自定义 `RateLimiter`）。

参考资料：
- Effective Go（Errors、Defer、Panic、Recover、Interfaces）
- Go by Example（Errors、Interfaces、Readers/Writers） 
