package core

// ScanConfig represents user-configurable parameters for a scan session.
type ScanConfig struct {
	IncludePaths    []string
	ExcludePatterns []string
	Mode            string // basic | video | text | image
	Concurrency     int
	// Filters
	MinSizeBytes int64 // 0 = no min
	MaxSizeBytes int64 // 0 = no max
	// Hashing / similarity
	HashAlgorithm       string  // sha1 | sha256 | md5 (占位)
	SimilarityThreshold float64 // 0.0-1.0 (媒体模式占位)
	// Optional progress callback
	OnProgress func(Progress)
	// Future: hash algorithm, similarity threshold, size filters, presets
}

// FileInfo represents a single file discovered by the scanner.
type FileInfo struct {
	Path         string
	SizeBytes    int64
	ModifiedUnix int64
	Hash         string
	Type         string // mime or coarse type
}

// DuplicateGroup represents a logical group of duplicate files.
type DuplicateGroup struct {
	GroupID string
	Files   []FileInfo
}

// Progress provides lightweight telemetry from scanner to UI/CLI.
type Progress struct {
	Stage        string // "walking" | "hashing" | "grouping" | "done"
	FilesScanned int
	GroupsFound  int
}
