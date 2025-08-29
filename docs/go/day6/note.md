# Day6 全面解析笔记：文件 I/O、JSON 处理与时间操作

> 代码参考：`code/go/day6`
> 示例涵盖：文件操作、JSON 序列化、时间处理、实践项目

## 1. 文件 I/O 操作

### 1.1 基本文件操作
Go 提供了多种文件操作方式，从简单的一次性读写到流式处理：

```go
// 一次性读取整个文件
content, err := os.ReadFile("file.txt")
if err != nil {
    log.Fatal(err)
}

// 一次性写入文件
err = os.WriteFile("output.txt", []byte("Hello World"), 0644)
if err != nil {
    log.Fatal(err)
}
```

### 1.2 文件句柄操作
```go
// 打开文件
file, err := os.Open("input.txt")
if err != nil {
    log.Fatal(err)
}
defer file.Close() // 重要：确保文件关闭

// 创建文件
file, err := os.Create("output.txt")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

// 以追加模式打开
file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
```

### 1.3 缓冲 I/O
使用 `bufio` 包提高 I/O 性能：

```go
// 缓冲读取
file, _ := os.Open("large.txt")
defer file.Close()

scanner := bufio.NewScanner(file)
for scanner.Scan() {
    line := scanner.Text()
    // 处理每一行
}

// 缓冲写入
file, _ := os.Create("output.txt")
defer file.Close()

writer := bufio.NewWriter(file)
defer writer.Flush() // 确保缓冲区内容写入

writer.WriteString("Hello\n")
writer.WriteString("World\n")
```

---

## 2. JSON 处理

### 2.1 结构体标签
Go 使用结构体标签控制 JSON 序列化行为：

```go
type User struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email,omitempty"`    // 空值时省略
    Password string `json:"-"`                  // 永不序列化
    Age      int    `json:"age,string"`         // 转换为字符串
}
```

### 2.2 序列化与反序列化
```go
// 序列化（编码）
user := User{ID: 1, Name: "张三", Email: "zhang@example.com"}
jsonData, err := json.Marshal(user)
if err != nil {
    log.Fatal(err)
}

// 美化输出
jsonData, err := json.MarshalIndent(user, "", "  ")

// 反序列化（解码）
var user User
err := json.Unmarshal(jsonData, &user)
if err != nil {
    log.Fatal(err)
}
```

### 2.3 流式 JSON 处理
对于大型 JSON 数据，使用流式处理：

```go
// 编码到流
file, _ := os.Create("users.json")
defer file.Close()

encoder := json.NewEncoder(file)
encoder.SetIndent("", "  ")
encoder.Encode(users)

// 从流解码
file, _ := os.Open("users.json")
defer file.Close()

decoder := json.NewDecoder(file)
var users []User
decoder.Decode(&users)
```

### 2.4 处理动态 JSON
```go
// 使用 interface{} 处理未知结构
var data interface{}
json.Unmarshal(jsonData, &data)

// 使用 map[string]interface{} 处理对象
var obj map[string]interface{}
json.Unmarshal(jsonData, &obj)

// 类型断言访问值
if name, ok := obj["name"].(string); ok {
    fmt.Println("Name:", name)
}
```

---

## 3. 时间处理

### 3.1 时间创建与格式化
Go 使用独特的时间格式化方式：

```go
// 当前时间
now := time.Now()

// 创建特定时间
t := time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC)

// 解析时间字符串
layout := "2006-01-02 15:04:05"
t, err := time.Parse(layout, "2024-01-01 12:00:00")

// 格式化时间
formatted := now.Format("2006-01-02 15:04:05")
```

### 3.2 时间运算
```go
// 时间加减
future := now.Add(24 * time.Hour)        // 加一天
past := now.Add(-1 * time.Hour)          // 减一小时

// 时间差
duration := future.Sub(now)
fmt.Printf("相差 %v\n", duration)

// 时间比较
if now.Before(future) {
    fmt.Println("now 在 future 之前")
}
```

### 3.3 时区处理
```go
// 加载时区
loc, err := time.LoadLocation("Asia/Shanghai")
if err != nil {
    log.Fatal(err)
}

// 在指定时区创建时间
t := time.Date(2024, 1, 1, 12, 0, 0, 0, loc)

// 转换时区
utc := t.UTC()
local := t.Local()
```

---

## 4. 错误处理最佳实践

### 4.1 错误包装
```go
import "fmt"

// 包装错误，添加上下文
if err != nil {
    return fmt.Errorf("读取配置文件失败: %w", err)
}

// 检查特定错误
if errors.Is(err, os.ErrNotExist) {
    // 文件不存在
}

// 检查错误类型
var pathErr *os.PathError
if errors.As(err, &pathErr) {
    fmt.Printf("路径错误: %s\n", pathErr.Path)
}
```

### 4.2 自定义错误类型
```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("验证失败 [%s]: %s", e.Field, e.Message)
}

