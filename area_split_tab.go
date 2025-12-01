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

// 创建区号拆分标签页
func (a *App) createAreaSplitTab() *fyne.Container {
	// 创建拖拽区域
	areaSplitDropArea := a.createAreaSplitDropArea()

	// 文件选择
	a.areaSplitFileLabel = widget.NewLabel("未选择文件")
	selectFileBtn := widget.NewButtonWithIcon("📁 选择文件", nil, func() {
		file, err := nativeDialog.File().Filter("文本文件", "txt").Title("选择要按区号拆分的文件").Load()
		if err != nil {
			if err.Error() != "Cancelled" {
				dialog.ShowError(err, a.window)
			}
			return
		}

		if file != "" {
			if err := a.validateFileContainsPhoneNumbers(file); err != nil {
				dialog.ShowError(err, a.window)
				fmt.Printf("❌ 区号拆分文件验证失败: %s - %v\n", filepath.Base(file), err)
				return
			}
			a.areaSplitFile = file
			a.areaSplitFileLabel.SetText(filepath.Base(file))
			fmt.Printf("✅ 选择区号拆分文件: %s\n", filepath.Base(file))
		}
	})

	// 国家选择下拉框
	a.areaSplitCountrySelect = widget.NewSelect(getAvailableCountries(), func(selected string) {
		a.areaSplitSelectedCountry = selected
		fmt.Printf("✅ 选择国家: %s\n", selected)
	})
	a.areaSplitCountrySelect.PlaceHolder = "请选择国家"

	// 开始拆分按钮
	splitBtn := widget.NewButtonWithIcon("🗺️ 开始拆分", nil, func() {
		if a.areaSplitFile == "" {
			dialog.ShowInformation("提示", "请先选择要拆分的文件", a.window)
			return
		}
		if a.areaSplitSelectedCountry == "" {
			dialog.ShowInformation("提示", "请先选择国家", a.window)
			return
		}
		a.startAreaSplit()
	})
	splitBtn.Importance = widget.HighImportance

	// 进度区域
	a.areaSplitProgress = widget.NewProgressBar()
	a.areaSplitStatus = widget.NewLabel("📋 就绪")
	a.areaSplitStatus.TextStyle = fyne.TextStyle{Italic: true}

	// 主布局
	topSection := container.NewVBox(
		widget.NewRichTextFromMarkdown("## 🗺️ 按地区拆分号码\n上传txt文件，选择国家，按地区拆分号码"),
		container.NewPadded(areaSplitDropArea),
	)

	middleSection := container.NewVBox(
		widget.NewSeparator(),
		widget.NewLabel("📄 选择的文件:"),
		a.areaSplitFileLabel,
		selectFileBtn,
		widget.NewSeparator(),
		widget.NewLabel("🌍 选择国家:"),
		a.areaSplitCountrySelect,
		widget.NewSeparator(),
		widget.NewLabel("⚙️ 拆分说明:"),
		widget.NewLabel("• 识别该国家的号码并按地区分类"),
		widget.NewLabel("• 非该国家号码 → 未知国家.txt"),
		widget.NewLabel("• 无法识别地区 → 未知地区.txt"),
		widget.NewLabel("• 已识别地区 → 对应地区.txt"),
	)

	bottomSection := container.NewVBox(
		widget.NewSeparator(),
		container.NewHBox(widget.NewLabel(""), splitBtn),
		widget.NewSeparator(),
		widget.NewLabel("📊 进度状态:"),
		a.areaSplitProgress,
		a.areaSplitStatus,
	)

	return container.NewVBox(
		topSection,
		middleSection,
		bottomSection,
	)
}

// 创建区号拆分专用拖拽区域
func (a *App) createAreaSplitDropArea() *fyne.Container {
	dropIcon := widget.NewLabel("🗺️")
	dropIcon.Alignment = fyne.TextAlignCenter
	dropIcon.TextStyle = fyne.TextStyle{Bold: true}

	dropLabel := widget.NewLabel("拖拽文件到此处或点击选择")
	dropLabel.Alignment = fyne.TextAlignCenter
	dropLabel.TextStyle = fyne.TextStyle{Bold: true}

	dropHint := widget.NewLabel("选择包含手机号的文件进行地区拆分")
	dropHint.Alignment = fyne.TextAlignCenter
	dropHint.TextStyle = fyne.TextStyle{Italic: true}

	dropContent := container.NewVBox(
		dropIcon,
		dropLabel,
		dropHint,
	)

	dropButton := widget.NewButton("", func() {
		file, err := nativeDialog.File().Filter("文本文件", "txt").Title("选择要按区号拆分的文件").Load()
		if err != nil {
			if err.Error() != "Cancelled" {
				dialog.ShowError(err, a.window)
			}
			return
		}

		if file != "" {
			if err := a.validateFileContainsPhoneNumbers(file); err != nil {
				dialog.ShowError(err, a.window)
				fmt.Printf("❌ 区号拆分文件验证失败: %s - %v\n", filepath.Base(file), err)
				return
			}
			a.areaSplitFile = file
			a.areaSplitFileLabel.SetText(filepath.Base(file))
			fmt.Printf("✅ 选择区号拆分文件: %s\n", filepath.Base(file))
		}
	})

	dropButton.Resize(fyne.NewSize(500, 120))
	dropButton.Importance = widget.LowImportance

	overlayContainer := container.NewStack(dropButton, dropContent)

	return container.NewPadded(overlayContainer)
}

