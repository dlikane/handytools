package frame

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
)

var namedColors = map[string]color.NRGBA{
	"white": {R: 255, G: 255, B: 255, A: 255},
	"black": {R: 0, G: 0, B: 0, A: 255},
	"cream": {R: 255, G: 253, B: 231, A: 255},
	"ivory": {R: 255, G: 255, B: 240, A: 255},
}

func runFrame(cfg Config) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	fc, err := parseColor(cfg.Color)
	if err != nil {
		logger.WithError(err).Error("Invalid frame color")
		return
	}

	if cfg.OutputDir != "." {
		if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
			logger.WithError(err).Errorf("Cannot create output directory: %s", cfg.OutputDir)
			return
		}
	}

	for _, inputPath := range cfg.InputFiles {
		if err := processFile(inputPath, cfg, fc, rng); err != nil {
			logger.WithError(err).Errorf("Failed: %s", inputPath)
		}
	}
}

func processFile(inputPath string, cfg Config, fc color.NRGBA, rng *rand.Rand) error {
	src, err := imaging.Open(inputPath)
	if err != nil {
		return err
	}

	bounds := src.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	frameW := max(int(float64(w)*cfg.FramePct/100.0), 1)

	newW := w + 2*frameW
	newH := h + 2*frameW

	dst := image.NewNRGBA(image.Rect(0, 0, newW, newH))

	// Fill canvas with frame color.
	for y := 0; y < newH; y++ {
		for x := 0; x < newW; x++ {
			dst.SetNRGBA(x, y, fc)
		}
	}

	if cfg.Torn {
		tearDepth := max(int(float64(frameW)*cfg.TornDepth/100.0), 2)
		topEdge := tornCurve(w, tearDepth, rng)
		bottomEdge := tornCurve(w, tearDepth, rng)
		leftEdge := tornCurve(h, tearDepth, rng)
		rightEdge := tornCurve(h, tearDepth, rng)

		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				// Only paint the photo pixel if it's inside all four torn boundaries.
				if x >= leftEdge[y] && x < w-rightEdge[y] &&
					y >= topEdge[x] && y < h-bottomEdge[x] {
					dst.Set(frameW+x, frameW+y, src.At(bounds.Min.X+x, bounds.Min.Y+y))
				}
			}
		}
	} else {
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				dst.Set(frameW+x, frameW+y, src.At(bounds.Min.X+x, bounds.Min.Y+y))
			}
		}
	}

	var outputPath string
	if cfg.OutputDir == "." {
		outputPath = inputPath
	} else {
		outputPath = filepath.Join(cfg.OutputDir, filepath.Base(inputPath))
	}

	logger.Infof("%s → %s (frame=%dpx torn=%v)", filepath.Base(inputPath), outputPath, frameW, cfg.Torn)
	return imaging.Save(dst, outputPath, imaging.JPEGQuality(85))
}

// tornCurve returns a non-repeating torn-edge offset curve of `length` values
// in [0, maxDepth]. Three scales of random control-point interpolation replace
// sine waves, eliminating any periodicity:
//   - coarse  (~12 pts): broad undulations — where the tear runs deep vs shallow
//   - medium  (~60 pts): secondary bumps and notches between the coarse features
//   - fine    (~350 pts): ragged paper-fibre micro-texture
func tornCurve(length, maxDepth int, rng *rand.Rand) []int {
	d := float64(maxDepth)
	coarse := randomInterp(length, 10+rng.Intn(8), 0, d, rng)
	medium := randomInterp(length, 50+rng.Intn(30), -d*0.25, d*0.25, rng)
	fine := randomInterp(length, 300+rng.Intn(100), -d*0.12, d*0.12, rng)

	curve := make([]int, length)
	for i := range curve {
		v := coarse[i] + medium[i] + fine[i]
		// Occasional sharp spike: a sudden deep pierce into the photo.
		if rng.Float64() < 0.006 {
			v += d * (0.25 + rng.Float64()*0.5)
		}
		curve[i] = int(math.Max(0, math.Min(d, v)))
	}
	return curve
}

// randomInterp produces a smooth curve of `length` values by interpolating
// `n` random control points uniformly sampled from [lo, hi], using smoothstep
// so there are no kinks at control-point boundaries.
func randomInterp(length, n int, lo, hi float64, rng *rand.Rand) []float64 {
	pts := make([]float64, n+1)
	for i := range pts {
		pts[i] = lo + rng.Float64()*(hi-lo)
	}
	out := make([]float64, length)
	for i := range out {
		t := float64(i) / float64(max(length-1, 1)) * float64(n)
		j := int(t)
		if j >= n {
			j = n - 1
		}
		f := t - float64(j)
		f = f * f * (3 - 2*f) // smoothstep
		out[i] = pts[j]*(1-f) + pts[j+1]*f
	}
	return out
}

func parseColor(s string) (color.NRGBA, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	if c, ok := namedColors[s]; ok {
		return c, nil
	}
	hex := strings.TrimPrefix(s, "#")
	if len(hex) != 6 {
		return color.NRGBA{}, fmt.Errorf("unknown color %q — use white/black/cream/ivory or #RRGGBB", s)
	}
	r, e1 := strconv.ParseUint(hex[0:2], 16, 8)
	g, e2 := strconv.ParseUint(hex[2:4], 16, 8)
	b, e3 := strconv.ParseUint(hex[4:6], 16, 8)
	if e1 != nil || e2 != nil || e3 != nil {
		return color.NRGBA{}, fmt.Errorf("invalid hex color %q", s)
	}
	return color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}, nil
}
