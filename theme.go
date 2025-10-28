package main

import (
	"image/color"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type modernTheme struct{}

func (t *modernTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNamePrimary:
		return color.RGBA{R: 0x00, G: 0x7A, B: 0xFF, A: 0xff} // 现代蓝色
	case theme.ColorNameBackground:
		if variant == theme.VariantDark {
			return color.RGBA{R: 0x1a, G: 0x1a, B: 0x1a, A: 0xff}
		}
		return color.RGBA{R: 0xfa, G: 0xfa, B: 0xfa, A: 0xff} // 浅灰背景
	case theme.ColorNameForeground:
		if variant == theme.VariantDark {
			return color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
		}
		return color.RGBA{R: 0x2d, G: 0x3e, B: 0x50, A: 0xff} // 深灰文字
	case theme.ColorNameButton:
		return color.RGBA{R: 0x00, G: 0x7A, B: 0xFF, A: 0xff} // 按钮蓝色
	case theme.ColorNameDisabledButton:
		return color.RGBA{R: 0xcc, G: 0xcc, B: 0xcc, A: 0xff}
	case theme.ColorNameHover:
		return color.RGBA{R: 0x00, G: 0x5a, B: 0xcc, A: 0x30} // 悬停效果
	case theme.ColorNameFocus:
		return color.RGBA{R: 0x00, G: 0x7A, B: 0xFF, A: 0x80}
	case theme.ColorNameSelection:
		return color.RGBA{R: 0x00, G: 0x7A, B: 0xFF, A: 0x30}
	case theme.ColorNameShadow:
		return color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x20}
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (t *modernTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *modernTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *modernTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return 14
	case theme.SizeNameCaptionText:
		return 12
	case theme.SizeNameHeadingText:
		return 20
	case theme.SizeNameSubHeadingText:
		return 16
	case theme.SizeNamePadding:
		return 8
	case theme.SizeNameInlineIcon:
		return 20
	case theme.SizeNameScrollBar:
		return 16
	case theme.SizeNameSeparatorThickness:
		return 1
	}
	return theme.DefaultTheme().Size(name)
}
