package gui

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"goduplicate/internal/core"
)

// buildConfigPage returns the configuration form tab
func buildConfigPage(state *AppState, onStart func(cfg core.ScanConfig)) fyne.CanvasObject {
	includeEntry := widget.NewEntry()
	includeEntry.SetPlaceHolder("示例: D\\;E\\docs")
	includeEntry.OnChanged = func(v string) {
		state.mu.Lock()
		state.IncludePathsInput = v
		state.mu.Unlock()
	}

	excludeEntry := widget.NewEntry()
	excludeEntry.SetPlaceHolder("示例: *.tmp;node_modules;*.bak")
	excludeEntry.OnChanged = func(v string) {
		state.mu.Lock()
		state.ExcludePatternsInput = v
		state.mu.Unlock()
	}

	modeSelect := widget.NewSelect([]string{"basic", "video", "text", "image"}, func(v string) {
		state.mu.Lock()
		state.Mode = v
		state.mu.Unlock()
	})
	modeSelect.Selected = state.Mode

	hashSelect := widget.NewSelect([]string{"sha1", "sha256", "md5"}, func(v string) {
		state.mu.Lock()
		state.HashAlgorithm = v
		state.mu.Unlock()
	})
	hashSelect.Selected = state.HashAlgorithm

	minEntry := widget.NewEntry()
	minEntry.SetPlaceHolder("最小大小(字节)")
	minEntry.OnChanged = func(v string) {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			state.mu.Lock()
			state.MinSizeBytes = n
			state.mu.Unlock()
		}
	}
	maxEntry := widget.NewEntry()
	maxEntry.SetPlaceHolder("最大大小(字节,0不限)")
	maxEntry.OnChanged = func(v string) {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			state.mu.Lock()
			state.MaxSizeBytes = n
			state.mu.Unlock()
		}
	}

	concurrency := widget.NewSlider(1, 16)
	concurrency.Step = 1
	cLabel := widget.NewLabel(fmt.Sprintf("并发度: %d", int(concurrency.Value)))
	concurrency.OnChanged = func(v float64) {
		iv := int(v)
		cLabel.SetText(fmt.Sprintf("并发度: %d", iv))
		state.mu.Lock()
		state.Concurrency = iv
		state.mu.Unlock()
	}
	concurrency.SetValue(float64(state.Concurrency))

	// similarity slider 0.50~0.99
	simSlider := widget.NewSlider(0.5, 0.99)
	simSlider.Step = 0.01
	simLabel := widget.NewLabel(fmt.Sprintf("相似度阈值: %.2f", state.SimilarityThreshold))
	simSlider.OnChanged = func(v float64) {
		state.mu.Lock()
		state.SimilarityThreshold = v
		state.mu.Unlock()
		simLabel.SetText(fmt.Sprintf("相似度阈值: %.2f", v))
	}
	simSlider.SetValue(state.SimilarityThreshold)

	startBtn := widget.NewButton("开始扫描", func() {
		onStart(state.ToScanConfig())
	})

	pickDirBtn := widget.NewButton("选择目录", func() {
		dialog.ShowFolderOpen(func(u fyne.ListableURI, err error) {
			if err == nil && u != nil {
				path := u.Path()
				state.mu.Lock()
				if state.IncludePathsInput == "" {
					state.IncludePathsInput = path
				} else {
					state.IncludePathsInput += ";" + path
				}
				state.mu.Unlock()
				includeEntry.SetText(state.IncludePathsInput)
			}
		}, fyne.CurrentApp().Driver().AllWindows()[0])
	})

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "扫描路径(;)分隔", Widget: includeEntry},
			{Text: "排除模式(;)分隔", Widget: excludeEntry},
			{Text: "模式", Widget: modeSelect},
			{Text: "哈希算法", Widget: hashSelect},
			{Text: "最小大小", Widget: minEntry},
			{Text: "最大大小", Widget: maxEntry},
			{Text: "并发度", Widget: container.NewHBox(concurrency, cLabel)},
			{Text: "相似度", Widget: container.NewHBox(simSlider, simLabel)},
		},
		OnSubmit: func() { onStart(state.ToScanConfig()) },
	}
	return container.NewBorder(nil, container.NewHBox(pickDirBtn, startBtn), nil, nil, form)
}

// buildMonitorPage shows minimal statistics
func buildMonitorPage(state *AppState) fyne.CanvasObject {
	files := widget.NewLabel(fmt.Sprintf("%s 0", t(state, "label_files")))
	groups := widget.NewLabel(fmt.Sprintf("%s 0", t(state, "label_groups")))
	status := widget.NewLabel(t(state, "status_idle"))
	speed := widget.NewLabel(fmt.Sprintf("%s -", t(state, "label_speed")))

	box := container.NewVBox(files, groups, status, speed)

	go func() {
		var lastFiles int
		var lastTime = time.Now()
		for {
			time.Sleep(500 * time.Millisecond)
			state.mu.RLock()
			f := state.FilesScanned
			g := state.GroupsFound
			scanning := state.IsScanning
			state.mu.RUnlock()

			files.SetText(fmt.Sprintf("%s %d", t(state, "label_files"), f))
			groups.SetText(fmt.Sprintf("%s %d", t(state, "label_groups"), g))
			if scanning {
				status.SetText(t(state, "status_scanning"))
			} else {
				status.SetText(t(state, "status_idle"))
			}

			now := time.Now()
			dt := now.Sub(lastTime).Seconds()
			if dt > 0 {
				spd := float64(f-lastFiles) / dt
				speed.SetText(fmt.Sprintf("%s %.1f 文件/秒", t(state, "label_speed"), spd))
			}
			lastFiles = f
			lastTime = now
		}
	}()

	return box
}

