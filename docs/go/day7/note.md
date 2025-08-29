# Day7 全面解析笔记：CLI 工具开发与项目交付

> 代码参考：`code/go/day7`
> 示例涵盖：命令行参数解析、Cobra 框架、用户交互、配置管理、项目打包

## 1. 命令行参数解析基础

### 1.1 标准库 flag 包
Go 标准库提供了基础的命令行参数解析功能：

```go
import "flag"

var (
    name    = flag.String("name", "World", "要问候的名字")
    age     = flag.Int("age", 0, "年龄")
    verbose = flag.Bool("verbose", false, "详细输出")
)

func main() {
    flag.Parse()
    
    if *verbose {
        fmt.Printf("详细模式：问候 %s，年龄 %d\n", *name, *age)
    } else {
        fmt.Printf("Hello, %s!\n", *name)
    }
}
```

### 1.2 参数类型与验证
```go
// 自定义参数类型
type LogLevel string

func (l *LogLevel) String() string {
    return string(*l)
}

func (l *LogLevel) Set(value string) error {
    switch value {
    case "debug", "info", "warn", "error":
        *l = LogLevel(value)
        return nil
    default:
        return fmt.Errorf("无效的日志级别: %s", value)
    }
}

var logLevel LogLevel
flag.Var(&logLevel, "log-level", "日志级别 (debug|info|warn|error)")
```

### 1.3 位置参数处理
```go
func main() {
    flag.Parse()
    
    // 获取非标志参数
    args := flag.Args()
    if len(args) == 0 {
        fmt.Println("请提供至少一个文件名")
        os.Exit(1)
    }
    
    for _, filename := range args {
        processFile(filename)
    }
}
```

---

## 2. Cobra 框架深入

### 2.1 Cobra 基础概念
Cobra 是 Go 生态中最流行的 CLI 框架，被 kubectl、docker、git 等工具使用：

- **Command**：命令，如 `git commit`
- **Flag**：标志，如 `--verbose`
- **Argument**：参数，如文件名

### 2.2 基本结构
```go
import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
    Use:   "myapp",
    Short: "一个简单的 CLI 应用",
    Long:  `这是一个使用 Cobra 构建的示例应用程序`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Hello from myapp!")
    },
}

func main() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "错误: %v\n", err)
        os.Exit(1)
    }
}
```

### 2.3 子命令设计
```go
var createCmd = &cobra.Command{
    Use:   "create [name]",
    Short: "创建新项目",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        name := args[0]
        template, _ := cmd.Flags().GetString("template")
        force, _ := cmd.Flags().GetBool("force")
        
        createProject(name, template, force)
    },
}

func init() {
    createCmd.Flags().StringP("template", "t", "basic", "项目模板")
    createCmd.Flags().BoolP("force", "f", false, "强制覆盖")
    rootCmd.AddCommand(createCmd)
}
```

### 2.4 持久化标志
```go
func init() {
    // 全局标志，所有子命令都可用
    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "配置文件路径")
    rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "详细输出")
    
    // 本地标志，仅当前命令可用
    rootCmd.Flags().BoolP("toggle", "t", false, "帮助信息")
}
```

---

## 3. 配置管理

### 3.1 Viper 配置库
Viper 提供了强大的配置管理功能：

```go
import "github.com/spf13/viper"

func initConfig() {
    if cfgFile != "" {
        viper.SetConfigFile(cfgFile)
    } else {
        home, err := os.UserHomeDir()
        cobra.CheckErr(err)
        
        viper.AddConfigPath(home)
        viper.AddConfigPath(".")
        viper.SetConfigType("yaml")
        viper.SetConfigName(".myapp")
    }
    
    viper.AutomaticEnv() // 自动读取环境变量
    
    if err := viper.ReadInConfig(); err == nil {
        fmt.Println("使用配置文件:", viper.ConfigFileUsed())
    }
}
```

### 3.2 配置优先级
Viper 的配置优先级（从高到低）：
1. 显式调用 Set
2. 命令行标志
3. 环境变量
4. 配置文件
5. 键值存储
6. 默认值

```go
// 绑定标志到配置
viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
viper.BindEnv("author", "MYAPP_AUTHOR")
viper.SetDefault("author", "Anonymous")
```

---

## 4. 用户交互设计

### 4.1 进度条显示
```go
import "github.com/schollz/progressbar/v3"

func processFiles(files []string) {
    bar := progressbar.Default(int64(len(files)))
    
    for _, file := range files {
        processFile(file)
        bar.Add(1)
        time.Sleep(100 * time.Millisecond) // 模拟处理时间
    }
}
```

