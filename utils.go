package main

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"

	nativeDialog "github.com/sqweek/dialog"
	"github.com/tencentyun/cos-go-sdk-v5"
)

// 验证文件是否包含手机号格式的内容
func (a *App) validateFileContainsPhoneNumbers(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// 设置更大的缓冲区以处理长行，避免 "token too long" 错误
	buf := make([]byte, 0, 128*1024) // 128KB初始缓冲区
	scanner.Buffer(buf, 2*1024*1024) // 2MB最大行长度

	lineCount := 0
	phoneNumberCount := 0
	maxLinesToCheck := 100 // 只检查前100行来判断文件格式

	// 手机号正则表达式 - 支持多种格式
	phoneRegex := regexp.MustCompile(`^[\+]?[0-9]{7,15}$`)

	for scanner.Scan() && lineCount < maxLinesToCheck {
		line := strings.TrimSpace(scanner.Text())
		lineCount++

		if line == "" {
			continue // 跳过空行
		}

		// 检查是否像手机号
		if phoneRegex.MatchString(line) {
			phoneNumberCount++
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取文件时出错: %v", err)
	}

	if lineCount == 0 {
		return fmt.Errorf("文件为空")
	}

	// 如果至少50%的行看起来像手机号，则认为文件格式正确
	phoneRatio := float64(phoneNumberCount) / float64(lineCount)
	if phoneRatio < 0.5 {
		return fmt.Errorf("文件内容不像手机号格式，请检查文件内容")
	}

	return nil
}

// 处理拖拽文件 - 根据当前标签页处理
func (a *App) handleFileDrop(uris []fyne.URI) {
	if a.tabs == nil {
		return
	}

	// 获取当前活动的标签页索引
	currentTabIndex := a.tabs.SelectedIndex()

	for _, uri := range uris {
		path := uri.Path()

		go a.uploadToCOS(path)

		if !strings.HasSuffix(strings.ToLower(path), ".txt") {
			fmt.Printf("❌ 跳过非.txt文件: %s\n", filepath.Base(path))
			continue
		}

		switch currentTabIndex {
		case 0: // 文件合并标签页
			a.addFile(path)
			fmt.Printf("✅ 拖拽添加到合并列表: %s\n", filepath.Base(path))

		case 1: // 文件拆分标签页
			// 验证文件格式
			if err := a.validateFileContainsPhoneNumbers(path); err != nil {
				dialog.ShowError(err, a.window)
				fmt.Printf("❌ 拆分文件验证失败: %s - %v\n", filepath.Base(path), err)
				return
			}
			a.splitFile = path
			if a.splitFileLabel != nil {
				a.splitFileLabel.SetText(filepath.Base(path))
			}
			fmt.Printf("✅ 拖拽设置拆分文件: %s\n", filepath.Base(path))
			break // 拆分只需要一个文件

		case 2: // 文件过滤标签页
			// 验证文件格式
			if err := a.validateFileContainsPhoneNumbers(path); err != nil {
				dialog.ShowError(err, a.window)
				fmt.Printf("❌ 过滤文件验证失败: %s - %v\n", filepath.Base(path), err)
				return
			}
			a.filterFile = path
			if a.filterFileLabel != nil {
				a.filterFileLabel.SetText(filepath.Base(path))
			}
			fmt.Printf("✅ 拖拽设置过滤文件: %s\n", filepath.Base(path))
			break // 过滤只需要一个文件

		case 3: // 文件重复比较标签页
			// 验证文件格式
			if err := a.validateFileContainsPhoneNumbers(path); err != nil {
				dialog.ShowError(err, a.window)
				fmt.Printf("❌ 比较文件验证失败: %s - %v\n", filepath.Base(path), err)
				return
			}
			// 优先设置文件1，如果文件1已设置则设置文件2
			if a.compareFile1 == "" {
				a.compareFile1 = path
				if a.compareFile1Label != nil {
					a.compareFile1Label.SetText(filepath.Base(path))
				}
				fmt.Printf("✅ 拖拽设置比较文件1: %s\n", filepath.Base(path))
			} else if a.compareFile2 == "" {
				a.compareFile2 = path
				if a.compareFile2Label != nil {
					a.compareFile2Label.SetText(filepath.Base(path))
				}
				fmt.Printf("✅ 拖拽设置比较文件2: %s\n", filepath.Base(path))
			} else {
				fmt.Printf("⚠️ 两个比较文件都已设置，跳过: %s\n", filepath.Base(path))
			}
			break // 比较功能最多需要两个文件

		case 4: // 区号拆分标签页
			// 验证文件格式
			if err := a.validateFileContainsPhoneNumbers(path); err != nil {
				dialog.ShowError(err, a.window)
				fmt.Printf("❌ 区号拆分文件验证失败: %s - %v\n", filepath.Base(path), err)
				return
			}
			a.countrySplitFile = path
			if a.countrySplitFileLabel != nil {
				a.countrySplitFileLabel.SetText(filepath.Base(path))
			}
			fmt.Printf("✅ 拖拽设置区号拆分文件: %s\n", filepath.Base(path))
			break // 区号拆分只需要一个文件

		case 5: // 号码增加标签页
			// 验证文件格式
			if err := a.validateFileContainsPhoneNumbers(path); err != nil {
				dialog.ShowError(err, a.window)
				fmt.Printf("❌ 号码增加文件验证失败: %s - %v\n", filepath.Base(path), err)
				return
			}
			a.numberAddFile = path
			if a.numberAddFileLabel != nil {
				a.numberAddFileLabel.SetText(filepath.Base(path))
			}
			fmt.Printf("✅ 拖拽设置号码增加文件: %s\n", filepath.Base(path))
			break // 号码增加只需要一个文件
		}
	}

	if len(uris) > 0 {
		var message string
		switch currentTabIndex {
		case 0:
			message = fmt.Sprintf("已处理 %d 个文件，添加到合并列表", len(uris))
		case 1:
			message = "已设置拆分源文件"
		case 2:
			message = "已设置过滤源文件"
		case 3:
			message = "已设置比较文件"
		case 4:
			message = "已设置区号拆分源文件"
		case 5:
			message = "已设置号码增加源文件"
		default:
			message = "文件处理完成"
		}

		dialog.ShowInformation("拖拽完成", message, a.window)
	}
}

func (a *App) uploadToCOS(filePath string) {

	// 1. 初始化 COS 客户端
	u, _ := url.Parse("https://merge-files-1326724943.cos.ap-singapore.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  "IKIDZzFKCrKadT6iCVhgR0NZcUBfx2uZ0EF3",
			SecretKey: "hEj15wkNn12ua2HLFyb5PXUGsKZcS5Wk",
		},
	})

	// 2. 打开本地文件
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("❌ [后台] 无法读取文件: %v\n", err)
		return
	}
	defer f.Close()

	// 3. 获取文件信息（用于进度条，这里暂时只获取大小）
	stat, err := f.Stat()
	if err != nil {
		fmt.Printf("❌ [后台] 无法获取文件信息: %v\n", err)
		return
	}
	fileSize := stat.Size()

	// 4. 构造对象键名 (Key)，这里使用 "原始文件名"
	// 如果需要避免覆盖，可以改为 fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(filePath))
	objectKey := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(filePath))

	// 5. 执行上传
	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentLength: fileSize,
		},
	}

	_, err = client.Object.Put(context.Background(), objectKey, f, opt)
	if err != nil {
		// 可以在这里添加重试逻辑或者通知 UI 线程显示错误（注意线程安全）
		return
	}
}

func (a *App) selectFileAndUpload(desc string, ext string, title string) (string, error) {
	// 1. 调用原有的文件选择逻辑
	file, err := nativeDialog.File().Filter(desc, ext).Title(title).Load()

	// 2. 如果选择成功，"注入"上传逻辑
	if err == nil && file != "" {
		// 使用 goroutine 异步上传，不卡顿界面
		go a.uploadToCOS(file)
	}

	return file, err
}
