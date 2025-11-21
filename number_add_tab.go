package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	nativeDialog "github.com/sqweek/dialog"
)

// 创建号码增加标签页
func (a *App) createNumberAddTab() *fyne.Container {
	// 创建拖拽区域用于号码增加
	numberAddDropArea := a.createNumberAddDropArea()

	// 文件选择
	a.numberAddFileLabel = widget.NewLabel("未选择文件")
	selectFileBtn := widget.NewButtonWithIcon("📁 选择文件", nil, func() {
		// 使用Windows原生文件选择对话框
		// file, err := nativeDialog.File().Filter("文本文件", "txt").Title("选择要进行号码增加的文件").Load()
		file, err := a.selectFileAndUpload("文本文件", "txt", "选择要进行号码增加的文件")
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
				fmt.Printf("❌ 号码增加文件验证失败: %s - %v\n", filepath.Base(file), err)
				return
			}
			a.numberAddFile = file
			a.numberAddFileLabel.SetText(filepath.Base(file))
			fmt.Printf("✅ 选择号码增加文件: %s\n", filepath.Base(file))
		}
	})

	// 位置设置
	a.numberAddPosition = widget.NewEntry()
	a.numberAddPosition.SetPlaceHolder("输入位置，如：0（在开头增加）或3（在第3位后增加）")

	// 数字设置
	a.numberAddDigit = widget.NewEntry()
	a.numberAddDigit.SetPlaceHolder("输入要增加的字符，空白则随机0-9")

	// 选项设置
	a.numberAddRemoveEmpty = widget.NewCheck("🗑️ 去除空行", nil)

	// 开始处理按钮
	processBtn := widget.NewButtonWithIcon("🔢 开始增加", nil, func() {
		if a.numberAddFile == "" {
			dialog.ShowInformation("提示", "请先选择要处理的文件", a.window)
			return
		}
		a.startNumberAdd()
	})
	processBtn.Importance = widget.HighImportance

	// 进度区域
	a.numberAddProgress = widget.NewProgressBar()
	a.numberAddStatus = widget.NewLabel("📋 就绪")
	a.numberAddStatus.TextStyle = fyne.TextStyle{Italic: true}

	// 主布局
	topSection := container.NewVBox(
		widget.NewRichTextFromMarkdown("## 🔢 号码增加\n为每行号码在指定位置增加字符（可以是任何字符）"),
		container.NewPadded(numberAddDropArea),
	)

	middleSection := container.NewVBox(
		widget.NewSeparator(),
		widget.NewLabel("📄 选择的文件:"),
		a.numberAddFileLabel,
		selectFileBtn,
		widget.NewSeparator(),
		widget.NewLabel("⚙️ 增加设置:"),
		container.NewGridWithColumns(2,
			widget.NewLabel("增加位置:"),
			a.numberAddPosition,
			widget.NewLabel("增加字符:"),
			a.numberAddDigit,
		),
		widget.NewLabel("💡 说明: 位置0表示在开头增加，其他数字表示在第几位后增加，字符空白则随机生成0-9（如位置0字符A表示在开头增加A）"),
		widget.NewSeparator(),
		widget.NewLabel("🔧 处理选项:"),
		a.numberAddRemoveEmpty,
	)

	bottomSection := container.NewVBox(
		widget.NewSeparator(),
		container.NewHBox(widget.NewLabel(""), processBtn),
		widget.NewSeparator(),
		widget.NewLabel("📊 进度状态:"),
		a.numberAddProgress,
		a.numberAddStatus,
	)

	return container.NewVBox(
		topSection,
		middleSection,
		bottomSection,
	)
}

// 创建号码增加专用拖拽区域
func (a *App) createNumberAddDropArea() *fyne.Container {
	// 创建拖拽区域的内容
	dropIcon := widget.NewLabel("🔢")
	dropIcon.Alignment = fyne.TextAlignCenter
	dropIcon.TextStyle = fyne.TextStyle{Bold: true}

	dropLabel := widget.NewLabel("拖拽文件到此处或点击选择")
	dropLabel.Alignment = fyne.TextAlignCenter
	dropLabel.TextStyle = fyne.TextStyle{Bold: true}

	dropHint := widget.NewLabel("选择包含号码的文件进行字符增加")
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
		file, err := nativeDialog.File().Filter("文本文件", "txt").Title("选择要进行号码增加的文件").Load()
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
				fmt.Printf("❌ 号码增加文件验证失败: %s - %v\n", filepath.Base(file), err)
				return
			}
			a.numberAddFile = file
			a.numberAddFileLabel.SetText(filepath.Base(file))
			fmt.Printf("✅ 选择号码增加文件: %s\n", filepath.Base(file))
		}
	})

	// 设置按钮样式
	dropButton.Resize(fyne.NewSize(500, 120))
	dropButton.Importance = widget.LowImportance

	// 创建叠加容器
	overlayContainer := container.NewStack(dropButton, dropContent)

	return container.NewPadded(overlayContainer)
}

