# Day8 全面解析笔记：Go 并发编程基础

> 代码参考：`code/go/day8`
> 示例涵盖：goroutine、channel、select、context、同步原语、内存模型

## 1. Goroutine 基础

### 1.1 什么是 Goroutine
Goroutine 是 Go 语言的轻量级线程，由 Go 运行时管理：

```go
// 启动 goroutine
go func() {
    fmt.Println("在 goroutine 中执行")
}()

// 带参数的 goroutine
go func(name string) {
    fmt.Printf("Hello, %s\n", name)
}("World")

// 调用函数作为 goroutine
go processData()
```

### 1.2 Goroutine 特点
- **轻量级**：初始栈大小仅 2KB，可动态增长
- **高并发**：单个程序可运行数十万个 goroutine
- **调度器**：Go 运行时的 M:N 调度器管理 goroutine
- **抢占式**：Go 1.14+ 支持异步抢占

### 1.3 Goroutine 生命周期
```go
func demonstrateGoroutineLifecycle() {
    fmt.Println("主 goroutine 开始")
    
    // 启动子 goroutine
    go func() {
        fmt.Println("子 goroutine 开始")
        time.Sleep(1 * time.Second)
        fmt.Println("子 goroutine 结束")
    }()
    
    // 主 goroutine 继续执行
    fmt.Println("主 goroutine 继续")
    time.Sleep(2 * time.Second) // 等待子 goroutine 完成
    fmt.Println("主 goroutine 结束")
}
```

---

## 2. Channel 通信

### 2.1 Channel 基础
Channel 是 goroutine 之间通信的管道：

```go
// 创建 channel
ch := make(chan int)        // 无缓冲 channel
buffered := make(chan int, 5) // 有缓冲 channel

// 发送数据
ch <- 42

// 接收数据
value := <-ch
value, ok := <-ch // ok 表示 channel 是否关闭

// 关闭 channel
close(ch)
```

### 2.2 无缓冲 vs 有缓冲 Channel
```go
// 无缓冲 channel：同步通信
unbuffered := make(chan int)
go func() {
    unbuffered <- 1 // 阻塞直到有接收者
}()
value := <-unbuffered // 阻塞直到有发送者

// 有缓冲 channel：异步通信
buffered := make(chan int, 3)
buffered <- 1 // 不阻塞，缓冲区未满
buffered <- 2
buffered <- 3
// buffered <- 4 // 会阻塞，缓冲区已满
```

### 2.3 Channel 方向
```go
// 只发送 channel
func sender(ch chan<- int) {
    ch <- 42
}

// 只接收 channel
func receiver(ch <-chan int) {
    value := <-ch
    fmt.Println(value)
}

// 双向 channel
func bidirectional(ch chan int) {
    ch <- 1
    value := <-ch
}
```

---

## 3. Select 语句

### 3.1 基本用法
Select 用于处理多个 channel 操作：

```go
select {
case msg1 := <-ch1:
    fmt.Println("从 ch1 接收:", msg1)
case msg2 := <-ch2:
    fmt.Println("从 ch2 接收:", msg2)
case ch3 <- 42:
    fmt.Println("向 ch3 发送 42")
default:
    fmt.Println("没有 channel 准备好")
}
```

### 3.2 超时处理
```go
select {
case result := <-ch:
    fmt.Println("接收到结果:", result)
case <-time.After(5 * time.Second):
    fmt.Println("操作超时")
}
```

### 3.3 非阻塞操作
```go
select {
case ch <- value:
    fmt.Println("发送成功")
default:
    fmt.Println("channel 已满，发送失败")
}

select {
case value := <-ch:
    fmt.Println("接收到:", value)
default:
    fmt.Println("channel 为空，接收失败")
}
```

---

## 4. Context 包

### 4.1 Context 基础
Context 用于传递取消信号、超时和请求范围的值：

```go
// 创建 context
ctx := context.Background()           // 根 context
ctx = context.TODO()                  // 占位符 context

// 带取消的 context
ctx, cancel := context.WithCancel(ctx)
defer cancel()

// 带超时的 context
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()

// 带截止时间的 context
deadline := time.Now().Add(10 * time.Second)
ctx, cancel := context.WithDeadline(ctx, deadline)
defer cancel()
```

### 4.2 Context 传值
```go
// 存储值
ctx := context.WithValue(context.Background(), "userID", 12345)

// 获取值
if userID, ok := ctx.Value("userID").(int); ok {
    fmt.Printf("用户 ID: %d\n", userID)
}
```

### 4.3 Context 最佳实践
```go
func processWithContext(ctx context.Context, data []int) error {
    for i, item := range data {
        // 检查是否被取消
        select {
        case <-ctx.Done():
            return ctx.Err() // 返回取消原因
        default:
        }
        
        // 处理数据
        if err := processItem(item); err != nil {
            return err
        }
        
        // 定期检查取消信号
        if i%100 == 0 {
            if ctx.Err() != nil {
                return ctx.Err()
            }
        }
    }
    return nil
}
```

---

## 5. 同步原语

### 5.1 WaitGroup
等待一组 goroutine 完成：

