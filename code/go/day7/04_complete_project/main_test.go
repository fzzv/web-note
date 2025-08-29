package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAnalyzeFile 测试文件分析功能
func TestAnalyzeFile(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expectedLines int
		expectedWords int
		expectedChars int
	}{
		{
			name:        "空文件",
			content:     "",
			expectedLines: 1,
			expectedWords: 0,
			expectedChars: 0,
		},
		{
			name:        "单行文件",
			content:     "Hello World",
			expectedLines: 1,
			expectedWords: 2,
			expectedChars: 11,
		},
		{
			name:        "多行文件",
			content:     "第一行\n第二行\n第三行",
			expectedLines: 3,
			expectedWords: 3,
			expectedChars: 9,
		},
		{
			name:        "包含空行",
			content:     "第一行\n\n第三行\n",
			expectedLines: 4,
			expectedWords: 2,
			expectedChars: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建临时文件
			tmpfile, err := os.CreateTemp("", "test_*.txt")
			require.NoError(t, err)
			defer os.Remove(tmpfile.Name())

			// 写入测试内容
			_, err = tmpfile.WriteString(tt.content)
			require.NoError(t, err)
			tmpfile.Close()

			// 分析文件
			info := analyzeFile(tmpfile.Name())
			require.NotNil(t, info)

			// 验证结果
			assert.Equal(t, tt.expectedLines, info.Lines, "行数不匹配")
			assert.Equal(t, tt.expectedWords, info.Words, "词数不匹配")
			assert.Equal(t, tt.expectedChars, info.Characters, "字符数不匹配")
			assert.Equal(t, int64(len(tt.content)), info.Size, "文件大小不匹配")
		})
	}
}

// TestReplaceInFile 测试文件替换功能
func TestReplaceInFile(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		searchPattern  string
		replacement    string
		useRegex       bool
		caseSensitive  bool
		expectedResult string
		expectedLines  int
	}{
		{
			name:           "简单替换",
			content:        "Hello World\nHello Go",
			searchPattern:  "Hello",
			replacement:    "Hi",
			useRegex:       false,
			caseSensitive:  true,
			expectedResult: "Hi World\nHi Go",
			expectedLines:  0,
		},
		{
			name:           "大小写不敏感替换",
			content:        "Hello World\nhello go",
			searchPattern:  "hello",
			replacement:    "Hi",
			useRegex:       false,
			caseSensitive:  false,
			expectedResult: "Hi World\nHi go",
			expectedLines:  0,
		},
		{
			name:           "正则表达式替换",
			content:        "test123\ntest456",
			searchPattern:  `test\d+`,
			replacement:    "result",
			useRegex:       true,
			caseSensitive:  true,
			expectedResult: "result\nresult",
			expectedLines:  0,
		},
		{
			name:           "无匹配内容",
			content:        "Hello World",
			searchPattern:  "xyz",
			replacement:    "abc",
			useRegex:       false,
			caseSensitive:  true,
			expectedResult: "Hello World",
			expectedLines:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建临时文件
			tmpfile, err := os.CreateTemp("", "test_*.txt")
			require.NoError(t, err)
			defer os.Remove(tmpfile.Name())

			// 写入测试内容
			_, err = tmpfile.WriteString(tt.content)
			require.NoError(t, err)
			tmpfile.Close()

			// 执行替换
			result := replaceInFile(tmpfile.Name(), tt.searchPattern, tt.replacement, 
				tt.useRegex, tt.caseSensitive, false)

			// 验证结果
			assert.True(t, result.Success, "替换应该成功")
			assert.Empty(t, result.Error, "不应该有错误")

			// 读取文件内容验证
			newContent, err := os.ReadFile(tmpfile.Name())
			require.NoError(t, err)
			assert.Equal(t, tt.expectedResult, string(newContent), "替换结果不匹配")
		})
	}
}

