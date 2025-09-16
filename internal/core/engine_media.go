package core

import (
    "path/filepath"
    "strings"
)

// MediaSimilarity groups images by perceptual hash threshold.
func MediaSimilarity(files []FileInfo, threshold int) []DuplicateGroup {
    buckets := map[string][]FileInfo{}
    hashes := map[string]string{}
    for _, f := range files {
        ext := strings.ToLower(filepath.Ext(f.Path))
        if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
            img, err := GenerateImageThumbnail(f.Path, 128)
            if err != nil { continue }
            h := PerceptualHash(img)
            hashes[f.Path] = h
        }
    }
    // simple clustering: compare to first in each bucket key set
    for _, f := range files {
        h, ok := hashes[f.Path]
        if !ok { continue }
        placed := false
        for key, list := range buckets {
            if len(list) == 0 { continue }
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
    // convert to groups (only groups with >=2 files)
    out := make([]DuplicateGroup, 0, len(buckets))
    for key, list := range buckets {
        if len(list) < 2 { continue }
        out = append(out, DuplicateGroup{ GroupID: key, Files: list })
    }
    return out
}


