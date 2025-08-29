// Package main 演示代码质量工具的使用
package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// User 用户结构体
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserService 用户服务
type UserService struct {
	users map[int]*User
}

// NewUserService 创建用户服务
func NewUserService() *UserService {
	return &UserService{
		users: make(map[int]*User),
	}
}

// AddUser 添加用户
func (s *UserService) AddUser(user *User) error {
	if user == nil {
		return errors.New("用户不能为空")
	}
	
	if user.Name == "" {
		return errors.New("用户名不能为空")
	}
	
	if user.Email == "" {
		return errors.New("邮箱不能为空")
	}
	
	// 检查邮箱格式（简单验证）
	if !strings.Contains(user.Email, "@") {
		return errors.New("邮箱格式无效")
	}
	
	s.users[user.ID] = user
	return nil
}

// GetUser 获取用户
func (s *UserService) GetUser(id int) (*User, error) {
	user, exists := s.users[id]
	if !exists {
		return nil, fmt.Errorf("用户 %d 不存在", id)
	}
	return user, nil
}

// GetAllUsers 获取所有用户
func (s *UserService) GetAllUsers() []*User {
	users := make([]*User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users
}

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) bool {
	// 简单的邮箱验证
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// ProcessNumbers 处理数字列表
func ProcessNumbers(numbers []int) (sum, avg float64, err error) {
	if len(numbers) == 0 {
		return 0, 0, errors.New("数字列表不能为空")
	}
	
	total := 0
	for _, num := range numbers {
		total += num
	}
	
	sum = float64(total)
	avg = sum / float64(len(numbers))
	
	return sum, avg, nil
}

// FileProcessor 文件处理器
type FileProcessor struct {
	filename string
}

// NewFileProcessor 创建文件处理器
func NewFileProcessor(filename string) *FileProcessor {
	return &FileProcessor{
		filename: filename,
	}
}

// ProcessFile 处理文件
func (fp *FileProcessor) ProcessFile() error {
	// 检查文件是否存在
	if _, err := os.Stat(fp.filename); os.IsNotExist(err) {
		return fmt.Errorf("文件 %s 不存在", fp.filename)
	}
	
	// 读取文件内容
	content, err := os.ReadFile(fp.filename)
	if err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}
	
	// 处理内容
	lines := strings.Split(string(content), "\n")
	fmt.Printf("文件 %s 包含 %d 行\n", fp.filename, len(lines))
	
	return nil
}

// Calculator 计算器
type Calculator struct{}

// Add 加法
func (c *Calculator) Add(a, b int) int {
	return a + b
}

// Subtract 减法
func (c *Calculator) Subtract(a, b int) int {
	return a - b
}

// Multiply 乘法
func (c *Calculator) Multiply(a, b int) int {
	return a * b
}

// Divide 除法
func (c *Calculator) Divide(a, b int) (float64, error) {
	if b == 0 {
		return 0, errors.New("除数不能为零")
	}
	return float64(a) / float64(b), nil
}

// StringUtils 字符串工具
type StringUtils struct{}

// Reverse 反转字符串
func (su *StringUtils) Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// IsPalindrome 检查是否为回文
func (su *StringUtils) IsPalindrome(s string) bool {
	s = strings.ToLower(s)
	return s == su.Reverse(s)
}

// WordCount 统计单词数量
func (su *StringUtils) WordCount(s string) int {
	words := strings.Fields(s)
	return len(words)
}

// 演示一些可能被 linter 检查的问题
func demonstrateCommonIssues() {
	fmt.Println("=== 常见代码质量问题演示 ===")
	
	// 1. 未使用的变量（会被 go vet 检查）
	// unusedVar := "这个变量没有被使用"
	
	// 2. 错误的格式化字符串（会被 go vet 检查）
	name := "张三"
	age := 25
	fmt.Printf("姓名: %s, 年龄: %d\n", name, age) // 正确
	// fmt.Printf("姓名: %d, 年龄: %s\n", name, age) // 错误：类型不匹配
	
	// 3. 可能的空指针解引用
	var user *User
	if user != nil {
		fmt.Printf("用户: %s\n", user.Name)
	}
	
	// 4. 错误处理
	if result, err := strconv.Atoi("123"); err != nil {
		log.Printf("转换失败: %v", err)
	} else {
		fmt.Printf("转换结果: %d\n", result)
	}
}

func main() {
	fmt.Println("代码质量工具演示")
	
	// 演示用户服务
	userService := NewUserService()
	
	user1 := &User{
		ID:    1,
		Name:  "张三",
		Email: "zhangsan@example.com",
	}
	
	if err := userService.AddUser(user1); err != nil {
		log.Printf("添加用户失败: %v", err)
	} else {
		fmt.Printf("用户添加成功: %+v\n", user1)
	}
	
	// 演示计算器
	calc := &Calculator{}
	fmt.Printf("加法: 10 + 5 = %d\n", calc.Add(10, 5))
	fmt.Printf("减法: 10 - 5 = %d\n", calc.Subtract(10, 5))
	
	if result, err := calc.Divide(10, 2); err != nil {
		log.Printf("除法失败: %v", err)
	} else {
		fmt.Printf("除法: 10 / 2 = %.2f\n", result)
	}
	
	// 演示字符串工具
	stringUtils := &StringUtils{}
	text := "hello"
	fmt.Printf("原字符串: %s\n", text)
	fmt.Printf("反转后: %s\n", stringUtils.Reverse(text))
	fmt.Printf("是否回文: %t\n", stringUtils.IsPalindrome(text))
	
	// 演示数字处理
	numbers := []int{1, 2, 3, 4, 5}
	if sum, avg, err := ProcessNumbers(numbers); err != nil {
		log.Printf("处理数字失败: %v", err)
	} else {
		fmt.Printf("数字总和: %.0f, 平均值: %.2f\n", sum, avg)
	}
	
	// 演示常见问题
	demonstrateCommonIssues()
	
	fmt.Println("\n代码质量检查命令:")
	fmt.Println("  go fmt ./...           # 格式化代码")
	fmt.Println("  go vet ./...           # 静态分析")
	fmt.Println("  golangci-lint run      # 综合检查")
	fmt.Println("  go test -race ./...    # 竞态检测")
}
