package main

import (
	"fyne.io/fyne/v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

// getIconResource 从文件系统加载图标资源
func getIconResource() fyne.Resource {
	// 获取当前执行文件的目录
	exePath, err := os.Executable()
	if err != nil {
		return nil
	}
	
	// 构建图标文件路径
	iconPath := filepath.Join(filepath.Dir(exePath), "icon.ico")
	
	// 检查文件是否存在
	if _, err := os.Stat(iconPath); os.IsNotExist(err) {
		// 如果exe目录没有，尝试当前工作目录
		iconPath = "icon.ico"
		if _, err := os.Stat(iconPath); os.IsNotExist(err) {
			return nil
		}
	}
	
	// 读取图标文件内容
	iconData, err := ioutil.ReadFile(iconPath)
	if err != nil {
		return nil
	}
	
	// 创建资源
	return fyne.NewStaticResource("icon.ico", iconData)
}
