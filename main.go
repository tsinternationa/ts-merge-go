package main

import (
	"fmt"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/flopp/go-findfont"
	"os"
	"strings"
)

var (
	version = "1.0.0" // 构建时会被替换
)

func init() {
	// 设置中文字体：解决中文乱码问题
	fontPaths := findfont.List()
	for _, path := range fontPaths {
		if strings.Contains(path, "msyh.ttf") || // 微软雅黑
			strings.Contains(path, "simhei.ttf") || // 黑体
			strings.Contains(path, "simsun.ttc") || // 宋体
			strings.Contains(path, "simkai.ttf") || // 楷体
			strings.Contains(path, "Microsoft YaHei") || // 微软雅黑
			strings.Contains(path, "SimSun") { // 宋体
			os.Setenv("FYNE_FONT", path)
			fmt.Printf("设置中文字体: %s\n", path)
			break
		}
	}
}

type App struct {
	window fyne.Window
	tabs   *container.AppTabs // 添加标签页引用

	// 合并相关
	mergeFiles    []string
	mergeList     *widget.List
	mergeDedup    *widget.Check
	mergeProgress *widget.ProgressBar
	mergeStatus   *widget.Label

	// 拆分相关
	splitFile      string
	splitFileLabel *widget.Label
	splitParts     *widget.Entry
	splitDedup     *widget.Check
	splitProgress  *widget.ProgressBar
	splitStatus    *widget.Label

	// 过滤相关
	filterFile      string
	filterFileLabel *widget.Label
	filterPrefix1   *widget.Entry // 第一个前缀输入框
	filterPrefix2   *widget.Entry // 第二个前缀输入框
	filterPrefix3   *widget.Entry // 第三个前缀输入框
	filterPrefix4   *widget.Entry // 第四个前缀输入框
	filterProgress  *widget.ProgressBar
	filterStatus    *widget.Label

	// 文件重复比较相关
	compareFile1      string
	compareFile1Label *widget.Label
	compareFile2      string
	compareFile2Label *widget.Label
	compareProgress   *widget.ProgressBar
	compareStatus     *widget.Label

	// 区号拆分相关
	countrySplitFile      string
	countrySplitFileLabel *widget.Label
	countrySplitProgress  *widget.ProgressBar
	countrySplitStatus    *widget.Label

	// 号码增加相关
	numberAddFile        string
	numberAddFileLabel   *widget.Label
	numberAddPosition    *widget.Entry
	numberAddDigit       *widget.Entry // 新增：用户输入要增加的数字
	numberAddRemoveEmpty *widget.Check
	numberAddProgress    *widget.ProgressBar
	numberAddStatus      *widget.Label
}

func main() {
	// 设置运行时参数以优化大文件处理
	runtime.GOMAXPROCS(runtime.NumCPU())

	myApp := app.New()

	// 设置应用程序图标
	if iconResource := getIconResource(); iconResource != nil {
		myApp.SetIcon(iconResource)
	}
	// 使用现代化主题
	myApp.Settings().SetTheme(&modernTheme{})

	window := myApp.NewWindow("TS-Merge v" + version + " - 高性能文件处理工具")
	window.Resize(fyne.NewSize(900, 800))
	window.CenterOnScreen()

	// 设置窗口最小尺寸
	window.SetFixedSize(false)

	app := &App{window: window}
	app.setupUI()

	// 设置拖拽功能 - 支持真正的文件拖拽
	window.SetOnDropped(func(position fyne.Position, uris []fyne.URI) {
		app.handleFileDrop(uris)
	})

	// 显示启动信息
	fmt.Printf("TS-Merge v%s 启动\n", version)
	fmt.Printf("Go版本: %s\n", runtime.Version())
	fmt.Printf("CPU核心数: %d\n", runtime.NumCPU())
	fmt.Printf("操作系统: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Println("✅ 拖拽功能已启用 - 可以直接拖拽.txt文件到窗口")

	window.ShowAndRun()
}

func (a *App) setupUI() {
	a.tabs = container.NewAppTabs(
		container.NewTabItem("📁 文件合并", a.createMergeTab()),
		container.NewTabItem("✂️ 文件拆分", a.createSplitTab()),
		container.NewTabItem("🔍 文件过滤", a.createFilterTab()),
		container.NewTabItem("🔄 文件重复", a.createCompareTab()),
		container.NewTabItem("🌍 区号拆分", a.createCountrySplitTab()),
		container.NewTabItem("🔢 号码增加", a.createNumberAddTab()),
	)
	a.tabs.SetTabLocation(container.TabLocationTop)

	// 添加边距和背景
	content := container.NewPadded(a.tabs)
	a.window.SetContent(content)
}
