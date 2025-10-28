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

// 创建文件合并标签页
func (a *App) createMergeTab() *fyne.Container {
	// 创建拖拽区域
	dropArea := a.createDropArea()
	
	// 文件列表 - 修复容器结构问题
	a.mergeList = widget.NewList(
		func() int { return len(a.mergeFiles) },
		func() fyne.CanvasObject {
			fileName := widget.NewLabel("")
			fileName.TextStyle = fyne.TextStyle{}
			removeBtn := widget.NewButton("×", nil)
			removeBtn.Resize(fyne.NewSize(30, 30))
			return container.NewHBox(fileName, removeBtn)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id >= len(a.mergeFiles) {
				return
			}
			container := obj.(*fyne.Container)
			fileName := container.Objects[0].(*widget.Label)
			removeBtn := container.Objects[1].(*widget.Button)
			
			fileName.SetText(filepath.Base(a.mergeFiles[id]))
			removeBtn.OnTapped = func() {
				a.removeFile(id)
			}
		},
	)
	
	// 操作按钮组
	clearBtn := widget.NewButtonWithIcon("🗑️ 清空列表", nil, func() {
		a.mergeFiles = []string{}
		a.mergeList.Refresh()
	})
	
	// 选项区域
	a.mergeDedup = widget.NewCheck("🔄 去除重复行", nil)
	mergeBtn := widget.NewButtonWithIcon("🚀 开始合并", nil, func() {
		if len(a.mergeFiles) == 0 {
			dialog.ShowInformation("提示", "请先选择要合并的文件", a.window)
			return
		}
		a.startMerge()
	})
	mergeBtn.Importance = widget.HighImportance
	
	// 进度区域
	a.mergeProgress = widget.NewProgressBar()
	a.mergeStatus = widget.NewLabel("📋 就绪")
	a.mergeStatus.TextStyle = fyne.TextStyle{Italic: true}
	
	// 顶部区域
	topSection := container.NewVBox(
		widget.NewRichTextFromMarkdown("## 📁 文件合并\n拖拽文件到下方区域或点击选择文件按钮"),
		container.NewPadded(dropArea),
	)
	
	// 左侧：文件列表区域（占更多空间）
	fileListContainer := container.NewScroll(a.mergeList)
	
	// 创建一个边框容器来包装文件列表，让它填满可用空间
	listWithBorder := container.NewBorder(
		widget.NewLabel("📋 已选择的文件:"), // 顶部
		clearBtn,                          // 底部
		nil,                               // 左侧
		nil,                               // 右侧
		fileListContainer,                 // 中心内容，会自动扩展
	)
	
	leftSection := listWithBorder
	
	// 右侧：控制区域
	rightSection := container.NewVBox(
		widget.NewLabel("⚙️ 合并选项:"),
		a.mergeDedup,
		widget.NewSeparator(),
		mergeBtn,
		widget.NewSeparator(),
		widget.NewLabel("📊 进度状态:"),
		a.mergeProgress,
		a.mergeStatus,
	)
	
	// 主要内容区域：左右布局
	mainSection := container.NewHSplit(leftSection, rightSection)
	mainSection.SetOffset(0.7) // 左边占70%，右边占30%
	
	return container.NewVBox(
		topSection,
		widget.NewSeparator(),
		mainSection,
	)
}

// 创建拖拽区域
func (a *App) createDropArea() *fyne.Container {
	// 创建一个可点击的按钮作为拖拽区域
	var dropButton *widget.Button
	dropButton = widget.NewButton("", func() {
		// 使用原生Windows文件选择对话框
		file, err := nativeDialog.File().Filter("文本文件", "txt").Title("选择要合并的文件").Load()
		if err != nil {
			if err.Error() != "Cancelled" {
				dialog.ShowError(err, a.window)
			}
			return
		}
		
		if file != "" {
			a.addFile(file)
			fmt.Printf("✅ 选择文件: %s\n", filepath.Base(file))
			
			// 询问是否继续添加更多文件
			dialog.ShowConfirm("继续添加文件？", 
				"文件已添加到合并列表！\n\n是否继续选择更多文件？\n(选择'是'可以继续添加文件)", 
				func(continue_adding bool) {
					if continue_adding {
						// 递归调用继续选择
						dropButton.OnTapped()
					}
				}, a.window)
		}
	})
	
	// 创建拖拽区域的内容
	dropIcon := widget.NewLabel("📁")
	dropIcon.Alignment = fyne.TextAlignCenter
	dropIcon.TextStyle = fyne.TextStyle{Bold: true}
	
	dropLabel := widget.NewLabel("拖拽文件到此处或点击选择")
	dropLabel.Alignment = fyne.TextAlignCenter
	dropLabel.TextStyle = fyne.TextStyle{Bold: true}
	
	dropHint := widget.NewLabel("点击可连续选择多个文件")
	dropHint.Alignment = fyne.TextAlignCenter
	dropHint.TextStyle = fyne.TextStyle{Italic: true}
	
	dropContent := container.NewVBox(
		dropIcon,
		dropLabel,
		dropHint,
	)
	
	// 创建带边框和背景的拖拽区域
	dropArea := container.NewBorder(nil, nil, nil, nil, dropContent)
	dropArea.Resize(fyne.NewSize(500, 120))
	
	// 使用按钮的样式但显示自定义内容
	dropButton.Resize(fyne.NewSize(500, 120))
	dropButton.Importance = widget.LowImportance
	
	// 创建一个叠加容器，将按钮和内容叠加
	overlayContainer := container.NewStack(dropButton, dropArea)
	
	return container.NewPadded(overlayContainer)
}

