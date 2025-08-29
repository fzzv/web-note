# FileProcessor Makefile
# 提供常用的开发和构建命令

# 应用信息
APP_NAME := fileprocessor
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "v1.0.0-dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")

# 构建标志
LDFLAGS := -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.commit=$(COMMIT)
BUILD_FLAGS := -ldflags "$(LDFLAGS)"

# 目录
BUILD_DIR := dist
COVERAGE_DIR := coverage

# Go 相关
GO := go
GOFMT := gofmt
GOLINT := golangci-lint

# 默认目标
.DEFAULT_GOAL := help

# 颜色定义
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[1;33m
RED := \033[0;31m
NC := \033[0m # No Color

# 帮助信息
.PHONY: help
help: ## 显示帮助信息
	@echo "$(BLUE)FileProcessor 构建工具$(NC)"
	@echo ""
	@echo "$(YELLOW)可用命令:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# 开发相关命令
.PHONY: dev
dev: ## 开发模式运行
	@echo "$(BLUE)启动开发模式...$(NC)"
	$(GO) run main.go

.PHONY: build
build: ## 构建当前平台的二进制文件
	@echo "$(BLUE)构建 $(APP_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(APP_NAME) main.go
	@echo "$(GREEN)构建完成: $(BUILD_DIR)/$(APP_NAME)$(NC)"

.PHONY: build-all
build-all: ## 构建所有平台的二进制文件
	@echo "$(BLUE)构建所有平台...$(NC)"
	@chmod +x build.sh
	./build.sh -a

.PHONY: install
install: build ## 安装到系统路径
	@echo "$(BLUE)安装 $(APP_NAME)...$(NC)"
	sudo cp $(BUILD_DIR)/$(APP_NAME) /usr/local/bin/
	@echo "$(GREEN)安装完成$(NC)"

# 测试相关命令
.PHONY: test
test: ## 运行所有测试
	@echo "$(BLUE)运行测试...$(NC)"
	$(GO) test -v ./...

.PHONY: test-cover
test-cover: ## 运行测试并生成覆盖率报告
	@echo "$(BLUE)运行测试覆盖率分析...$(NC)"
	@mkdir -p $(COVERAGE_DIR)
	$(GO) test -cover -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	$(GO) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "$(GREEN)覆盖率报告已生成: $(COVERAGE_DIR)/coverage.html$(NC)"

.PHONY: test-race
test-race: ## 运行竞态检测
	@echo "$(BLUE)运行竞态检测...$(NC)"
	$(GO) test -race ./...

.PHONY: bench
bench: ## 运行基准测试
	@echo "$(BLUE)运行基准测试...$(NC)"
	$(GO) test -bench=. -benchmem ./...

.PHONY: test-all
test-all: test test-race test-cover ## 运行所有测试（包括竞态检测和覆盖率）

# 代码质量相关命令
.PHONY: fmt
fmt: ## 格式化代码
	@echo "$(BLUE)格式化代码...$(NC)"
	$(GOFMT) -w .
	$(GO) mod tidy

.PHONY: lint
lint: ## 运行代码检查
	@echo "$(BLUE)运行代码检查...$(NC)"
	@if command -v $(GOLINT) >/dev/null 2>&1; then \
		$(GOLINT) run; \
	else \
		echo "$(YELLOW)golangci-lint 未安装，跳过代码检查$(NC)"; \
		echo "$(YELLOW)安装命令: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(NC)"; \
	fi

.PHONY: vet
vet: ## 运行 go vet
	@echo "$(BLUE)运行 go vet...$(NC)"
	$(GO) vet ./...

.PHONY: check
check: fmt vet lint ## 运行所有代码质量检查

# 清理相关命令
.PHONY: clean
clean: ## 清理构建文件
	@echo "$(BLUE)清理构建文件...$(NC)"
	rm -rf $(BUILD_DIR)
	rm -rf $(COVERAGE_DIR)
	rm -f coverage.out
	rm -f coverage.html
	@echo "$(GREEN)清理完成$(NC)"

.PHONY: clean-all
clean-all: clean ## 清理所有生成的文件
	@echo "$(BLUE)清理所有文件...$(NC)"
	$(GO) clean -cache
	$(GO) clean -testcache
	$(GO) clean -modcache

# 依赖管理
.PHONY: deps
deps: ## 下载依赖
	@echo "$(BLUE)下载依赖...$(NC)"
	$(GO) mod download
	$(GO) mod verify

.PHONY: deps-update
deps-update: ## 更新依赖
	@echo "$(BLUE)更新依赖...$(NC)"
	$(GO) get -u ./...
	$(GO) mod tidy

# 发布相关命令
.PHONY: release
release: clean test-all build-all ## 准备发布版本
	@echo "$(BLUE)准备发布版本 $(VERSION)...$(NC)"
	@cd $(BUILD_DIR) && sha256sum * > checksums.txt
	@echo "$(GREEN)发布准备完成$(NC)"

# 文档相关命令
.PHONY: docs
docs: ## 生成文档
	@echo "$(BLUE)生成文档...$(NC)"
	$(GO) doc -all . > docs.txt
	@echo "$(GREEN)文档已生成: docs.txt$(NC)"

# 示例和演示
.PHONY: demo
demo: build ## 运行演示
	@echo "$(BLUE)运行演示...$(NC)"
	@mkdir -p demo_files
	@echo "这是第一个测试文件" > demo_files/file1.txt
	@echo "这是第二个测试文件" > demo_files/file2.txt
	@echo "Hello World" > demo_files/hello.txt
	@echo ""
	@echo "$(YELLOW)1. 分析文件:$(NC)"
	./$(BUILD_DIR)/$(APP_NAME) analyze demo_files/
	@echo ""
	@echo "$(YELLOW)2. 搜索内容:$(NC)"
	./$(BUILD_DIR)/$(APP_NAME) search "测试" demo_files/
	@echo ""
	@echo "$(YELLOW)3. 替换内容（试运行）:$(NC)"
	./$(BUILD_DIR)/$(APP_NAME) replace -d "测试" "示例" demo_files/*.txt
	@echo ""
	@echo "$(GREEN)演示完成$(NC)"
	@rm -rf demo_files

# 性能分析
.PHONY: profile
profile: build ## 运行性能分析
	@echo "$(BLUE)运行性能分析...$(NC)"
	@mkdir -p profile_data
	@for i in {1..100}; do echo "Line $$i: This is a test line with some content" >> profile_data/large_file.txt; done
	./$(BUILD_DIR)/$(APP_NAME) analyze profile_data/ > /dev/null
	@rm -rf profile_data
	@echo "$(GREEN)性能分析完成$(NC)"

# 安装开发工具
.PHONY: install-tools
install-tools: ## 安装开发工具
	@echo "$(BLUE)安装开发工具...$(NC)"
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GO) install golang.org/x/tools/cmd/goimports@latest
	@echo "$(GREEN)开发工具安装完成$(NC)"

# 版本信息
.PHONY: version
version: ## 显示版本信息
	@echo "$(BLUE)版本信息:$(NC)"
	@echo "  应用名称: $(APP_NAME)"
	@echo "  版本: $(VERSION)"
	@echo "  构建时间: $(BUILD_TIME)"
	@echo "  提交: $(COMMIT)"
	@echo "  Go 版本: $(shell $(GO) version)"

# 快速开始
.PHONY: quick-start
quick-start: deps build test ## 快速开始（下载依赖、构建、测试）
	@echo "$(GREEN)快速开始完成！$(NC)"
	@echo "$(YELLOW)运行 'make demo' 查看演示$(NC)"

# CI/CD 相关
.PHONY: ci
ci: deps check test-all build ## CI 流水线
	@echo "$(GREEN)CI 流水线完成$(NC)"

# 监视文件变化（需要安装 fswatch）
.PHONY: watch
watch: ## 监视文件变化并自动测试
	@if command -v fswatch >/dev/null 2>&1; then \
		echo "$(BLUE)监视文件变化...$(NC)"; \
		fswatch -o . | xargs -n1 -I{} make test; \
	else \
		echo "$(RED)fswatch 未安装$(NC)"; \
		echo "$(YELLOW)macOS: brew install fswatch$(NC)"; \
		echo "$(YELLOW)Linux: apt-get install fswatch$(NC)"; \
	fi

# 显示项目统计信息
.PHONY: stats
stats: ## 显示项目统计信息
	@echo "$(BLUE)项目统计信息:$(NC)"
	@echo "  Go 文件数: $(shell find . -name '*.go' | wc -l)"
	@echo "  代码行数: $(shell find . -name '*.go' -exec wc -l {} + | tail -1 | awk '{print $$1}')"
	@echo "  测试文件数: $(shell find . -name '*_test.go' | wc -l)"
	@echo "  包数量: $(shell $(GO) list ./... | wc -l)"

# 特殊目标标记
.PHONY: all
all: clean deps check test-all build-all ## 执行完整的构建流程
