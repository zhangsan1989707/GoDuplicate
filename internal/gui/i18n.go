package gui

// Minimal i18n utility for two languages zh-CN/en-US

var zhCN = map[string]string{
	"tab_config":      "扫描配置",
	"tab_monitor":     "扫描监控",
	"tab_results":     "结果展示",
	"tab_strategy":    "处理策略",
	"tab_execute":     "执行处理",
	"tab_settings":    "系统设置",
	"label_sort":      "排序:",
	"label_thumbwall": "缩略图预览:",
	"status_scanning": "状态: 扫描中",
	"status_idle":     "状态: 空闲",
	"label_files":     "已扫描文件:",
	"label_groups":    "发现重复组:",
	"label_speed":     "速度:",
}

var enUS = map[string]string{
	"tab_config":      "Config",
	"tab_monitor":     "Monitor",
	"tab_results":     "Results",
	"tab_strategy":    "Strategy",
	"tab_execute":     "Execute",
	"tab_settings":    "Settings",
	"label_sort":      "Sort:",
	"label_thumbwall": "Thumbnails:",
	"status_scanning": "Status: Scanning",
	"status_idle":     "Status: Idle",
	"label_files":     "Files scanned:",
	"label_groups":    "Duplicate groups:",
	"label_speed":     "Speed:",
}

func t(state *AppState, key string) string {
	state.mu.RLock()
	lang := state.Language
	state.mu.RUnlock()
	if lang == "zh-CN" {
		if v, ok := zhCN[key]; ok {
			return v
		}
	} else {
		if v, ok := enUS[key]; ok {
			return v
		}
	}
	return key
}