// buildResultsPage lists duplicate groups with sorting and details
func buildResultsPage(state *AppState) fyne.CanvasObject {
	sortSelect := widget.NewSelect([]string{"默认", "按文件数降序", "按大小降序", "按相似度降序"}, nil)
	sortSelect.Selected = "默认"

	groupTitle := widget.NewLabel("选择一个重复组以查看详情")
	var thumbImg *canvas.Image
	filesList := widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.ListItemID, o fyne.CanvasObject) {},
	)

	groupsList := widget.NewList(
		func() int { state.mu.RLock(); defer state.mu.RUnlock(); return len(state.Results) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			state.mu.RLock()
			defer state.mu.RUnlock()
			if i < 0 || i >= len(state.Results) {
				return
			}
			g := state.Results[i]
			sim := core.EstimateGroupSimilarity(g.Files)
			o.(*widget.Label).SetText(fmt.Sprintf("组 %d | 文件数 %d | 相似度 %.0f%%", i+1, len(g.Files), sim))
		},
	)

	thumbGrid := container.NewGridWrap(fyne.NewSize(96, 96))
	thumbError := widget.NewLabel("")

	refreshLists := func() {
		state.mu.Lock()
		switch sortSelect.Selected {
		case "按文件数降序":
			sort.Slice(state.Results, func(i, j int) bool { return len(state.Results[i].Files) > len(state.Results[j].Files) })
		case "按大小降序":
			size := func(g core.DuplicateGroup) int64 {
				var s int64
				for _, f := range g.Files {
					s += f.SizeBytes
				}
				return s
			}
			sort.Slice(state.Results, func(i, j int) bool { return size(state.Results[i]) > size(state.Results[j]) })
		case "按相似度降序":
			sort.Slice(state.Results, func(i, j int) bool {
				si := core.EstimateGroupSimilarity(state.Results[i].Files)
				sj := core.EstimateGroupSimilarity(state.Results[j].Files)
				return si > sj
			})
		}
		state.mu.Unlock()
		groupsList.Refresh()
	}
	sortSelect.OnChanged = func(string) { refreshLists() }

	groupsList.OnSelected = func(id widget.ListItemID) {
		state.mu.RLock()
		if id < 0 || id >= len(state.Results) {
			state.mu.RUnlock()
			return
		}
		g := state.Results[id]
		state.mu.RUnlock()
		sim := core.EstimateGroupSimilarity(g.Files)
		groupTitle.SetText(fmt.Sprintf("组 %d 详情： 相似度≈%.0f%%", id+1, sim))
		files := g.Files
		filesList.Length = func() int { return len(files) }
		filesList.UpdateItem = func(i widget.ListItemID, o fyne.CanvasObject) {
			if i < 0 || i >= len(files) {
				return
			}
			f := files[i]
			o.(*widget.Label).SetText(fmt.Sprintf("%s | %dB", f.Path, f.SizeBytes))
		}
		filesList.Refresh()
		thumbError.SetText("")
		if len(files) > 0 {
			path := files[0].Path
			state.mu.RLock()
			img := state.ThumbCache[path]
			state.mu.RUnlock()
			if img == nil {
				if timg, err := core.GetMediaThumbnail(path, 160); err == nil {
					state.mu.Lock()
					state.ThumbCache[path] = timg
					state.mu.Unlock()
					img = timg
				} else {
					thumbError.SetText("缩略图生成失败")
				}
			}
			if img != nil {
				thumbImg = canvas.NewImageFromImage(img)
				thumbImg.SetMinSize(fyne.NewSize(160, 160))
			}
		}
		thumbGrid.Objects = nil
		sem := make(chan struct{}, 4)
		for _, f := range files {
			p := f.Path
			state.mu.RLock()
			ti := state.ThumbCache[p]
			state.mu.RUnlock()
			if ti != nil {
				thumbGrid.Add(canvas.NewImageFromImage(ti))
				continue
			}
			thumbGrid.Add(widget.NewLabel("生成中…"))
			go func(path string) {
				sem <- struct{}{}
				defer func() { <-sem }()
				if timg, err := core.GetMediaThumbnail(path, 96); err == nil && timg != nil {
					state.mu.Lock()
					state.ThumbCache[path] = timg
					state.mu.Unlock()
					thumbGrid.Refresh()
				}
			}(p)
		}
		thumbGrid.Refresh()
	}

	thumbBox := container.NewMax()
	if thumbImg != nil {
		thumbBox.Add(thumbImg)
	}

	left := container.NewBorder(container.NewHBox(widget.NewLabel(t(state, "label_sort")), sortSelect), nil, nil, nil, groupsList)
	right := container.NewBorder(container.NewVBox(groupTitle, thumbBox, thumbError, widget.NewLabel(t(state, "label_thumbwall")), thumbGrid), nil, nil, nil, filesList)
	return container.NewHSplit(left, right)
}

func filepathExt(p string) string {
	for i := len(p) - 1; i >= 0; i-- {
		if p[i] == '.' {
			return p[i:]
		}
		if p[i] == '/' || p[i] == '\\' {
			break
		}
	}
	return ""
}
