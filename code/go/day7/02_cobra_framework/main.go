// Package main 演示 Cobra 框架的使用
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// 全局配置变量
	cfgFile string
	verbose bool
	output  string
)

// rootCmd 根命令
var rootCmd = &cobra.Command{
	Use:   "filetools",
	Short: "一个强大的文件处理工具",
	Long: `FileTools 是一个使用 Cobra 构建的文件处理工具。
它提供了多种文件操作功能，包括搜索、统计、转换等。

示例:
  filetools search -p "*.go" -t "func"
  filetools stats file1.txt file2.txt
  filetools convert -f json -t yaml input.json`,
	
	// 全局前置钩子
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if verbose {
			fmt.Printf("执行命令: %s\n", cmd.CommandPath())
			fmt.Printf("参数: %v\n", args)
		}
	},
}

// searchCmd 搜索命令
var searchCmd = &cobra.Command{
	Use:   "search [目录]",
	Short: "在文件中搜索内容",
	Long: `在指定目录中搜索文件内容。
支持正则表达式和多种文件类型过滤。`,
	Args: cobra.MaximumNArgs(1),
	Run:  runSearch,
}

// statsCmd 统计命令
var statsCmd = &cobra.Command{
	Use:   "stats [文件...]",
	Short: "统计文件信息",
	Long:  `统计文件的行数、字数、字符数等信息。`,
	Args:  cobra.MinimumNArgs(1),
	Run:   runStats,
}

// convertCmd 转换命令
var convertCmd = &cobra.Command{
	Use:   "convert [文件]",
	Short: "转换文件格式",
	Long:  `支持在 JSON、YAML、TOML 等格式之间转换。`,
	Args:  cobra.ExactArgs(1),
	Run:   runConvert,
}

// versionCmd 版本命令
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("FileTools v1.0.0\n")
		fmt.Printf("构建时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
		fmt.Printf("Go 版本: %s\n", "go1.22")
	},
}

func main() {
	// 执行根命令
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// 初始化 Cobra
	cobra.OnInitialize(initConfig)

	// 全局持久化标志
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "配置文件路径 (默认: $HOME/.filetools.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "详细输出")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "输出文件路径")

	// 搜索命令标志
	searchCmd.Flags().StringP("pattern", "p", "*", "文件名模式")
	searchCmd.Flags().StringP("text", "t", "", "搜索文本")
	searchCmd.Flags().BoolP("regex", "r", false, "使用正则表达式")
	searchCmd.Flags().BoolP("recursive", "R", false, "递归搜索")
	searchCmd.Flags().IntP("context", "C", 0, "显示上下文行数")

	// 统计命令标志
	statsCmd.Flags().BoolP("detailed", "d", false, "显示详细统计")
	statsCmd.Flags().StringP("format", "f", "table", "输出格式 (table|json|csv)")

	// 转换命令标志
	convertCmd.Flags().StringP("from", "f", "", "源格式 (json|yaml|toml)")
	convertCmd.Flags().StringP("to", "t", "", "目标格式 (json|yaml|toml)")
	convertCmd.Flags().BoolP("pretty", "p", false, "美化输出")

	// 标记必需标志
	convertCmd.MarkFlagRequired("from")
	convertCmd.MarkFlagRequired("to")

	// 绑定标志到 Viper
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))

	// 添加子命令
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(statsCmd)
	rootCmd.AddCommand(convertCmd)
	rootCmd.AddCommand(versionCmd)
}

// initConfig 初始化配置
func initConfig() {
	if cfgFile != "" {
		// 使用指定的配置文件
		viper.SetConfigFile(cfgFile)
	} else {
		// 查找主目录
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// 搜索配置文件
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".filetools")
	}

	// 自动读取环境变量
	viper.AutomaticEnv()
	viper.SetEnvPrefix("FILETOOLS")

	// 读取配置文件
	if err := viper.ReadInConfig(); err == nil && verbose {
		fmt.Printf("使用配置文件: %s\n", viper.ConfigFileUsed())
	}
}

// runSearch 执行搜索命令
func runSearch(cmd *cobra.Command, args []string) {
	// 获取标志值
	pattern, _ := cmd.Flags().GetString("pattern")
	text, _ := cmd.Flags().GetString("text")
	regex, _ := cmd.Flags().GetBool("regex")
	recursive, _ := cmd.Flags().GetBool("recursive")
	context, _ := cmd.Flags().GetInt("context")

	// 确定搜索目录
	searchDir := "."
	if len(args) > 0 {
		searchDir = args[0]
	}

	if verbose {
		fmt.Printf("搜索配置:\n")
		fmt.Printf("  目录: %s\n", searchDir)
		fmt.Printf("  文件模式: %s\n", pattern)
		fmt.Printf("  搜索文本: %s\n", text)
		fmt.Printf("  正则表达式: %t\n", regex)
		fmt.Printf("  递归搜索: %t\n", recursive)
		fmt.Printf("  上下文行数: %d\n", context)
		fmt.Printf("\n")
	}

	// 验证参数
	if text == "" {
		fmt.Fprintf(os.Stderr, "错误: 必须指定搜索文本 (-t 或 --text)\n")
		cmd.Usage()
		os.Exit(1)
	}

	// 执行搜索
	results, err := performSearch(searchDir, pattern, text, recursive, regex, context)
	if err != nil {
		fmt.Fprintf(os.Stderr, "搜索失败: %v\n", err)
		os.Exit(1)
	}

	// 输出结果
	if len(results) == 0 {
		fmt.Println("未找到匹配的内容")
	} else {
		fmt.Printf("找到 %d 个匹配结果:\n\n", len(results))
		for _, result := range results {
			fmt.Println(result)
		}
	}

	// 保存到文件
	if output != "" {
		saveResults(output, results)
	}
}

