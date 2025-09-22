package gui

import (
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

// 重写字体方法，确保使用支持中文的字体
func (c customTheme) Font(s fyne.TextStyle) fyne.Resource {
	// 在Windows上，返回nil让Fyne自动在系统中查找支持中文的字体
	// 但我们需要确保系统能正确识别中文字符
	// Windows会优先使用系统字体如微软雅黑、宋体、黑体等
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
