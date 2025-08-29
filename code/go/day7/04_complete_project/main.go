// Package main 完整的文件处理 CLI 工具
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// 版本信息（构建时注入）
	version   = "1.0.0"
	buildTime = "2024-01-01"
	commit    = "dev"

	// 全局变量
	logger *logrus.Logger
	config *Config
)

// Config 应用程序配置
type Config struct {
	LogLevel     string `mapstructure:"log_level"`
	OutputDir    string `mapstructure:"output_dir"`
	MaxWorkers   int    `mapstructure:"max_workers"`
	BufferSize   int    `mapstructure:"buffer_size"`
	ProgressBar  bool   `mapstructure:"progress_bar"`
	BackupFiles  bool   `mapstructure:"backup_files"`
	IgnoreHidden bool   `mapstructure:"ignore_hidden"`
}

// FileInfo 文件信息
type FileInfo struct {
	Path         string    `json:"path"`
	Size         int64     `json:"size"`
	Lines        int       `json:"lines"`
	Words        int       `json:"words"`
	Characters   int       `json:"characters"`
	ModTime      time.Time `json:"mod_time"`
	IsDirectory  bool      `json:"is_directory"`
	Extension    string    `json:"extension"`
	Permissions  string    `json:"permissions"`
}

// ProcessResult 处理结果
type ProcessResult struct {
	File         string        `json:"file"`
	Success      bool          `json:"success"`
	Error        string        `json:"error,omitempty"`
	Duration     time.Duration `json:"duration"`
	BytesRead    int64         `json:"bytes_read"`
	BytesWritten int64         `json:"bytes_written"`
	LinesChanged int           `json:"lines_changed"`
}

// Statistics 统计信息
type Statistics struct {
	TotalFiles     int           `json:"total_files"`
	ProcessedFiles int           `json:"processed_files"`
	FailedFiles    int           `json:"failed_files"`
	TotalBytes     int64         `json:"total_bytes"`
	TotalDuration  time.Duration `json:"total_duration"`
	StartTime      time.Time     `json:"start_time"`
	EndTime        time.Time     `json:"end_time"`
}

var rootCmd = &cobra.Command{
	Use:   "fileprocessor",
	Short: "强大的文件处理工具",
	Long: `FileProcessor 是一个功能丰富的文件处理工具，支持：
- 文件内容搜索和替换
- 批量文件重命名
- 文件统计分析
- 并发处理
- 进度显示
- 备份功能`,
	Version: fmt.Sprintf("%s (built: %s, commit: %s)", version, buildTime, commit),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initializeApp()
	},
}

var analyzeCmd = &cobra.Command{
	Use:   "analyze [目录或文件...]",
	Short: "分析文件统计信息",
	Long:  "分析指定文件或目录的统计信息，包括行数、字数、大小等",
	Run:   runAnalyze,
}

var searchCmd = &cobra.Command{
	Use:   "search [模式] [目录...]",
	Short: "搜索文件内容",
	Long:  "在指定目录中搜索匹配的文件内容，支持正则表达式",
	Args:  cobra.MinimumNArgs(1),
	Run:   runSearch,
}

var replaceCmd = &cobra.Command{
	Use:   "replace [搜索模式] [替换文本] [文件...]",
	Short: "替换文件内容",
	Long:  "在指定文件中搜索并替换内容，支持正则表达式和备份",
	Args:  cobra.MinimumNArgs(3),
	Run:   runReplace,
}

var renameCmd = &cobra.Command{
	Use:   "rename [模式] [替换] [目录]",
	Short: "批量重命名文件",
	Long:  "根据指定模式批量重命名文件",
	Args:  cobra.ExactArgs(3),
	Run:   runRename,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		if logger != nil {
			logger.Fatal(err)
		} else {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			os.Exit(1)
		}
	}
}

