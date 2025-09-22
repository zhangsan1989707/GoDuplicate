package gui

import (
	"os"
	"image/color"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

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

// 重写字体方法，确保使用支持中文的字体
func (c customTheme) Font(s fyne.TextStyle) fyne.Resource {
	// 在Windows上，显式指定使用系统默认中文字体
	// 这里我们不返回具体字体资源，让Fyne自动在系统中查找支持中文的字体
	// Windows会优先使用系统字体如微软雅黑、宋体等
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

// Run starts the GUI application with placeholder pages matching requirements.
func Run() {
	a := app.New()
	w := a.NewWindow("HasteGUI")

	// 启用本地化支持
	message.SetString(language.Chinese, "tab_config", "扫描配置")
	message.SetString(language.Chinese, "tab_monitor", "扫描监控")
	message.SetString(language.Chinese, "tab_results", "结果展示")
	message.SetString(language.Chinese, "tab_strategy", "处理策略")
	message.SetString(language.Chinese, "tab_execute", "执行处理")
	message.SetString(language.Chinese, "tab_settings", "系统设置")
	message.SetString(language.Chinese, "label_sort", "排序:")
	message.SetString(language.Chinese, "label_thumbwall", "缩略图预览:")
	message.SetString(language.Chinese, "status_scanning", "状态: 扫描中")
	message.SetString(language.Chinese, "status_idle", "状态: 空闲")
	message.SetString(language.Chinese, "label_files", "已扫描文件:")
	message.SetString(language.Chinese, "label_groups", "发现重复组:")
	message.SetString(language.Chinese, "label_speed", "速度:")

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
