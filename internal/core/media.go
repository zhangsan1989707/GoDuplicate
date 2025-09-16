package core

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

// GenerateImageThumbnail decodes an image file and returns a small RGBA image thumbnail (max 160px)
func GenerateImageThumbnail(path string, maxSide int) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	b := img.Bounds()
	w := b.Dx()
	h := b.Dy()
	// compute scale
	scale := float64(maxSide) / float64(w)
	if h > w {
		scale = float64(maxSide) / float64(h)
	}
	if scale >= 1 {
		return img, nil
	}
	nw := int(float64(w) * scale)
	nh := int(float64(h) * scale)
	// nearest resize (fast placeholder)
	out := image.NewRGBA(image.Rect(0, 0, nw, nh))
	for y := 0; y < nh; y++ {
		for x := 0; x < nw; x++ {
			sx := int(float64(x) / float64(nw) * float64(w))
			sy := int(float64(y) / float64(nh) * float64(h))
			out.Set(x, y, img.At(sx, sy))
		}
	}
	return out, nil
}

// PerceptualHash computes a simple average hash (aHash) 8x8 -> 64-bit string hex
func PerceptualHash(img image.Image) string {
	// downscale to 8x8 gray and compute average
	w, h := 8, 8
	total := uint64(0)
	vals := make([]uint8, 0, 64)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			// sample center of cell
			px := (x*img.Bounds().Dx() + img.Bounds().Dx()/16) / w
			py := (y*img.Bounds().Dy() + img.Bounds().Dy()/16) / h
			r, g, b, _ := img.At(px, py).RGBA()
			g8 := uint8((r*299 + g*587 + b*114) / 1000 >> 8)
			vals = append(vals, g8)
			total += uint64(g8)
		}
	}
	avg := uint8(total / 64)
	var bits uint64
	for i, v := range vals {
		if v >= avg {
			bits |= 1 << uint(i)
		}
	}
	// return as 16-hex chars
	const hexdigits = "0123456789abcdef"
	out := make([]byte, 16)
	for i := 0; i < 16; i++ {
		nibble := (bits >> uint((15-i)*4)) & 0xF
		out[i] = hexdigits[nibble]
	}
	return string(out)
}

// HammingDistanceHex between two equal-length hex strings (64-bit represented as 16 hex)
func HammingDistanceHex(a, b string) int {
	if len(a) != len(b) {
		return 64
	}
	// convert hex chars to 4-bit and compare
	toNibble := func(c byte) uint8 {
		switch {
		case c >= '0' && c <= '9':
			return c - '0'
		case c >= 'a' && c <= 'f':
			return 10 + c - 'a'
		case c >= 'A' && c <= 'F':
			return 10 + c - 'A'
		}
		return 0
	}
	dist := 0
	for i := 0; i < len(a); i++ {
		na := toNibble(a[i])
		nb := toNibble(b[i])
		x := na ^ nb
		for j := 0; j < 4; j++ {
			if (x>>uint(j))&1 == 1 {
				dist++
			}
		}
	}
	return dist
}

// HashMedia generates/loads a thumbnail for the given path and computes its perceptual hash.
func HashMedia(path string, maxSide int) (string, error) {
	thumb, err := GenerateImageThumbnail(path, maxSide)
	if err != nil {
		return "", err
	}
	return PerceptualHash(thumb), nil
}

// EstimateGroupSimilarity returns an approximate similarity percent [0,100] using the first file as reference.
// Only meaningful for image/video groups.
func EstimateGroupSimilarity(files []FileInfo) float64 {
	if len(files) <= 1 {
		return 100
	}
	refHash, err := HashMedia(files[0].Path, 128)
	if err != nil {
		return 0
	}
	total := 0.0
	count := 0
	for i := 1; i < len(files); i++ {
		h, err := HashMedia(files[i].Path, 128)
		if err != nil {
			continue
		}
		d := HammingDistanceHex(refHash, h)
		// map distance [0,64] -> similarity percent
		sim := 100.0 * (1.0 - float64(d)/64.0)
		if sim < 0 {
			sim = 0
		}
		if sim > 100 {
			sim = 100
		}
		total += sim
		count++
	}
	if count == 0 {
		return 0
	}
	return total / float64(count)
}
