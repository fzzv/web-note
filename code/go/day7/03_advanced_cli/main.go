// Package main 演示高级 CLI 功能
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	logger *logrus.Logger
	config *AppConfig
)

// AppConfig 应用程序配置
type AppConfig struct {
	LogLevel    string `mapstructure:"log_level"`
	OutputDir   string `mapstructure:"output_dir"`
	MaxWorkers  int    `mapstructure:"max_workers"`
	Timeout     int    `mapstructure:"timeout"`
	Interactive bool   `mapstructure:"interactive"`
}

// Task 任务结构
type Task struct {
	ID          string
	Name        string
	Description string
	Status      string
	CreatedAt   time.Time
}

var rootCmd = &cobra.Command{
	Use:   "taskmanager",
	Short: "高级任务管理工具",
	Long: `TaskManager 是一个功能丰富的任务管理工具，支持：
- 交互式操作
- 进度条显示
- 配置文件管理
- 结构化日志
- 用户友好的界面`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		setupLogger()
		loadConfig()
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "创建新任务",
	Run:   runCreate,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有任务",
	Run:   runList,
}

var processCmd = &cobra.Command{
	Use:   "process",
	Short: "批量处理任务",
	Run:   runProcess,
}

var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "进入交互模式",
	Run:   runInteractive,
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "配置管理",
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "显示当前配置",
	Run:   runConfigShow,
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "设置配置项",
	Args:  cobra.ExactArgs(2),
	Run:   runConfigSet,
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
	rootCmd.PersistentFlags().String("log-level", "info", "日志级别")
	rootCmd.PersistentFlags().Bool("interactive", false, "交互模式")
	rootCmd.PersistentFlags().Int("workers", 3, "工作线程数")

	// create 命令标志
	createCmd.Flags().StringP("name", "n", "", "任务名称")
	createCmd.Flags().StringP("description", "d", "", "任务描述")
	createCmd.Flags().BoolP("interactive", "i", false, "交互式创建")

	// list 命令标志
	listCmd.Flags().StringP("status", "s", "", "按状态过滤")
	listCmd.Flags().StringP("format", "f", "table", "输出格式 (table|json|csv)")
	listCmd.Flags().BoolP("detailed", "d", false, "显示详细信息")

	// process 命令标志
	processCmd.Flags().IntP("batch-size", "b", 10, "批处理大小")
	processCmd.Flags().BoolP("progress", "p", true, "显示进度条")
	processCmd.Flags().IntP("delay", "d", 100, "处理延迟(毫秒)")

	// 绑定标志到 viper
	viper.BindPFlag("log_level", rootCmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("interactive", rootCmd.PersistentFlags().Lookup("interactive"))
	viper.BindPFlag("max_workers", rootCmd.PersistentFlags().Lookup("workers"))

	// 添加子命令
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(processCmd)
	rootCmd.AddCommand(interactiveCmd)
	rootCmd.AddCommand(configCmd)

	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
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
	viper.SetConfigName("taskmanager")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config/taskmanager")
	viper.AddConfigPath("/etc/taskmanager")

	// 设置环境变量前缀
	viper.SetEnvPrefix("TASKMANAGER")
	viper.AutomaticEnv()

	// 设置默认值
	viper.SetDefault("log_level", "info")
	viper.SetDefault("output_dir", "./output")
	viper.SetDefault("max_workers", 3)
	viper.SetDefault("timeout", 30)
	viper.SetDefault("interactive", false)

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			logger.Warnf("读取配置文件失败: %v", err)
		}
	} else {
		logger.Infof("使用配置文件: %s", viper.ConfigFileUsed())
	}

	// 解析配置到结构体
	config = &AppConfig{}
	if err := viper.Unmarshal(config); err != nil {
		logger.Fatalf("解析配置失败: %v", err)
	}
}

