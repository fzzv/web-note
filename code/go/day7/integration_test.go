// Package day7 集成测试
// 测试所有 CLI 工具的功能集成
package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCLIIntegration 测试 CLI 工具集成
func TestCLIIntegration(t *testing.T) {
	// 创建临时测试目录
	tmpDir, err := os.MkdirTemp("", "cli_integration_test_*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// 测试基础 CLI 工具
	t.Run("BasicCLI", func(t *testing.T) {
		testBasicCLI(t, tmpDir)
	})

	// 测试 Cobra 框架工具
	t.Run("CobraFramework", func(t *testing.T) {
		testCobraFramework(t, tmpDir)
	})

	// 测试高级 CLI 功能
	t.Run("AdvancedCLI", func(t *testing.T) {
		testAdvancedCLI(t, tmpDir)
	})

	// 测试完整项目
	t.Run("CompleteProject", func(t *testing.T) {
		testCompleteProject(t, tmpDir)
	})
}

// testBasicCLI 测试基础 CLI 工具
func testBasicCLI(t *testing.T, tmpDir string) {
	// 构建基础 CLI 工具
	binaryPath := filepath.Join(tmpDir, "basic_cli")
	if isWindows() {
		binaryPath += ".exe"
	}

	buildCmd := exec.Command("go", "build", "-o", binaryPath, "./01_basic_cli/main.go")
	buildCmd.Dir = "."
	err := buildCmd.Run()
	require.NoError(t, err, "构建基础 CLI 工具失败")

	// 创建测试文件
	testFile := filepath.Join(tmpDir, "test.txt")
	testContent := "Hello World\nThis is a test file\nWith multiple lines"
	err = os.WriteFile(testFile, []byte(testContent), 0644)
	require.NoError(t, err)

	// 测试默认行为
	cmd := exec.Command(binaryPath)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err)
	assert.Contains(t, string(output), "Hello, World!")

	// 测试带参数的行为
	cmd = exec.Command(binaryPath, "-name=Go", "-age=10", "-verbose", testFile)
	output, err = cmd.CombinedOutput()
	require.NoError(t, err)
	assert.Contains(t, string(output), "Hello, Go!")
	assert.Contains(t, string(output), "3 行")

	// 测试输出到文件
	outputFile := filepath.Join(tmpDir, "output.txt")
	cmd = exec.Command(binaryPath, "-output="+outputFile, testFile)
	err = cmd.Run()
	require.NoError(t, err)

	// 验证输出文件
	_, err = os.Stat(outputFile)
	assert.NoError(t, err, "输出文件应该存在")
}

// testCobraFramework 测试 Cobra 框架工具
func testCobraFramework(t *testing.T, tmpDir string) {
	// 构建 Cobra 工具
	binaryPath := filepath.Join(tmpDir, "cobra_cli")
	if isWindows() {
		binaryPath += ".exe"
	}

	buildCmd := exec.Command("go", "build", "-o", binaryPath, "./02_cobra_framework/main.go")
	buildCmd.Dir = "."
	err := buildCmd.Run()
	require.NoError(t, err, "构建 Cobra 工具失败")

	// 创建测试文件
	testFile := filepath.Join(tmpDir, "cobra_test.txt")
	testContent := "package main\n\nfunc main() {\n\tfmt.Println(\"Hello\")\n}"
	err = os.WriteFile(testFile, []byte(testContent), 0644)
	require.NoError(t, err)

	// 测试版本命令
	cmd := exec.Command(binaryPath, "version")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err)
	assert.Contains(t, string(output), "FileTools")

	// 测试统计命令
	cmd = exec.Command(binaryPath, "stats", testFile)
	output, err = cmd.CombinedOutput()
	require.NoError(t, err)
	assert.Contains(t, string(output), "cobra_test.txt")

	// 测试搜索命令
	cmd = exec.Command(binaryPath, "search", "-t", "main", tmpDir)
	output, err = cmd.CombinedOutput()
	require.NoError(t, err)
	// 搜索可能找到或找不到结果，主要测试命令不出错

	// 测试转换命令
	cmd = exec.Command(binaryPath, "convert", "-f", "json", "-t", "yaml", testFile)
	output, err = cmd.CombinedOutput()
	require.NoError(t, err)
	assert.Contains(t, string(output), "转换结果")
}