// runStats 执行统计命令
func runStats(cmd *cobra.Command, args []string) {
	detailed, _ := cmd.Flags().GetBool("detailed")
	format, _ := cmd.Flags().GetString("format")

	if verbose {
		fmt.Printf("统计 %d 个文件，格式: %s\n\n", len(args), format)
	}

	var allStats []FileStats
	for _, filename := range args {
		stats, err := getFileStats(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "统计文件 %s 失败: %v\n", filename, err)
			continue
		}
		allStats = append(allStats, stats)
	}

	// 输出统计结果
	outputStats(allStats, format, detailed)

	// 保存到文件
	if output != "" {
		saveStatsToFile(output, allStats, format)
	}
}

// runConvert 执行转换命令
func runConvert(cmd *cobra.Command, args []string) {
	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")
	pretty, _ := cmd.Flags().GetBool("pretty")

	filename := args[0]

	if verbose {
		fmt.Printf("转换文件: %s (%s -> %s)\n", filename, from, to)
	}

	// 执行转换
	result, err := convertFile(filename, from, to, pretty)
	if err != nil {
		fmt.Fprintf(os.Stderr, "转换失败: %v\n", err)
		os.Exit(1)
	}

	// 输出结果
	if output != "" {
		err := os.WriteFile(output, []byte(result), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "保存文件失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("转换结果已保存到: %s\n", output)
	} else {
		fmt.Println(result)
	}
}

// SearchResult 搜索结果
type SearchResult struct {
	File    string
	Line    int
	Content string
	Context []string
}

// FileStats 文件统计信息
type FileStats struct {
	Filename  string
	Size      int64
	Lines     int
	Words     int
	Chars     int
	Extension string
}

// performSearch 执行搜索
func performSearch(dir, pattern, text string, recursive, regex bool, context int) ([]string, error) {
	// 这里是搜索逻辑的简化实现
	var results []string
	
	// 模拟搜索结果
	results = append(results, fmt.Sprintf("在文件 example.go 第 10 行找到: %s", text))
	results = append(results, fmt.Sprintf("在文件 main.go 第 25 行找到: %s", text))
	
	return results, nil
}

// getFileStats 获取文件统计信息
func getFileStats(filename string) (FileStats, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return FileStats{}, err
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		return FileStats{}, err
	}

	lines := strings.Count(string(content), "\n") + 1
	words := len(strings.Fields(string(content)))
	chars := len(content)

	return FileStats{
		Filename:  filename,
		Size:      info.Size(),
		Lines:     lines,
		Words:     words,
		Chars:     chars,
		Extension: filepath.Ext(filename),
	}, nil
}

// outputStats 输出统计结果
func outputStats(stats []FileStats, format string, detailed bool) {
	switch format {
	case "table":
		outputStatsTable(stats, detailed)
	case "json":
		outputStatsJSON(stats)
	case "csv":
		outputStatsCSV(stats)
	default:
		fmt.Fprintf(os.Stderr, "不支持的格式: %s\n", format)
	}
}

// outputStatsTable 表格格式输出
func outputStatsTable(stats []FileStats, detailed bool) {
	fmt.Printf("%-20s %10s %8s %8s %8s\n", "文件名", "大小", "行数", "词数", "字符数")
	fmt.Println(strings.Repeat("-", 60))
	
	for _, stat := range stats {
		fmt.Printf("%-20s %10d %8d %8d %8d\n",
			filepath.Base(stat.Filename), stat.Size, stat.Lines, stat.Words, stat.Chars)
	}
}

// outputStatsJSON JSON 格式输出
func outputStatsJSON(stats []FileStats) {
	fmt.Println("JSON 格式输出 (简化)")
	for _, stat := range stats {
		fmt.Printf(`{"file": "%s", "size": %d, "lines": %d}`, 
			stat.Filename, stat.Size, stat.Lines)
		fmt.Println()
	}
}

// outputStatsCSV CSV 格式输出
func outputStatsCSV(stats []FileStats) {
	fmt.Println("文件名,大小,行数,词数,字符数")
	for _, stat := range stats {
		fmt.Printf("%s,%d,%d,%d,%d\n",
			stat.Filename, stat.Size, stat.Lines, stat.Words, stat.Chars)
	}
}

// convertFile 转换文件格式
func convertFile(filename, from, to string, pretty bool) (string, error) {
	// 这里是格式转换的简化实现
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	// 模拟转换过程
	result := fmt.Sprintf("# 转换结果 (%s -> %s)\n%s", from, to, string(content))
	
	if pretty {
		result = "# 美化后的输出\n" + result
	}

	return result, nil
}

// saveResults 保存搜索结果
func saveResults(filename string, results []string) {
	content := strings.Join(results, "\n")
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "保存结果失败: %v\n", err)
	} else {
		fmt.Printf("搜索结果已保存到: %s\n", filename)
	}
}

// saveStatsToFile 保存统计结果到文件
func saveStatsToFile(filename string, stats []FileStats, format string) {
	// 简化实现
	content := fmt.Sprintf("统计结果 (%s 格式)\n", format)
	for _, stat := range stats {
		content += fmt.Sprintf("%s: %d 字节\n", stat.Filename, stat.Size)
	}
	
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "保存统计结果失败: %v\n", err)
	} else {
		fmt.Printf("统计结果已保存到: %s\n", filename)
	}
}
