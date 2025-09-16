package gui

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"goduplicate/internal/core"
)

type execState struct {
	plan       []core.PlanItem
	logs       []string
	lastResult *core.ExecResult
}

func buildExecutePage(state *AppState) fyne.CanvasObject {
	st := &execState{}
	logList := widget.NewList(
		func() int { return len(st.logs) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if i >= 0 && i < len(st.logs) {
				o.(*widget.Label).SetText(st.logs[i])
			}
		},
	)

	refreshPlan := func() {
		state.mu.RLock()
		st.plan = append([]core.PlanItem(nil), state.Plan...)
		state.mu.RUnlock()
		st.logs = append(st.logs, fmt.Sprintf("载入计划条目: %d", len(st.plan)))
		logList.Refresh()
	}

	dryRun := widget.NewCheck("Dry-Run", func(bool) {})
	dryRun.SetChecked(true)

	policySelect := widget.NewSelect([]string{"skip", "overwrite", "rename"}, nil)
	policySelect.Selected = "rename"

	previewBtn := widget.NewButton("刷新预览", func() { refreshPlan() })
	executeBtn := widget.NewButton("执行并保存日志", func() {
		if len(st.plan) == 0 {
			st.logs = append(st.logs, "无可执行项")
			logList.Refresh()
			return
		}
		opts := core.ExecuteOptions{DryRun: dryRun.Checked, ConflictPolicy: core.ConflictPolicy(policySelect.Selected)}
		res := core.Execute(st.plan, opts)
		st.lastResult = &res
		path, err := core.PersistExecLog(res)
		if err != nil {
			st.logs = append(st.logs, fmt.Sprintf("保存日志失败: %v", err))
		} else {
			st.logs = append(st.logs, fmt.Sprintf("日志已保存: %s", path))
		}
		logList.Refresh()
	})

	undoBtn := widget.NewButton("撤销上次", func() {
		if st.lastResult == nil {
			st.logs = append(st.logs, "无可撤销记录")
			logList.Refresh()
			return
		}
		res := core.Undo(*st.lastResult)
		path, err := core.PersistExecLog(res)
		if err != nil {
			st.logs = append(st.logs, fmt.Sprintf("撤销日志保存失败: %v", err))
		} else {
			st.logs = append(st.logs, fmt.Sprintf("撤销日志已保存: %s", path))
		}
		logList.Refresh()
	})

	openLogsBtn := widget.NewButton("打开日志目录", func() {
		dir := os.TempDir()
		st.logs = append(st.logs, fmt.Sprintf("日志目录: %s", dir))
		logList.Refresh()
	})

	controls := container.NewHBox(dryRun, widget.NewLabel("冲突策略:"), policySelect, previewBtn, executeBtn, undoBtn, openLogsBtn)
	return container.NewBorder(controls, nil, nil, nil, logList)
}
