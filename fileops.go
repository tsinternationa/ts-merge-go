package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// 合并文件实现 - 高效处理大文件
func (a *App) mergeFiles_impl(inputFiles []string, outputFile string, dedup bool) error {
	output, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %v", err)
	}
	defer output.Close()

	writer := bufio.NewWriter(output)
	defer writer.Flush()

	if dedup {
		return a.mergeWithDedup(inputFiles, writer)
	} else {
		return a.mergeSimple(inputFiles, writer)
	}
}

// 简单合并（不去重）- 极快速度，最小内存占用
func (a *App) mergeSimple(inputFiles []string, writer *bufio.Writer) error {
	totalFiles := len(inputFiles)
	
	for i, inputFile := range inputFiles {
		// 更新进度
		progress := float64(i) / float64(totalFiles)
		a.mergeProgress.SetValue(progress)
		a.mergeStatus.SetText(fmt.Sprintf("处理文件 %d/%d: %s", i+1, totalFiles, filepath.Base(inputFile)))

		file, err := os.Open(inputFile)
		if err != nil {
			fmt.Printf("警告: 无法打开文件 %s: %v\n", inputFile, err)
			continue
		}

		// 使用更大的缓冲区提高性能 - 优化大文件处理
		scanner := bufio.NewScanner(file)
		buf := make([]byte, 0, 128*1024) // 增加到128KB
		scanner.Buffer(buf, 2*1024*1024) // 增加到2MB最大行长度

		lineCount := 0
		batchLines := make([]string, 0, 1000) // 批量写入优化
		
		for scanner.Scan() {
			line := scanner.Text()
			batchLines = append(batchLines, line)
			lineCount++
			
			// 批量写入，减少系统调用
			if len(batchLines) >= 1000 {
				for _, batchLine := range batchLines {
					writer.WriteString(batchLine + "\n")
				}
				batchLines = batchLines[:0] // 重置切片
				writer.Flush()
			}
		}
		
		// 写入剩余的行
		for _, batchLine := range batchLines {
			writer.WriteString(batchLine + "\n")
		}
		writer.Flush()

		if err := scanner.Err(); err != nil {
			fmt.Printf("警告: 读取文件 %s 时出错: %v\n", inputFile, err)
		}

		file.Close()
		fmt.Printf("完成文件 %s，共处理 %d 行\n", filepath.Base(inputFile), lineCount)
	}

	return nil
}

// 带去重的合并 - 使用高效的map实现
func (a *App) mergeWithDedup(inputFiles []string, writer *bufio.Writer) error {
	seen := make(map[string]bool)
	totalFiles := len(inputFiles)
	totalLines := 0
	uniqueLines := 0

	for i, inputFile := range inputFiles {
		// 更新进度
		progress := float64(i) / float64(totalFiles)
		a.mergeProgress.SetValue(progress)
		a.mergeStatus.SetText(fmt.Sprintf("处理文件 %d/%d: %s", i+1, totalFiles, filepath.Base(inputFile)))

		file, err := os.Open(inputFile)
		if err != nil {
			fmt.Printf("警告: 无法打开文件 %s: %v\n", inputFile, err)
			continue
		}

		scanner := bufio.NewScanner(file)
        buf := make([]byte, 0, 128*1024) // 增加到128KB
        scanner.Buffer(buf, 2*1024*1024) // 增加到2MB最大行长度

		fileLines := 0
		fileUnique := 0

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			totalLines++
			fileLines++

			if line != "" && !seen[line] {
				seen[line] = true
				_, err := writer.WriteString(line + "\n")
				if err != nil {
					file.Close()
					return fmt.Errorf("写入失败: %v", err)
				}
				uniqueLines++
				fileUnique++
			}

			// 每10万行刷新一次
			if fileLines%100000 == 0 {
				writer.Flush()
				fmt.Printf("  已处理 %d 行，唯一行 %d\n", fileLines, fileUnique)
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("警告: 读取文件 %s 时出错: %v\n", inputFile, err)
		}

		file.Close()
		fmt.Printf("完成文件 %s: 总行数 %d，唯一行 %d\n", filepath.Base(inputFile), fileLines, fileUnique)
	}

	fmt.Printf("合并完成: 总行数 %d，唯一行 %d，去重率 %.2f%%\n", 
		totalLines, uniqueLines, float64(totalLines-uniqueLines)/float64(totalLines)*100)

	return nil
}

