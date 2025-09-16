package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"goduplicate/internal/core"
)

func buildSettingsPage(state *AppState) fyne.CanvasObject {
	theme := widget.NewSelect([]string{"light", "dark"}, func(v string) {
		state.mu.Lock()
		state.Theme = v
		state.mu.Unlock()
		if v == "dark" {
			fyne.CurrentApp().Settings().SetTheme(themePkgDark())
		} else {
			fyne.CurrentApp().Settings().SetTheme(themePkgLight())
		}
	})
	theme.Selected = state.Theme
	lang := widget.NewSelect([]string{"zh-CN", "en-US"}, func(v string) {
		state.mu.Lock()
		state.Language = v
		state.mu.Unlock()
		ws := fyne.CurrentApp().Driver().AllWindows()
		if len(ws) > 0 {
			if v == "zh-CN" {
				ws[0].SetTitle("HasteGUI 重复文件管理")
			} else {
				ws[0].SetTitle("HasteGUI Duplicate Manager")
			}
		}
	})
	lang.Selected = state.Language
	ffmpegEntry := widget.NewEntry()
	ffmpegEntry.SetPlaceHolder("可选：手动指定 ffmpeg 可执行路径")
	state.mu.RLock()
	ffmpegEntry.SetText(state.FfmpegPath)
	state.mu.RUnlock()
	ffmpegEntry.OnChanged = func(v string) { state.mu.Lock(); state.FfmpegPath = v; state.mu.Unlock() }

	presetName := widget.NewEntry()
	presetName.SetPlaceHolder("预设名称")
	savePresetBtn := widget.NewButton("保存扫描预设", func() {
		state.mu.RLock()
		cfg := state.ToScanConfig()
		state.mu.RUnlock()
		_, _ = core.SaveScanPreset(core.ScanPreset{Name: presetName.Text, Config: cfg})
	})
	loadPresetBtn := widget.NewButton("加载扫描预设", func() {
		p, err := core.LoadScanPreset(presetName.Text)
		if err == nil {
			state.mu.Lock()
			state.IncludePathsInput = joinWithSemicolon(p.Config.IncludePaths)
			state.ExcludePatternsInput = joinWithSemicolon(p.Config.ExcludePatterns)
			state.Mode = p.Config.Mode
			state.Concurrency = p.Config.Concurrency
			state.MinSizeBytes = p.Config.MinSizeBytes
			state.MaxSizeBytes = p.Config.MaxSizeBytes
			state.HashAlgorithm = p.Config.HashAlgorithm
			state.SimilarityThreshold = p.Config.SimilarityThreshold
			state.mu.Unlock()
		}
	})

	presetList := widget.NewSelect([]string{}, func(v string) { presetName.SetText(v) })
	refreshListBtn := widget.NewButton("刷新预设列表", func() {
		if names, err := core.ListScanPresets(); err == nil {
			presetList.Options = names
			presetList.Refresh()
		}
	})
	deletePresetBtn := widget.NewButton("删除选中预设", func() { _ = core.DeleteScanPreset(presetList.Selected) })

	return container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("主题", theme),
			widget.NewFormItem("语言", lang),
			widget.NewFormItem("ffmpeg 路径", ffmpegEntry),
		),
		widget.NewForm(
			widget.NewFormItem("预设名称", presetName),
			widget.NewFormItem("保存", savePresetBtn),
			widget.NewFormItem("加载", loadPresetBtn),
		),
		widget.NewForm(
			widget.NewFormItem("预设列表", presetList),
			widget.NewFormItem("刷新", refreshListBtn),
			widget.NewFormItem("删除", deletePresetBtn),
		),
	)
}

// local helpers to avoid import name clash
func themePkgDark() fyne.Theme  { return theme.DarkTheme() }
func themePkgLight() fyne.Theme { return theme.LightTheme() }

func joinWithSemicolon(items []string) string {
	if len(items) == 0 {
		return ""
	}
	out := items[0]
	for i := 1; i < len(items); i++ {
		out += ";" + items[i]
	}
	return out
}
