# Go 自学与实战 4 周路线图

## 目标
- 1 个月内掌握 Go 基础、并发、Web/后端与系统编程。
- 产出：
  - 完整的「Web 笔记系统（web-note）」：REST API + DB + 鉴权 + 日志 + 配置 + 测试 + Docker + CI。
  - 一个 CLI 工具与一个系统级小工具（文件监控/并发处理）。

## 环境建议（Windows 10）
- 推荐 WSL2 + Ubuntu（系统编程更顺滑）
  - PowerShell（管理员）：`wsl --install -d Ubuntu`
  - Ubuntu 内安装 Go：`sudo tar -C /usr/local -xzf go1.22.x.linux-amd64.tar.gz`
- 或原生 Windows：`winget install -e --id GoLang.Go`
- 工具：VS Code + Go 插件、Docker Desktop、`air`、`golangci-lint`、`goose`
- 验证：`go version`，`go env GOPATH`，按需设置 `GOPROXY`。

## 目录建议
```
web-note/
  cmd/api/main.go
  cmd/worker/main.go
  internal/config/
  internal/http/
  internal/note/
  internal/storage/
  internal/worker/
  pkg/logger/
  migrations/
  go.mod
```

## 每周计划

### 第 1 周：语言基础与工程习惯
- Day1：环境与工具、变量/常量、基础类型、切片/映射、控制流。
- Day2：函数/方法、指针、结构体、接收者、可见性与包。
- Day3：接口、错误处理、`defer`/`panic`/`recover`。
- Day4：泛型（1.18+）、`constraints`、常见容器范式。
- Day5：`go mod`、工程结构（`cmd/`、`internal/`）、`fmt/vet/lint`、表驱动测试。
- Day6：文件/JSON/时间；练习：JSON 日志清洗器、并发 ping。
- Day7：交付 CLI 小工具：参数解析、读写文件、单测、`-race`。

### 第 2 周：并发与网络 I/O
- Day8：goroutine、channel、`select`、`context`、`WaitGroup`、`Mutex`、内存模型。
- Day9：并发模式（worker pool、fan-in/out、pipeline、背压与取消）。
- Day10：HTTP 客户端、重试/超时、连接池、流式读写。
- Day11：`net/http` 服务、路由与中间件思想。
- Day12：综合练习：并发爬取 + 去重 + 限流 + 持久化（SQLite）。
- Day13：测试/基准：`testing`、`testify`、`-race`、`pprof` 入门。
- Day14：周整合与总结。

### 第 3 周：Web/后端工程化
- Day15：Gin/Echo、恢复/日志/CORS/限流中间件。
- Day16：Postgres（`pgx`/`gorm`）、迁移（`goose`）、事务、索引与 N+1。
- Day17：分层（handler/service/repo）、DTO 校验、统一错误码。
- Day18：鉴权（JWT/Session）、密码学、`http.Server` 超时配置。
- Day19：日志（`zap`/`zerolog`）、配置（`viper`/env）、可观测性（OTEL 可选）。
- Day20：集成测试（docker 起 Postgres）、`seed` 数据、`Makefile/Taskfile`。
- Day21：交付最小可用 REST 服务。

### 第 4 周：系统编程与发布
- Day22：CLI（`cobra`）与 API 打通；二进制配置与子命令设计。
- Day23：文件系统/通知（`fsnotify`）、信号与优雅退出（`os/signal`）。
- Day24：性能与稳定性：`pprof`、`-race`、基准；连接池/对象池；资源泄漏扫描。
- Day25：容器化（多阶段构建、无根镜像）、多环境配置、交叉编译。
- Day26：CI（GitHub Actions：lint、test、build）；少量 e2e。
- Day27：项目收尾（文档、README、脚本、一键启动）。
- Day28：演示与答辩（压测、指标、代码走查、Q&A）。

## 核心项目（贯穿式）
- Web 笔记系统（web-note）：REST API + 背景任务 + CLI + 系统监控
- 技术：Gin（或 Echo）、`net/http`、Postgres（`pgx`/`gorm`）、`goose`、`zap`、`viper`、`docker`、CI
- 并发：`errgroup`、Channel、Worker Pool、`context` 取消

## JS → Go 思维迁移
- 强类型、无隐式转换；组合优于继承；接口隐式实现；错误为返回值；goroutine+channel；`go mod`。

## 验收标准（第 4 周末）
- API：CRUD+鉴权；`go test -race -cover` ≥ 60%。
- 并发：可取消、有背压；`pprof` 无明显热点/泄漏。
- 工程：`golangci-lint` 零阻塞；CI 绿；容器镜像可运行。
- 文档：完整 README、接口文档与一键启动脚本。

---

> 提示：按 `go/day1/README.md` 开始第一天学习与练习。 
