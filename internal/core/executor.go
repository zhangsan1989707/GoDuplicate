package core

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ExecLogEntry represents a single dry-run/real execution log item for persistence and undo.
type ExecLogEntry struct {
	TimeUnix int64
	Action   ActionType
	Source   string
	Target   string
	Status   string // success|fail|skipped
	Message  string
}

// ExecResult wraps the execution log and potential undo info placeholder.
type ExecResult struct {
	Entries []ExecLogEntry
}

// DryRunExecute simulates executing a plan and returns logs without modifying files.
func DryRunExecute(plan []PlanItem) ExecResult {
	logs := make([]ExecLogEntry, 0, len(plan))
	now := time.Now().Unix()
	for _, p := range plan {
		logs = append(logs, ExecLogEntry{
			TimeUnix: now,
			Action:   p.Action,
			Source:   p.Source.Path,
			Target:   p.Target,
			Status:   "success",
			Message:  "dry-run",
		})
	}
	return ExecResult{Entries: logs}
}

// ConflictPolicy defines how to handle existing destination files.
type ConflictPolicy string

const (
	ConflictSkip      ConflictPolicy = "skip"
	ConflictOverwrite ConflictPolicy = "overwrite"
	ConflictRename    ConflictPolicy = "rename"
)

// ExecuteOptions controls real execution behavior.
type ExecuteOptions struct {
	DryRun         bool
	ConflictPolicy ConflictPolicy
}

// Execute performs real operations according to the plan. Best-effort; continues on error, logging each.
func Execute(plan []PlanItem, opts ExecuteOptions) ExecResult {
	if opts.DryRun {
		return DryRunExecute(plan)
	}
	logs := make([]ExecLogEntry, 0, len(plan))
	for _, p := range plan {
		entry := ExecLogEntry{TimeUnix: time.Now().Unix(), Action: p.Action, Source: p.Source.Path, Target: p.Target}
		var err error
		switch p.Action {
		case ActionDelete:
			err = os.Remove(p.Source.Path)
		case ActionMove:
			if p.Target == "" {
				err = os.ErrInvalid
				break
			}
			_ = os.MkdirAll(p.Target, 0o755)
			targetPath := filepath.Join(p.Target, filepath.Base(p.Source.Path))
			targetPath, err = resolveConflict(targetPath, opts.ConflictPolicy)
			if err != nil && err != os.ErrExist {
				break
			}
			if err == os.ErrExist { // skipped
				entry.Status = "skipped"
				entry.Message = "conflict: exists"
				logs = append(logs, entry)
				continue
			}
			err = os.Rename(p.Source.Path, targetPath)
			if err == nil {
				entry.Target = targetPath
			}
		case ActionCopy:
			if p.Target == "" {
				err = os.ErrInvalid
				break
			}
			_ = os.MkdirAll(p.Target, 0o755)
			targetPath := filepath.Join(p.Target, filepath.Base(p.Source.Path))
			targetPath, err = resolveConflict(targetPath, opts.ConflictPolicy)
			if err != nil && err != os.ErrExist {
				break
			}
			if err == os.ErrExist {
				entry.Status = "skipped"
				entry.Message = "conflict: exists"
				logs = append(logs, entry)
				continue
			}
			err = copyFile(p.Source.Path, targetPath)
			if err == nil {
				entry.Target = targetPath
			}
		case ActionRename:
			if p.Target == "" {
				err = os.ErrInvalid
				break
			}
			// apply conflict on rename too
			targetPath, e2 := resolveConflict(p.Target, opts.ConflictPolicy)
			if e2 != nil && e2 != os.ErrExist {
				err = e2
				break
			}
			if e2 == os.ErrExist {
				entry.Status = "skipped"
				entry.Message = "conflict: exists"
				logs = append(logs, entry)
				continue
			}
			err = os.Rename(p.Source.Path, targetPath)
			if err == nil {
				entry.Target = targetPath
			}
		default:
			entry.Status = "skipped"
			entry.Message = "unsupported action"
			logs = append(logs, entry)
			continue
		}
		if err != nil {
			entry.Status = "fail"
			entry.Message = err.Error()
		} else {
			entry.Status = "success"
		}
		logs = append(logs, entry)
	}
	return ExecResult{Entries: logs}
}

// resolveConflict returns a path to use according to policy.
// If policy=skip and exists, returns os.ErrExist.
// If policy=overwrite, removes existing file.
// If policy=rename, returns a unique suffixed filename.
func resolveConflict(target string, policy ConflictPolicy) (string, error) {
	if _, err := os.Stat(target); err != nil {
		return target, nil // not exists
	}
	switch policy {
	case ConflictSkip:
		return target, os.ErrExist
	case ConflictOverwrite:
		if err := os.Remove(target); err != nil {
			return target, err
		}
		return target, nil
	case ConflictRename:
		dir := filepath.Dir(target)
		base := filepath.Base(target)
		dot := strings.LastIndexByte(base, '.')
		name, ext := base, ""
		if dot > 0 {
			name, ext = base[:dot], base[dot:]
		}
		for i := 1; i < 10000; i++ {
			cand := filepath.Join(dir, name+" ("+fmt.Sprintf("%d", i)+")"+ext)
			if _, err := os.Stat(cand); os.IsNotExist(err) {
				return cand, nil
			}
		}
		return target, os.ErrExist
	default:
		return target, nil
	}
}

func copyFile(src, dst string) error {
	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()
	df, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer df.Close()
	_, err = io.Copy(df, sf)
	return err
}

// PersistExecLog saves the execution log to a JSON file in temp dir and returns the path.
func PersistExecLog(res ExecResult) (string, error) {
	dir := filepath.Join(os.TempDir(), "haste_logs")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	name := filepath.Join(dir, time.Now().Format("20060102_150405")+"_exec.json")
	f, err := os.Create(name)
	if err != nil {
		return "", err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(res); err != nil {
		return "", err
	}
	return name, nil
}

// Undo best-effort based on a previous ExecResult (reverse actions where possible).
func Undo(res ExecResult) ExecResult {
	logs := make([]ExecLogEntry, 0, len(res.Entries))
	for i := len(res.Entries) - 1; i >= 0; i-- {
		e := res.Entries[i]
		entry := ExecLogEntry{TimeUnix: time.Now().Unix(), Action: e.Action, Source: e.Source, Target: e.Target}
		var err error
		switch e.Action {
		case ActionMove:
			if e.Target != "" {
				err = os.Rename(e.Target, e.Source)
			}
		case ActionRename:
			if e.Target != "" {
				err = os.Rename(e.Target, e.Source)
			}
		case ActionCopy:
			if e.Target != "" {
				err = os.Remove(e.Target)
			}
		case ActionDelete:
			entry.Status = "skipped"
			entry.Message = "cannot undo delete"
			logs = append(logs, entry)
			continue
		default:
			entry.Status = "skipped"
			entry.Message = "unsupported"
			logs = append(logs, entry)
			continue
		}
		if err != nil {
			entry.Status = "fail"
			entry.Message = err.Error()
		} else {
			entry.Status = "success"
		}
		logs = append(logs, entry)
	}
	return ExecResult{Entries: logs}
}