// 拆分文件实现
func (a *App) splitFile_impl(inputFile string, parts int, dedup bool) ([]string, error) {
	// 第一步：计算总行数
	a.splitStatus.SetText("正在计算文件行数...")
	totalLines, err := a.countLines(inputFile)
	if err != nil {
		return nil, fmt.Errorf("计算行数失败: %v", err)
	}

	fmt.Printf("文件总行数: %d\n", totalLines)

	// 创建输出文件
	outputFiles := make([]string, parts)
	writers := make([]*bufio.Writer, parts)
	files := make([]*os.File, parts)

	timestamp := time.Now().Unix()
	for i := 0; i < parts; i++ {
		filename := fmt.Sprintf("split_%d_%d.txt", i+1, timestamp)
		outputFiles[i] = filename
		
		file, err := os.Create(filename)
		if err != nil {
			// 清理已创建的文件
			for j := 0; j < i; j++ {
				files[j].Close()
				os.Remove(outputFiles[j])
			}
			return nil, fmt.Errorf("创建输出文件失败: %v", err)
		}
		files[i] = file
		writers[i] = bufio.NewWriter(file)
	}

	defer func() {
		for i := 0; i < parts; i++ {
			writers[i].Flush()
			files[i].Close()
		}
	}()

	// 开始拆分
	a.splitStatus.SetText("正在拆分文件...")
	
	if dedup {
		err = a.splitWithDedup(inputFile, writers, totalLines)
	} else {
		err = a.splitSimple(inputFile, writers, totalLines)
	}

	if err != nil {
		// 清理文件
		for i := 0; i < parts; i++ {
			writers[i].Flush()
			files[i].Close()
			os.Remove(outputFiles[i])
		}
		return nil, err
	}

	return outputFiles, nil
}

// 简单拆分（不去重）
func (a *App) splitSimple(inputFile string, writers []*bufio.Writer, totalLines int) error {
	file, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("打开输入文件失败: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
    buf := make([]byte, 0, 128*1024) // 增加到128KB
    scanner.Buffer(buf, 2*1024*1024) // 增加到2MB最大行长度

	linesPerPart := totalLines / len(writers)
	currentPart := 0
	linesInCurrentPart := 0
	processedLines := 0

	for scanner.Scan() {
		line := scanner.Text()
		
		// 写入当前部分
		writers[currentPart].WriteString(line + "\n")
		linesInCurrentPart++
		processedLines++

		// 检查是否需要切换到下一部分
		if linesInCurrentPart >= linesPerPart && currentPart < len(writers)-1 {
			currentPart++
			linesInCurrentPart = 0
		}

		// 更新进度
		if processedLines%10000 == 0 {
			progress := float64(processedLines) / float64(totalLines)
			a.splitProgress.SetValue(progress)
		}
	}

	return scanner.Err()
}

// 带去重的拆分
func (a *App) splitWithDedup(inputFile string, writers []*bufio.Writer, totalLines int) error {
	file, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("打开输入文件失败: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
    buf := make([]byte, 0, 128*1024) // 增加到128KB
    scanner.Buffer(buf, 2*1024*1024) // 增加到2MB最大行长度

	seen := make(map[string]bool)
	uniqueLines := 0
	processedLines := 0
	
	// 先计算唯一行数
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !seen[line] {
			seen[line] = true
			uniqueLines++
		}
		processedLines++
		
		if processedLines%50000 == 0 {
			progress := float64(processedLines) / float64(totalLines) * 0.5 // 前50%用于去重计算
			a.splitProgress.SetValue(progress)
		}
	}

	// 重新打开文件进行拆分
	file.Close()
	file, err = os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("重新打开输入文件失败: %v", err)
	}
	defer file.Close()

	scanner = bufio.NewScanner(file)
    scanner.Buffer(buf, 2*1024*1024) // 增加到2MB最大行长度

	linesPerPart := uniqueLines / len(writers)
	currentPart := 0
	linesInCurrentPart := 0
	processedLines = 0
	writtenLines := 0

	seen = make(map[string]bool) // 重置去重map

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		processedLines++

		if line != "" && !seen[line] {
			seen[line] = true
			
			// 写入当前部分
			writers[currentPart].WriteString(line + "\n")
			linesInCurrentPart++
			writtenLines++

			// 检查是否需要切换到下一部分
			if linesInCurrentPart >= linesPerPart && currentPart < len(writers)-1 {
				currentPart++
				linesInCurrentPart = 0
			}
		}

		// 更新进度
		if processedLines%10000 == 0 {
			progress := 0.5 + float64(processedLines)/float64(totalLines)*0.5 // 后50%用于写入
			a.splitProgress.SetValue(progress)
		}
	}

	fmt.Printf("拆分完成: 原始行数 %d，唯一行数 %d\n", totalLines, writtenLines)
	return scanner.Err()
}

