package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"goduplicate/internal/core"
)

type strategyUI struct {
	policy core.Policy
	plan   []core.PlanItem
}

func buildStrategyPage(state *AppState, onPlan func([]core.PlanItem)) fyne.CanvasObject {
	ui := &strategyUI{policy: core.Policy{
		Name:        "安全删除",
		Description: "每组保留一个，删除其他（预览模式）",
		Rule:        core.PolicyRule{KeepNewest: true},
		Action:      core.Action{Type: core.ActionDelete, DryRun: true},
	}}

	templateSelect := widget.NewSelect([]string{"安全删除", "移动到目录", "重命名加后缀"}, func(v string) {
		switch v {
		case "安全删除":
			ui.policy = core.Policy{Name: v, Description: "每组保留一个，删除其他（预览）", Rule: core.PolicyRule{KeepNewest: true}, Action: core.Action{Type: core.ActionDelete, DryRun: true}}
		case "移动到目录":
			ui.policy = core.Policy{Name: v, Description: "将重复文件移动到指定目录（预览）", Rule: core.PolicyRule{KeepNewest: true}, Action: core.Action{Type: core.ActionMove, DestinationDir: "D:/DuplicateArchive", DryRun: true}}
		case "重命名加后缀":
			ui.policy = core.Policy{Name: v, Description: "为重复文件添加 .dup 后缀（预览）", Rule: core.PolicyRule{KeepNewest: true}, Action: core.Action{Type: core.ActionRename, RenameSuffix: ".dup", DryRun: true}}
		}
	})
	templateSelect.Selected = "安全删除"

	list := widget.NewList(
		func() int { return len(ui.plan) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if i < 0 || i >= len(ui.plan) {
				return
			}
			p := ui.plan[i]
			o.(*widget.Label).SetText(fmt.Sprintf("[%s] %s -> %s", p.Action, p.Source.Path, p.Target))
		},
	)

	genBtn := widget.NewButton("生成预览计划", func() {
		state.mu.RLock()
		groups := state.Results
		state.mu.RUnlock()
		ui.plan = core.BuildPlan(groups, ui.policy)
		onPlan(ui.plan)
		// write to shared state
		state.mu.Lock()
		state.Plan = ui.plan
		state.Logs = append(state.Logs, fmt.Sprintf("生成计划: %d 项", len(ui.plan)))
		state.mu.Unlock()
		list.Refresh()
	})

	presetName := widget.NewEntry()
	presetName.SetPlaceHolder("策略预设名称")
	saveBtn := widget.NewButton("保存策略预设", func() { _, _ = core.SavePolicyPreset(core.PolicyPreset{Name: presetName.Text, Policy: ui.policy}) })
	loadBtn := widget.NewButton("加载策略预设", func() {
		p, err := core.LoadPolicyPreset(presetName.Text)
		if err == nil {
			ui.policy = p.Policy
		}
	})

	header := container.NewVBox(
		container.NewHBox(widget.NewLabel("策略模板:"), templateSelect, genBtn),
		container.NewHBox(widget.NewLabel("策略预设:"), presetName, saveBtn, loadBtn),
	)
	return container.NewBorder(header, nil, nil, nil, list)
}
