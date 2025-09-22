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

// 创建本地化函数，用于获取正确的翻译文本
func t(state *AppState, key string) string {
	state.mu.RLock()
	lang := state.Language
	state.mu.RUnlock()
	
	// 如果是中文环境，直接返回中文文本
	if lang == "zh-CN" {
		switch key {
		case "tab_config":
			return "扫描配置"
		case "tab_monitor":
			return "扫描监控"
		case "tab_results":
			return "结果展示"
		case "tab_strategy":
			return "处理策略"
		case "tab_execute":
			return "执行处理"
		case "tab_settings":
			return "系统设置"
		case "label_sort":
			return "排序:"
		case "label_thumbwall":
			return "缩略图预览:"
		case "status_scanning":
			return "状态: 扫描中"
		case "status_idle":
			return "状态: 空闲"
		case "label_files":
			return "已扫描文件:"
		case "label_groups":
			return "发现重复组:"
		case "label_speed":
			return "速度:"
		case "placeholder_include":
			return "输入要扫描的路径，多个路径用;分隔"
		case "placeholder_exclude":
			return "输入要排除的模式，多个模式用;分隔"
		case "placeholder_min_size":
			return "最小文件大小(字节)"
		case "placeholder_max_size":
			return "最大文件大小(字节)"
		case "label_concurrency":
			return "并发度: %d"
		case "label_similarity":
			return "相似度阈值: %.2f"
		case "btn_start_scan":
			return "开始扫描"
		case "btn_pick_dir":
			return "选择目录"
		case "sort_default":
			return "默认排序"
		case "sort_files_desc":
			return "文件数量降序"
		case "sort_size_desc":
			return "文件大小降序"
		case "sort_similarity_desc":
			return "相似度降序"
		case "msg_select_group_for_details":
			return "选择一个组查看详情"
		case "msg_thumbnail_generation_failed":
			return "生成缩略图失败"
		case "msg_generating":
			return "生成中..."
		case "placeholder_ffmpeg_path":
			return "可选：手动指定 ffmpeg 可执行路径"
		case "placeholder_preset_name":
			return "预设名称"
		case "btn_save_preset":
			return "保存扫描预设"
		case "btn_load_preset":
			return "加载扫描预设"
		case "btn_refresh_presets":
			return "刷新预设列表"
		case "btn_delete_preset":
			return "删除选中预设"
		case "label_theme":
			return "主题"
		case "label_language":
			return "语言"
		case "label_ffmpeg_path":
			return "ffmpeg 路径"
		case "label_preset_name":
			return "预设名称"
		case "label_save":
			return "保存"
		case "label_load":
			return "加载"
		case "label_preset_list":
			return "预设列表"
		case "label_refresh":
			return "刷新"
		case "label_delete":
			return "删除"
		default:
			return key
		}
	}
	// 英文环境返回默认文本
	return key
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
