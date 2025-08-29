// Package main 演示文件 I/O 操作
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println("=== 文件 I/O 操作演示 ===")

	// 1. 基本文件读写
	demonstrateBasicFileOps()

	// 2. 缓冲 I/O
	demonstrateBufferedIO()

	// 3. 文件信息和目录操作
	demonstrateFileInfo()

	// 4. 文件复制
	demonstrateFileCopy()

	// 5. 逐行处理大文件
	demonstrateLineByLine()
}

// demonstrateBasicFileOps 演示基本文件操作
func demonstrateBasicFileOps() {
	fmt.Println("\n--- 基本文件操作 ---")

	// 写入文件
	content := "Hello, World!\n这是一个测试文件。\n"
	filename := "test_output.txt"

	// 一次性写入
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		log.Printf("写入文件失败: %v", err)
		return
	}
	fmt.Printf("已写入文件: %s\n", filename)

	// 一次性读取
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("读取文件失败: %v", err)
		return
	}
	fmt.Printf("文件内容:\n%s", string(data))

	// 清理
	defer os.Remove(filename)
}

// demonstrateBufferedIO 演示缓冲 I/O
func demonstrateBufferedIO() {
	fmt.Println("\n--- 缓冲 I/O ---")

	filename := "buffered_test.txt"

	// 缓冲写入
	file, err := os.Create(filename)
	if err != nil {
		log.Printf("创建文件失败: %v", err)
		return
	}
	defer file.Close()
	defer os.Remove(filename)

	writer := bufio.NewWriter(file)
	defer writer.Flush() // 确保缓冲区内容写入

	lines := []string{
		"第一行内容",
		"第二行内容",
		"第三行内容",
		"第四行内容",
	}

	for i, line := range lines {
		_, err := writer.WriteString(fmt.Sprintf("%d: %s\n", i+1, line))
		if err != nil {
			log.Printf("写入失败: %v", err)
			return
		}
	}

	// 强制刷新缓冲区
	writer.Flush()
	file.Close()

	// 缓冲读取
	file, err = os.Open(filename)
	if err != nil {
		log.Printf("打开文件失败: %v", err)
		return
	}
	defer file.Close()

	fmt.Println("逐行读取:")
	scanner := bufio.NewScanner(file)
	lineNum := 1
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("  行 %d: %s\n", lineNum, line)
		lineNum++
	}

	if err := scanner.Err(); err != nil {
		log.Printf("读取过程中出错: %v", err)
	}
}

// demonstrateFileInfo 演示文件信息和目录操作
func demonstrateFileInfo() {
	fmt.Println("\n--- 文件信息和目录操作 ---")

	// 创建测试文件
	filename := "info_test.txt"
	content := "测试文件信息"
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		log.Printf("创建测试文件失败: %v", err)
		return
	}
	defer os.Remove(filename)

	// 获取文件信息
	info, err := os.Stat(filename)
	if err != nil {
		log.Printf("获取文件信息失败: %v", err)
		return
	}

	fmt.Printf("文件名: %s\n", info.Name())
	fmt.Printf("文件大小: %d 字节\n", info.Size())
	fmt.Printf("修改时间: %s\n", info.ModTime().Format("2006-01-02 15:04:05"))
	fmt.Printf("是否为目录: %t\n", info.IsDir())
	fmt.Printf("文件权限: %s\n", info.Mode())

	// 检查文件是否存在
	if _, err := os.Stat("不存在的文件.txt"); os.IsNotExist(err) {
		fmt.Println("文件不存在检查: 正确")
	}

	// 创建目录
	dirName := "test_dir"
	err = os.Mkdir(dirName, 0755)
	if err != nil {
		log.Printf("创建目录失败: %v", err)
		return
	}
	defer os.RemoveAll(dirName)

	// 在目录中创建文件
	subFile := filepath.Join(dirName, "sub_file.txt")
	err = os.WriteFile(subFile, []byte("子目录文件"), 0644)
	if err != nil {
		log.Printf("创建子文件失败: %v", err)
		return
	}

	// 遍历目录
	fmt.Printf("\n目录 %s 的内容:\n", dirName)
	entries, err := os.ReadDir(dirName)
	if err != nil {
		log.Printf("读取目录失败: %v", err)
		return
	}

	for _, entry := range entries {
		fmt.Printf("  %s (目录: %t)\n", entry.Name(), entry.IsDir())
	}
}

// demonstrateFileCopy 演示文件复制
func demonstrateFileCopy() {
	fmt.Println("\n--- 文件复制 ---")

	// 创建源文件
	srcFile := "source.txt"
	srcContent := "这是源文件的内容\n包含多行文本\n用于测试文件复制功能"
	err := os.WriteFile(srcFile, []byte(srcContent), 0644)
	if err != nil {
		log.Printf("创建源文件失败: %v", err)
		return
	}
	defer os.Remove(srcFile)

	// 复制文件
	dstFile := "destination.txt"
	err = copyFile(srcFile, dstFile)
	if err != nil {
		log.Printf("复制文件失败: %v", err)
		return
	}
	defer os.Remove(dstFile)

	// 验证复制结果
	srcData, _ := os.ReadFile(srcFile)
	dstData, _ := os.ReadFile(dstFile)

	if string(srcData) == string(dstData) {
		fmt.Printf("文件复制成功: %s -> %s\n", srcFile, dstFile)
		fmt.Printf("文件大小: %d 字节\n", len(dstData))
	} else {
		fmt.Println("文件复制失败: 内容不匹配")
	}
}

// copyFile 复制文件
func copyFile(src, dst string) error {
	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("打开源文件失败: %w", err)
	}
	defer srcFile.Close()

	// 创建目标文件
	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer dstFile.Close()

	// 复制内容
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("复制内容失败: %w", err)
	}

	// 同步到磁盘
	err = dstFile.Sync()
	if err != nil {
		return fmt.Errorf("同步文件失败: %w", err)
	}

	return nil
}

// demonstrateLineByLine 演示逐行处理大文件
func demonstrateLineByLine() {
	fmt.Println("\n--- 逐行处理大文件 ---")

	// 创建测试文件
	filename := "large_test.txt"
	file, err := os.Create(filename)
	if err != nil {
		log.Printf("创建测试文件失败: %v", err)
		return
	}

	// 写入多行数据
	writer := bufio.NewWriter(file)
	for i := 1; i <= 1000; i++ {
		line := fmt.Sprintf("这是第 %d 行数据，包含一些测试内容\n", i)
		writer.WriteString(line)
	}
	writer.Flush()
	file.Close()
	defer os.Remove(filename)

	// 逐行处理
	err = processLargeFile(filename)
	if err != nil {
		log.Printf("处理大文件失败: %v", err)
	}
}

// processLargeFile 处理大文件
func processLargeFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	wordCount := 0
	charCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++
		wordCount += len(strings.Fields(line))
		charCount += len(line)

		// 只显示前几行和最后几行
		if lineCount <= 3 || lineCount%200 == 0 {
			fmt.Printf("  处理第 %d 行: %s\n", lineCount, 
				truncateString(line, 50))
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取文件时出错: %w", err)
	}

	fmt.Printf("\n文件统计:\n")
	fmt.Printf("  总行数: %d\n", lineCount)
	fmt.Printf("  总单词数: %d\n", wordCount)
	fmt.Printf("  总字符数: %d\n", charCount)

	return nil
}

// truncateString 截断字符串
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
