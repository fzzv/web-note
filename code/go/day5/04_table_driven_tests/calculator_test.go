package main

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCalculator_Add 测试加法运算
func TestCalculator_Add(t *testing.T) {
	calc := &Calculator{}
	
	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"正数相加", 2.5, 3.5, 6.0},
		{"负数相加", -1.5, -2.5, -4.0},
		{"正负数相加", 5.0, -3.0, 2.0},
		{"零值相加", 0.0, 5.0, 5.0},
		{"小数相加", 0.1, 0.2, 0.3},
		{"大数相加", 1000000.0, 2000000.0, 3000000.0},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.Add(tt.a, tt.b)
			assert.InDelta(t, tt.expected, result, 0.0001, "加法结果不正确")
		})
	}
}

// TestCalculator_Divide 测试除法运算
func TestCalculator_Divide(t *testing.T) {
	calc := &Calculator{}
	
	tests := []struct {
		name      string
		a, b      float64
		expected  float64
		expectErr bool
	}{
		{"正常除法", 10.0, 2.0, 5.0, false},
		{"小数除法", 7.5, 2.5, 3.0, false},
		{"除以1", 5.0, 1.0, 5.0, false},
		{"0除以数", 0.0, 5.0, 0.0, false},
		{"除以0", 5.0, 0.0, 0.0, true},
		{"负数除法", -10.0, 2.0, -5.0, false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calc.Divide(tt.a, tt.b)
			
			if tt.expectErr {
				assert.Error(t, err, "应该返回错误")
				assert.Equal(t, 0.0, result, "错误时结果应为0")
			} else {
				assert.NoError(t, err, "不应该返回错误")
				assert.InDelta(t, tt.expected, result, 0.0001, "除法结果不正确")
			}
		})
	}
}

// TestCalculator_Sqrt 测试平方根运算
func TestCalculator_Sqrt(t *testing.T) {
	calc := &Calculator{}
	
	tests := []struct {
		name      string
		input     float64
		expected  float64
		expectErr bool
	}{
		{"正数平方根", 9.0, 3.0, false},
		{"0的平方根", 0.0, 0.0, false},
		{"1的平方根", 1.0, 1.0, false},
		{"小数平方根", 0.25, 0.5, false},
		{"负数平方根", -4.0, 0.0, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calc.Sqrt(tt.input)
			
			if tt.expectErr {
				assert.Error(t, err, "应该返回错误")
			} else {
				assert.NoError(t, err, "不应该返回错误")
				assert.InDelta(t, tt.expected, result, 0.0001, "平方根结果不正确")
			}
		})
	}
}

// TestStringProcessor_Reverse 测试字符串反转
func TestStringProcessor_Reverse(t *testing.T) {
	sp := &StringProcessor{}
	
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"英文字符串", "hello", "olleh"},
		{"中文字符串", "你好", "好你"},
		{"空字符串", "", ""},
		{"单个字符", "a", "a"},
		{"数字字符串", "12345", "54321"},
		{"混合字符串", "hello世界", "界世olleh"},
		{"回文字符串", "aba", "aba"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sp.Reverse(tt.input)
			assert.Equal(t, tt.expected, result, "字符串反转结果不正确")
		})
	}
}

// TestStringProcessor_IsPalindrome 测试回文检查
func TestStringProcessor_IsPalindrome(t *testing.T) {
	sp := &StringProcessor{}
	
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"简单回文", "aba", true},
		{"非回文", "hello", false},
		{"空字符串", "", true},
		{"单个字符", "a", true},
		{"忽略大小写", "Aba", true},
		{"带空格的回文", "a b a", true},
		{"数字回文", "12321", true},
		{"长回文", "racecar", true},
		{"非回文数字", "12345", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sp.IsPalindrome(tt.input)
			assert.Equal(t, tt.expected, result, "回文检查结果不正确")
		})
	}
}

