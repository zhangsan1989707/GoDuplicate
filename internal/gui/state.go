package gui

import (
	"image"
	"sync"

	"goduplicate/internal/core"
)

// AppState holds shared state between GUI tabs.
type AppState struct {
	mu sync.RWMutex

	// Live configuration fields bound to the config form
	IncludePathsInput    string // semicolon separated
	ExcludePatternsInput string // semicolon separated
	Mode                 string
	Concurrency          int
	MinSizeBytes         int64
	MaxSizeBytes         int64
	HashAlgorithm        string
	SimilarityThreshold  float64

	// Scan results and stats
	Results       []core.DuplicateGroup
	LastScanError error

	// Strategy & execution
	Plan []core.PlanItem
	Logs []string

	// Monitoring snapshot
	FilesScanned int
	GroupsFound  int
	IsScanning   bool

	// Settings (placeholder)
	Theme      string // light|dark
	Language   string // zh-CN|en-US
	FfmpegPath string // optional custom ffmpeg path

	// Caches
	ThumbCache map[string]image.Image
}

func NewAppState() *AppState {
	return &AppState{
		Mode:                "basic",
		Concurrency:         4,
		HashAlgorithm:       "sha1",
		SimilarityThreshold: 0.85,
		Theme:               "light",
		Language:            "zh-CN",
		ThumbCache:          make(map[string]image.Image),
	}
}

func (s *AppState) ToScanConfig() core.ScanConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return core.ScanConfig{
		IncludePaths:        splitSemicolon(s.IncludePathsInput),
		ExcludePatterns:     splitSemicolon(s.ExcludePatternsInput),
		Mode:                s.Mode,
		Concurrency:         s.Concurrency,
		MinSizeBytes:        s.MinSizeBytes,
		MaxSizeBytes:        s.MaxSizeBytes,
		HashAlgorithm:       s.HashAlgorithm,
		SimilarityThreshold: s.SimilarityThreshold,
	}
}

func splitSemicolon(v string) []string {
	out := make([]string, 0, 8)
	start := 0
	for i := 0; i <= len(v); i++ {
		if i == len(v) || v[i] == ';' {
			token := v[start:i]
			// trim spaces
			token = trimSpaces(token)
			if token != "" {
				out = append(out, token)
			}
			start = i + 1
		}
	}
	return out
}

func trimSpaces(s string) string {
	// simple local trim to avoid bringing strings here
	i := 0
	j := len(s)
	for i < j && (s[i] == ' ' || s[i] == '\t' || s[i] == '\n' || s[i] == '\r') {
		i++
	}
	for j > i && (s[j-1] == ' ' || s[j-1] == '\t' || s[j-1] == '\n' || s[j-1] == '\r') {
		j--
	}
	return s[i:j]
}