// 使用
if email == "" {
    return &ValidationError{
        Field:   "email",
        Message: "邮箱不能为空",
    }
}
```

---

## 5. 资源管理与 defer

### 5.1 defer 的使用模式
```go
func processFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close() // 确保文件关闭

    // 处理文件...
    return nil
}
```

### 5.2 defer 的执行顺序
```go
func example() {
    defer fmt.Println("第一个 defer")
    defer fmt.Println("第二个 defer")
    defer fmt.Println("第三个 defer")
    fmt.Println("函数体")
}
// 输出：
// 函数体
// 第三个 defer
// 第二个 defer
// 第一个 defer
```

### 5.3 defer 与错误处理
```go
func processWithCleanup() (err error) {
    resource, err := acquireResource()
    if err != nil {
        return err
    }
    
    defer func() {
        if cleanupErr := resource.Close(); cleanupErr != nil {
            if err == nil {
                err = cleanupErr
            }
        }
    }()
    
    // 使用资源...
    return nil
}
```

---

## 6. 性能优化技巧

### 6.1 缓冲 I/O
```go
// 低效：多次系统调用
for _, line := range lines {
    file.WriteString(line + "\n")
}

// 高效：使用缓冲
writer := bufio.NewWriter(file)
for _, line := range lines {
    writer.WriteString(line + "\n")
}
writer.Flush()
```

### 6.2 JSON 性能优化
```go
// 预分配切片容量
users := make([]User, 0, expectedCount)

// 使用 json.RawMessage 延迟解析
type Response struct {
    Status string          `json:"status"`
    Data   json.RawMessage `json:"data"`
}
```

### 6.3 字符串构建
```go
// 低效：字符串拼接
var result string
for _, item := range items {
    result += item + "\n"
}

// 高效：使用 strings.Builder
var builder strings.Builder
builder.Grow(estimatedSize) // 预分配容量
for _, item := range items {
    builder.WriteString(item)
    builder.WriteString("\n")
}
result := builder.String()
```

---

## 7. 实践项目：JSON 日志清洗器

### 7.1 项目需求
创建一个工具，用于：
- 读取 JSON 格式的日志文件
- 过滤特定级别的日志
- 格式化输出
- 统计日志信息

### 7.2 核心功能
```go
type LogEntry struct {
    Timestamp time.Time `json:"timestamp"`
    Level     string    `json:"level"`
    Message   string    `json:"message"`
    Source    string    `json:"source"`
}

type LogProcessor struct {
    minLevel LogLevel
    output   io.Writer
}

func (lp *LogProcessor) ProcessFile(filename string) error {
    // 实现日志处理逻辑
}
```

---

## 8. 常见陷阱与解决方案

### 8.1 文件句柄泄漏
```go
// 错误：忘记关闭文件
func badExample() {
    file, _ := os.Open("file.txt")
    // 忘记 defer file.Close()
    // 文件句柄泄漏
}

// 正确：总是使用 defer
func goodExample() {
    file, err := os.Open("file.txt")
    if err != nil {
        return err
    }
    defer file.Close() // 立即添加 defer
}
```

### 8.2 JSON 字段名不匹配
```go
// 问题：字段名不匹配
type User struct {
    UserName string `json:"username"` // JSON 中是 "user_name"
}

// 解决：正确的标签
type User struct {
    UserName string `json:"user_name"`
}
```

### 8.3 时区问题
```go
// 问题：忽略时区
t, _ := time.Parse("2006-01-02 15:04:05", "2024-01-01 12:00:00")
// t 是 UTC 时间，可能不是预期的

// 解决：明确指定时区
loc, _ := time.LoadLocation("Asia/Shanghai")
t, _ := time.ParseInLocation("2006-01-02 15:04:05", "2024-01-01 12:00:00", loc)
```

---

## 9. 测试策略

### 9.1 文件操作测试
```go
func TestFileProcessor(t *testing.T) {
    // 创建临时文件
    tmpfile, err := os.CreateTemp("", "test")
    require.NoError(t, err)
    defer os.Remove(tmpfile.Name())
    
    // 写入测试数据
    testData := "test content"
    _, err = tmpfile.WriteString(testData)
    require.NoError(t, err)
    tmpfile.Close()
    
    // 测试处理逻辑
    result, err := ProcessFile(tmpfile.Name())
    assert.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

### 9.2 时间相关测试
```go
func TestTimeProcessing(t *testing.T) {
    // 固定时间进行测试
    fixedTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
    
    // 使用依赖注入传入时间
    processor := NewProcessor(func() time.Time {
        return fixedTime
    })
    
    result := processor.Process()
    assert.Equal(t, expected, result)
}
```

---

## 10. 下一步学习

Day7 将学习：
- CLI 工具开发
- 命令行参数解析
- 并发编程基础
- 项目实战：完整的 CLI 工具

重点关注：
- 用户体验设计
- 错误处理和用户反馈
- 性能优化
- 代码组织和测试
