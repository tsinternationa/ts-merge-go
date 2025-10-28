package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	nativeDialog "github.com/sqweek/dialog"
)

// 创建文件重复比较标签页
func (a *App) createCompareTab() *fyne.Container {
	// 左侧文件选择区域
	leftDropArea := a.createCompareDropArea(1)
	a.compareFile1Label = widget.NewLabel("未选择文件")
	selectFile1Btn := widget.NewButtonWithIcon("📁 选择文件1", nil, func() {
		file, err := nativeDialog.File().Filter("文本文件", "txt").Title("选择第一个文件").Load()
		if err != nil {
			if err.Error() != "Cancelled" {
				dialog.ShowError(err, a.window)
			}
			return
		}
		if file != "" {
			// 验证文件格式
			if err := a.validateFileContainsPhoneNumbers(file); err != nil {
				dialog.ShowError(err, a.window)
				fmt.Printf("❌ 比较文件1验证失败: %s - %v\n", filepath.Base(file), err)
				return
			}
			a.compareFile1 = file
			a.compareFile1Label.SetText(filepath.Base(file))
			fmt.Printf("✅ 选择文件1: %s\n", filepath.Base(file))
		}
	})
	
	// 右侧文件选择区域
	rightDropArea := a.createCompareDropArea(2)
	a.compareFile2Label = widget.NewLabel("未选择文件")
	selectFile2Btn := widget.NewButtonWithIcon("📁 选择文件2", nil, func() {
		file, err := nativeDialog.File().Filter("文本文件", "txt").Title("选择第二个文件").Load()
		if err != nil {
			if err.Error() != "Cancelled" {
				dialog.ShowError(err, a.window)
			}
			return
		}
		if file != "" {
			// 验证文件格式
			if err := a.validateFileContainsPhoneNumbers(file); err != nil {
				dialog.ShowError(err, a.window)
				fmt.Printf("❌ 比较文件2验证失败: %s - %v\n", filepath.Base(file), err)
				return
			}
			a.compareFile2 = file
			a.compareFile2Label.SetText(filepath.Base(file))
			fmt.Printf("✅ 选择文件2: %s\n", filepath.Base(file))
		}
	})
	
	// 开始比较按钮
	compareBtn := widget.NewButtonWithIcon("🔄 开始比较", nil, func() {
		if a.compareFile1 == "" || a.compareFile2 == "" {
			dialog.ShowInformation("提示", "请先选择两个要比较的文件", a.window)
			return
		}
		a.startCompare()
	})
	compareBtn.Importance = widget.HighImportance
	
	// 进度区域
	a.compareProgress = widget.NewProgressBar()
	a.compareStatus = widget.NewLabel("📋 就绪")
	a.compareStatus.TextStyle = fyne.TextStyle{Italic: true}
	
	// 顶部说明
	topSection := container.NewVBox(
		widget.NewRichTextFromMarkdown("## 🔄 文件重复比较\n比较两个文件，生成相同内容和不同内容的文件"),
	)
	
	// 左右文件选择区域
	leftSection := container.NewVBox(
		widget.NewLabel("📄 文件1:"),
		container.NewPadded(leftDropArea),
		a.compareFile1Label,
		selectFile1Btn,
	)
	
	rightSection := container.NewVBox(
		widget.NewLabel("📄 文件2:"),
		container.NewPadded(rightDropArea),
		a.compareFile2Label,
		selectFile2Btn,
	)
	
	// 左右布局
	middleSection := container.NewHSplit(leftSection, rightSection)
	middleSection.SetOffset(0.5) // 50-50分割
	
	// 底部控制区域
	bottomSection := container.NewVBox(
		widget.NewSeparator(),
		container.NewHBox(widget.NewLabel(""), compareBtn),
		widget.NewSeparator(),
		widget.NewLabel("📊 进度状态:"),
		a.compareProgress,
		a.compareStatus,
	)
	
	return container.NewVBox(
		topSection,
		middleSection,
		bottomSection,
	)
}

// 创建文件比较专用拖拽区域
func (a *App) createCompareDropArea(fileNum int) *fyne.Container {
	// 创建拖拽区域的内容
	dropIcon := widget.NewLabel("📁")
	dropIcon.Alignment = fyne.TextAlignCenter
	dropIcon.TextStyle = fyne.TextStyle{Bold: true}
	
	var dropLabel *widget.Label
	if fileNum == 1 {
		dropLabel = widget.NewLabel("拖拽文件1到此处")
	} else {
		dropLabel = widget.NewLabel("拖拽文件2到此处")
	}
	dropLabel.Alignment = fyne.TextAlignCenter
	dropLabel.TextStyle = fyne.TextStyle{Bold: true}
	
	dropContent := container.NewVBox(
		dropIcon,
		dropLabel,
	)
	
	// 创建带边框和背景的拖拽区域
	dropArea := container.NewBorder(nil, nil, nil, nil, dropContent)
	dropArea.Resize(fyne.NewSize(300, 100))
	
	return dropArea
}