func init() {
	// 全局标志
	rootCmd.PersistentFlags().String("config", "", "配置文件路径")
	rootCmd.PersistentFlags().String("log-level", "info", "日志级别 (debug|info|warn|error)")
	rootCmd.PersistentFlags().String("output-dir", "./output", "输出目录")
	rootCmd.PersistentFlags().Int("workers", 4, "并发工作线程数")
	rootCmd.PersistentFlags().Bool("progress", true, "显示进度条")
	rootCmd.PersistentFlags().Bool("backup", false, "创建备份文件")

	// analyze 命令标志
	analyzeCmd.Flags().StringP("format", "f", "table", "输出格式 (table|json|csv)")
	analyzeCmd.Flags().BoolP("recursive", "r", false, "递归处理子目录")
	analyzeCmd.Flags().StringP("include", "i", "*", "包含文件模式")
	analyzeCmd.Flags().StringP("exclude", "e", "", "排除文件模式")
	analyzeCmd.Flags().BoolP("detailed", "d", false, "显示详细信息")

	// search 命令标志
	searchCmd.Flags().BoolP("regex", "r", false, "使用正则表达式")
	searchCmd.Flags().BoolP("case-sensitive", "c", false, "区分大小写")
	searchCmd.Flags().IntP("context", "C", 0, "显示上下文行数")
	searchCmd.Flags().StringP("include", "i", "*", "包含文件模式")
	searchCmd.Flags().StringP("exclude", "e", "", "排除文件模式")
	searchCmd.Flags().BoolP("line-numbers", "n", true, "显示行号")

	// replace 命令标志
	replaceCmd.Flags().BoolP("regex", "r", false, "使用正则表达式")
	replaceCmd.Flags().BoolP("case-sensitive", "c", false, "区分大小写")
	replaceCmd.Flags().BoolP("dry-run", "d", false, "试运行（不实际修改文件）")
	replaceCmd.Flags().StringP("include", "i", "*", "包含文件模式")
	replaceCmd.Flags().StringP("exclude", "e", "", "排除文件模式")

	// rename 命令标志
	renameCmd.Flags().BoolP("regex", "r", false, "使用正则表达式")
	renameCmd.Flags().BoolP("dry-run", "d", false, "试运行（不实际重命名）")
	renameCmd.Flags().BoolP("recursive", "R", false, "递归处理子目录")

	// 绑定标志到 viper
	viper.BindPFlag("log_level", rootCmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("output_dir", rootCmd.PersistentFlags().Lookup("output-dir"))
	viper.BindPFlag("max_workers", rootCmd.PersistentFlags().Lookup("workers"))
	viper.BindPFlag("progress_bar", rootCmd.PersistentFlags().Lookup("progress"))
	viper.BindPFlag("backup_files", rootCmd.PersistentFlags().Lookup("backup"))

	// 添加子命令
	rootCmd.AddCommand(analyzeCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(replaceCmd)
	rootCmd.AddCommand(renameCmd)
}

// initializeApp 初始化应用程序
func initializeApp() {
	setupLogger()
	loadConfig()
	createOutputDir()
}

// setupLogger 设置日志器
func setupLogger() {
	logger = logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel(viper.GetString("log_level"))
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// 设置格式
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
}

// loadConfig 加载配置
func loadConfig() {
	// 设置配置文件搜索路径
	viper.SetConfigName("fileprocessor")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config/fileprocessor")

	// 设置环境变量前缀
	viper.SetEnvPrefix("FILEPROCESSOR")
	viper.AutomaticEnv()

	// 设置默认值
	viper.SetDefault("log_level", "info")
	viper.SetDefault("output_dir", "./output")
	viper.SetDefault("max_workers", 4)
	viper.SetDefault("buffer_size", 4096)
	viper.SetDefault("progress_bar", true)
	viper.SetDefault("backup_files", false)
	viper.SetDefault("ignore_hidden", true)

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			logger.Warnf("读取配置文件失败: %v", err)
		}
	}

	// 解析配置到结构体
	config = &Config{}
	if err := viper.Unmarshal(config); err != nil {
		logger.Fatalf("解析配置失败: %v", err)
	}
}

// createOutputDir 创建输出目录
func createOutputDir() {
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		logger.Fatalf("创建输出目录失败: %v", err)
	}
}

