package core

import (
	"context"
	"fmt"
	"image"
	_ "image/png"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// GenerateVideoThumbnail extracts a frame using ffmpeg and returns a resized image.Image.
// It requires ffmpeg to be available in PATH. This is a cross-platform placeholder implementation.
func GenerateVideoThumbnail(path string, maxSide int) (image.Image, error) {
	// temp png path
	base := filepath.Base(path)
	tmp := filepath.Join(os.TempDir(), fmt.Sprintf("haste_thumb_%d_%s.png", time.Now().UnixNano(), base))
	// resolve ffmpeg binary
	bin := os.Getenv("HASTE_FFMPEG_PATH")
	if bin == "" {
		bin = "ffmpeg"
	}
	// pick 1s as default seek position; add timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin, "-y", "-ss", "00:00:01.000", "-i", path, "-frames:v", "1", "-f", "image2", "-vcodec", "png", tmp)
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	defer os.Remove(tmp)
	// load and resize via existing image thumbnail util
	return GenerateImageThumbnail(tmp, maxSide)
}

// VideoSimilarity groups videos by perceptual hash of an extracted frame.
func VideoSimilarity(files []FileInfo, threshold int) []DuplicateGroup {
	buckets := map[string][]FileInfo{}
	hashes := map[string]string{}
	for _, f := range files {
		img, err := GenerateVideoThumbnail(f.Path, 128)
		if err != nil {
			continue
		}
		h := PerceptualHash(img)
		hashes[f.Path] = h
	}
	for _, f := range files {
		h, ok := hashes[f.Path]
		if !ok {
			continue
		}
		placed := false
		for key, list := range buckets {
			if len(list) == 0 {
				continue
			}
			ref := list[0]
			dist := HammingDistanceHex(hashes[ref.Path], h)
			if dist <= threshold {
				buckets[key] = append(buckets[key], f)
				placed = true
				break
			}
		}
		if !placed {
			buckets[f.Path] = []FileInfo{f}
		}
	}
	out := make([]DuplicateGroup, 0, len(buckets))
	for key, list := range buckets {
		if len(list) < 2 {
			continue
		}
		out = append(out, DuplicateGroup{GroupID: key, Files: list})
	}
	return out
}
