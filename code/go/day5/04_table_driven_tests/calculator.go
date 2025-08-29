// Package main 演示表驱动测试
package main

import (
	"errors"
	"math"
	"strings"
)

// Calculator 计算器结构体
type Calculator struct{}

// Add 加法运算
func (c *Calculator) Add(a, b float64) float64 {
	return a + b
}

// Subtract 减法运算
func (c *Calculator) Subtract(a, b float64) float64 {
	return a - b
}

// Multiply 乘法运算
func (c *Calculator) Multiply(a, b float64) float64 {
	return a * b
}

// Divide 除法运算
func (c *Calculator) Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("除数不能为零")
	}
	return a / b, nil
}

// Power 幂运算
func (c *Calculator) Power(base, exponent float64) float64 {
	return math.Pow(base, exponent)
}

// Sqrt 平方根运算
func (c *Calculator) Sqrt(x float64) (float64, error) {
	if x < 0 {
		return 0, errors.New("负数不能开平方根")
	}
	return math.Sqrt(x), nil
}

// StringProcessor 字符串处理器
type StringProcessor struct{}

// Reverse 反转字符串
func (sp *StringProcessor) Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// IsPalindrome 检查是否为回文
func (sp *StringProcessor) IsPalindrome(s string) bool {
	// 转换为小写并移除空格
	cleaned := strings.ToLower(strings.ReplaceAll(s, " ", ""))
	return cleaned == sp.Reverse(cleaned)
}

// WordCount 统计单词数量
func (sp *StringProcessor) WordCount(s string) int {
	if strings.TrimSpace(s) == "" {
		return 0
	}
	words := strings.Fields(s)
	return len(words)
}

// Capitalize 首字母大写
func (sp *StringProcessor) Capitalize(s string) string {
	if s == "" {
		return s
	}
	
	words := strings.Fields(s)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, " ")
}

// ValidateEmail 验证邮箱格式
func (sp *StringProcessor) ValidateEmail(email string) bool {
	// 简单的邮箱验证
	if email == "" {
		return false
	}
	
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	
	localPart := parts[0]
	domainPart := parts[1]
	
	if localPart == "" || domainPart == "" {
		return false
	}
	
	if !strings.Contains(domainPart, ".") {
		return false
	}
	
	return true
}

// SliceUtils 切片工具
type SliceUtils struct{}

// Sum 计算整数切片的和
func (su *SliceUtils) Sum(numbers []int) int {
	sum := 0
	for _, num := range numbers {
		sum += num
	}
	return sum
}

// Average 计算平均值
func (su *SliceUtils) Average(numbers []int) (float64, error) {
	if len(numbers) == 0 {
		return 0, errors.New("空切片无法计算平均值")
	}
	
	sum := su.Sum(numbers)
	return float64(sum) / float64(len(numbers)), nil
}

// Max 找到最大值
func (su *SliceUtils) Max(numbers []int) (int, error) {
	if len(numbers) == 0 {
		return 0, errors.New("空切片无法找到最大值")
	}
	
	max := numbers[0]
	for _, num := range numbers[1:] {
		if num > max {
			max = num
		}
	}
	return max, nil
}

// Min 找到最小值
func (su *SliceUtils) Min(numbers []int) (int, error) {
	if len(numbers) == 0 {
		return 0, errors.New("空切片无法找到最小值")
	}
	
	min := numbers[0]
	for _, num := range numbers[1:] {
		if num < min {
			min = num
		}
	}
	return min, nil
}

// Contains 检查切片是否包含指定元素
func (su *SliceUtils) Contains(numbers []int, target int) bool {
	for _, num := range numbers {
		if num == target {
			return true
		}
	}
	return false
}

// Remove 移除指定元素的第一个出现
func (su *SliceUtils) Remove(numbers []int, target int) []int {
	for i, num := range numbers {
		if num == target {
			return append(numbers[:i], numbers[i+1:]...)
		}
	}
	return numbers
}

// Unique 去重
func (su *SliceUtils) Unique(numbers []int) []int {
	seen := make(map[int]bool)
	var result []int
	
	for _, num := range numbers {
		if !seen[num] {
			seen[num] = true
			result = append(result, num)
		}
	}
	
	return result
}

func main() {
	// 这个文件主要用于测试，实际的演示在测试文件中
	calc := &Calculator{}
	result := calc.Add(2, 3)
	println("2 + 3 =", result)
	
	sp := &StringProcessor{}
	reversed := sp.Reverse("hello")
	println("hello 反转后:", reversed)
	
	su := &SliceUtils{}
	numbers := []int{1, 2, 3, 4, 5}
	sum := su.Sum(numbers)
	println("数组和:", sum)
}
