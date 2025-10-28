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
func (a *App) createCountrySplitTab() *fyne.Container {
	// 创建拖拽区域用于区号拆分
	countrySplitDropArea := a.createCountrySplitDropArea()
	
	// 文件选择
	a.countrySplitFileLabel = widget.NewLabel("未选择文件")
	selectFileBtn := widget.NewButtonWithIcon("📁 选择文件", nil, func() {
		// 使用Windows原生文件选择对话框
		file, err := nativeDialog.File().Filter("文本文件", "txt").Title("选择要按区号拆分的文件").Load()
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
				fmt.Printf("❌ 区号拆分文件验证失败: %s - %v\n", filepath.Base(file), err)
				return
			}
			a.countrySplitFile = file
			a.countrySplitFileLabel.SetText(filepath.Base(file))
			fmt.Printf("✅ 选择区号拆分文件: %s\n", filepath.Base(file))
		}
	})
	
	// 开始拆分按钮
	splitBtn := widget.NewButtonWithIcon("🌍 开始拆分", nil, func() {
		if a.countrySplitFile == "" {
			dialog.ShowInformation("提示", "请先选择要拆分的文件", a.window)
			return
		}
		a.startCountrySplit()
	})
	splitBtn.Importance = widget.HighImportance
	
	// 进度区域
	a.countrySplitProgress = widget.NewProgressBar()
	a.countrySplitStatus = widget.NewLabel("📋 就绪")
	a.countrySplitStatus.TextStyle = fyne.TextStyle{Italic: true}
	
	// 主布局
	topSection := container.NewVBox(
		widget.NewRichTextFromMarkdown("## 🌍 按国家区号拆分\n拖拽文件到下方区域或点击选择文件按钮"),
		container.NewPadded(countrySplitDropArea),
	)
	
	middleSection := container.NewVBox(
		widget.NewSeparator(),
		widget.NewLabel("📄 选择的文件:"),
		a.countrySplitFileLabel,
		selectFileBtn,
		widget.NewSeparator(),
		widget.NewLabel("⚙️ 拆分说明:"),
		widget.NewLabel("• 自动识别手机号的国家区号"),
		widget.NewLabel("• 按国家分组生成独立文件"),
		widget.NewLabel("• 支持美国、英国等主要国家"),
		widget.NewLabel("• 输出文件格式: 国家名.txt"),
	)
	
	bottomSection := container.NewVBox(
		widget.NewSeparator(),
		container.NewHBox(widget.NewLabel(""), splitBtn),
		widget.NewSeparator(),
		widget.NewLabel("📊 进度状态:"),
		a.countrySplitProgress,
		a.countrySplitStatus,
	)
	
	return container.NewVBox(
		topSection,
		middleSection,
		bottomSection,
	)
}

// 创建区号拆分专用拖拽区域
func (a *App) createCountrySplitDropArea() *fyne.Container {
	// 创建拖拽区域的内容
	dropIcon := widget.NewLabel("🌍")
	dropIcon.Alignment = fyne.TextAlignCenter
	dropIcon.TextStyle = fyne.TextStyle{Bold: true}
	
	dropLabel := widget.NewLabel("拖拽文件到此处或点击选择")
	dropLabel.Alignment = fyne.TextAlignCenter
	dropLabel.TextStyle = fyne.TextStyle{Bold: true}
	
	dropHint := widget.NewLabel("选择包含手机号的文件进行区号拆分")
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
		file, err := nativeDialog.File().Filter("文本文件", "txt").Title("选择要按区号拆分的文件").Load()
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
				fmt.Printf("❌ 区号拆分文件验证失败: %s - %v\n", filepath.Base(file), err)
				return
			}
			a.countrySplitFile = file
			a.countrySplitFileLabel.SetText(filepath.Base(file))
			fmt.Printf("✅ 选择区号拆分文件: %s\n", filepath.Base(file))
		}
	})
	
	// 设置按钮样式
	dropButton.Resize(fyne.NewSize(500, 120))
	dropButton.Importance = widget.LowImportance
	
	// 创建叠加容器
	overlayContainer := container.NewStack(dropButton, dropContent)
	
	return container.NewPadded(overlayContainer)
}

