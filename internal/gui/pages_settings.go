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
			fyne.CurrentApp().Settings().SetTheme(newChineseTheme(themePkgDark()))
		} else {
			fyne.CurrentApp().Settings().SetTheme(newChineseTheme(themePkgLight()))
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
		// 触发语言变更回调，刷新所有界面元素
		state.TriggerLanguageChangedCallbacks()
	})
	lang.Selected = state.Language
	ffmpegEntry := widget.NewEntry()
	ffmpegEntry.SetPlaceHolder(t(state, "placeholder_ffmpeg_path"))
	state.mu.RLock()
	ffmpegEntry.SetText(state.FfmpegPath)
	state.mu.RUnlock()
	ffmpegEntry.OnChanged = func(v string) { state.mu.Lock(); state.FfmpegPath = v; state.mu.Unlock() }

	presetName := widget.NewEntry()
	presetName.SetPlaceHolder(t(state, "placeholder_preset_name"))
	savePresetBtn := widget.NewButton(t(state, "btn_save_preset"), func() {
		state.mu.RLock()
		cfg := state.ToScanConfig()
		state.mu.RUnlock()
		_, _ = core.SaveScanPreset(core.ScanPreset{Name: presetName.Text, Config: cfg})
	})
	loadPresetBtn := widget.NewButton(t(state, "btn_load_preset"), func() {
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
	refreshListBtn := widget.NewButton(t(state, "btn_refresh_presets"), func() {
		if names, err := core.ListScanPresets(); err == nil {
			presetList.Options = names
			presetList.Refresh()
		}
	})
	deletePresetBtn := widget.NewButton(t(state, "btn_delete_preset"), func() { _ = core.DeleteScanPreset(presetList.Selected) })

	return container.NewVBox(
		widget.NewForm(
			widget.NewFormItem(t(state, "label_theme"), theme),
			widget.NewFormItem(t(state, "label_language"), lang),
			widget.NewFormItem(t(state, "label_ffmpeg_path"), ffmpegEntry),
		),
		widget.NewForm(
			widget.NewFormItem(t(state, "label_preset_name"), presetName),
			widget.NewFormItem(t(state, "label_save"), savePresetBtn),
			widget.NewFormItem(t(state, "label_load"), loadPresetBtn),
		),
		widget.NewForm(
			widget.NewFormItem(t(state, "label_preset_list"), presetList),
			widget.NewFormItem(t(state, "label_refresh"), refreshListBtn),
			widget.NewFormItem(t(state, "label_delete"), deletePresetBtn),
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