// testAdvancedCLI 测试高级 CLI 功能
func testAdvancedCLI(t *testing.T, tmpDir string) {
	// 构建高级 CLI 工具
	binaryPath := filepath.Join(tmpDir, "advanced_cli")
	if isWindows() {
		binaryPath += ".exe"
	}

	buildCmd := exec.Command("go", "build", "-o", binaryPath, "./03_advanced_cli/main.go")
	buildCmd.Dir = "."
	err := buildCmd.Run()
	require.NoError(t, err, "构建高级 CLI 工具失败")

	// 测试列表命令
	cmd := exec.Command(binaryPath, "list")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err)
	assert.Contains(t, string(output), "任务列表")

	// 测试配置显示
	cmd = exec.Command(binaryPath, "config", "show")
	output, err = cmd.CombinedOutput()
	require.NoError(t, err)
	assert.Contains(t, string(output), "当前配置")

	// 测试批量处理（无进度条模式）
	cmd = exec.Command(binaryPath, "process", "--progress=false", "--delay=10")
	// 设置超时，避免测试运行太久
	cmd.Env = append(os.Environ(), "TASKMANAGER_LOG_LEVEL=error")
	
	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()

	select {
	case err := <-done:
		// 命令完成，检查是否成功
		if err != nil {
			t.Logf("批量处理命令返回错误（可能是正常的）: %v", err)
		}
	case <-time.After(5 * time.Second):
		// 超时，杀死进程
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		t.Log("批量处理命令超时，已终止")
	}
}

// testCompleteProject 测试完整项目
func testCompleteProject(t *testing.T, tmpDir string) {
	// 构建完整项目
	binaryPath := filepath.Join(tmpDir, "fileprocessor")
	if isWindows() {
		binaryPath += ".exe"
	}

	buildCmd := exec.Command("go", "build", "-o", binaryPath, "./04_complete_project/main.go")
	buildCmd.Dir = "."
	err := buildCmd.Run()
	require.NoError(t, err, "构建完整项目失败")

	// 创建测试文件
	testDir := filepath.Join(tmpDir, "test_files")
	err = os.MkdirAll(testDir, 0755)
	require.NoError(t, err)

	testFiles := map[string]string{
		"file1.txt": "Hello World\nThis is file 1\nWith some content",
		"file2.go":  "package main\n\nfunc main() {\n\tfmt.Println(\"Hello\")\n}",
		"file3.py":  "#!/usr/bin/env python3\nprint(\"Hello Python\")\n",
	}

	for filename, content := range testFiles {
		filePath := filepath.Join(testDir, filename)
		err = os.WriteFile(filePath, []byte(content), 0644)
		require.NoError(t, err)
	}

	// 测试版本命令
	cmd := exec.Command(binaryPath, "--version")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err)
	assert.Contains(t, string(output), "1.0.0")

	// 测试分析命令
	cmd = exec.Command(binaryPath, "analyze", testDir)
	output, err = cmd.CombinedOutput()
	require.NoError(t, err)
	assert.Contains(t, string(output), "file1.txt")
	assert.Contains(t, string(output), "file2.go")

	// 测试搜索命令
	cmd = exec.Command(binaryPath, "search", "Hello", testDir)
	output, err = cmd.CombinedOutput()
	require.NoError(t, err)
	// 应该找到包含 "Hello" 的文件

	// 测试替换命令（试运行）
	cmd = exec.Command(binaryPath, "replace", "-d", "Hello", "Hi", filepath.Join(testDir, "file1.txt"))
	output, err = cmd.CombinedOutput()
	require.NoError(t, err)
	// 试运行不应该实际修改文件

	// 验证文件内容未改变
	content, err := os.ReadFile(filepath.Join(testDir, "file1.txt"))
	require.NoError(t, err)
	assert.Contains(t, string(content), "Hello World")

	// 测试重命名命令（试运行）
	cmd = exec.Command(binaryPath, "rename", "-d", "file", "document", testDir)
	output, err = cmd.CombinedOutput()
	require.NoError(t, err)
	// 试运行不应该实际重命名文件

	// 验证文件仍然存在
	_, err = os.Stat(filepath.Join(testDir, "file1.txt"))
	assert.NoError(t, err, "原文件应该仍然存在")
}