```go
var wg sync.WaitGroup

for i := 0; i < 5; i++ {
    wg.Add(1) // 增加计数
    go func(id int) {
        defer wg.Done() // 完成时减少计数
        fmt.Printf("Worker %d 完成\n", id)
    }(i)
}

wg.Wait() // 等待所有 goroutine 完成
fmt.Println("所有 worker 完成")
```

### 5.2 Mutex（互斥锁）
保护共享资源：

```go
type Counter struct {
    mu    sync.Mutex
    value int
}

func (c *Counter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

func (c *Counter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}
```

### 5.3 RWMutex（读写锁）
允许多个读者或一个写者：

```go
type SafeMap struct {
    mu   sync.RWMutex
    data map[string]int
}

func (sm *SafeMap) Get(key string) (int, bool) {
    sm.mu.RLock()
    defer sm.mu.RUnlock()
    value, ok := sm.data[key]
    return value, ok
}

func (sm *SafeMap) Set(key string, value int) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    sm.data[key] = value
}
```

### 5.4 Once
确保函数只执行一次：

```go
var once sync.Once
var instance *Singleton

func GetInstance() *Singleton {
    once.Do(func() {
        instance = &Singleton{}
    })
    return instance
}
```

---

## 6. 并发模式

### 6.1 生产者-消费者模式
```go
func producerConsumer() {
    ch := make(chan int, 10)
    
    // 生产者
    go func() {
        defer close(ch)
        for i := 0; i < 100; i++ {
            ch <- i
        }
    }()
    
    // 消费者
    for value := range ch {
        fmt.Printf("消费: %d\n", value)
    }
}
```

### 6.2 扇出-扇入模式
```go
func fanOutFanIn(input <-chan int) <-chan int {
    // 扇出：启动多个 worker
    workers := make([]<-chan int, 3)
    for i := 0; i < 3; i++ {
        worker := make(chan int)
        workers[i] = worker
        
        go func(ch chan<- int) {
            defer close(ch)
            for value := range input {
                ch <- value * value // 处理数据
            }
        }(worker)
    }
    
    // 扇入：合并结果
    output := make(chan int)
    var wg sync.WaitGroup
    
    for _, worker := range workers {
        wg.Add(1)
        go func(ch <-chan int) {
            defer wg.Done()
            for value := range ch {
                output <- value
            }
        }(worker)
    }
    
    go func() {
        wg.Wait()
        close(output)
    }()
    
    return output
}
```

---

## 7. 内存模型与竞态条件

### 7.1 数据竞争
```go
// 错误：数据竞争
var counter int
go func() { counter++ }()
go func() { counter++ }()

// 正确：使用同步
var mu sync.Mutex
var counter int
go func() {
    mu.Lock()
    counter++
    mu.Unlock()
}()
```

### 7.2 Happens-Before 关系
Go 内存模型定义了内存操作的顺序：

- Channel 操作
- Mutex 操作
- Once 操作
- 原子操作

### 7.3 竞态检测
```bash
# 使用竞态检测器
go run -race main.go
go test -race ./...
go build -race
```

---

## 8. 错误处理与并发

### 8.1 errgroup 包
```go
import "golang.org/x/sync/errgroup"

func processWithErrGroup(ctx context.Context, items []string) error {
    g, ctx := errgroup.WithContext(ctx)
    
    for _, item := range items {
        item := item // 避免闭包陷阱
        g.Go(func() error {
            return processItem(ctx, item)
        })
    }
    
    return g.Wait() // 等待所有 goroutine 完成
}
```

### 8.2 错误收集
```go
func collectErrors(items []string) []error {
    errCh := make(chan error, len(items))
    var wg sync.WaitGroup
    
    for _, item := range items {
        wg.Add(1)
        go func(item string) {
            defer wg.Done()
            if err := processItem(item); err != nil {
                errCh <- err
            }
        }(item)
    }
    
    go func() {
        wg.Wait()
        close(errCh)
    }()
    
    var errors []error
    for err := range errCh {
        errors = append(errors, err)
    }
    
    return errors
}
```

---

## 9. 性能考虑

### 9.1 Goroutine 池
```go
type WorkerPool struct {
    workers int
    jobs    chan Job
    results chan Result
}

func NewWorkerPool(workers int) *WorkerPool {
    return &WorkerPool{
        workers: workers,
        jobs:    make(chan Job, 100),
        results: make(chan Result, 100),
    }
}

func (wp *WorkerPool) Start() {
    for i := 0; i < wp.workers; i++ {
        go wp.worker()
    }
}

func (wp *WorkerPool) worker() {
    for job := range wp.jobs {
        result := job.Process()
        wp.results <- result
    }
}
```

### 9.2 避免 Goroutine 泄漏
```go
func avoidGoroutineLeak() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel() // 确保取消
    
    ch := make(chan int)
    go func() {
        select {
        case ch <- computeValue():
        case <-ctx.Done():
            return // 避免泄漏
        }
    }()
    
    select {
    case result := <-ch:
        fmt.Println("结果:", result)
    case <-ctx.Done():
        fmt.Println("超时")
    }
}
```

---

## 10. 下一步学习

Day9 将学习：
- 高级并发模式
- Worker Pool 实现
- Pipeline 模式
- 背压处理
- 取消传播

重点关注：
- 并发安全设计
- 性能优化
- 错误处理策略
- 资源管理
