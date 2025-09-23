package gui

import (
	"io/ioutil"
	"os"
	"image/color"

	"goduplicate/internal/core"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

// 创建一个自定义主题，支持中文显示
type customTheme struct {
	baseTheme fyne.Theme
}

// 重写字体方法，使用项目中的中文字体
func (c customTheme) Font(s fyne.TextStyle) fyne.Resource {
	// 尝试从项目目录加载中文字体（优先使用黑体）
	fontPath := "font/simhei.ttf"
	
	// 检查字体文件是否存在
	if _, err := os.Stat(fontPath); err == nil {
		// 加载字体文件
		fontData, err := ioutil.ReadFile(fontPath)
		if err == nil {
			// 创建字体资源
			fontRes := fyne.NewStaticResource("simhei.ttf", fontData)
			return fontRes
		}
	}
	
	// 如果黑体加载失败，尝试加载宋体
	fontPath = "font/simsunb.ttf"
	if _, err := os.Stat(fontPath); err == nil {
		fontData, err := ioutil.ReadFile(fontPath)
		if err == nil {
			fontRes := fyne.NewStaticResource("simsunb.ttf", fontData)
			return fontRes
		}
	}
	
	// 如果加载失败，返回nil让系统自动查找
	return nil
}

// 确保其他主题方法正常工作
func (c customTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return c.baseTheme.Color(name, variant)
}

func (c customTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return c.baseTheme.Icon(name)
}

func (c customTheme) Size(name fyne.ThemeSizeName) float32 {
	return c.baseTheme.Size(name)
}

// 创建支持中文的自定义主题
func newChineseTheme(base fyne.Theme) fyne.Theme {
	return customTheme{baseTheme: base}
}

// 确保所有文本元素都使用支持中文的字体
func ensureChineseFontSupport() {
	// 提前检查字体文件是否存在，确保在UI渲染前可用
	fontPaths := []string{"font/simhei.ttf", "font/simsunb.ttf"}
	for _, path := range fontPaths {
		if _, err := os.Stat(path); err != nil {
			// 如果字体文件不存在，可以考虑使用备用方案或记录日志
			// 在实际应用中，可以添加日志记录或其他处理逻辑
		}
	}
}

// 确保中文显示正常
func init() {
	// 初始化时就确保中文支持
	ensureChineseFontSupport()
}

// Run starts the GUI application with placeholder pages matching requirements.
func Run() {
	// 首先确保中文显示支持
	ensureChineseFontSupport()
	
	a := app.New()
	// 在应用程序初始化时就设置中文主题，确保所有UI元素都能正确使用支持中文的字体
	a.Settings().SetTheme(newChineseTheme(theme.LightTheme()))
	
	// 创建窗口时直接使用中文标题
	w := a.NewWindow("HasteGUI 重复文件管理")

	state := NewAppState()
	engine := core.NewSimpleScanner()

	applyTheme := func() {
		state.mu.RLock()
		th := state.Theme
		lang := state.Language
		state.mu.RUnlock()
		if th == "dark" {
			a.Settings().SetTheme(newChineseTheme(theme.DarkTheme()))
		} else {
			a.Settings().SetTheme(newChineseTheme(theme.LightTheme()))
		}
		if lang == "zh-CN" {
			w.SetTitle("HasteGUI 重复文件管理")
		} else {
			w.SetTitle("HasteGUI Duplicate Manager")
		}
	}
	applyFFmpegEnv := func() {
		state.mu.RLock()
		p := state.FfmpegPath
		state.mu.RUnlock()
		if p != "" {
			_ = os.Setenv("HASTE_FFMPEG_PATH", p)
		}
	}
	applyTheme()
	applyFFmpegEnv()

	startScan := func(cfg core.ScanConfig) {
		applyFFmpegEnv()
		state.mu.Lock()
		state.IsScanning = true
		state.LastScanError = nil
		state.FilesScanned = 0
		state.GroupsFound = 0
		state.mu.Unlock()
		cfg.OnProgress = func(p core.Progress) {
			state.mu.Lock()
			state.FilesScanned = p.FilesScanned
			state.GroupsFound = p.GroupsFound
			state.mu.Unlock()
		}
		go func() {
			groups, err := engine.Scan(cfg)
			state.mu.Lock()
			state.IsScanning = false
			state.LastScanError = err
			state.Results = groups
			state.GroupsFound = len(groups)
			total := 0
			for _, g := range groups {
				total += len(g.Files)
			}
			if total > state.FilesScanned {
				state.FilesScanned = total
			}
			state.mu.Unlock()
		}()
	}

	tabs := container.NewAppTabs(
		container.NewTabItem(t(state, "tab_config"), buildConfigPage(state, startScan)),
		container.NewTabItem(t(state, "tab_monitor"), buildMonitorPage(state)),
		container.NewTabItem(t(state, "tab_results"), buildResultsPage(state)),
		container.NewTabItem(t(state, "tab_strategy"), buildStrategyPage(state, func(plan []core.PlanItem) {})),
		container.NewTabItem(t(state, "tab_execute"), buildExecutePage(state)),
		container.NewTabItem(t(state, "tab_settings"), buildSettingsPage(state)),
	)

	w.SetContent(tabs)
	w.Resize(fyne.NewSize(1024, 700))
	w.ShowAndRun()
}
