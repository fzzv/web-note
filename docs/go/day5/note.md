# Day5 全面解析笔记：Go 模块管理与工程结构

> 代码参考：`code/go/day5`
> 示例涵盖：go mod、工程结构、代码质量工具、表驱动测试

## 1. Go 模块系统（go mod）

### 1.1 基本概念
- **模块（Module）**：Go 代码的版本化单元，由 `go.mod` 文件定义
- **包（Package）**：同一目录下的 `.go` 文件集合，共享包名
- **依赖管理**：通过 `go.mod` 和 `go.sum` 管理外部依赖

### 1.2 go.mod 文件结构
```go
module github.com/username/project-name

go 1.22

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/stretchr/testify v1.8.4
)

require (
    // 间接依赖
    github.com/bytedance/sonic v1.9.1 // indirect
)

replace github.com/old/package => github.com/new/package v1.0.0
exclude github.com/bad/package v1.2.3
```

### 1.3 常用命令
```bash
# 初始化模块
go mod init module-name

# 添加依赖
go get github.com/gin-gonic/gin@latest
go get github.com/gin-gonic/gin@v1.9.1

# 更新依赖
go get -u ./...                    # 更新所有依赖
go get -u github.com/gin-gonic/gin # 更新特定依赖

# 清理依赖
go mod tidy                        # 移除未使用的依赖，添加缺失的依赖
go mod download                    # 下载依赖到本地缓存

# 查看依赖
go list -m all                     # 列出所有依赖
go mod graph                       # 显示依赖图
go mod why github.com/gin-gonic/gin # 解释为什么需要某个依赖

# 验证依赖
go mod verify                      # 验证依赖完整性
```

---

## 2. 标准工程结构

### 2.1 推荐目录布局
```
project/
├── cmd/                    # 主应用程序入口
│   ├── api/               # API 服务器
│   │   └── main.go
│   ├── worker/            # 后台任务处理器
│   │   └── main.go
│   └── cli/               # 命令行工具
│       └── main.go
├── internal/              # 私有应用程序代码
│   ├── config/           # 配置管理
│   ├── handler/          # HTTP 处理器
│   ├── service/          # 业务逻辑
│   ├── repository/       # 数据访问层
│   └── model/            # 数据模型
├── pkg/                   # 可被外部应用使用的库代码
│   ├── logger/           # 日志工具
│   ├── validator/        # 验证工具
│   └── utils/            # 通用工具
├── api/                   # API 定义文件
│   └── openapi.yaml
├── web/                   # Web 资源
│   ├── static/
│   └── templates/
├── scripts/               # 构建、安装、分析等脚本
├── test/                  # 额外的测试数据
├── docs/                  # 设计和用户文档
├── examples/              # 应用程序或公共库的示例
├── deployments/           # 部署配置
├── .gitignore
├── README.md
├── Makefile
└── go.mod
```

### 2.2 目录说明

