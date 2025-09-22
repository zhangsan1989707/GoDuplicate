package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

// 创建一个支持中文的自定义主题
type customChineseTheme struct {
	baseTheme fyne.Theme
}

// 重写字体方法，确保使用支持中文的字体
func (c customChineseTheme) Font(s fyne.TextStyle) fyne.Resource {
	// 在Windows系统上，返回nil让Fyne自动查找支持中文的系统字体
	return nil
}

// 其他主题方法保持不变
func (c customChineseTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return c.baseTheme.Color(name, variant)
}

func (c customChineseTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return c.baseTheme.Icon(name)
}

func (c customChineseTheme) Size(name fyne.ThemeSizeName) float32 {
	return c.baseTheme.Size(name)
}

func main() {
	// 创建应用程序
	a := app.New()
	// 创建主窗口
	w := a.NewWindow("中文显示测试")
	
	// 设置自定义主题以支持中文显示
	a.Settings().SetTheme(customChineseTheme{baseTheme: theme.LightTheme()})
	
	// 创建一个包含中文文本的标签
	chineseText := widget.NewLabel("这是一个中文显示测试\n中文应该能正常显示，而不是显示为问号")
	chineseText.TextStyle.Monospace = false // 使用非等宽字体，通常更好地支持中文
	
	// 创建一个按钮，也包含中文
	chineseButton := widget.NewButton("点击测试", func() {
		chineseText.SetText("按钮被点击了！\n中文显示正常吗？")
	})
	
	// 将标签和按钮放入一个垂直容器
	content := container.NewVBox(
		chineseText,
		chineseButton,
	)
	
	// 设置窗口内容
	w.SetContent(content)
	// 设置窗口大小
	w.Resize(fyne.NewSize(400, 200))
	// 显示并运行应用程序
	w.ShowAndRun()
}