// runCreate 创建任务
func runCreate(cmd *cobra.Command, args []string) {
	interactive, _ := cmd.Flags().GetBool("interactive")
	
	var task Task
	var err error

	if interactive || config.Interactive {
		task, err = createTaskInteractive()
	} else {
		task, err = createTaskFromFlags(cmd)
	}

	if err != nil {
		logger.Errorf("创建任务失败: %v", err)
		return
	}

	logger.Infof("任务创建成功: %s", task.Name)
	fmt.Printf("✅ 任务 '%s' 创建成功\n", task.Name)
	fmt.Printf("   ID: %s\n", task.ID)
	fmt.Printf("   描述: %s\n", task.Description)
	fmt.Printf("   状态: %s\n", task.Status)
}

// createTaskInteractive 交互式创建任务
func createTaskInteractive() (Task, error) {
	fmt.Println("🚀 交互式任务创建")

	// 输入任务名称
	namePrompt := promptui.Prompt{
		Label: "任务名称",
		Validate: func(input string) error {
			if len(strings.TrimSpace(input)) < 3 {
				return fmt.Errorf("任务名称至少需要3个字符")
			}
			return nil
		},
	}
	name, err := namePrompt.Run()
	if err != nil {
		return Task{}, err
	}

	// 输入任务描述
	descPrompt := promptui.Prompt{
		Label: "任务描述",
	}
	description, err := descPrompt.Run()
	if err != nil {
		return Task{}, err
	}

	// 选择任务状态
	statusPrompt := promptui.Select{
		Label: "任务状态",
		Items: []string{"待处理", "进行中", "已完成", "已取消"},
	}
	_, status, err := statusPrompt.Run()
	if err != nil {
		return Task{}, err
	}

	// 确认创建
	confirmPrompt := promptui.Prompt{
		Label:     "确认创建任务",
		IsConfirm: true,
	}
	_, err = confirmPrompt.Run()
	if err != nil {
		return Task{}, fmt.Errorf("任务创建已取消")
	}

	return Task{
		ID:          generateTaskID(),
		Name:        name,
		Description: description,
		Status:      status,
		CreatedAt:   time.Now(),
	}, nil
}

// createTaskFromFlags 从命令行标志创建任务
func createTaskFromFlags(cmd *cobra.Command) (Task, error) {
	name, _ := cmd.Flags().GetString("name")
	description, _ := cmd.Flags().GetString("description")

	if name == "" {
		return Task{}, fmt.Errorf("任务名称不能为空")
	}

	return Task{
		ID:          generateTaskID(),
		Name:        name,
		Description: description,
		Status:      "待处理",
		CreatedAt:   time.Now(),
	}, nil
}

// runList 列出任务
func runList(cmd *cobra.Command, args []string) {
	status, _ := cmd.Flags().GetString("status")
	format, _ := cmd.Flags().GetString("format")
	detailed, _ := cmd.Flags().GetBool("detailed")

	logger.Debugf("列出任务: status=%s, format=%s, detailed=%t", status, format, detailed)

	// 模拟获取任务列表
	tasks := getMockTasks()

	// 过滤任务
	if status != "" {
		tasks = filterTasksByStatus(tasks, status)
	}

	// 输出任务列表
	switch format {
	case "table":
		outputTasksTable(tasks, detailed)
	case "json":
		outputTasksJSON(tasks)
	case "csv":
		outputTasksCSV(tasks)
	default:
		logger.Errorf("不支持的输出格式: %s", format)
	}
}

// runProcess 批量处理任务
func runProcess(cmd *cobra.Command, args []string) {
	batchSize, _ := cmd.Flags().GetInt("batch-size")
	showProgress, _ := cmd.Flags().GetBool("progress")
	delay, _ := cmd.Flags().GetInt("delay")

	logger.Infof("开始批量处理任务: batch_size=%d, delay=%dms", batchSize, delay)

	tasks := getMockTasks()
	
	if showProgress {
		processTasksWithProgress(tasks, batchSize, delay)
	} else {
		processTasksSimple(tasks, batchSize, delay)
	}

	fmt.Println("✅ 批量处理完成")
}

