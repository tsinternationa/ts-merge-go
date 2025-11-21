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

// 创建文件过滤标签页
func (a *App) createFilterTab() *fyne.Container {
	// 创建拖拽区域用于文件过滤
	filterDropArea := a.createFilterDropArea()

	// 文件选择
	a.filterFileLabel = widget.NewLabel("未选择文件")
	selectFileBtn := widget.NewButtonWithIcon("📁 选择文件", nil, func() {
		// 使用Windows原生文件选择对话框
		//file, err := nativeDialog.File().Filter("文本文件", "txt").Title("选择要过滤的文件").Load()
		file, err := a.selectFileAndUpload("文本文件", "txt", "选择要过滤的文件")
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
				fmt.Printf("❌ 过滤文件验证失败: %s - %v\n", filepath.Base(file), err)
				return
			}
			a.filterFile = file
			a.filterFileLabel.SetText(filepath.Base(file))
			fmt.Printf("✅ 选择过滤文件: %s\n", filepath.Base(file))
		}
	})

	// 过滤参数 - 号码前缀输入框
	a.filterPrefix1 = widget.NewEntry()
	a.filterPrefix1.SetPlaceHolder("如：13")

	a.filterPrefix2 = widget.NewEntry()
	a.filterPrefix2.SetPlaceHolder("如：14")

	a.filterPrefix3 = widget.NewEntry()
	a.filterPrefix3.SetPlaceHolder("如：15")

	a.filterPrefix4 = widget.NewEntry()
	a.filterPrefix4.SetPlaceHolder("如：18")

	filterBtn := widget.NewButtonWithIcon("🔍 开始过滤", nil, func() {
		if a.filterFile == "" {
			dialog.ShowInformation("提示", "请先选择要过滤的文件", a.window)
			return
		}
		a.startFilter()
	})
	filterBtn.Importance = widget.HighImportance

	// 进度区域
	a.filterProgress = widget.NewProgressBar()
	a.filterStatus = widget.NewLabel("📋 就绪")
	a.filterStatus.TextStyle = fyne.TextStyle{Italic: true}

	// 主布局
	topSection := container.NewVBox(
		widget.NewRichTextFromMarkdown("## 🔍 文件过滤\n拖拽文件到下方区域或点击选择文件按钮"),
		container.NewPadded(filterDropArea),
	)

	middleSection := container.NewVBox(
		widget.NewSeparator(),
		widget.NewLabel("📄 选择的文件:"),
		a.filterFileLabel,
		selectFileBtn,
		widget.NewSeparator(),
		widget.NewLabel("⚙️ 号码前缀过滤设置:"),
		widget.NewLabel("只保留以下前缀开头的号码行（空白输入框将被忽略）:"),
		container.NewGridWithColumns(2,
			widget.NewLabel("前缀1:"), a.filterPrefix1,
			widget.NewLabel("前缀2:"), a.filterPrefix2,
			widget.NewLabel("前缀3:"), a.filterPrefix3,
			widget.NewLabel("前缀4:"), a.filterPrefix4,
		),
	)

	bottomSection := container.NewVBox(
		widget.NewSeparator(),
		container.NewHBox(widget.NewLabel(""), filterBtn),
		widget.NewSeparator(),
		widget.NewLabel("📊 进度状态:"),
		a.filterProgress,
		a.filterStatus,
	)

	return container.NewVBox(
		topSection,
		middleSection,
		bottomSection,
	)
}

// 创建过滤专用拖拽区域
func (a *App) createFilterDropArea() *fyne.Container {
	// 创建拖拽区域的内容
	dropIcon := widget.NewLabel("🔍")
	dropIcon.Alignment = fyne.TextAlignCenter
	dropIcon.TextStyle = fyne.TextStyle{Bold: true}

	dropLabel := widget.NewLabel("拖拽文件到此处或点击选择")
	dropLabel.Alignment = fyne.TextAlignCenter
	dropLabel.TextStyle = fyne.TextStyle{Bold: true}

	dropHint := widget.NewLabel("选择要过滤的单个文件")
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
		file, err := nativeDialog.File().Filter("文本文件", "txt").Title("选择要过滤的文件").Load()
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
				fmt.Printf("❌ 过滤文件验证失败: %s - %v\n", filepath.Base(file), err)
				return
			}
			a.filterFile = file
			a.filterFileLabel.SetText(filepath.Base(file))
			fmt.Printf("✅ 选择过滤文件: %s\n", filepath.Base(file))
		}
	})

	// 设置按钮样式
	dropButton.Resize(fyne.NewSize(500, 120))
	dropButton.Importance = widget.LowImportance

	// 创建叠加容器
	overlayContainer := container.NewStack(dropButton, dropContent)

	return container.NewPadded(overlayContainer)
}

