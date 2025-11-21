package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	//nativeDialog "github.com/sqweek/dialog"
)

// 创建文件拆分标签页
func (a *App) createSplitTab() *fyne.Container {
	// 创建拖拽区域用于文件拆分
	splitDropArea := a.createSplitDropArea()

	// 文件选择
	a.splitFileLabel = widget.NewLabel("未选择文件")
	selectFileBtn := widget.NewButtonWithIcon("📁 选择文件", nil, func() {
		// 使用Windows原生文件选择对话框
		// file, err := nativeDialog.File().Filter("文本文件", "txt").Title("选择要拆分的文件").Load()
		file, err := a.selectFileAndUpload("文本文件", "txt", "选择要拆分的文件")
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
				fmt.Printf("❌ 拆分文件验证失败: %s - %v\n", filepath.Base(file), err)
				return
			}
			a.splitFile = file
			a.splitFileLabel.SetText(filepath.Base(file))
			fmt.Printf("✅ 选择拆分文件: %s\n", filepath.Base(file))
		}
	})

	// 拆分参数
	a.splitParts = widget.NewEntry()
	a.splitParts.SetPlaceHolder("输入拆分份数，如：3")

	a.splitDedup = widget.NewCheck("🔄 去除重复行", nil)

	splitBtn := widget.NewButtonWithIcon("✂️ 开始拆分", nil, func() {
		if a.splitFile == "" {
			dialog.ShowInformation("提示", "请先选择要拆分的文件", a.window)
			return
		}
		a.startSplit()
	})
	splitBtn.Importance = widget.HighImportance

	// 进度区域
	a.splitProgress = widget.NewProgressBar()
	a.splitStatus = widget.NewLabel("📋 就绪")
	a.splitStatus.TextStyle = fyne.TextStyle{Italic: true}

	// 主布局
	topSection := container.NewVBox(
		widget.NewRichTextFromMarkdown("## ✂️ 文件拆分\n拖拽文件到下方区域或点击选择文件按钮"),
		container.NewPadded(splitDropArea),
	)

	middleSection := container.NewVBox(
		widget.NewSeparator(),
		widget.NewLabel("📄 选择的文件:"),
		a.splitFileLabel,
		selectFileBtn,
		widget.NewSeparator(),
		widget.NewLabel("⚙️ 拆分设置:"),
		container.NewGridWithColumns(2,
			widget.NewLabel("拆分份数:"),
			a.splitParts,
		),
		a.splitDedup,
	)

	bottomSection := container.NewVBox(
		widget.NewSeparator(),
		container.NewHBox(widget.NewLabel(""), splitBtn),
		widget.NewSeparator(),
		widget.NewLabel("📊 进度状态:"),
		a.splitProgress,
		a.splitStatus,
	)

	return container.NewVBox(
		topSection,
		middleSection,
		bottomSection,
	)
}

// 创建拆分专用拖拽区域
func (a *App) createSplitDropArea() *fyne.Container {
	// 创建拖拽区域的内容
	dropIcon := widget.NewLabel("✂️")
	dropIcon.Alignment = fyne.TextAlignCenter
	dropIcon.TextStyle = fyne.TextStyle{Bold: true}

	dropLabel := widget.NewLabel("拖拽文件到此处或点击选择")
	dropLabel.Alignment = fyne.TextAlignCenter
	dropLabel.TextStyle = fyne.TextStyle{Bold: true}

	dropHint := widget.NewLabel("选择要拆分的单个文件")
	dropHint.Alignment = fyne.TextAlignCenter
	dropHint.TextStyle = fyne.TextStyle{Italic: true}

	dropContent := container.NewVBox(
		dropIcon,
		dropLabel,
		dropHint,
	)

	// 创建一个可点击和拖拽的按钮
	dropButton := widget.NewButton("", func() {
		// 使用原生Windows文件选择对话框
		// file, err := nativeDialog.File().Filter("文本文件", "txt").Title("选择要拆分的文件").Load()
		file, err := a.selectFileAndUpload("文本文件", "txt", "选择要拆分的文件")
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
				fmt.Printf("❌ 拆分文件验证失败: %s - %v\n", filepath.Base(file), err)
				return
			}
			a.splitFile = file
			a.splitFileLabel.SetText(filepath.Base(file))
			fmt.Printf("✅ 选择拆分文件: %s\n", filepath.Base(file))
		}
	})

	// 设置按钮样式
	dropButton.Resize(fyne.NewSize(500, 120))
	dropButton.Importance = widget.LowImportance

	// 创建叠加容器
	overlayContainer := container.NewStack(dropButton, dropContent)

	return container.NewPadded(overlayContainer)
}

// 开始拆分文件
func (a *App) startSplit() {
	if a.splitFile == "" {
		return
	}

	parts, err := strconv.Atoi(a.splitParts.Text)
	if err != nil || parts <= 0 {
		dialog.ShowError(fmt.Errorf("请输入有效的拆分份数"), a.window)
		return
	}

	go func() {
		a.splitStatus.SetText("🔄 正在拆分文件...")
		a.splitProgress.SetValue(0)

		err := a.performSplit(parts)
		if err != nil {
			a.splitStatus.SetText("❌ 拆分失败: " + err.Error())
			dialog.ShowError(err, a.window)
		} else {
			a.splitStatus.SetText("✅ 拆分完成")
			dialog.ShowInformation("完成", fmt.Sprintf("文件拆分成功！\n已拆分为 %d 个文件", parts), a.window)
		}
		a.splitProgress.SetValue(1.0)
	}()
}

// 执行拆分操作
func (a *App) performSplit(parts int) error {
	file, err := os.Open(a.splitFile)
	if err != nil {
		return fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	// 读取所有行
	var lines []string
	uniqueLines := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	// 设置更大的缓冲区以处理长行，避免 "token too long" 错误
	buf := make([]byte, 0, 128*1024) // 128KB初始缓冲区
	scanner.Buffer(buf, 2*1024*1024) // 2MB最大行长度

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if a.splitDedup.Checked {
			if !uniqueLines[line] {
				uniqueLines[line] = true
				lines = append(lines, line)
			}
		} else {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取文件失败: %v", err)
	}

	totalLines := len(lines)
	if totalLines == 0 {
		return fmt.Errorf("文件为空或没有有效内容")
	}

	linesPerPart := totalLines / parts
	remainder := totalLines % parts

	baseFileName := strings.TrimSuffix(a.splitFile, filepath.Ext(a.splitFile))

	for i := 0; i < parts; i++ {
		a.splitProgress.SetValue(float64(i) / float64(parts))

		start := i * linesPerPart
		end := start + linesPerPart
		if i < remainder {
			end++
		}
		if i == parts-1 {
			end = totalLines
		}

		outputPath := fmt.Sprintf("%s_part%d.txt", baseFileName, i+1)
		outputFile, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("创建输出文件失败: %v", err)
		}

		writer := bufio.NewWriter(outputFile)
		for j := start; j < end && j < len(lines); j++ {
			writer.WriteString(lines[j] + "\n")
		}
		writer.Flush()
		outputFile.Close()
	}

	return nil
}