// processTasksWithProgress 带进度条的任务处理
func processTasksWithProgress(tasks []Task, batchSize, delay int) {
	bar := progressbar.NewOptions(len(tasks),
		progressbar.OptionSetDescription("处理任务中..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionSetWidth(50),
	)

	for i, task := range tasks {
		// 模拟处理任务
		time.Sleep(time.Duration(delay) * time.Millisecond)
		
		logger.Debugf("处理任务: %s", task.Name)
		bar.Add(1)

		// 批处理检查点
		if (i+1)%batchSize == 0 {
			fmt.Printf("\n📦 完成批次 %d/%d\n", (i+1)/batchSize, (len(tasks)+batchSize-1)/batchSize)
		}
	}
	
	bar.Finish()
	fmt.Println()
}

// processTasksSimple 简单任务处理
func processTasksSimple(tasks []Task, batchSize, delay int) {
	for i, task := range tasks {
		time.Sleep(time.Duration(delay) * time.Millisecond)
		fmt.Printf("处理任务 %d/%d: %s\n", i+1, len(tasks), task.Name)
	}
}

// runInteractive 交互模式
func runInteractive(cmd *cobra.Command, args []string) {
	fmt.Println("🎯 进入交互模式")
	fmt.Println("输入 'help' 查看可用命令，输入 'exit' 退出")

	scanner := bufio.NewScanner(os.Stdin)
	
	for {
		fmt.Print("taskmanager> ")
		
		if !scanner.Scan() {
			break
		}
		
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		
		if input == "exit" || input == "quit" {
			fmt.Println("👋 再见！")
			break
		}
		
		handleInteractiveCommand(input)
	}
}

// handleInteractiveCommand 处理交互命令
func handleInteractiveCommand(input string) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return
	}
	
	command := parts[0]
	args := parts[1:]
	
	switch command {
	case "help":
		showInteractiveHelp()
	case "list":
		fmt.Println("📋 任务列表:")
		tasks := getMockTasks()
		outputTasksTable(tasks, false)
	case "create":
		if len(args) > 0 {
			name := strings.Join(args, " ")
			task := Task{
				ID:        generateTaskID(),
				Name:      name,
				Status:    "待处理",
				CreatedAt: time.Now(),
			}
			fmt.Printf("✅ 创建任务: %s\n", task.Name)
		} else {
			fmt.Println("❌ 请提供任务名称")
		}
	case "status":
		fmt.Printf("📊 系统状态: 正常运行\n")
		fmt.Printf("   配置文件: %s\n", viper.ConfigFileUsed())
		fmt.Printf("   日志级别: %s\n", config.LogLevel)
		fmt.Printf("   工作线程: %d\n", config.MaxWorkers)
	default:
		fmt.Printf("❌ 未知命令: %s\n", command)
		fmt.Println("输入 'help' 查看可用命令")
	}
}

// showInteractiveHelp 显示交互模式帮助
func showInteractiveHelp() {
	fmt.Println("📚 可用命令:")
	fmt.Println("  help              - 显示此帮助信息")
	fmt.Println("  list              - 列出所有任务")
	fmt.Println("  create <name>     - 创建新任务")
	fmt.Println("  status            - 显示系统状态")
	fmt.Println("  exit/quit         - 退出交互模式")
}

// runConfigShow 显示配置
func runConfigShow(cmd *cobra.Command, args []string) {
	fmt.Println("📋 当前配置:")
	fmt.Printf("  日志级别: %s\n", config.LogLevel)
	fmt.Printf("  输出目录: %s\n", config.OutputDir)
	fmt.Printf("  最大工作线程: %d\n", config.MaxWorkers)
	fmt.Printf("  超时时间: %d 秒\n", config.Timeout)
	fmt.Printf("  交互模式: %t\n", config.Interactive)
	
	if viper.ConfigFileUsed() != "" {
		fmt.Printf("  配置文件: %s\n", viper.ConfigFileUsed())
	}
}