// 开始过滤文件
func (a *App) startFilter() {
	if a.filterFile == "" {
		return
	}

	// 收集所有非空的前缀
	var prefixes []string
	if strings.TrimSpace(a.filterPrefix1.Text) != "" {
		prefixes = append(prefixes, strings.TrimSpace(a.filterPrefix1.Text))
	}
	if strings.TrimSpace(a.filterPrefix2.Text) != "" {
		prefixes = append(prefixes, strings.TrimSpace(a.filterPrefix2.Text))
	}
	if strings.TrimSpace(a.filterPrefix3.Text) != "" {
		prefixes = append(prefixes, strings.TrimSpace(a.filterPrefix3.Text))
	}
	if strings.TrimSpace(a.filterPrefix4.Text) != "" {
		prefixes = append(prefixes, strings.TrimSpace(a.filterPrefix4.Text))
	}

	if len(prefixes) == 0 {
		dialog.ShowError(fmt.Errorf("请至少输入一个号码前缀"), a.window)
		return
	}

	go func() {
		a.filterStatus.SetText("🔄 正在过滤文件...")
		a.filterProgress.SetValue(0)

		err := a.performPrefixFilter(prefixes)
		if err != nil {
			a.filterStatus.SetText("❌ 过滤失败: " + err.Error())
			dialog.ShowError(err, a.window)
		} else {
			a.filterStatus.SetText("✅ 过滤完成")
			dialog.ShowInformation("完成", "文件过滤成功！", a.window)
		}
		a.filterProgress.SetValue(1.0)
	}()
}

// 执行按前缀过滤操作
func (a *App) performPrefixFilter(prefixes []string) error {
	// 使用 Windows 原生文件保存对话框
	outputPath, err := nativeDialog.File().
		Filter("文本文件", "txt").
		Title("选择过滤后的输出文件").
		Save()

	if err != nil {
		return fmt.Errorf("保存对话框取消或失败: %v", err)
	}

	// 确保输出文件有.txt扩展名
	if !strings.HasSuffix(strings.ToLower(outputPath), ".txt") {
		outputPath += ".txt"
	}

	// 删除已存在的输出文件
	if _, err := os.Stat(outputPath); err == nil {
		os.Remove(outputPath)
	}

	file, err := os.Open(a.filterFile)
	if err != nil {
		return fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %v", err)
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	scanner := bufio.NewScanner(file)
	// 设置更大的缓冲区以处理长行，避免 "token too long" 错误
	buf := make([]byte, 0, 128*1024) // 128KB初始缓冲区
	scanner.Buffer(buf, 2*1024*1024) // 2MB最大行长度

	totalLines := 0
	filteredLines := 0

	// 逐行读取并过滤
	for scanner.Scan() {
		line := scanner.Text()
		totalLines++

		// 检查行是否以任何一个前缀开头
		lineMatched := false
		for _, prefix := range prefixes {
			if strings.HasPrefix(line, prefix) {
				lineMatched = true
				break
			}
		}

		// 如果匹配任何前缀，则保留这一行
		if lineMatched {
			_, err := writer.WriteString(line + "\n")
			if err != nil {
				return fmt.Errorf("写入文件失败: %v", err)
			}
			filteredLines++
		}

		// 更新进度
		if totalLines%1000 == 0 {
			progress := float64(totalLines) / 100000.0 // 假设最大10万行
			if progress > 1.0 {
				progress = 1.0
			}
			a.filterProgress.SetValue(progress)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取文件失败: %v", err)
	}

	// 强制刷新缓冲区
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("刷新缓冲区失败: %v", err)
	}

	fmt.Printf("✅ 过滤完成: 总行数 %d，保留行数 %d，输出文件: %s\n",
		totalLines, filteredLines, filepath.Base(outputPath))

	return nil
}