### 4.2 交互式输入
```go
import "github.com/manifoldco/promptui"

func getUserInput() (string, error) {
    prompt := promptui.Prompt{
        Label:    "请输入项目名称",
        Validate: validateProjectName,
    }
    
    return prompt.Run()
}

func validateProjectName(input string) error {
    if len(input) < 3 {
        return errors.New("项目名称至少需要3个字符")
    }
    return nil
}

func selectTemplate() (string, error) {
    prompt := promptui.Select{
        Label: "选择项目模板",
        Items: []string{"Web API", "CLI Tool", "Library", "Microservice"},
    }
    
    _, result, err := prompt.Run()
    return result, err
}
```

### 4.3 确认对话框
```go
func confirmAction(message string) bool {
    prompt := promptui.Prompt{
        Label:     message,
        IsConfirm: true,
    }
    
    result, err := prompt.Run()
    return err == nil && (result == "y" || result == "Y")
}
```

---

## 5. 错误处理与日志

### 5.1 优雅的错误处理
```go
type CLIError struct {
    Code    int
    Message string
    Cause   error
}

func (e *CLIError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Cause)
    }
    return e.Message
}

func handleError(err error) {
    if cliErr, ok := err.(*CLIError); ok {
        fmt.Fprintf(os.Stderr, "错误: %s\n", cliErr.Message)
        if verbose && cliErr.Cause != nil {
            fmt.Fprintf(os.Stderr, "详细信息: %v\n", cliErr.Cause)
        }
        os.Exit(cliErr.Code)
    } else {
        fmt.Fprintf(os.Stderr, "未知错误: %v\n", err)
        os.Exit(1)
    }
}
```

### 5.2 结构化日志
```go
import "github.com/sirupsen/logrus"

func setupLogger() *logrus.Logger {
    logger := logrus.New()
    
    if verbose {
        logger.SetLevel(logrus.DebugLevel)
    } else {
        logger.SetLevel(logrus.InfoLevel)
    }
    
    logger.SetFormatter(&logrus.TextFormatter{
        FullTimestamp: true,
        ForceColors:   true,
    })
    
    return logger
}
```

---

## 6. 文件操作与数据处理

### 6.1 安全的文件操作
```go
func processFile(filename string) error {
    // 检查文件是否存在
    if _, err := os.Stat(filename); os.IsNotExist(err) {
        return &CLIError{
            Code:    2,
            Message: fmt.Sprintf("文件不存在: %s", filename),
            Cause:   err,
        }
    }
    
    // 检查文件权限
    file, err := os.Open(filename)
    if err != nil {
        return &CLIError{
            Code:    3,
            Message: fmt.Sprintf("无法打开文件: %s", filename),
            Cause:   err,
        }
    }
    defer file.Close()
    
    // 处理文件内容
    return processFileContent(file)
}
```

### 6.2 并发文件处理
```go
func processFilesParallel(files []string, workers int) error {
    jobs := make(chan string, len(files))
    results := make(chan error, len(files))
    
    // 启动工作协程
    for w := 0; w < workers; w++ {
        go func() {
            for filename := range jobs {
                results <- processFile(filename)
            }
        }()
    }
    
    // 发送任务
    for _, filename := range files {
        jobs <- filename
    }
    close(jobs)
    
    // 收集结果
    var errors []error
    for i := 0; i < len(files); i++ {
        if err := <-results; err != nil {
            errors = append(errors, err)
        }
    }
    
    if len(errors) > 0 {
        return fmt.Errorf("处理失败的文件数: %d", len(errors))
    }
    
    return nil
}
```

---

## 7. 测试策略

### 7.1 命令测试
```go
func TestRootCommand(t *testing.T) {
    tests := []struct {
        name     string
        args     []string
        wantErr  bool
        wantCode int
    }{
        {
            name:     "无参数",
            args:     []string{},
            wantErr:  false,
            wantCode: 0,
        },
        {
            name:     "帮助标志",
            args:     []string{"--help"},
            wantErr:  false,
            wantCode: 0,
        },
        {
            name:     "无效标志",
            args:     []string{"--invalid"},
            wantErr:  true,
            wantCode: 1,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            cmd := newRootCommand()
            cmd.SetArgs(tt.args)
            
            err := cmd.Execute()
            
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### 7.2 集成测试
```go
func TestCLIIntegration(t *testing.T) {
    // 创建临时目录
    tmpDir, err := os.MkdirTemp("", "cli-test")
    require.NoError(t, err)
    defer os.RemoveAll(tmpDir)
    
    // 创建测试文件
    testFile := filepath.Join(tmpDir, "test.txt")
    err = os.WriteFile(testFile, []byte("test content"), 0644)
    require.NoError(t, err)
    
    // 执行命令
    cmd := exec.Command("./myapp", "process", testFile)
    output, err := cmd.CombinedOutput()
    
    assert.NoError(t, err)
    assert.Contains(t, string(output), "处理完成")
}
```

---

## 8. 项目打包与分发

### 8.1 构建脚本
```bash
#!/bin/bash
# build.sh

