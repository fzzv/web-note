# Day4 全面解析笔记：Go 泛型（Generics）与约束（Constraints）

> 代码参考：`code/go/day4`
> 示例涵盖：泛型函数、泛型类型、约束与类型集、Map/Filter/Reduce、Set

## 1. 基本术语与模型
- **类型参数（Type Parameter）**：函数/类型定义处的尖括号参数，如 `func F[T any](...)`、`type Box[T any] ...`。
- **类型实参（Type Argument）**：调用/实例化时的具体类型，如 `F[int](...)`、`Box[string]{...}`。多数情况下可依靠类型推断省略。
- **约束（Constraint）**：对类型参数可用操作的描述，如 `any`、`~int|~float64`、`comparable`。
- **类型集（Type Set）**：满足某约束的所有类型集合。运算符可用性取决于类型集（例如仅数值类型可 `+`）。
- **近似约束 `~`**：表示“底层类型”为某类型的命名类型也被包含，如 `type MyInt int` 满足 `~int`。

常见内置约束：
- `any`：任意类型（`interface{}` 的别名）
- `comparable`：可作 map 键/可比较的类型（支持 `==`/`!=`）

---

## 2. 泛型函数（01_generic_functions）
文件：`code/go/day4/01_generic_functions/main.go`

要点：
- 在约束允许的前提下使用运算符，例如 `Min` 中使用 `<`，需类型集包含可比较且有序的类型。
- `Map[T,R]` 的返回切片应按输入容量预分配，降低扩容开销。

示例要点：
```go
// 有序最小值（演示 purpose）：
func Min[T ~int | ~int64 | ~float64 | ~string](a, b T) T {
    if a < b { return a }
    return b
}

// Map（保持容量）：
func Map[T any, R any](in []T, f func(T) R) []R {
    out := make([]R, 0, len(in))
    for _, v := range in { out = append(out, f(v)) }
    return out
}
```
实践提示：
- 若需更通用的“可排序”集合，可自定义 `Ordered` 约束，或在外层提供比较器函数 `less(a,b T) bool`。
- 当类型推断失败（参数类型不一致）时，显式写出类型实参：`Min[int](a, b)`。

---

## 3. 泛型类型（02_generic_types）
文件：`code/go/day4/02_generic_types/main.go`

要点：
- 泛型类型可在方法接收者上继续使用其类型参数：`func (p Pair[T,U]) String() string`。
- 实例化：`Pair[int,string]{...}`；类型推断不适用于字面量，需显式标注。

示例要点：
```go
type Pair[T any, U any] struct { First T; Second U }
func (p Pair[T, U]) String() string { return fmt.Sprintf("(%v,%v)", p.First, p.Second) }
```
实践提示：
- 命名约定：`T`、`U`、`K`、`V` 表示元素/键值；`S` 表示切片元素；`E` 表示 Error/Element，根据语义选择更清晰的名字。

---

## 4. 约束与类型集（03_constraints）
文件：`code/go/day4/03_constraints/main.go`

要点：
- 通过类型集限定可用运算符，示例 `Number` 允许 `+`。
- 使用 `~` 支持拥有指定底层类型的命名类型（例如 `type Age int` 也能参与 `Sum`）。

示例要点：
```go
type Number interface { ~int | ~int64 | ~float64 }
func Sum[T Number](xs []T) T { var s T; for _, v := range xs { s += v }; return s }
```
实践提示：
- 约束过宽可能导致在函数体内无法使用需要的运算；约束过窄会降低复用。优先按“需要的操作集合”设计约束。
- 若需比较大小（`<`/`>`），约束需包含有序类型。

---

## 5. 通用集合算法（04_map_filter_reduce）
文件：`code/go/day4/04_map_filter_reduce/main.go`

要点：
- Map/Filter/Reduce 三件套的签名需简洁、零反射。
- 预分配容量；闭包函数 `f`/`pred` 可能捕获外部变量，注意逃逸与分配。

示例要点：
```go
func Map[T any, R any](in []T, f func(T) R) []R
func Filter[T any](in []T, pred func(T) bool) []T
func Reduce[T any, R any](in []T, init R, f func(R, T) R) R
```
扩展建议：
- `Find[T any](in []T, pred func(T) bool) (T, bool)`
- `GroupBy[T any, K comparable](in []T, key func(T) K) map[K][]T`

---

## 6. 泛型 Set（05_set）
文件：`code/go/day4/05_set/main.go`

要点：
- `T` 必须 `comparable` 才能作为 map 键。
- 典型方法：`Add`、`Has`、`Union`，可扩展 `Intersect`、`Diff`、`IsSubset`、`Size`。

示例要点：
```go
type Set[T comparable] map[T]struct{}
func NewSet[T comparable]() Set[T] { return make(Set[T]) }
func (s Set[T]) Add(v T)          { s[v] = struct{}{} }
func (s Set[T]) Has(v T) bool     { _, ok := s[v]; return ok }
```
注意事项：
- `float64` 的 `NaN` 不等于自身，放入集合时需谨慎；业务上可规避或做归一化。
- map 非并发安全，读写需加锁或使用外层保护；遍历顺序不保证稳定。

---

## 7. 性能与工程实践
- 预分配切片容量：`make([]T, 0, len(in))`。
- 约束精确化：为需要的操作留出空间（如 `+`/`<`），避免 `any` 导致无法使用运算符。
- 避免不必要的装箱：若函数内部将 `T` 转为 `any`，可能触发逃逸与额外分配。
- 基准测试：用 `testing.B` 比较泛型实现与特定类型实现的差距；结合 `pprof` 定位热点。

---

## 8. 常见坑位
- 缺少 `comparable` 却在 map 中用作键 → 编译错误。
- 约束未包含所需运算（如 `<`），在函数体内使用相关运算 → 编译错误。
- 类型推断失败：不同类型混用（如 `Min(1, 1.0)`），需显式类型参数或统一类型。
- 使用 `~` 时误解“底层类型”：只有底层类型匹配的命名类型才被包含。

---

## 9. 实战建议与可复用原型
- 比较器风格（更通用的 Min/Max）：
```go
func MinBy[T any](a, b T, less func(T, T) bool) T {
    if less(a, b) { return a }
    return b
}
```
- 可排序约束（自行维护）：
```go
type Ordered interface{ ~int | ~int64 | ~float64 | ~string }
```
- Set 扩展：
```go
func (s Set[T]) Intersect(o Set[T]) Set[T] { r := NewSet[T](); for v := range s { if o.Has(v) { r.Add(v) } } ; return r }
func (s Set[T]) Diff(o Set[T]) Set[T]      { r := NewSet[T](); for v := range s { if !o.Has(v) { r.Add(v) } } ; return r }
```

---

## 10. 命令清单
```bash
cd code/go/day4
(go mod init web-note/day4) 2> NUL || echo ok
# 逐个运行示例
go run ./01_generic_functions
go run ./02_generic_types
go run ./03_constraints
go run ./04_map_filter_reduce
go run ./05_set
```

参考资料：
- Go 官方博客：Generics 入门与约束（`https://go.dev/blog/intro-generics`，`https://go.dev/blog/constraints`）
- Go by Example：Generics（`https://gobyexample.com/generics`） 
