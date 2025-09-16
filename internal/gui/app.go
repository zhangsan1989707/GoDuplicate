package gui

import (
	"os"

	"goduplicate/internal/core"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

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
			a.Settings().SetTheme(theme.DarkTheme())
		} else {
			a.Settings().SetTheme(theme.LightTheme())
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