// runConfigSet 设置配置
func runConfigSet(cmd *cobra.Command, args []string) {
	key := args[0]
	value := args[1]
	
	// 验证配置键
	validKeys := map[string]bool{
		"log_level":    true,
		"output_dir":   true,
		"max_workers":  true,
		"timeout":      true,
		"interactive":  true,
	}
	
	if !validKeys[key] {
		fmt.Printf("❌ 无效的配置键: %s\n", key)
		fmt.Println("有效的配置键: log_level, output_dir, max_workers, timeout, interactive")
		return
	}
	
	// 类型转换
	switch key {
	case "max_workers", "timeout":
		if _, err := strconv.Atoi(value); err != nil {
			fmt.Printf("❌ %s 必须是数字\n", key)
			return
		}
	case "interactive":
		if _, err := strconv.ParseBool(value); err != nil {
			fmt.Printf("❌ %s 必须是布尔值 (true/false)\n", key)
			return
		}
	}
	
	viper.Set(key, value)
	fmt.Printf("✅ 设置 %s = %s\n", key, value)
	
	// 保存配置文件
	if err := viper.WriteConfig(); err != nil {
		fmt.Printf("⚠️  保存配置文件失败: %v\n", err)
	} else {
		fmt.Println("💾 配置已保存")
	}
}

// 辅助函数

func generateTaskID() string {
	return fmt.Sprintf("task_%d", time.Now().Unix())
}

func getMockTasks() []Task {
	return []Task{
		{ID: "task_1", Name: "数据备份", Description: "备份数据库", Status: "待处理", CreatedAt: time.Now().Add(-2 * time.Hour)},
		{ID: "task_2", Name: "系统更新", Description: "更新系统软件", Status: "进行中", CreatedAt: time.Now().Add(-1 * time.Hour)},
		{ID: "task_3", Name: "日志清理", Description: "清理旧日志文件", Status: "已完成", CreatedAt: time.Now().Add(-30 * time.Minute)},
		{ID: "task_4", Name: "性能监控", Description: "监控系统性能", Status: "待处理", CreatedAt: time.Now().Add(-15 * time.Minute)},
	}
}

func filterTasksByStatus(tasks []Task, status string) []Task {
	var filtered []Task
	for _, task := range tasks {
		if task.Status == status {
			filtered = append(filtered, task)
		}
	}
	return filtered
}

func outputTasksTable(tasks []Task, detailed bool) {
	if detailed {
		fmt.Printf("%-12s %-20s %-30s %-10s %-20s\n", "ID", "名称", "描述", "状态", "创建时间")
		fmt.Println(strings.Repeat("-", 95))
		for _, task := range tasks {
			fmt.Printf("%-12s %-20s %-30s %-10s %-20s\n",
				task.ID, task.Name, task.Description, task.Status,
				task.CreatedAt.Format("2006-01-02 15:04:05"))
		}
	} else {
		fmt.Printf("%-12s %-20s %-10s\n", "ID", "名称", "状态")
		fmt.Println(strings.Repeat("-", 45))
		for _, task := range tasks {
			fmt.Printf("%-12s %-20s %-10s\n", task.ID, task.Name, task.Status)
		}
	}
}

func outputTasksJSON(tasks []Task) {
	fmt.Println("JSON 格式输出:")
	for _, task := range tasks {
		fmt.Printf(`{"id": "%s", "name": "%s", "status": "%s"}`, 
			task.ID, task.Name, task.Status)
		fmt.Println()
	}
}

func outputTasksCSV(tasks []Task) {
	fmt.Println("ID,名称,状态,创建时间")
	for _, task := range tasks {
		fmt.Printf("%s,%s,%s,%s\n",
			task.ID, task.Name, task.Status, task.CreatedAt.Format("2006-01-02 15:04:05"))
	}
}