APP_NAME="myapp"
VERSION=$(git describe --tags --always --dirty)
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT=$(git rev-parse HEAD)

# 设置构建标志
LDFLAGS="-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.commit=${COMMIT}"

# 构建不同平台的二进制文件
GOOS=linux GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o dist/${APP_NAME}-linux-amd64
GOOS=darwin GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o dist/${APP_NAME}-darwin-amd64
GOOS=windows GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o dist/${APP_NAME}-windows-amd64.exe
```

### 8.2 版本信息
```go
var (
    version   = "dev"
    buildTime = "unknown"
    commit    = "unknown"
)

var versionCmd = &cobra.Command{
    Use:   "version",
    Short: "显示版本信息",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Printf("%s version %s\n", rootCmd.Use, version)
        fmt.Printf("Built: %s\n", buildTime)
        fmt.Printf("Commit: %s\n", commit)
    },
}
```

### 8.3 安装脚本
```bash
#!/bin/bash
# install.sh

set -e

# 检测操作系统和架构
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64) ARCH="amd64" ;;
    arm64) ARCH="arm64" ;;
    *) echo "不支持的架构: $ARCH"; exit 1 ;;
esac

# 下载并安装
BINARY_NAME="myapp-${OS}-${ARCH}"
DOWNLOAD_URL="https://github.com/user/myapp/releases/latest/download/${BINARY_NAME}"

echo "下载 ${BINARY_NAME}..."
curl -L -o myapp "${DOWNLOAD_URL}"
chmod +x myapp

# 安装到系统路径
sudo mv myapp /usr/local/bin/
echo "安装完成！"
```

---

## 9. 实践项目：文件处理工具

### 9.1 项目需求
创建一个文件处理工具，支持：
- 文件内容搜索和替换
- 批量文件重命名
- 文件统计信息
- 并发处理
- 进度显示

### 9.2 核心功能设计
```go
type FileProcessor struct {
    logger    *logrus.Logger
    config    *Config
    stats     *ProcessStats
}

type Config struct {
    Workers     int
    Pattern     string
    Replacement string
    DryRun      bool
    Recursive   bool
}

type ProcessStats struct {
    FilesProcessed int64
    BytesProcessed int64
    ErrorCount     int64
    StartTime      time.Time
}
```

---

## 10. 性能优化与监控

### 10.1 内存使用优化
```go
// 使用对象池减少内存分配
var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 4096)
    },
}

func processLargeFile(filename string) error {
    buffer := bufferPool.Get().([]byte)
    defer bufferPool.Put(buffer)
    
    // 使用缓冲区处理文件
    return processWithBuffer(filename, buffer)
}
```

### 10.2 性能监控
```go
func (fp *FileProcessor) reportProgress() {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            elapsed := time.Since(fp.stats.StartTime)
            rate := float64(fp.stats.FilesProcessed) / elapsed.Seconds()
            
            fmt.Printf("\r处理进度: %d 文件, %.2f 文件/秒",
                fp.stats.FilesProcessed, rate)
        case <-fp.done:
            return
        }
    }
}
```

---

## 11. 常见陷阱与解决方案

### 11.1 命令行参数陷阱
```go
// 错误：忘记调用 flag.Parse()
func badExample() {
    name := flag.String("name", "World", "名字")
    // 忘记调用 flag.Parse()
    fmt.Printf("Hello, %s\n", *name) // 总是输出默认值
}

// 正确：记得解析参数
func goodExample() {
    name := flag.String("name", "World", "名字")
    flag.Parse()
    fmt.Printf("Hello, %s\n", *name)
}
```

### 11.2 并发安全问题
```go
// 错误：并发访问共享变量
var counter int
func unsafeIncrement() {
    counter++ // 竞态条件
}

// 正确：使用同步原语
var (
    counter int
    mu      sync.Mutex
)

func safeIncrement() {
    mu.Lock()
    defer mu.Unlock()
    counter++
}
```

---

## 12. 下一步学习

Day8 将学习：
- Goroutine 和 Channel 深入
- 并发模式设计
- Context 使用
- 同步原语

重点关注：
- 并发安全设计
- 性能优化技巧
- 错误处理策略
- 资源管理最佳实践

CLI 工具开发为并发编程奠定了基础，特别是在文件处理、网络请求等场景中，合理使用并发可以显著提升性能。