// 开始按地区拆分
func (a *App) startAreaSplit() {
	if a.areaSplitFile == "" || a.areaSplitSelectedCountry == "" {
		return
	}

	go func() {
		a.areaSplitStatus.SetText("🔄 正在按地区拆分文件...")
		a.areaSplitProgress.SetValue(0)

		// 选择输出目录
		outputDir, err := nativeDialog.Directory().Title("选择拆分文件的输出文件夹").Browse()
		if err != nil {
			a.areaSplitStatus.SetText("❌ 拆分已取消")
			return
		}

		err = a.performAreaSplit(outputDir)
		if err != nil {
			a.areaSplitStatus.SetText("❌ 拆分失败: " + err.Error())
			dialog.ShowError(err, a.window)
		} else {
			a.areaSplitStatus.SetText("✅ 拆分完成")
			dialog.ShowInformation("完成", "按地区拆分成功！\n已生成各地区的独立文件", a.window)
		}
		a.areaSplitProgress.SetValue(1.0)
	}()
}

// 执行按地区拆分操作
func (a *App) performAreaSplit(outputDir string) error {
	file, err := os.Open(a.areaSplitFile)
	if err != nil {
		return fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	// 获取选中国家的配置
	countryConfig := getCountryConfig(a.areaSplitSelectedCountry)
	if countryConfig == nil {
		return fmt.Errorf("未找到国家配置: %s", a.areaSplitSelectedCountry)
	}

	// 用于存储不同分类的号码
	areaPhones := make(map[string][]string)   // 地区号码
	unknownAreaPhones := make([]string, 0)    // 未知地区号码
	unknownCountryPhones := make([]string, 0) // 未知国家号码

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 128*1024)
	scanner.Buffer(buf, 2*1024*1024)

	totalLines := 0
	processedLines := 0

	// 第一遍：计算总行数
	a.areaSplitStatus.SetText("🔄 正在计算文件行数...")
	for scanner.Scan() {
		totalLines++
	}

	// 重新打开文件进行处理
	file.Close()
	file, err = os.Open(a.areaSplitFile)
	if err != nil {
		return fmt.Errorf("重新打开文件失败: %v", err)
	}
	defer file.Close()

	scanner = bufio.NewScanner(file)
	scanner.Buffer(buf, 2*1024*1024)

	a.areaSplitStatus.SetText(fmt.Sprintf("🔄 正在识别 %s 的地区...", a.areaSplitSelectedCountry))

	// 第二遍：分类号码
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		processedLines++

		if line != "" {
			// 判断是否属于该国家
			if !countryConfig.IsCountryNumber(line) {
				unknownCountryPhones = append(unknownCountryPhones, line)
			} else {
				// 识别地区
				area := countryConfig.IdentifyArea(line)
				if area == "" {
					unknownAreaPhones = append(unknownAreaPhones, line)
				} else {
					if areaPhones[area] == nil {
						areaPhones[area] = make([]string, 0)
					}
					areaPhones[area] = append(areaPhones[area], line)
				}
			}
		}

		// 更新进度
		if processedLines%10000 == 0 {
			progress := float64(processedLines) / float64(totalLines) * 0.7
			a.areaSplitProgress.SetValue(progress)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取文件失败: %v", err)
	}

	a.areaSplitProgress.SetValue(0.7)
	a.areaSplitStatus.SetText("🔄 正在生成文件...")

	// 写入未知国家文件
	if len(unknownCountryPhones) > 0 {
		if err := writePhonesToFile(filepath.Join(outputDir, "未知国家.txt"), unknownCountryPhones); err != nil {
			return err
		}
		fmt.Printf("✅ 生成文件: 未知国家.txt (%d个号码)\n", len(unknownCountryPhones))
	}

	// 写入未知地区文件
	if len(unknownAreaPhones) > 0 {
		if err := writePhonesToFile(filepath.Join(outputDir, "未知地区.txt"), unknownAreaPhones); err != nil {
			return err
		}
		fmt.Printf("✅ 生成文件: 未知地区.txt (%d个号码)\n", len(unknownAreaPhones))
	}

	// 写入各地区文件
	totalAreas := len(areaPhones)
	currentArea := 0

	for area, phones := range areaPhones {
		if len(phones) == 0 {
			continue
		}

		fileName := filepath.Join(outputDir, fmt.Sprintf("%s.txt", area))
		if err := writePhonesToFile(fileName, phones); err != nil {
			return err
		}

		currentArea++
		progress := 0.7 + float64(currentArea)/float64(totalAreas)*0.3
		a.areaSplitProgress.SetValue(progress)

		fmt.Printf("✅ 生成文件: %s.txt (%d个号码)\n", area, len(phones))
	}

	// 输出统计信息
	fmt.Printf("✅ 按地区拆分完成:\n")
	fmt.Printf("   未知国家: %d个号码\n", len(unknownCountryPhones))
	fmt.Printf("   未知地区: %d个号码\n", len(unknownAreaPhones))
	for area, phones := range areaPhones {
		if len(phones) > 0 {
			fmt.Printf("   %s: %d个号码\n", area, len(phones))
		}
	}

	return nil
}

// 写入号码到文件
func writePhonesToFile(fileName string, phones []string) error {
	outputFile, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("创建文件 %s 失败: %v", fileName, err)
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	for _, phone := range phones {
		if _, err := writer.WriteString(phone + "\n"); err != nil {
			return fmt.Errorf("写入文件 %s 失败: %v", fileName, err)
		}
	}

	return nil
}
