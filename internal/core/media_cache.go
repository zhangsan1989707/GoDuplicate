package core

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

// thumbnailCacheDir returns a directory for storing thumbnails.
func thumbnailCacheDir() string {
	base := os.TempDir()
	return filepath.Join(base, "haste_thumbs")
}

// ensureDir ensures a directory exists.
func ensureDir(dir string) { _ = os.MkdirAll(dir, 0o755) }

// ThumbnailCachePath computes a cache file path for an original path and size.
func ThumbnailCachePath(originalPath string, maxSide int) string {
	name := strings.ReplaceAll(originalPath, string(os.PathSeparator), "__")
	name = strings.ReplaceAll(name, ":", "_")
	return filepath.Join(thumbnailCacheDir(), name+"_"+itoa(maxSide)+".png")
}

// LoadThumbnail tries to load a cached thumbnail from disk.
func LoadThumbnail(originalPath string, maxSide int) (image.Image, error) {
	path := ThumbnailCachePath(originalPath, maxSide)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
}

// SaveThumbnail writes a PNG thumbnail to disk.
func SaveThumbnail(originalPath string, maxSide int, img image.Image) error {
	ensureDir(thumbnailCacheDir())
	path := ThumbnailCachePath(originalPath, maxSide)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

// GetMediaThumbnail returns a thumbnail for image/video based on extension.
// It uses disk cache if available.
func GetMediaThumbnail(path string, maxSide int) (image.Image, error) {
	if img, err := LoadThumbnail(path, maxSide); err == nil && img != nil {
		return img, nil
	}
	ext := strings.ToLower(filepath.Ext(path))
	var (
		img image.Image
		err error
	)
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp":
		img, err = GenerateImageThumbnail(path, maxSide)
	case ".mp4", ".mov", ".avi", ".mkv", ".wmv":
		img, err = GenerateVideoThumbnail(path, maxSide)
	default:
		err = os.ErrInvalid
	}
	if err == nil && img != nil {
		_ = SaveThumbnail(path, maxSide, img)
	}
	return img, err
}

func itoa(i int) string { return fmtInt(int64(i)) }

func fmtInt(i int64) string {
	// simple base-10 without strconv to avoid extra imports in this file
	if i == 0 {
		return "0"
	}
	neg := false
	if i < 0 {
		neg = true
		i = -i
	}
	var buf [20]byte
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}