**cmd/**：每个应用程序的主入口点
- 保持 main.go 简洁，主要逻辑放在 internal/ 或 pkg/
- 按应用程序组织子目录

**internal/**：私有应用程序代码
- 其他项目无法导入 internal/ 下的代码
- 放置不希望被外部使用的代码

**pkg/**：可被外部项目使用的库代码
- 需要谨慎设计 API，因为可能被其他项目依赖
- 通常包含通用工具和库

---

## 3. 代码质量工具

### 3.1 格式化工具
```bash
# gofmt - 标准格式化工具
gofmt -w .                         # 格式化并写入文件
gofmt -d .                         # 显示格式化差异

# goimports - 自动管理导入
go install golang.org/x/tools/cmd/goimports@latest
goimports -w .                     # 格式化并整理导入
```

### 3.2 代码检查工具
```bash
# go vet - 静态分析工具
go vet ./...                       # 检查常见错误

# golangci-lint - 综合代码检查工具
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
golangci-lint run                  # 运行所有启用的检查器
golangci-lint run --enable-all     # 启用所有检查器
```

### 3.3 .golangci.yml 配置示例
```yaml
run:
  timeout: 5m
  tests: true

linters:
  enable:
    - gofmt
    - goimports
    - govet
    - errcheck
    - staticcheck
    - unused
    - gosimple
    - ineffassign
    - typecheck
    - gocritic
    - revive

linters-settings:
  gocritic:
    enabled-tags:
      - diagnostic
      - style
      - performance
  revive:
    rules:
      - name: exported
        disabled: true
```

---

## 4. 表驱动测试（Table-Driven Tests）

### 4.1 基本模式
表驱动测试是 Go 中常用的测试模式，通过定义测试用例表来减少重复代码：

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"正数相加", 2, 3, 5},
        {"负数相加", -1, -2, -3},
        {"零值相加", 0, 5, 5},
        {"大数相加", 1000000, 2000000, 3000000},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Add(%d, %d) = %d; expected %d", 
                    tt.a, tt.b, result, tt.expected)
            }
        })
    }
}
```

### 4.2 高级模式
```go
func TestComplexFunction(t *testing.T) {
    tests := []struct {
        name    string
        input   Input
        want    Output
        wantErr bool
        setup   func() // 测试前的设置
        cleanup func() // 测试后的清理
    }{
        {
            name: "成功案例",
            input: Input{Value: "test"},
            want: Output{Result: "processed"},
            wantErr: false,
            setup: func() {
                // 设置测试环境
            },
            cleanup: func() {
                // 清理测试环境
            },
        },
        {
            name: "错误案例",
            input: Input{Value: ""},
            want: Output{},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if tt.setup != nil {
                tt.setup()
            }
            if tt.cleanup != nil {
                defer tt.cleanup()
            }

            got, err := ComplexFunction(tt.input)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("ComplexFunction() error = %v, wantErr %v", 
                    err, tt.wantErr)
                return
            }
            
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("ComplexFunction() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

---

## 5. 测试最佳实践

### 5.1 测试命名规范
- 测试函数：`TestFunctionName`
- 基准测试：`BenchmarkFunctionName`
- 示例函数：`ExampleFunctionName`

### 5.2 常用测试命令
```bash
# 运行测试
go test ./...                      # 运行所有测试
go test -v ./...                   # 详细输出
go test -run TestSpecific          # 运行特定测试
go test -short                     # 跳过长时间运行的测试

# 测试覆盖率
go test -cover ./...               # 显示覆盖率
go test -coverprofile=coverage.out # 生成覆盖率文件
go tool cover -html=coverage.out   # 生成 HTML 覆盖率报告

# 竞态检测
go test -race ./...                # 检测竞态条件

# 基准测试
go test -bench=.                   # 运行基准测试
go test -bench=. -benchmem         # 包含内存分配统计
```

### 5.3 测试辅助工具
```go
// testify 库示例
import (
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestWithTestify(t *testing.T) {
    result := SomeFunction()
    
    assert.Equal(t, expected, result)           // 断言相等
    assert.NotNil(t, result)                   // 断言非空
    assert.Contains(t, slice, element)         // 断言包含
    
    require.NoError(t, err)                    // 要求无错误，失败时停止测试
}
```

---

## 6. 实践要点

### 6.1 模块管理最佳实践
1. **语义化版本**：使用 `v1.2.3` 格式的标签
2. **最小版本选择**：Go 选择满足要求的最小版本
3. **vendor 目录**：可选，用于离线构建
4. **私有模块**：配置 GOPRIVATE 环境变量

### 6.2 工程结构最佳实践
1. **单一职责**：每个包有明确的职责
2. **依赖方向**：高层依赖低层，避免循环依赖
3. **接口设计**：在使用方定义接口，而非实现方
4. **错误处理**：统一错误处理策略

### 6.3 测试最佳实践
1. **测试覆盖率**：目标 70-80%，关键路径 100%
2. **测试隔离**：每个测试独立，可并行运行
3. **测试数据**：使用 testdata/ 目录存放测试数据
4. **Mock 和 Stub**：使用接口便于测试

---

## 7. 常见问题与解决方案

### 7.1 依赖管理问题
**问题**：依赖版本冲突
**解决**：使用 `go mod graph` 分析依赖关系，必要时使用 `replace` 指令

**问题**：私有仓库访问
**解决**：配置 Git 凭据或使用 `GOPRIVATE` 环境变量

### 7.2 测试问题
**问题**：测试运行缓慢
**解决**：使用 `-short` 标志跳过长时间测试，或使用并行测试

**问题**：竞态条件
**解决**：始终使用 `-race` 标志运行测试

---

## 8. 下一步学习

Day6 将学习：
- 文件 I/O 操作
- JSON 处理
- 时间处理
- 实践项目：JSON 日志清洗器

重点关注：
- 错误处理模式
- 资源管理（defer 的使用）
- 性能优化（缓冲 I/O）
