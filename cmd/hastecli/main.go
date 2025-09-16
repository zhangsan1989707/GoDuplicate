package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"goduplicate/internal/core"
)

// entrypoint for CLI mode. Implements minimal argument parsing and delegates to core engine later.
func main() {
	var includePathsArg string
	var excludePatternsArg string
	var mode string
	var concurrency int
	var minSize int64
	var maxSize int64
	var hashAlg string
	var sim float64

	flag.StringVar(&includePathsArg, "paths", "", "要扫描的路径，使用;分隔多个路径")
	flag.StringVar(&excludePatternsArg, "exclude", "", "排除的通配符模式，使用;分隔")
	flag.StringVar(&mode, "mode", "basic", "扫描模式：basic|video|text|image")
	flag.IntVar(&concurrency, "concurrency", 4, "并发度")
	flag.Int64Var(&minSize, "min-size", 0, "最小文件大小(字节)")
	flag.Int64Var(&maxSize, "max-size", 0, "最大文件大小(字节，0为不限)")
	flag.StringVar(&hashAlg, "hash", "sha1", "哈希算法：sha1|sha256|md5(占位)")
	flag.Float64Var(&sim, "similarity", 0.0, "相似度阈值(0.0-1.0，占位)")
	flag.Parse()

	if includePathsArg == "" {
		fmt.Println("请使用 --paths 指定至少一个路径（使用;分隔）")
		os.Exit(2)
	}

	includePaths := splitAndTrim(includePathsArg)
	excludePatterns := splitAndTrim(excludePatternsArg)

	cfg := core.ScanConfig{
		IncludePaths:        includePaths,
		ExcludePatterns:     excludePatterns,
		Mode:                strings.ToLower(mode),
		Concurrency:         concurrency,
		MinSizeBytes:        minSize,
		MaxSizeBytes:        maxSize,
		HashAlgorithm:       strings.ToLower(hashAlg),
		SimilarityThreshold: sim,
	}

	engine := core.NewSimpleScanner()
	groups, err := engine.Scan(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "扫描失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("发现重复组数: %d\n", len(groups))
	for i, g := range groups {
		if i >= 10 {
			fmt.Println("...更多结果已省略")
			break
		}
		fmt.Printf("组 %d (id=%s, 文件数=%d)\n", i+1, g.GroupID[:8], len(g.Files))
	}
}

func splitAndTrim(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ";")
	var out []string
	for _, p := range parts {
		t := strings.TrimSpace(p)
		if t != "" {
			out = append(out, t)
		}
	}
	return out
}