// runAnalyze 执行分析命令
func runAnalyze(cmd *cobra.Command, args []string) {
	format, _ := cmd.Flags().GetString("format")
	recursive, _ := cmd.Flags().GetBool("recursive")
	include, _ := cmd.Flags().GetString("include")
	exclude, _ := cmd.Flags().GetString("exclude")
	detailed, _ := cmd.Flags().GetBool("detailed")

	if len(args) == 0 {
		args = []string{"."}
	}

	logger.Infof("开始分析文件，格式: %s", format)

	// 收集文件列表
	files, err := collectFiles(args, recursive, include, exclude)
	if err != nil {
		logger.Fatalf("收集文件失败: %v", err)
	}

	logger.Infof("找到 %d 个文件", len(files))

	// 分析文件
	results := analyzeFiles(files)

	// 输出结果
	switch format {
	case "json":
		outputJSON(results)
	case "csv":
		outputCSV(results)
	default:
		outputTable(results, detailed)
	}

	// 保存结果到文件
	saveResults(results, format)
}

// runSearch 执行搜索命令
func runSearch(cmd *cobra.Command, args []string) {
	pattern := args[0]
	dirs := args[1:]
	if len(dirs) == 0 {
		dirs = []string{"."}
	}

	useRegex, _ := cmd.Flags().GetBool("regex")
	caseSensitive, _ := cmd.Flags().GetBool("case-sensitive")
	context, _ := cmd.Flags().GetInt("context")
	include, _ := cmd.Flags().GetString("include")
	exclude, _ := cmd.Flags().GetString("exclude")
	lineNumbers, _ := cmd.Flags().GetBool("line-numbers")

	logger.Infof("搜索模式: %s", pattern)

	// 编译搜索模式
	var regex *regexp.Regexp
	var err error

	if useRegex {
		flags := ""
		if !caseSensitive {
			flags = "(?i)"
		}
		regex, err = regexp.Compile(flags + pattern)
		if err != nil {
			logger.Fatalf("无效的正则表达式: %v", err)
		}
	}

	// 收集文件
	files, err := collectFiles(dirs, true, include, exclude)
	if err != nil {
		logger.Fatalf("收集文件失败: %v", err)
	}

	// 搜索文件
	searchFiles(files, pattern, regex, caseSensitive, context, lineNumbers)
}

// runReplace 执行替换命令
func runReplace(cmd *cobra.Command, args []string) {
	searchPattern := args[0]
	replacement := args[1]
	files := args[2:]

	useRegex, _ := cmd.Flags().GetBool("regex")
	caseSensitive, _ := cmd.Flags().GetBool("case-sensitive")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	include, _ := cmd.Flags().GetString("include")
	exclude, _ := cmd.Flags().GetString("exclude")

	logger.Infof("替换操作: '%s' -> '%s'", searchPattern, replacement)

	if dryRun {
		logger.Info("试运行模式，不会实际修改文件")
	}

	// 收集文件
	targetFiles, err := collectFiles(files, false, include, exclude)
	if err != nil {
		logger.Fatalf("收集文件失败: %v", err)
	}

	// 执行替换
	results := replaceInFiles(targetFiles, searchPattern, replacement, useRegex, caseSensitive, dryRun)

	// 输出结果
	outputReplaceResults(results)
}

// runRename 执行重命名命令
func runRename(cmd *cobra.Command, args []string) {
	pattern := args[0]
	replacement := args[1]
	directory := args[2]

	useRegex, _ := cmd.Flags().GetBool("regex")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	recursive, _ := cmd.Flags().GetBool("recursive")

	logger.Infof("重命名操作: '%s' -> '%s' 在目录 %s", pattern, replacement, directory)

	if dryRun {
		logger.Info("试运行模式，不会实际重命名文件")
	}

	// 收集文件
	files, err := collectFiles([]string{directory}, recursive, "*", "")
	if err != nil {
		logger.Fatalf("收集文件失败: %v", err)
	}

	// 执行重命名
	renameFiles(files, pattern, replacement, useRegex, dryRun)
}