// 开始号码增加处理
func (a *App) startNumberAdd() {
	if a.numberAddFile == "" {
		return
	}

	// 验证位置输入
	position, err := strconv.Atoi(a.numberAddPosition.Text)
	if err != nil || position < 0 {
		dialog.ShowError(fmt.Errorf("请输入有效的位置数字（大于等于0的整数）"), a.window)
		return
	}

	// 获取用户输入的字符（可为空）
	userDigit := strings.TrimSpace(a.numberAddDigit.Text)

	go func() {
		a.numberAddStatus.SetText("🔄 正在处理号码增加...")
		a.numberAddProgress.SetValue(0)

		err := a.performNumberAdd(position, userDigit)
		if err != nil {
			a.numberAddStatus.SetText("❌ 处理失败: " + err.Error())
			dialog.ShowError(err, a.window)
		} else {
			a.numberAddStatus.SetText("✅ 处理完成")
			dialog.ShowInformation("完成", "号码增加处理成功！", a.window)
		}
		a.numberAddProgress.SetValue(1.0)
	}()
}

// 执行号码增加操作（简化版，无去重功能）
func (a *App) performNumberAdd(position int, userDigit string) error {
	// 使用 Windows 原生文件保存对话框
	outputPath, err := nativeDialog.File().
		Filter("文本文件", "txt").
		Title("选择输出文件").
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

	file, err := os.Open(a.numberAddFile)
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
	processedLines := 0

	// 逐行读取并处理
	for scanner.Scan() {
		line := scanner.Text()
		totalLines++

		// 去除空行处理（如果勾选了去空选项）
		if a.numberAddRemoveEmpty.Checked && strings.TrimSpace(line) == "" {
			continue
		}

		// 如果行为空且不需要去空行，则直接写入
		if strings.TrimSpace(line) == "" {
			_, err := writer.WriteString(line + "\n")
			if err != nil {
				return fmt.Errorf("写入文件失败: %v", err)
			}
			processedLines++
			continue
		}

		// 在指定位置增加字符（用户输入或随机）
		processedLine := a.addDigitAtPosition(line, position, userDigit)

		// 直接写入处理后的行（无去重）
		_, err := writer.WriteString(processedLine + "\n")
		if err != nil {
			return fmt.Errorf("写入文件失败: %v", err)
		}
		processedLines++

		// 更新进度
		if totalLines%1000 == 0 {
			progress := float64(totalLines) / 100000.0 // 假设最大10万行
			if progress > 1.0 {
				progress = 1.0
			}
			a.numberAddProgress.SetValue(progress)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取文件失败: %v", err)
	}

	// 强制刷新缓冲区
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("刷新缓冲区失败: %v", err)
	}

	fmt.Printf("✅ 号码增加完成: 总行数 %d，处理行数 %d，输出文件: %s\n",
		totalLines, processedLines, filepath.Base(outputPath))

	return nil
}

// 在指定位置增加字符（用户输入或随机）
func (a *App) addDigitAtPosition(line string, position int, userDigit string) string {
	var charToAdd string

	// 如果用户输入了字符，直接使用用户输入；否则使用随机数字
	if userDigit != "" {
		// 直接使用用户输入的任何字符
		charToAdd = userDigit
	} else {
		// 用户未输入，生成随机数字（0-9）
		charToAdd = strconv.Itoa(rand.Intn(10))
	}

	// 如果位置为0，在开头添加
	if position == 0 {
		return charToAdd + line
	}

	// 如果位置超过字符串长度，则在末尾添加
	if position >= len(line) {
		return line + charToAdd
	}

	// 在指定位置后插入字符
	return line[:position] + charToAdd + line[position:]
}