// 开始按国家区号拆分
func (a *App) startCountrySplit() {
	if a.countrySplitFile == "" {
		return
	}
	
	go func() {
		a.countrySplitStatus.SetText("🔄 正在按区号拆分文件...")
		a.countrySplitProgress.SetValue(0)
		
		// 选择输出目录
		outputDir, err := nativeDialog.Directory().Title("选择拆分文件的输出文件夹").Browse()
		if err != nil {
			a.countrySplitStatus.SetText("❌ 拆分已取消")
			return
		}
		
		err = a.performCountrySplit(outputDir)
		if err != nil {
			a.countrySplitStatus.SetText("❌ 拆分失败: " + err.Error())
			dialog.ShowError(err, a.window)
		} else {
			a.countrySplitStatus.SetText("✅ 拆分完成")
			dialog.ShowInformation("完成", "按国家区号拆分成功！\n已生成各国家的独立文件", a.window)
		}
		a.countrySplitProgress.SetValue(1.0)
	}()
}

// 国家区号映射表
type CountryCode struct {
	Name     string   // 国家名称
	Prefixes []string // 手机号前缀列表
}

// 获取国家区号映射表
func getCountryCodes() []CountryCode {
	return []CountryCode{
		// 美国
		{
			Name: "美国",
			Prefixes: []string{
				"1201", "1202", "1203", "1205", "1206", "1207", "1208", "1209", "1210",
				"1212", "1213", "1214", "1215", "1216", "1217", "1218", "1219", "1224",
				"1225", "1228", "1229", "1231", "1234", "1239", "1240", "1248", "1251",
				"1252", "1253", "1254", "1256", "1260", "1262", "1267", "1269", "1270",
				"1276", "1281", "1301", "1302", "1303", "1304", "1305", "1307", "1308",
				"1309", "1310", "1312", "1313", "1314", "1315", "1316", "1317", "1318",
				"1319", "1320", "1321", "1323", "1325", "1330", "1331", "1334", "1336",
				"1337", "1339", "1341", "1347", "1351", "1352", "1360", "1361",
				"1364", "1380", "1385", "1386", "1401", "1402", "1404", "1405", "1406",
				"1407", "1408", "1409", "1410", "1412", "1413", "1414", "1415", "1417",
				"1419", "1423", "1424", "1425", "1430", "1432", "1434", "1435", "1440",
				"1442", "1443", "1458", "1463", "1464", "1469", "1470", "1475", "1478",
				"1479", "1480", "1484", "1501", "1502", "1503", "1504", "1505", "1507",
				"1508", "1509", "1510", "1512", "1513", "1515", "1516", "1517", "1518",
				"1520", "1530", "1531", "1534", "1539", "1540", "1541", "1551", "1559",
				"1561", "1562", "1563", "1564", "1567", "1570", "1571", "1573", "1574",
				"1575", "1580", "1585", "1586", "1601", "1602", "1603", "1605", "1606",
				"1607", "1608", "1609", "1610", "1612", "1614", "1615", "1616", "1617",
				"1618", "1619", "1620", "1623", "1626", "1628", "1629", "1630", "1631",
				"1636", "1641", "1646", "1650", "1651", "1657", "1660", "1661", "1662",
				"1667", "1669", "1678", "1681", "1682", "1689", "1701", "1702",
				"1703", "1704", "1706", "1707", "1708", "1712", "1713", "1714", "1715",
				"1716", "1717", "1718", "1719", "1720", "1724", "1725", "1727", "1731",
				"1732", "1734", "1737", "1740", "1743", "1747", "1754", "1757", "1760",
				"1762", "1763", "1764", "1765", "1769", "1770", "1772", "1773", "1774",
				"1775", "1779", "1781", "1785", "1786", "1801", "1802", "1803",
				"1804", "1805", "1806", "1810", "1812", "1813", "1814", "1815",
				"1816", "1817", "1818", "1828", "1830", "1831", "1832", "1843", "1845",
				"1847", "1848", "1850", "1856", "1857", "1858", "1859", "1860", "1862",
				"1863", "1864", "1865", "1870", "1872", "1878", "1901", "1903", "1904",
				"1906", "1907", "1908", "1909", "1910", "1912", "1913", "1914", "1915",
				"1916", "1917", "1918", "1919", "1920", "1925", "1928", "1929", "1930",
				"1931", "1934", "1936", "1937", "1940", "1941", "1947", "1949", "1951",
				"1952", "1954", "1956", "1959", "1970", "1971", "1972", "1973", "1978",
				"1979", "1980", "1984", "1985", "1989",
			},
		},
		// 加拿大
		{
			Name: "加拿大",
			Prefixes: []string{
				"1403", "1587", "1825", // 阿尔伯塔省
				"1236", "1250", "1604", "1672", "1778", // 不列颠哥伦比亚省
				"1204", "1431", // 马尼托巴省
				"1506", // 新不伦瑞克省
				"1709", // 纽芬兰和拉布拉多省
				"1782", "1902", // 新斯科舍省
				"1226", "1249", "1289", "1343", "1365", "1416", "1437", "1519", "1548", "1613", "1647", "1705", "1807", "1905", // 安大略省
				"1418", "1438", "1450", "1514", "1579", "1581", "1819", "1873", // 魁北克省
				"1306", "1639", // 萨斯喀彻温省
				"1867", // 西北地区、努纳武特地区、育空地区
			},
		},
		// 俄罗斯/哈萨克斯坦 +7
		{Name: "俄罗斯", Prefixes: []string{"7"}},
		// 埃及 +20
		{Name: "埃及", Prefixes: []string{"20"}},
		// 南非 +27
		{Name: "南非", Prefixes: []string{"27"}},
		// 希腊 +30
		{Name: "希腊", Prefixes: []string{"30"}},
		// 荷兰 +31
		{Name: "荷兰", Prefixes: []string{"31"}},
		// 比利时 +32
		{Name: "比利时", Prefixes: []string{"32"}},
		// 法国 +33
		{Name: "法国", Prefixes: []string{"33"}},
		// 西班牙 +34
		{Name: "西班牙", Prefixes: []string{"34"}},
		// 匈牙利 +36
		{Name: "匈牙利", Prefixes: []string{"36"}},
		// 意大利 +39
		{Name: "意大利", Prefixes: []string{"39"}},
		// 罗马尼亚 +40
		{Name: "罗马尼亚", Prefixes: []string{"40"}},
		// 瑞士 +41
		{Name: "瑞士", Prefixes: []string{"41"}},
		// 奥地利 +43
		{Name: "奥地利", Prefixes: []string{"43"}},
		// 英国 +44
		{Name: "英国", Prefixes: []string{"44"}},
		// 丹麦 +45
		{Name: "丹麦", Prefixes: []string{"45"}},
		// 瑞典 +46
		{Name: "瑞典", Prefixes: []string{"46"}},
		// 挪威 +47
		{Name: "挪威", Prefixes: []string{"47"}},
		// 波兰 +48
		{Name: "波兰", Prefixes: []string{"48"}},
		// 德国 +49
		{Name: "德国", Prefixes: []string{"49"}},
		// 秘鲁 +51
		{Name: "秘鲁", Prefixes: []string{"51"}},
		// 墨西哥 +52
		{Name: "墨西哥", Prefixes: []string{"52"}},
		// 古巴 +53
		{Name: "古巴", Prefixes: []string{"53"}},
		// 阿根廷 +54
		{Name: "阿根廷", Prefixes: []string{"54"}},
		// 巴西 +55
		{Name: "巴西", Prefixes: []string{"55"}},
		// 智利 +56
		{Name: "智利", Prefixes: []string{"56"}},
		// 哥伦比亚 +57
		{Name: "哥伦比亚", Prefixes: []string{"57"}},
		// 委内瑞拉 +58
		{Name: "委内瑞拉", Prefixes: []string{"58"}},
		// 马来西亚 +60
		{Name: "马来西亚", Prefixes: []string{"60"}},
		// 澳大利亚 +61
		{Name: "澳大利亚", Prefixes: []string{"61"}},
		// 印度尼西亚 +62
		{Name: "印度尼西亚", Prefixes: []string{"62"}},
		// 菲律宾 +63
		{Name: "菲律宾", Prefixes: []string{"63"}},
		// 新西兰 +64
		{Name: "新西兰", Prefixes: []string{"64"}},
		// 新加坡 +65
		{Name: "新加坡", Prefixes: []string{"65"}},
		// 泰国 +66
		{Name: "泰国", Prefixes: []string{"66"}},
		// 日本 +81
		{Name: "日本", Prefixes: []string{"81"}},
		// 韩国 +82
		{Name: "韩国", Prefixes: []string{"82"}},
		// 越南 +84
		{Name: "越南", Prefixes: []string{"84"}},
		// 土耳其 +90
		{Name: "土耳其", Prefixes: []string{"90"}},
		// 印度 +91
		{Name: "印度", Prefixes: []string{"91"}},
		// 巴基斯坦 +92
		{Name: "巴基斯坦", Prefixes: []string{"92"}},
		// 阿富汗 +93
		{Name: "阿富汗", Prefixes: []string{"93"}},
		// 斯里兰卡 +94
		{Name: "斯里兰卡", Prefixes: []string{"94"}},
		// 缅甸 +95
		{Name: "缅甸", Prefixes: []string{"95"}},
		// 伊朗 +98
		{Name: "伊朗", Prefixes: []string{"98"}},
		// 摩洛哥 +212
		{Name: "摩洛哥", Prefixes: []string{"212"}},
		// 阿尔及利亚 +213
		{Name: "阿尔及利亚", Prefixes: []string{"213"}},
		// 突尼斯 +216
		{Name: "突尼斯", Prefixes: []string{"216"}},
		// 利比亚 +218
		{Name: "利比亚", Prefixes: []string{"218"}},
		// 尼日利亚 +234
		{Name: "尼日利亚", Prefixes: []string{"234"}},
		// 肯尼亚 +254
		{Name: "肯尼亚", Prefixes: []string{"254"}},
		// 坦桑尼亚 +255
		{Name: "坦桑尼亚", Prefixes: []string{"255"}},
		// 乌干达 +256
		{Name: "乌干达", Prefixes: []string{"256"}},
		// 津巴布韦 +263
		{Name: "津巴布韦", Prefixes: []string{"263"}},
		// 葡萄牙 +351
		{Name: "葡萄牙", Prefixes: []string{"351"}},
		// 卢森堡 +352
		{Name: "卢森堡", Prefixes: []string{"352"}},
		// 爱尔兰 +353
		{Name: "爱尔兰", Prefixes: []string{"353"}},
		// 冰岛 +354
		{Name: "冰岛", Prefixes: []string{"354"}},
		// 阿尔巴尼亚 +355
		{Name: "阿尔巴尼亚", Prefixes: []string{"355"}},
		// 马耳他 +356
		{Name: "马耳他", Prefixes: []string{"356"}},
		// 芬兰 +358
		{Name: "芬兰", Prefixes: []string{"358"}},
		// 保加利亚 +359
		{Name: "保加利亚", Prefixes: []string{"359"}},
		// 立陶宛 +370
		{Name: "立陶宛", Prefixes: []string{"370"}},
		// 拉脱维亚 +371
		{Name: "拉脱维亚", Prefixes: []string{"371"}},
		// 爱沙尼亚 +372
		{Name: "爱沙尼亚", Prefixes: []string{"372"}},
		// 摩尔多瓦 +373
		{Name: "摩尔多瓦", Prefixes: []string{"373"}},
		// 白俄罗斯 +375
		{Name: "白俄罗斯", Prefixes: []string{"375"}},
		// 乌克兰 +380
		{Name: "乌克兰", Prefixes: []string{"380"}},
		// 塞尔维亚 +381
		{Name: "塞尔维亚", Prefixes: []string{"381"}},
		// 黑山 +382
		{Name: "黑山", Prefixes: []string{"382"}},
		// 克罗地亚 +385
		{Name: "克罗地亚", Prefixes: []string{"385"}},
		// 斯洛文尼亚 +386
		{Name: "斯洛文尼亚", Prefixes: []string{"386"}},
		// 波黑 +387
		{Name: "波黑", Prefixes: []string{"387"}},
		// 马其顿 +389
		{Name: "马其顿", Prefixes: []string{"389"}},
		// 捷克 +420
		{Name: "捷克", Prefixes: []string{"420"}},
		// 斯洛伐克 +421
		{Name: "斯洛伐克", Prefixes: []string{"421"}},
		// 以色列 +972
		{Name: "以色列", Prefixes: []string{"972"}},
		// 阿联酋 +971
		{Name: "阿联酋", Prefixes: []string{"971"}},
		// 沙特阿拉伯 +966
		{Name: "沙特阿拉伯", Prefixes: []string{"966"}},
	}
}