// 开始文件比较
func (a *App) startCompare() {
	if a.compareFile1 == "" || a.compareFile2 == "" {
		return
	}
	
	go func() {
		a.compareStatus.SetText("🔄 正在比较文件...")
		a.compareProgress.SetValue(0)
		
		err := a.performCompare()
		if err != nil {
			a.compareStatus.SetText("❌ 比较失败: " + err.Error())
			dialog.ShowError(err, a.window)
		} else {
			a.compareStatus.SetText("✅ 比较完成")
			dialog.ShowInformation("完成", "文件比较成功！\n已生成相同内容和不同内容的文件", a.window)
		}
		a.compareProgress.SetValue(1.0)
	}()
}

// 执行文件比较操作
func (a *App) performCompare() error {
	// 选择输出目录
	outputDir, err := nativeDialog.Directory().Title("选择输出文件夹").Browse()
	if err != nil {
		return fmt.Errorf("选择输出目录失败: %v", err)
	}
	
	// 读取第一个文件
	file1Lines, err := a.readFileLines(a.compareFile1)
	if err != nil {
		return fmt.Errorf("读取文件1失败: %v", err)
	}
	
	// 读取第二个文件
	file2Lines, err := a.readFileLines(a.compareFile2)
	if err != nil {
		return fmt.Errorf("读取文件2失败: %v", err)
	}
	
	// 使用高效的集合算法进行比较
	file1Set := make(map[string]bool)
	file2Set := make(map[string]bool)
	sameSet := make(map[string]bool)  // 用于去重相同内容
	
	totalLines := len(file1Lines) + len(file2Lines)
	processedLines := 0
	
	// 构建文件1的集合
	for _, line := range file1Lines {
		file1Set[line] = true
		processedLines++
		if processedLines%1000 == 0 {
			progress := float64(processedLines) / float64(totalLines) * 0.3 // 30%用于构建集合
			a.compareProgress.SetValue(progress)
		}
	}
	
	// 构建文件2的集合
	for _, line := range file2Lines {
		file2Set[line] = true
		processedLines++
		if processedLines%1000 == 0 {
			progress := float64(processedLines) / float64(totalLines) * 0.3
			a.compareProgress.SetValue(progress)
		}
	}
	
	a.compareProgress.SetValue(0.3) // 集合构建完成
	
	// 找出相同和不同的内容
	var sameLines []string
	var diffLines []string
	
	// 检查文件1中的每一行
	processedLines = 0
	for _, line := range file1Lines {
		if file2Set[line] {
			// 相同内容（使用map去重，O(1)复杂度）
			if !sameSet[line] {
				sameSet[line] = true
				sameLines = append(sameLines, line)
			}
		} else {
			// 文件1独有的内容
			diffLines = append(diffLines, line)
		}
		
		processedLines++
		if processedLines%1000 == 0 {
			progress := 0.3 + float64(processedLines)/float64(len(file1Lines))*0.35 // 30%-65%
			a.compareProgress.SetValue(progress)
		}
	}
	
	a.compareProgress.SetValue(0.65) // 文件1处理完成
	
	// 检查文件2中独有的内容
	processedLines = 0
	for _, line := range file2Lines {
		if !file1Set[line] {
			diffLines = append(diffLines, line)
		}
		
		processedLines++
		if processedLines%1000 == 0 {
			progress := 0.65 + float64(processedLines)/float64(len(file2Lines))*0.25 // 65%-90%
			a.compareProgress.SetValue(progress)
		}
	}
	
	a.compareProgress.SetValue(0.9) // 比较完成，准备写入文件
	
	// 生成输出文件名
	baseFileName1 := strings.TrimSuffix(filepath.Base(a.compareFile1), filepath.Ext(a.compareFile1))
	baseFileName2 := strings.TrimSuffix(filepath.Base(a.compareFile2), filepath.Ext(a.compareFile2))
	
	sameFileName := filepath.Join(outputDir, fmt.Sprintf("%s_%s_相同内容.txt", baseFileName1, baseFileName2))
	diffFileName := filepath.Join(outputDir, fmt.Sprintf("%s_%s_不同内容.txt", baseFileName1, baseFileName2))
	
	// 写入相同内容文件
	err = a.writeLinesToFile(sameFileName, sameLines)
	if err != nil {
		return fmt.Errorf("写入相同内容文件失败: %v", err)
	}
	
	// 写入不同内容文件
	err = a.writeLinesToFile(diffFileName, diffLines)
	if err != nil {
		return fmt.Errorf("写入不同内容文件失败: %v", err)
	}
	
	fmt.Printf("✅ 比较完成:\n")
	fmt.Printf("   相同内容: %d 行 -> %s\n", len(sameLines), filepath.Base(sameFileName))
	fmt.Printf("   不同内容: %d 行 -> %s\n", len(diffLines), filepath.Base(diffFileName))
	
	return nil
}

// 读取文件所有行
func (a *App) readFileLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	var lines []string
	scanner := bufio.NewScanner(file)
	
	// 设置更大的缓冲区以处理长行，避免 "token too long" 错误
	buf := make([]byte, 0, 128*1024) // 128KB初始缓冲区
	scanner.Buffer(buf, 2*1024*1024) // 2MB最大行长度
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" { // 跳过空行
			lines = append(lines, line)
		}
	}
	
	return lines, scanner.Err()
}

// 将行写入文件
func (a *App) writeLinesToFile(filePath string, lines []string) error {
	// 删除已存在的文件
	if _, err := os.Stat(filePath); err == nil {
		os.Remove(filePath)
	}
	
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	writer := bufio.NewWriter(file)
	defer writer.Flush()
	
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	
	return writer.Flush()
}