// collectFiles 收集文件列表
func collectFiles(paths []string, recursive bool, include, exclude string) ([]string, error) {
	var files []string
	var mu sync.Mutex

	includePattern, err := filepath.Match(include, "")
	if err != nil && include != "*" {
		return nil, fmt.Errorf("无效的包含模式: %v", err)
	}

	var excludePattern *regexp.Regexp
	if exclude != "" {
		excludePattern, err = regexp.Compile(exclude)
		if err != nil {
			return nil, fmt.Errorf("无效的排除模式: %v", err)
		}
	}

	for _, path := range paths {
		err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				logger.Warnf("访问文件失败 %s: %v", filePath, err)
				return nil
			}

			// 跳过目录（除非需要递归处理）
			if info.IsDir() && !recursive && filePath != path {
				return filepath.SkipDir
			}

			// 跳过隐藏文件
			if config.IgnoreHidden && strings.HasPrefix(info.Name(), ".") {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}

			// 检查包含模式
			if include != "*" {
				matched, _ := filepath.Match(include, info.Name())
				if !matched {
					return nil
				}
			}

			// 检查排除模式
			if excludePattern != nil && excludePattern.MatchString(filePath) {
				return nil
			}

			if !info.IsDir() {
				mu.Lock()
				files = append(files, filePath)
				mu.Unlock()
			}

			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	return files, nil
}

// analyzeFiles 分析文件
func analyzeFiles(files []string) []FileInfo {
	var results []FileInfo
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 创建工作池
	jobs := make(chan string, len(files))

	// 启动工作协程
	for i := 0; i < config.MaxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range jobs {
				info := analyzeFile(file)
				if info != nil {
					mu.Lock()
					results = append(results, *info)
					mu.Unlock()
				}
			}
		}()
	}

	// 发送任务
	go func() {
		defer close(jobs)
		for _, file := range files {
			jobs <- file
		}
	}()

	// 等待完成
	wg.Wait()

	return results
}

// analyzeFile 分析单个文件
func analyzeFile(filePath string) *FileInfo {
	info, err := os.Stat(filePath)
	if err != nil {
		logger.Warnf("获取文件信息失败 %s: %v", filePath, err)
		return nil
	}

	fileInfo := &FileInfo{
		Path:        filePath,
		Size:        info.Size(),
		ModTime:     info.ModTime(),
		IsDirectory: info.IsDir(),
		Extension:   filepath.Ext(filePath),
		Permissions: info.Mode().String(),
	}

	if !info.IsDir() {
		// 读取文件内容进行分析
		content, err := os.ReadFile(filePath)
		if err != nil {
			logger.Warnf("读取文件失败 %s: %v", filePath, err)
			return fileInfo
		}

		text := string(content)
		fileInfo.Lines = strings.Count(text, "\n") + 1
		fileInfo.Words = len(strings.Fields(text))
		fileInfo.Characters = len(text)
	}

	return fileInfo
}

// searchFiles 搜索文件内容
func searchFiles(files []string, pattern string, regex *regexp.Regexp, caseSensitive bool, context int, lineNumbers bool) {
	var totalMatches int

	for _, file := range files {
		matches := searchInFile(file, pattern, regex, caseSensitive, context, lineNumbers)
		totalMatches += matches
	}

	logger.Infof("搜索完成，共找到 %d 个匹配", totalMatches)
}

// searchInFile 在单个文件中搜索
func searchInFile(filePath, pattern string, regex *regexp.Regexp, caseSensitive bool, context int, lineNumbers bool) int {
	file, err := os.Open(filePath)
	if err != nil {
		logger.Warnf("打开文件失败 %s: %v", filePath, err)
		return 0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0
	matches := 0
	var contextLines []string

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		var matched bool
		if regex != nil {
			matched = regex.MatchString(line)
		} else {
			searchText := pattern
			lineText := line
			if !caseSensitive {
				searchText = strings.ToLower(searchText)
				lineText = strings.ToLower(lineText)
			}
			matched = strings.Contains(lineText, searchText)
		}

		if matched {
			matches++
			if matches == 1 {
				fmt.Printf("\n📁 %s:\n", filePath)
			}

			if lineNumbers {
				fmt.Printf("%4d: %s\n", lineNum, line)
			} else {
				fmt.Printf("%s\n", line)
			}

			// 显示上下文
			if context > 0 {
				// 这里可以实现上下文显示逻辑
			}
		}

		// 保存上下文行
		if context > 0 {
			contextLines = append(contextLines, line)
			if len(contextLines) > context*2+1 {
				contextLines = contextLines[1:]
			}
		}
	}

	return matches
}

