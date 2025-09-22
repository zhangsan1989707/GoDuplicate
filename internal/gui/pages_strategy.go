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
		Name:        t(state, "strategy_safe_delete"),
		Description: t(state, "strategy_safe_delete_desc"),
		Rule:        core.PolicyRule{KeepNewest: true},
		Action:      core.Action{Type: core.ActionDelete, DryRun: true},
	}}

	templateSelect := widget.NewSelect([]string{t(state, "strategy_safe_delete"), t(state, "strategy_move_to_dir"), t(state, "strategy_rename_suffix")}, func(v string) {
		switch v {
		case t(state, "strategy_safe_delete"):
			ui.policy = core.Policy{Name: v, Description: t(state, "strategy_safe_delete_desc_short"), Rule: core.PolicyRule{KeepNewest: true}, Action: core.Action{Type: core.ActionDelete, DryRun: true}}
		case t(state, "strategy_move_to_dir"):
			ui.policy = core.Policy{Name: v, Description: t(state, "strategy_move_to_dir_desc"), Rule: core.PolicyRule{KeepNewest: true}, Action: core.Action{Type: core.ActionMove, DestinationDir: "D:/DuplicateArchive", DryRun: true}}
		case t(state, "strategy_rename_suffix"):
			ui.policy = core.Policy{Name: v, Description: t(state, "strategy_rename_suffix_desc"), Rule: core.PolicyRule{KeepNewest: true}, Action: core.Action{Type: core.ActionRename, RenameSuffix: ".dup", DryRun: true}}
		}
	})
	templateSelect.Selected = t(state, "strategy_safe_delete")

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

	genBtn := widget.NewButton(t(state, "btn_generate_preview_plan"), func() {
		state.mu.RLock()
		groups := state.Results
		state.mu.RUnlock()
		ui.plan = core.BuildPlan(groups, ui.policy)
		onPlan(ui.plan)
		// write to shared state
		state.mu.Lock()
		state.Plan = ui.plan
		state.Logs = append(state.Logs, fmt.Sprintf(t(state, "msg_plan_generated"), len(ui.plan)))
		state.mu.Unlock()
		list.Refresh()
	})

	presetName := widget.NewEntry()
	presetName.SetPlaceHolder(t(state, "placeholder_policy_preset_name"))
	saveBtn := widget.NewButton(t(state, "btn_save_policy_preset"), func() { _, _ = core.SavePolicyPreset(core.PolicyPreset{Name: presetName.Text, Policy: ui.policy}) })
	loadBtn := widget.NewButton(t(state, "btn_load_policy_preset"), func() {
		p, err := core.LoadPolicyPreset(presetName.Text)
		if err == nil {
			ui.policy = p.Policy
		}
	})

	header := container.NewVBox(
		container.NewHBox(widget.NewLabel(t(state, "label_strategy_template")+"："), templateSelect, genBtn),
		container.NewHBox(widget.NewLabel(t(state, "label_policy_preset")+"："), presetName, saveBtn, loadBtn),
	)
	return container.NewBorder(header, nil, nil, nil, list)
}