// 过滤文件实现
func (a *App) filterFile_impl(inputFile string, numbers []int, outputFile string) error {
	// 创建数字集合以提高查找效率
	numberSet := make(map[int]bool)
	for _, num := range numbers {
		numberSet[num] = true
	}

	// 计算总行数
	a.filterStatus.SetText("正在计算文件行数...")
	totalLines, err := a.countLines(inputFile)
	if err != nil {
		return fmt.Errorf("计算行数失败: %v", err)
	}

	// 开始过滤
	a.filterStatus.SetText("正在过滤文件...")
	
	input, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("打开输入文件失败: %v", err)
	}
	defer input.Close()

	output, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %v", err)
	}
	defer output.Close()

	scanner := bufio.NewScanner(input)
    writer := bufio.NewWriter(output)
    defer writer.Flush()

    buf := make([]byte, 0, 128*1024) // 增加到128KB
    scanner.Buffer(buf, 2*1024*1024) // 增加到2MB最大行长度

	processedLines := 0
	filteredCount := 0
	invalidLines := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		processedLines++

		if line != "" {
			// 安全的数字转换，避免崩溃
			if num, err := strconv.Atoi(line); err == nil {
				if numberSet[num] {
					writer.WriteString(line + "\n")
					filteredCount++
				}
			} else {
				// 记录无效行但不中断处理
				invalidLines++
				if invalidLines <= 10 { // 只显示前10个无效行的警告
					fmt.Printf("⚠️ 第%d行包含非数字内容，已跳过: %s\n", processedLines, line)
				}
			}
		}

		// 更新进度
		if processedLines%10000 == 0 {
			progress := float64(processedLines) / float64(totalLines)
			a.filterProgress.SetValue(progress)
		}
	}

	if invalidLines > 10 {
		fmt.Printf("⚠️ 总共跳过了%d行非数字内容\n", invalidLines)
	}

	fmt.Printf("过滤完成: 处理 %d 行，匹配 %d 行，跳过无效行 %d 行\n", processedLines, filteredCount, invalidLines)
	return scanner.Err()
}

// 计算文件行数
func (a *App) countLines(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
    buf := make([]byte, 0, 128*1024) // 增加到128KB
    scanner.Buffer(buf, 2*1024*1024) // 增加到2MB最大行长度

	lines := 0
	for scanner.Scan() {
		lines++
	}

	return lines, scanner.Err()
}

// 复制文件
func (a *App) copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