// replaceInFiles 在文件中执行替换
func replaceInFiles(files []string, searchPattern, replacement string, useRegex, caseSensitive, dryRun bool) []ProcessResult {
	var results []ProcessResult
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 创建进度条
	var bar *progressbar.ProgressBar
	if config.ProgressBar {
		bar = progressbar.Default(int64(len(files)))
	}

	// 创建工作池
	jobs := make(chan string, len(files))

	// 启动工作协程
	for i := 0; i < config.MaxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range jobs {
				result := replaceInFile(file, searchPattern, replacement, useRegex, caseSensitive, dryRun)

				mu.Lock()
				results = append(results, result)
				if bar != nil {
					bar.Add(1)
				}
				mu.Unlock()
			}
		}()
	}

	// 发送任务
	go func() {
		defer close(jobs)
		for _, file := range files {
			jobs <- file
		}
	}()

	// 等待完成
	wg.Wait()

	if bar != nil {
		bar.Finish()
	}

	return results
}

// replaceInFile 在单个文件中执行替换
func replaceInFile(filePath, searchPattern, replacement string, useRegex, caseSensitive, dryRun bool) ProcessResult {
	start := time.Now()
	result := ProcessResult{
		File:    filePath,
		Success: false,
	}

	// 读取文件
	content, err := os.ReadFile(filePath)
	if err != nil {
		result.Error = err.Error()
		result.Duration = time.Since(start)
		return result
	}

	result.BytesRead = int64(len(content))
	text := string(content)
	originalLines := strings.Count(text, "\n") + 1

	// 执行替换
	var newText string
	if useRegex {
		flags := ""
		if !caseSensitive {
			flags = "(?i)"
		}
		regex, err := regexp.Compile(flags + searchPattern)
		if err != nil {
			result.Error = fmt.Sprintf("无效的正则表达式: %v", err)
			result.Duration = time.Since(start)
			return result
		}
		newText = regex.ReplaceAllString(text, replacement)
	} else {
		if caseSensitive {
			newText = strings.ReplaceAll(text, searchPattern, replacement)
		} else {
			// 大小写不敏感的替换
			newText = replaceAllIgnoreCase(text, searchPattern, replacement)
		}
	}

	newLines := strings.Count(newText, "\n") + 1
	result.LinesChanged = abs(newLines - originalLines)

	// 如果内容没有变化，跳过写入
	if newText == text {
		result.Success = true
		result.Duration = time.Since(start)
		return result
	}

	if !dryRun {
		// 创建备份
		if config.BackupFiles {
			backupPath := filePath + ".bak"
			if err := os.WriteFile(backupPath, content, 0644); err != nil {
				logger.Warnf("创建备份文件失败 %s: %v", backupPath, err)
			}
		}

		// 写入新内容
		err = os.WriteFile(filePath, []byte(newText), 0644)
		if err != nil {
			result.Error = err.Error()
			result.Duration = time.Since(start)
			return result
		}
	}

	result.Success = true
	result.BytesWritten = int64(len(newText))
	result.Duration = time.Since(start)

	return result
}