// TestStringProcessor_ValidateEmail 测试邮箱验证
func TestStringProcessor_ValidateEmail(t *testing.T) {
	sp := &StringProcessor{}
	
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{"有效邮箱", "user@example.com", true},
		{"有效邮箱2", "test.email@domain.org", true},
		{"空字符串", "", false},
		{"缺少@", "userexample.com", false},
		{"缺少域名", "user@", false},
		{"缺少用户名", "@example.com", false},
		{"缺少点", "user@example", false},
		{"多个@", "user@@example.com", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sp.ValidateEmail(tt.email)
			assert.Equal(t, tt.expected, result, "邮箱验证结果不正确")
		})
	}
}

// TestSliceUtils_Sum 测试切片求和
func TestSliceUtils_Sum(t *testing.T) {
	su := &SliceUtils{}
	
	tests := []struct {
		name     string
		numbers  []int
		expected int
	}{
		{"正数求和", []int{1, 2, 3, 4, 5}, 15},
		{"包含负数", []int{-1, 2, -3, 4}, 2},
		{"全为负数", []int{-1, -2, -3}, -6},
		{"包含零", []int{0, 1, 2}, 3},
		{"空切片", []int{}, 0},
		{"单个元素", []int{42}, 42},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := su.Sum(tt.numbers)
			assert.Equal(t, tt.expected, result, "求和结果不正确")
		})
	}
}

// TestSliceUtils_Average 测试平均值计算
func TestSliceUtils_Average(t *testing.T) {
	su := &SliceUtils{}
	
	tests := []struct {
		name      string
		numbers   []int
		expected  float64
		expectErr bool
	}{
		{"正数平均值", []int{1, 2, 3, 4, 5}, 3.0, false},
		{"包含负数", []int{-2, 0, 2}, 0.0, false},
		{"单个元素", []int{10}, 10.0, false},
		{"空切片", []int{}, 0.0, true},
		{"偶数个元素", []int{1, 2, 3, 4}, 2.5, false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := su.Average(tt.numbers)
			
			if tt.expectErr {
				assert.Error(t, err, "应该返回错误")
			} else {
				assert.NoError(t, err, "不应该返回错误")
				assert.InDelta(t, tt.expected, result, 0.0001, "平均值计算不正确")
			}
		})
	}
}

// TestSliceUtils_MaxMin 测试最大值和最小值
func TestSliceUtils_MaxMin(t *testing.T) {
	su := &SliceUtils{}
	
	tests := []struct {
		name      string
		numbers   []int
		expectedMax int
		expectedMin int
		expectErr bool
	}{
		{"正数数组", []int{1, 5, 3, 9, 2}, 9, 1, false},
		{"包含负数", []int{-5, -1, -10, 0}, 0, -10, false},
		{"单个元素", []int{42}, 42, 42, false},
		{"相同元素", []int{5, 5, 5}, 5, 5, false},
		{"空切片", []int{}, 0, 0, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			max, maxErr := su.Max(tt.numbers)
			min, minErr := su.Min(tt.numbers)
			
			if tt.expectErr {
				assert.Error(t, maxErr, "Max应该返回错误")
				assert.Error(t, minErr, "Min应该返回错误")
			} else {
				assert.NoError(t, maxErr, "Max不应该返回错误")
				assert.NoError(t, minErr, "Min不应该返回错误")
				assert.Equal(t, tt.expectedMax, max, "最大值不正确")
				assert.Equal(t, tt.expectedMin, min, "最小值不正确")
			}
		})
	}
}

// BenchmarkCalculator_Add 基准测试：加法运算
func BenchmarkCalculator_Add(b *testing.B) {
	calc := &Calculator{}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.Add(float64(i), float64(i+1))
	}
}

// BenchmarkStringProcessor_Reverse 基准测试：字符串反转
func BenchmarkStringProcessor_Reverse(b *testing.B) {
	sp := &StringProcessor{}
	testString := "这是一个用于基准测试的字符串"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sp.Reverse(testString)
	}
}

// ExampleCalculator_Add 示例测试：加法运算
func ExampleCalculator_Add() {
	calc := &Calculator{}
	result := calc.Add(2.5, 3.5)
	println(result)
	// Output: 6
}

// TestMain 测试主函数，可以进行测试前后的设置和清理
func TestMain(m *testing.M) {
	// 测试前的设置
	println("开始运行测试...")
	
	// 运行测试
	code := m.Run()
	
	// 测试后的清理
	println("测试运行完成")
	
	// 退出
	os.Exit(code)
}
