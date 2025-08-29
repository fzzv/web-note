#!/bin/bash

# FileProcessor 构建脚本
# 支持多平台交叉编译和版本信息注入

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 应用信息
APP_NAME="fileprocessor"
VERSION=${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo "v1.0.0-dev")}
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT=${COMMIT:-$(git rev-parse HEAD 2>/dev/null || echo "unknown")}

# 构建目录
BUILD_DIR="dist"
LDFLAGS="-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.commit=${COMMIT}"

# 支持的平台
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
    "windows/arm64"
)

# 打印信息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 显示帮助信息
show_help() {
    cat << EOF
FileProcessor 构建脚本

用法: $0 [选项]

选项:
    -h, --help          显示此帮助信息
    -v, --version       显示版本信息
    -p, --platform      指定构建平台 (例如: linux/amd64)
    -a, --all           构建所有支持的平台
    -c, --clean         清理构建目录
    -t, --test          运行测试
    -r, --race          运行竞态检测
    -b, --bench         运行基准测试
    --no-compress       不压缩二进制文件

示例:
    $0                  # 构建当前平台
    $0 -a               # 构建所有平台
    $0 -p linux/amd64   # 构建指定平台
    $0 -t               # 运行测试
    $0 -c               # 清理构建目录

支持的平台:
EOF
    for platform in "${PLATFORMS[@]}"; do
        echo "    $platform"
    done
}

# 清理构建目录
clean_build() {
    print_info "清理构建目录..."
    rm -rf "$BUILD_DIR"
    print_success "构建目录已清理"
}

# 运行测试
run_tests() {
    print_info "运行测试..."
    
    # 单元测试
    print_info "运行单元测试..."
    go test -v ./...
    
    # 测试覆盖率
    print_info "生成测试覆盖率报告..."
    go test -cover -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    print_success "测试覆盖率报告已生成: coverage.html"
}

# 运行竞态检测
run_race_tests() {
    print_info "运行竞态检测..."
    go test -race ./...
    print_success "竞态检测完成"
}

# 运行基准测试
run_benchmarks() {
    print_info "运行基准测试..."
    go test -bench=. -benchmem ./...
    print_success "基准测试完成"
}

# 构建单个平台
build_platform() {
    local platform=$1
    local goos=$(echo $platform | cut -d'/' -f1)
    local goarch=$(echo $platform | cut -d'/' -f2)
    
    local output_name="${APP_NAME}-${goos}-${goarch}"
    if [ "$goos" = "windows" ]; then
        output_name="${output_name}.exe"
    fi
    
    local output_path="${BUILD_DIR}/${output_name}"
    
    print_info "构建 ${platform}..."
    
    # 设置环境变量并构建
    GOOS=$goos GOARCH=$goarch CGO_ENABLED=0 go build \
        -ldflags "${LDFLAGS}" \
        -o "$output_path" \
        main.go
    
    if [ $? -eq 0 ]; then
        # 获取文件大小
        local size=$(du -h "$output_path" | cut -f1)
        print_success "构建完成: $output_path ($size)"
        
        # 压缩二进制文件（可选）
        if [ "$COMPRESS" = "true" ] && command -v upx >/dev/null 2>&1; then
            print_info "压缩二进制文件..."
            upx --best "$output_path" >/dev/null 2>&1
            local compressed_size=$(du -h "$output_path" | cut -f1)
            print_success "压缩完成: $compressed_size"
        fi
    else
        print_error "构建失败: $platform"
        return 1
    fi
}

# 构建所有平台
build_all() {
    print_info "开始构建所有平台..."
    
    local success_count=0
    local total_count=${#PLATFORMS[@]}
    
    for platform in "${PLATFORMS[@]}"; do
        if build_platform "$platform"; then
            ((success_count++))
        fi
    done
    
    print_info "构建统计: $success_count/$total_count 成功"
    
    if [ $success_count -eq $total_count ]; then
        print_success "所有平台构建完成！"
    else
        print_warning "部分平台构建失败"
        return 1
    fi
}

# 构建当前平台
build_current() {
    local goos=$(go env GOOS)
    local goarch=$(go env GOARCH)
    local platform="${goos}/${goarch}"
    
    print_info "构建当前平台: $platform"
    build_platform "$platform"
}

# 生成校验和
generate_checksums() {
    if [ -d "$BUILD_DIR" ]; then
        print_info "生成校验和..."
        cd "$BUILD_DIR"
        sha256sum * > checksums.txt
        cd - >/dev/null
        print_success "校验和已生成: ${BUILD_DIR}/checksums.txt"
    fi
}

# 显示构建信息
show_build_info() {
    print_info "构建信息:"
    echo "  应用名称: $APP_NAME"
    echo "  版本: $VERSION"
    echo "  构建时间: $BUILD_TIME"
    echo "  提交: $COMMIT"
    echo "  Go 版本: $(go version)"
}

# 主函数
main() {
    # 检查 Go 环境
    if ! command -v go >/dev/null 2>&1; then
        print_error "Go 未安装或不在 PATH 中"
        exit 1
    fi
    
    # 创建构建目录
    mkdir -p "$BUILD_DIR"
    
    # 解析命令行参数
    local build_all_platforms=false
    local run_tests_flag=false
    local run_race_flag=false
    local run_bench_flag=false
    local clean_flag=false
    local specific_platform=""
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -v|--version)
                echo "$VERSION"
                exit 0
                ;;
            -a|--all)
                build_all_platforms=true
                shift
                ;;
            -p|--platform)
                specific_platform="$2"
                shift 2
                ;;
            -c|--clean)
                clean_flag=true
                shift
                ;;
            -t|--test)
                run_tests_flag=true
                shift
                ;;
            -r|--race)
                run_race_flag=true
                shift
                ;;
            -b|--bench)
                run_bench_flag=true
                shift
                ;;
            --no-compress)
                COMPRESS=false
                shift
                ;;
            *)
                print_error "未知选项: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    # 设置默认压缩选项
    COMPRESS=${COMPRESS:-true}
    
    # 显示构建信息
    show_build_info
    
    # 执行清理
    if [ "$clean_flag" = true ]; then
        clean_build
        exit 0
    fi
    
    # 运行测试
    if [ "$run_tests_flag" = true ]; then
        run_tests
    fi
    
    # 运行竞态检测
    if [ "$run_race_flag" = true ]; then
        run_race_tests
    fi
    
    # 运行基准测试
    if [ "$run_bench_flag" = true ]; then
        run_benchmarks
    fi
    
    # 如果只是运行测试，不进行构建
    if [ "$run_tests_flag" = true ] || [ "$run_race_flag" = true ] || [ "$run_bench_flag" = true ]; then
        if [ "$build_all_platforms" = false ] && [ -z "$specific_platform" ]; then
            exit 0
        fi
    fi
    
    # 执行构建
    if [ "$build_all_platforms" = true ]; then
        build_all
    elif [ -n "$specific_platform" ]; then
        build_platform "$specific_platform"
    else
        build_current
    fi
    
    # 生成校验和
    generate_checksums
    
    print_success "构建完成！输出目录: $BUILD_DIR"
}

# 执行主函数
main "$@"