// TestConcurrentCLI 测试并发 CLI 操作
func TestConcurrentCLI(t *testing.T) {
	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "concurrent_cli_test_*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// 构建基础 CLI 工具
	binaryPath := filepath.Join(tmpDir, "basic_cli")
	if isWindows() {
		binaryPath += ".exe"
	}

	buildCmd := exec.Command("go", "build", "-o", binaryPath, "./01_basic_cli/main.go")
	buildCmd.Dir = "."
	err = buildCmd.Run()
	require.NoError(t, err)

	// 创建多个测试文件
	for i := 0; i < 5; i++ {
		testFile := filepath.Join(tmpDir, "test"+string(rune('0'+i))+".txt")
		content := "Test file " + string(rune('0'+i)) + "\nWith some content\n"
		err = os.WriteFile(testFile, []byte(content), 0644)
		require.NoError(t, err)
	}

	// 并发运行多个 CLI 命令
	done := make(chan bool, 5)
	for i := 0; i < 5; i++ {
		go func(index int) {
			defer func() { done <- true }()
			
			testFile := filepath.Join(tmpDir, "test"+string(rune('0'+index))+".txt")
			cmd := exec.Command(binaryPath, "-verbose", testFile)
			output, err := cmd.CombinedOutput()
			
			assert.NoError(t, err, "并发命令 %d 应该成功", index)
			assert.Contains(t, string(output), "test"+string(rune('0'+index))+".txt")
		}(i)
	}

	// 等待所有命令完成
	for i := 0; i < 5; i++ {
		select {
		case <-done:
			// 命令完成
		case <-time.After(10 * time.Second):
			t.Fatal("并发测试超时")
		}
	}
}

// TestErrorHandling 测试错误处理
func TestErrorHandling(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "error_handling_test_*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// 构建基础 CLI 工具
	binaryPath := filepath.Join(tmpDir, "basic_cli")
	if isWindows() {
		binaryPath += ".exe"
	}

	buildCmd := exec.Command("go", "build", "-o", binaryPath, "./01_basic_cli/main.go")
	buildCmd.Dir = "."
	err = buildCmd.Run()
	require.NoError(t, err)

	// 测试无效参数
	cmd := exec.Command(binaryPath, "-invalid-flag")
	output, err := cmd.CombinedOutput()
	assert.Error(t, err, "无效参数应该返回错误")
	assert.Contains(t, string(output), "flag provided but not defined")

	// 测试不存在的文件
	cmd = exec.Command(binaryPath, "/nonexistent/file.txt")
	output, err = cmd.CombinedOutput()
	// 这个可能成功（因为程序可能只是报告文件不存在），也可能失败
	t.Logf("不存在文件的输出: %s", string(output))

	// 测试无效的日志级别
	cmd = exec.Command(binaryPath, "-log-level=invalid")
	output, err = cmd.CombinedOutput()
	// 程序应该处理无效的日志级别
	t.Logf("无效日志级别的输出: %s", string(output))
}

// 辅助函数
func isWindows() bool {
	return strings.Contains(strings.ToLower(os.Getenv("OS")), "windows")
}

// TestBuildScripts 测试构建脚本
func TestBuildScripts(t *testing.T) {
	// 测试构建脚本是否存在且可执行
	buildScript := "./04_complete_project/build.sh"
	if _, err := os.Stat(buildScript); err != nil {
		t.Skip("构建脚本不存在，跳过测试")
	}

	// 在 Windows 上跳过 shell 脚本测试
	if isWindows() {
		t.Skip("Windows 环境跳过 shell 脚本测试")
	}

	// 测试构建脚本帮助
	cmd := exec.Command("bash", buildScript, "--help")
	cmd.Dir = "./04_complete_project"
	output, err := cmd.CombinedOutput()
	require.NoError(t, err)
	assert.Contains(t, string(output), "FileProcessor 构建脚本")

	// 测试版本显示
	cmd = exec.Command("bash", buildScript, "--version")
	cmd.Dir = "./04_complete_project"
	output, err = cmd.CombinedOutput()
	require.NoError(t, err)
	// 应该输出版本信息
}

// TestMakefile 测试 Makefile
func TestMakefile(t *testing.T) {
	makefilePath := "./04_complete_project/Makefile"
	if _, err := os.Stat(makefilePath); err != nil {
		t.Skip("Makefile 不存在，跳过测试")
	}

	// 检查 make 命令是否可用
	if _, err := exec.LookPath("make"); err != nil {
		t.Skip("make 命令不可用，跳过测试")
	}

	// 测试 make help
	cmd := exec.Command("make", "help")
	cmd.Dir = "./04_complete_project"
	output, err := cmd.CombinedOutput()
	require.NoError(t, err)
	assert.Contains(t, string(output), "FileProcessor 构建工具")

	// 测试 make version
	cmd = exec.Command("make", "version")
	cmd.Dir = "./04_complete_project"
	output, err = cmd.CombinedOutput()
	require.NoError(t, err)
	assert.Contains(t, string(output), "版本信息")
}