// TestReplaceAllIgnoreCase 测试大小写不敏感替换
func TestReplaceAllIgnoreCase(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		old      string
		new      string
		expected string
	}{
		{
			name:     "基本替换",
			text:     "Hello World Hello",
			old:      "hello",
			new:      "hi",
			expected: "hi World hi",
		},
		{
			name:     "混合大小写",
			text:     "Hello HELLO hello HeLLo",
			old:      "hello",
			new:      "hi",
			expected: "hi hi hi hi",
		},
		{
			name:     "无匹配",
			text:     "Hello World",
			old:      "xyz",
			new:      "abc",
			expected: "Hello World",
		},
		{
			name:     "空字符串",
			text:     "",
			old:      "hello",
			new:      "hi",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceAllIgnoreCase(tt.text, tt.old, tt.new)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestCollectFiles 测试文件收集功能
func TestCollectFiles(t *testing.T) {
	// 创建临时目录结构
	tmpDir, err := os.MkdirTemp("", "test_collect_*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// 创建测试文件
	testFiles := []string{
		"file1.txt",
		"file2.go",
		"subdir/file3.txt",
		"subdir/file4.py",
		".hidden.txt",
	}

	for _, file := range testFiles {
		fullPath := filepath.Join(tmpDir, file)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		require.NoError(t, err)
		
		err = os.WriteFile(fullPath, []byte("test content"), 0644)
		require.NoError(t, err)
	}

	// 初始化配置
	config = &Config{IgnoreHidden: true}

	tests := []struct {
		name      string
		paths     []string
		recursive bool
		include   string
		exclude   string
		expected  int
	}{
		{
			name:      "非递归收集",
			paths:     []string{tmpDir},
			recursive: false,
			include:   "*",
			exclude:   "",
			expected:  2, // file1.txt, file2.go (忽略隐藏文件)
		},
		{
			name:      "递归收集",
			paths:     []string{tmpDir},
			recursive: true,
			include:   "*",
			exclude:   "",
			expected:  4, // 所有非隐藏文件
		},
		{
			name:      "按扩展名过滤",
			paths:     []string{tmpDir},
			recursive: true,
			include:   "*.txt",
			exclude:   "",
			expected:  2, // file1.txt, subdir/file3.txt
		},
		{
			name:      "排除模式",
			paths:     []string{tmpDir},
			recursive: true,
			include:   "*",
			exclude:   ".*\\.py$",
			expected:  3, // 排除 .py 文件
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files, err := collectFiles(tt.paths, tt.recursive, tt.include, tt.exclude)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, len(files), "收集的文件数量不匹配")
		})
	}
}

// TestTruncateString 测试字符串截断功能
func TestTruncateString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{
			name:     "短字符串",
			input:    "Hello",
			maxLen:   10,
			expected: "Hello",
		},
		{
			name:     "长字符串",
			input:    "This is a very long string",
			maxLen:   10,
			expected: "This is...",
		},
		{
			name:     "正好等于长度",
			input:    "Hello",
			maxLen:   5,
			expected: "Hello",
		},
		{
			name:     "空字符串",
			input:    "",
			maxLen:   5,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateString(tt.input, tt.maxLen)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestAbs 测试绝对值函数
func TestAbs(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"正数", 5, 5},
		{"负数", -5, 5},
		{"零", 0, 0},
		{"大正数", 1000, 1000},
		{"大负数", -1000, 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := abs(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// BenchmarkAnalyzeFile 基准测试文件分析
func BenchmarkAnalyzeFile(b *testing.B) {
	// 创建测试文件
	tmpfile, err := os.CreateTemp("", "bench_*.txt")
	require.NoError(b, err)
	defer os.Remove(tmpfile.Name())

	// 写入大量内容
	content := strings.Repeat("这是一行测试内容，包含多个单词和字符。\n", 1000)
	_, err = tmpfile.WriteString(content)
	require.NoError(b, err)
	tmpfile.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analyzeFile(tmpfile.Name())
	}
}

// BenchmarkReplaceInFile 基准测试文件替换
func BenchmarkReplaceInFile(b *testing.B) {
	// 创建测试文件
	tmpfile, err := os.CreateTemp("", "bench_*.txt")
	require.NoError(b, err)
	defer os.Remove(tmpfile.Name())

	// 写入大量内容
	content := strings.Repeat("Hello World! This is a test line.\n", 1000)
	_, err = tmpfile.WriteString(content)
	require.NoError(b, err)
	tmpfile.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		replaceInFile(tmpfile.Name(), "Hello", "Hi", false, true, true)
	}
}

// TestMain 测试主函数
func TestMain(m *testing.M) {
	// 测试前设置
	config = &Config{
		LogLevel:     "error", // 测试时减少日志输出
		OutputDir:    "./test_output",
		MaxWorkers:   2,
		BufferSize:   1024,
		ProgressBar:  false,
		BackupFiles:  false,
		IgnoreHidden: true,
	}

	// 创建测试输出目录
	os.MkdirAll(config.OutputDir, 0755)

	// 运行测试
	code := m.Run()

	// 测试后清理
	os.RemoveAll(config.OutputDir)

	os.Exit(code)
}