// renameFiles 重命名文件
func renameFiles(files []string, pattern, replacement string, useRegex, dryRun bool) {
	var renamed int

	for _, file := range files {
		dir := filepath.Dir(file)
		oldName := filepath.Base(file)
		var newName string

		if useRegex {
			regex, err := regexp.Compile(pattern)
			if err != nil {
				logger.Errorf("无效的正则表达式: %v", err)
				continue
			}
			newName = regex.ReplaceAllString(oldName, replacement)
		} else {
			newName = strings.ReplaceAll(oldName, pattern, replacement)
		}

		if newName == oldName {
			continue
		}

		newPath := filepath.Join(dir, newName)

		if dryRun {
			fmt.Printf("重命名: %s -> %s\n", file, newPath)
		} else {
			if err := os.Rename(file, newPath); err != nil {
				logger.Errorf("重命名失败 %s -> %s: %v", file, newPath, err)
			} else {
				fmt.Printf("✅ 重命名: %s -> %s\n", file, newPath)
				renamed++
			}
		}
	}

	if dryRun {
		fmt.Printf("试运行完成，将重命名 %d 个文件\n", renamed)
	} else {
		fmt.Printf("重命名完成，成功重命名 %d 个文件\n", renamed)
	}
}

// 输出函数
func outputTable(results []FileInfo, detailed bool) {
	if detailed {
		fmt.Printf("%-40s %10s %8s %8s %10s %20s\n", "文件路径", "大小", "行数", "词数", "字符数", "修改时间")
		fmt.Println(strings.Repeat("-", 100))
		for _, info := range results {
			fmt.Printf("%-40s %10d %8d %8d %10d %20s\n",
				truncateString(info.Path, 40), info.Size, info.Lines, info.Words, info.Characters,
				info.ModTime.Format("2006-01-02 15:04:05"))
		}
	} else {
		fmt.Printf("%-50s %10s %8s\n", "文件路径", "大小", "行数")
		fmt.Println(strings.Repeat("-", 70))
		for _, info := range results {
			fmt.Printf("%-50s %10d %8d\n", truncateString(info.Path, 50), info.Size, info.Lines)
		}
	}
}

func outputJSON(results []FileInfo) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.Encode(results)
}

func outputCSV(results []FileInfo) {
	fmt.Println("路径,大小,行数,词数,字符数,修改时间")
	for _, info := range results {
		fmt.Printf("%s,%d,%d,%d,%d,%s\n",
			info.Path, info.Size, info.Lines, info.Words, info.Characters,
			info.ModTime.Format("2006-01-02 15:04:05"))
	}
}

func outputReplaceResults(results []ProcessResult) {
	successful := 0
	failed := 0
	totalDuration := time.Duration(0)

	for _, result := range results {
		totalDuration += result.Duration
		if result.Success {
			successful++
			if result.LinesChanged > 0 {
				fmt.Printf("✅ %s (修改了 %d 行)\n", result.File, result.LinesChanged)
			}
		} else {
			failed++
			fmt.Printf("❌ %s: %s\n", result.File, result.Error)
		}
	}

	fmt.Printf("\n📊 替换统计:\n")
	fmt.Printf("  成功: %d 个文件\n", successful)
	fmt.Printf("  失败: %d 个文件\n", failed)
	fmt.Printf("  总耗时: %v\n", totalDuration)
}

// 辅助函数
func saveResults(results []FileInfo, format string) {
	filename := filepath.Join(config.OutputDir, fmt.Sprintf("analysis_%s.%s",
		time.Now().Format("20060102_150405"), format))

	file, err := os.Create(filename)
	if err != nil {
		logger.Errorf("创建输出文件失败: %v", err)
		return
	}
	defer file.Close()

	switch format {
	case "json":
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		encoder.Encode(results)
	case "csv":
		fmt.Fprintln(file, "路径,大小,行数,词数,字符数,修改时间")
		for _, info := range results {
			fmt.Fprintf(file, "%s,%d,%d,%d,%d,%s\n",
				info.Path, info.Size, info.Lines, info.Words, info.Characters,
				info.ModTime.Format("2006-01-02 15:04:05"))
		}
	}

	logger.Infof("结果已保存到: %s", filename)
}

func replaceAllIgnoreCase(text, old, new string) string {
	// 简化的大小写不敏感替换
	lowerText := strings.ToLower(text)
	lowerOld := strings.ToLower(old)

	result := ""
	start := 0

	for {
		index := strings.Index(lowerText[start:], lowerOld)
		if index == -1 {
			result += text[start:]
			break
		}

		actualIndex := start + index
		result += text[start:actualIndex] + new
		start = actualIndex + len(old)
	}

	return result
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
