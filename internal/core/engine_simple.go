package core

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// SimpleScanner is a minimal, synchronous implementation to get things working end-to-end.
// It walks include paths, applies exclude patterns, computes a fast hash (sha1) for small files,
// and groups files by identical hash. This is a placeholder to be replaced by an optimized engine.
type SimpleScanner struct{}

func NewSimpleScanner() *SimpleScanner { return &SimpleScanner{} }

func (s *SimpleScanner) Scan(config ScanConfig) ([]DuplicateGroup, error) {
	// collect all files
	var mu sync.Mutex
	files := make([]FileInfo, 0, 1024)
	var count int

	report := func(stage string) {
		if config.OnProgress != nil {
			config.OnProgress(Progress{Stage: stage, FilesScanned: count})
		}
	}

	shouldExclude := func(path string) bool {
		for _, pat := range config.ExcludePatterns {
			if pat == "" {
				continue
			}
			matched, _ := filepath.Match(pat, filepath.Base(path))
			if matched {
				return true
			}
			// allow glob on full path too
			matched, _ = filepath.Match(pat, path)
			if matched {
				return true
			}
		}
		return false
	}

	report("walking")
	walker := func(root string) error {
		return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil // skip unreadable entries
			}
			if d.IsDir() {
				return nil
			}
			if shouldExclude(path) {
				return nil
			}
			info, err := d.Info()
			if err != nil {
				return nil
			}
			// size filters
			if config.MinSizeBytes > 0 && info.Size() < config.MinSizeBytes {
				return nil
			}
			if config.MaxSizeBytes > 0 && info.Size() > config.MaxSizeBytes {
				return nil
			}
			hash := quickHash(path, info.Size())
			mu.Lock()
			files = append(files, FileInfo{
				Path:         path,
				SizeBytes:    info.Size(),
				ModifiedUnix: info.ModTime().Unix(),
				Hash:         hash,
				Type:         strings.ToLower(filepath.Ext(path)),
			})
			count++
			mu.Unlock()
			if count%200 == 0 {
				report("walking")
			}
			return nil
		})
	}

	for _, root := range config.IncludePaths {
		if root == "" {
			continue
		}
		_ = walker(root)
	}

	report("grouping")
	// grouping strategy by mode
	var groups []DuplicateGroup
	switch config.Mode {
	case "image":
		threshold := 10
		if config.SimilarityThreshold > 0 {
			bits := int((1.0 - config.SimilarityThreshold) * 64.0)
			if bits < 0 {
				bits = 0
			}
			if bits > 64 {
				bits = 64
			}
			threshold = bits
		}
		groups = MediaSimilarity(files, threshold)
	case "video":
		threshold := 10
		if config.SimilarityThreshold > 0 {
			bits := int((1.0 - config.SimilarityThreshold) * 64.0)
			if bits < 0 {
				bits = 0
			}
			if bits > 64 {
				bits = 64
			}
			threshold = bits
		}
		// filter to likely video extensions
		videoFiles := make([]FileInfo, 0, len(files))
		for _, f := range files {
			ext := strings.ToLower(f.Type)
			if ext == ".mp4" || ext == ".mov" || ext == ".avi" || ext == ".mkv" || ext == ".wmv" {
				videoFiles = append(videoFiles, f)
			}
		}
		groups = VideoSimilarity(videoFiles, threshold)
	default:
		byHash := map[string][]FileInfo{}
		for _, f := range files {
			byHash[f.Hash] = append(byHash[f.Hash], f)
		}
		for h, fs := range byHash {
			if len(fs) < 2 {
				continue
			}
			groups = append(groups, DuplicateGroup{GroupID: h, Files: fs})
		}
	}

	report("done")
	return groups, nil
}

// quickHash computes a modest hash for small portion of the file for speed.
// For now, hash the whole file up to 1MB to keep it simple.
func quickHash(path string, size int64) string {
	f, err := os.Open(path)
	if err != nil {
		return fallbackHash(path, size)
	}
	defer f.Close()
	const limit int64 = 1 << 20
	h := sha1.New()
	_, _ = io.CopyN(h, f, limit)
	sum := h.Sum(nil)
	return hex.EncodeToString(sum)
}

func fallbackHash(path string, size int64) string {
	// deterministic fallback using metadata
	h := sha1.New()
	_, _ = io.WriteString(h, path)
	_, _ = io.WriteString(h, "|")
	_, _ = fmt.Fprintf(h, "%d", size)
	return hex.EncodeToString(h.Sum(nil))
}
