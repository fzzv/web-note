# FileProcessor - 强大的文件处理工具

FileProcessor 是一个功能丰富的命令行文件处理工具，使用 Go 语言开发，支持文件分析、搜索、替换和重命名等操作。

## 功能特性

- 📊 **文件分析**: 统计文件行数、字数、字符数等信息
- 🔍 **内容搜索**: 支持正则表达式的文件内容搜索
- 🔄 **内容替换**: 批量替换文件内容，支持备份
- 📝 **文件重命名**: 批量重命名文件，支持正则表达式
- ⚡ **并发处理**: 多线程并发处理，提高效率
- 📈 **进度显示**: 实时显示处理进度
- 🛡️ **安全操作**: 支持试运行和备份功能
- 🎯 **灵活过滤**: 支持文件包含/排除模式
- 📋 **多种输出**: 支持表格、JSON、CSV 格式输出

## 安装

### 从源码构建

```bash
git clone <repository-url>
cd fileprocessor
go build -o fileprocessor main.go
```

### 使用构建脚本

```bash
chmod +x build.sh
./build.sh
```

## 使用方法

### 基本命令

```bash
# 显示帮助信息
./fileprocessor --help

# 显示版本信息
./fileprocessor version
```

### 文件分析

```bash
# 分析当前目录的文件
./fileprocessor analyze

# 分析指定文件
./fileprocessor analyze file1.txt file2.txt

# 递归分析目录
./fileprocessor analyze -r /path/to/directory

# 输出为 JSON 格式
./fileprocessor analyze -f json

# 显示详细信息
./fileprocessor analyze -d

# 过滤文件类型
./fileprocessor analyze -i "*.go" -e "*.test.go"
```

### 内容搜索

```bash
# 搜索文本
./fileprocessor search "Hello World" /path/to/search

# 使用正则表达式搜索
./fileprocessor search -r "func\s+\w+" .

# 大小写敏感搜索
./fileprocessor search -c "Hello" .

# 显示上下文行
./fileprocessor search -C 3 "error" .

# 过滤文件类型
./fileprocessor search -i "*.go" "package" .
```

### 内容替换

```bash
# 简单替换
./fileprocessor replace "old text" "new text" file1.txt file2.txt

# 使用正则表达式替换
./fileprocessor replace -r "func\s+(\w+)" "function $1" *.go

# 试运行（不实际修改文件）
./fileprocessor replace -d "old" "new" *.txt

# 大小写不敏感替换
./fileprocessor replace "hello" "hi" *.txt

# 包含文件过滤
./fileprocessor replace -i "*.go" "old" "new" .
```

### 文件重命名

```bash
# 简单重命名
./fileprocessor rename "old" "new" /path/to/directory

# 使用正则表达式重命名
./fileprocessor rename -r "(\d+)\.txt" "file_$1.txt" .

# 试运行
./fileprocessor rename -d "old" "new" .

# 递归重命名
./fileprocessor rename -R "old" "new" .
```

## 配置文件

FileProcessor 支持 YAML 配置文件，默认搜索路径：
- `./fileprocessor.yaml`
- `$HOME/.config/fileprocessor/fileprocessor.yaml`

示例配置文件：

```yaml
# 日志级别
log_level: info

# 输出目录
output_dir: ./output

# 最大并发工作线程数
max_workers: 4

# 是否显示进度条
progress_bar: true

# 是否创建备份文件
backup_files: false

# 是否忽略隐藏文件
ignore_hidden: true
```

## 命令行选项

### 全局选项

- `--config`: 指定配置文件路径
- `--log-level`: 设置日志级别 (debug|info|warn|error)
- `--output-dir`: 设置输出目录
- `--workers`: 设置并发工作线程数
- `--progress`: 是否显示进度条
- `--backup`: 是否创建备份文件

### analyze 命令选项

- `-f, --format`: 输出格式 (table|json|csv)
- `-r, --recursive`: 递归处理子目录
- `-i, --include`: 包含文件模式
- `-e, --exclude`: 排除文件模式
- `-d, --detailed`: 显示详细信息

### search 命令选项

- `-r, --regex`: 使用正则表达式
- `-c, --case-sensitive`: 区分大小写
- `-C, --context`: 显示上下文行数
- `-i, --include`: 包含文件模式
- `-e, --exclude`: 排除文件模式
- `-n, --line-numbers`: 显示行号

### replace 命令选项

- `-r, --regex`: 使用正则表达式
- `-c, --case-sensitive`: 区分大小写
- `-d, --dry-run`: 试运行（不实际修改文件）
- `-i, --include`: 包含文件模式
- `-e, --exclude`: 排除文件模式

### rename 命令选项

- `-r, --regex`: 使用正则表达式
- `-d, --dry-run`: 试运行（不实际重命名）
- `-R, --recursive`: 递归处理子目录

## 使用示例

### 代码重构

```bash
# 将所有 Go 文件中的函数名从 camelCase 改为 snake_case
./fileprocessor replace -r "func\s+([a-z])([A-Z])" "func ${1}_${2}" *.go

# 更新包导入路径
./fileprocessor replace "old.package.com" "new.package.com" *.go
```

### 日志分析

```bash
# 搜索错误日志
./fileprocessor search -i "*.log" -C 2 "ERROR"

# 统计日志文件信息
./fileprocessor analyze -f json *.log > log_stats.json
```

### 批量文件整理

```bash
# 重命名图片文件
./fileprocessor rename -r "IMG_(\d+)" "photo_$1" /path/to/photos

# 清理临时文件（试运行）
./fileprocessor search -i "*.tmp" "." | wc -l
```

## 性能优化

- 使用 `--workers` 参数调整并发线程数
- 对于大文件，适当增加缓冲区大小
- 使用文件过滤减少处理的文件数量
- 在 SSD 上运行以获得更好的 I/O 性能

## 测试

```bash
# 运行所有测试
go test ./...

# 运行测试并显示覆盖率
go test -cover ./...

# 运行竞态检测
go test -race ./...

# 运行基准测试
go test -bench=. ./...
```

## 构建

```bash
# 本地构建
go build -o fileprocessor main.go

# 交叉编译
GOOS=linux GOARCH=amd64 go build -o fileprocessor-linux-amd64 main.go
GOOS=windows GOARCH=amd64 go build -o fileprocessor-windows-amd64.exe main.go
GOOS=darwin GOARCH=amd64 go build -o fileprocessor-darwin-amd64 main.go
```

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License

## 更新日志

### v1.0.0
- 初始版本发布
- 支持文件分析、搜索、替换、重命名功能
- 支持并发处理和进度显示
- 支持配置文件和多种输出格式
