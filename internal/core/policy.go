package core

// ActionType enumerates supported file operations for duplicates handling.
type ActionType string

const (
    ActionDelete    ActionType = "delete"   // 删除文件（占位）
    ActionRecycle   ActionType = "recycle"  // 移至回收站（占位）
    ActionMove      ActionType = "move"     // 移动到目录
    ActionCopy      ActionType = "copy"     // 复制到目录
    ActionRename    ActionType = "rename"   // 重命名/加后缀
    ActionMark      ActionType = "mark"     // 标记（元数据/DB，占位）
)

// PolicyRule defines one rule used to decide which files to keep or operate.
type PolicyRule struct {
    KeepNewest      bool
    KeepOldest      bool
    KeepShortestDir bool
    // more: by path contains, by extension, etc.
}

// Policy defines a high level strategy template.
type Policy struct {
    Name        string
    Description string
    Rule        PolicyRule
    Action      Action
}

// Action holds parameters for an operation to apply on selected files.
type Action struct {
    Type           ActionType
    DestinationDir string // for move/copy
    RenameSuffix   string // for rename
    DryRun         bool   // preview only
}

// PlanItem represents a single file operation in preview/execution.
type PlanItem struct {
    GroupID string
    Source  FileInfo
    Target  string // path or new name
    Action  ActionType
}

// BuildPlan creates a naive plan: keep first in each group, operate others by policy.Action.
func BuildPlan(groups []DuplicateGroup, p Policy) []PlanItem {
    var plan []PlanItem
    for _, g := range groups {
        if len(g.Files) <= 1 {
            continue
        }
        keeperIdx := 0 // TODO: apply Rule
        for i, f := range g.Files {
            if i == keeperIdx {
                continue
            }
            var target string
            switch p.Action.Type {
            case ActionMove, ActionCopy:
                target = p.Action.DestinationDir
            case ActionRename:
                target = f.Path + p.Action.RenameSuffix
            default:
                target = ""
            }
            plan = append(plan, PlanItem{
                GroupID: g.GroupID,
                Source:  f,
                Target:  target,
                Action:  p.Action.Type,
            })
        }
    }
    return plan
}