// 添加文件到列表
func (a *App) addFile(path string) {
	// 检查文件是否已存在
	for _, existing := range a.mergeFiles {
		if existing == path {
			fmt.Printf("⚠️ 文件已存在，跳过: %s\n", filepath.Base(path))
			return
		}
	}
	
	// 检查文件扩展名
	if !strings.HasSuffix(strings.ToLower(path), ".txt") {
		dialog.ShowError(fmt.Errorf("只支持 .txt 文件"), a.window)
		return
	}
	
	// 验证文件是否包含手机号格式的内容
	if err := a.validateFileContainsPhoneNumbers(path); err != nil {
		dialog.ShowError(err, a.window)
		fmt.Printf("❌ 文件验证失败: %s - %v\n", filepath.Base(path), err)
		return
	}
	
	a.mergeFiles = append(a.mergeFiles, path)
	fmt.Printf("📝 文件列表更新: 当前有 %d 个文件\n", len(a.mergeFiles))
	
	// 强制刷新列表
	if a.mergeList != nil {
		a.mergeList.Refresh()
		fmt.Printf("🔄 列表已刷新\n")
	} else {
		fmt.Printf("❌ 列表对象为空\n")
	}
}

// 从列表中移除文件
func (a *App) removeFile(index int) {
	if index >= 0 && index < len(a.mergeFiles) {
		a.mergeFiles = append(a.mergeFiles[:index], a.mergeFiles[index+1:]...)
		a.mergeList.Refresh()
	}
}

// 开始合并文件
func (a *App) startMerge() {
	if len(a.mergeFiles) == 0 {
		return
	}
	
	go func() {
		a.mergeStatus.SetText("🔄 正在合并文件...")
		a.mergeProgress.SetValue(0)
		
		// 使用 Windows 原生文件保存对话框
		outputPath, err := nativeDialog.File().
			Filter("文本文件", "txt").
			Title("选择合并后的输出文件").
			Save()
			
		if err != nil {
			a.mergeStatus.SetText("❌ 合并已取消")
			return
		}
		
		// 确保文件扩展名为 .txt
		if !strings.HasSuffix(strings.ToLower(outputPath), ".txt") {
			outputPath += ".txt"
		}
		
		err = a.performMerge(outputPath)
		if err != nil {
			a.mergeStatus.SetText("❌ 合并失败: " + err.Error())
			dialog.ShowError(err, a.window)
		} else {
			a.mergeStatus.SetText("✅ 合并完成")
			dialog.ShowInformation("完成", "文件合并成功！\n输出文件: "+filepath.Base(outputPath), a.window)
		}
		a.mergeProgress.SetValue(1.0)
	}()
}

// 执行合并操作
func (a *App) performMerge(outputPath string) error {
	// 删除可能存在的空文件
	os.Remove(outputPath)
	
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %v", err)
	}
	defer outputFile.Close()
	
	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()
	
	uniqueLines := make(map[string]bool)
	totalFiles := len(a.mergeFiles)
	linesWritten := 0
	
	for i, filePath := range a.mergeFiles {
		a.mergeProgress.SetValue(float64(i) / float64(totalFiles))
		a.mergeStatus.SetText(fmt.Sprintf("🔄 处理文件 %d/%d: %s", i+1, totalFiles, filepath.Base(filePath)))
		
		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("打开文件 %s 失败: %v", filePath, err)
		}
		
		scanner := bufio.NewScanner(file)
		// 设置更大的缓冲区以处理长行，避免 "token too long" 错误
		buf := make([]byte, 0, 128*1024) // 128KB初始缓冲区
		scanner.Buffer(buf, 2*1024*1024) // 2MB最大行长度
		
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			
			if a.mergeDedup.Checked {
				if !uniqueLines[line] {
					uniqueLines[line] = true
					_, err := writer.WriteString(line + "\n")
					if err != nil {
						file.Close()
						return fmt.Errorf("写入文件失败: %v", err)
					}
					linesWritten++
				}
			} else {
				_, err := writer.WriteString(line + "\n")
				if err != nil {
					file.Close()
					return fmt.Errorf("写入文件失败: %v", err)
				}
				linesWritten++
			}
		}
		
		file.Close()
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("读取文件 %s 失败: %v", filePath, err)
		}
	}
	
	// 强制刷新缓冲区
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("刷新缓冲区失败: %v", err)
	}
	
	fmt.Printf("✅ 合并完成，共写入 %d 行到文件: %s\n", linesWritten, outputPath)
	return nil
}