// 根据手机号前缀识别国家
func identifyCountry(phoneNumber string) string {
	countryCodes := getCountryCodes()
	
	// 按前缀长度从长到短排序，优先匹配更长的前缀
	for _, country := range countryCodes {
		for _, prefix := range country.Prefixes {
			if strings.HasPrefix(phoneNumber, prefix) {
				return country.Name
			}
		}
	}
	
	return "未知国家"
}

// 执行按国家区号拆分操作
func (a *App) performCountrySplit(outputDir string) error {
	file, err := os.Open(a.countrySplitFile)
	if err != nil {
		return fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	// 用于存储每个国家的手机号
	countryPhones := make(map[string][]string)
	
	scanner := bufio.NewScanner(file)
	// 设置更大的缓冲区以处理长行，避免 "token too long" 错误
	buf := make([]byte, 0, 128*1024) // 128KB初始缓冲区
	scanner.Buffer(buf, 2*1024*1024) // 2MB最大行长度
	
	totalLines := 0
	processedLines := 0
	
	// 第一遍：计算总行数
	a.countrySplitStatus.SetText("🔄 正在计算文件行数...")
	for scanner.Scan() {
		totalLines++
	}
	
	// 重新打开文件进行处理
	file.Close()
	file, err = os.Open(a.countrySplitFile)
	if err != nil {
		return fmt.Errorf("重新打开文件失败: %v", err)
	}
	defer file.Close()
	
	scanner = bufio.NewScanner(file)
	scanner.Buffer(buf, 2*1024*1024)
	
	a.countrySplitStatus.SetText("🔄 正在识别国家区号...")
	
	// 第二遍：按国家分类手机号
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		processedLines++
		
		if line != "" {
			// 识别国家
			country := identifyCountry(line)
			
			// 添加到对应国家的列表中
			if countryPhones[country] == nil {
				countryPhones[country] = make([]string, 0)
			}
			countryPhones[country] = append(countryPhones[country], line)
		}
		
		// 更新进度
		if processedLines%10000 == 0 {
			progress := float64(processedLines) / float64(totalLines) * 0.7 // 70%用于分类
			a.countrySplitProgress.SetValue(progress)
		}
	}
	
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取文件失败: %v", err)
	}
	
	a.countrySplitProgress.SetValue(0.7)
	a.countrySplitStatus.SetText("🔄 正在生成国家文件...")
	
	// 第三遍：为每个国家创建文件
	countryCount := len(countryPhones)
	currentCountry := 0
	
	for country, phones := range countryPhones {
		if len(phones) == 0 {
			continue
		}
		
		// 创建国家文件
		fileName := filepath.Join(outputDir, fmt.Sprintf("%s.txt", country))
		outputFile, err := os.Create(fileName)
		if err != nil {
			return fmt.Errorf("创建国家文件 %s 失败: %v", fileName, err)
		}
		
		writer := bufio.NewWriter(outputFile)
		
		// 写入该国家的所有手机号
		for _, phone := range phones {
			_, err := writer.WriteString(phone + "\n")
			if err != nil {
				writer.Flush()
				outputFile.Close()
				return fmt.Errorf("写入文件 %s 失败: %v", fileName, err)
			}
		}
		
		writer.Flush()
		outputFile.Close()
		
		currentCountry++
		progress := 0.7 + float64(currentCountry)/float64(countryCount)*0.3 // 剩余30%用于写入文件
		a.countrySplitProgress.SetValue(progress)
		
		fmt.Printf("✅ 生成文件: %s (%d个手机号)\n", fileName, len(phones))
	}
	
	// 输出统计信息
	fmt.Printf("✅ 按国家区号拆分完成:\n")
	for country, phones := range countryPhones {
		if len(phones) > 0 {
			fmt.Printf("   %s: %d个手机号\n", country, len(phones))
		}
	}
	
	return nil
}
