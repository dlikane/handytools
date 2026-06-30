package distort

import (
	"image"
	"math"
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
)

func runDistort(cfg Config) {
	r := newRand(cfg.Seed)
	logger.Infof("Distorting %s → %s (mode=%s intensity=%.2f)", cfg.InputFile, cfg.OutputFile, cfg.Mode, cfg.Intensity)

	var err error
	switch cfg.Mode {
	case "corrupt":
		err = blockDisplace(cfg.InputFile, cfg.OutputFile, cfg.Intensity, r)
	case "shift":
		err = shiftRows(cfg.InputFile, cfg.OutputFile, cfg.Intensity, r)
	case "melt":
		err = meltRows(cfg.InputFile, cfg.OutputFile, cfg.Intensity, r)
	default:
		logger.Errorf("Unknown mode: %s (use corrupt, shift, or melt)", cfg.Mode)
		return
	}

	if err != nil {
		logger.WithError(err).Error("Distortion failed")
		return
	}
	logger.Infof("Saved: %s", cfg.OutputFile)
}

// blockDisplace simulates I-frame removal by creating a few glitch zones, each
// with its own displacement vector (like a P-frame motion vector referencing the
// wrong position). All blocks within a zone shift in the same direction, so the
// effect reads as intentional regions of glitch rather than scattered noise.
func blockDisplace(input, output string, intensity float64, r *rand.Rand) error {
	src, err := imaging.Open(input)
	if err != nil {
		return err
	}

	bounds := src.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	const blockW, blockH = 24, 16

	dst := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dst.Set(x, y, src.At(bounds.Min.X+x, bounds.Min.Y+y))
		}
	}

	type zone struct{ cx, cy, radius, dx, dy int }

	numZones := 3 + r.Intn(2+int(intensity*4))
	maxDisplace := max(int(float64(w)*0.35), 80)
	zoneRadius := int(float64(min(w, h)) * (0.03 + intensity*0.05))

	zones := make([]zone, numZones)
	for i := range zones {
		zones[i] = zone{
			cx:     r.Intn(w),
			cy:     r.Intn(h),
			radius: zoneRadius + r.Intn(max(zoneRadius/2, 1)),
			dx:     r.Intn(maxDisplace*2+1) - maxDisplace,
			dy:     r.Intn(maxDisplace+1) - maxDisplace/2,
		}
	}

	for by := 0; by < h/blockH; by++ {
		for bx := 0; bx < w/blockW; bx++ {
			bcx := bx*blockW + blockW/2
			bcy := by*blockH + blockH/2

			// Find the first zone that covers this block.
			var dx, dy int
			influenced := false
			for _, z := range zones {
				dist := math.Sqrt(float64((bcx-z.cx)*(bcx-z.cx) + (bcy-z.cy)*(bcy-z.cy)))
				if dist <= float64(z.radius) {
					// Small per-block jitter so zone edges don't look perfectly uniform.
					dx = z.dx + r.Intn(blockW) - blockW/2
					dy = z.dy + r.Intn(blockH/2) - blockH/4
					influenced = true
					break
				}
			}
			if !influenced {
				continue
			}

			for py := 0; py < blockH; py++ {
				for px := 0; px < blockW; px++ {
					tx := bx*blockW + px
					ty := by*blockH + py
					sx := tx + dx
					sy := ty + dy
					if sx < 0 || sx >= w || sy < 0 || sy >= h {
						continue
					}
					dst.Set(tx, ty, src.At(bounds.Min.X+sx, bounds.Min.Y+sy))
				}
			}
		}
	}

	return saveOutput(dst, output)
}

// meltRows creates a downward smear by randomly freezing horizontal bands:
// rows repeat content from above, as if the image is "dripping" down.
func meltRows(input, output string, intensity float64, r *rand.Rand) error {
	src, err := imaging.Open(input)
	if err != nil {
		return err
	}

	bounds := src.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	dst := image.NewNRGBA(image.Rect(0, 0, w, h))

	frozenAt := -1
	for y := 0; y < h; y++ {
		// Chance to start a freeze event (repeat a row from above).
		if frozenAt < 0 && r.Float64() < intensity*0.4 {
			lookback := int(float64(h) * intensity * 0.3)
			if lookback < 1 {
				lookback = 1
			}
			frozenAt = y - r.Intn(lookback)
			if frozenAt < 0 {
				frozenAt = 0
			}
		}
		// Chance to end the freeze.
		if frozenAt >= 0 && r.Float64() < 0.25 {
			frozenAt = -1
		}

		srcY := y
		if frozenAt >= 0 {
			srcY = frozenAt
		}

		for x := 0; x < w; x++ {
			dst.Set(x, y, src.At(bounds.Min.X+x, bounds.Min.Y+srcY))
		}
	}

	return saveOutput(dst, output)
}

// shiftRows displaces rows of pixels horizontally by random amounts, producing
// a CRT scan-error or magnetic tape dropout look.
func shiftRows(input, output string, intensity float64, r *rand.Rand) error {
	src, err := imaging.Open(input)
	if err != nil {
		return err
	}

	bounds := src.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	maxShift := int(float64(w) * intensity * 0.4)
	if maxShift < 1 {
		maxShift = 1
	}

	dst := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		shift := 0
		if r.Float64() < intensity*2 {
			shift = r.Intn(maxShift*2+1) - maxShift
		}
		for x := 0; x < w; x++ {
			sx := x + shift
			if sx < 0 {
				sx = 0
			} else if sx >= w {
				sx = w - 1
			}
			dst.Set(x, y, src.At(bounds.Min.X+sx, bounds.Min.Y+y))
		}
	}

	return saveOutput(dst, output)
}

func saveOutput(img image.Image, path string) error {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".png":
		return imaging.Save(img, path)
	default:
		return imaging.Save(img, path, imaging.JPEGQuality(85))
	}
}

func newRand(seed int64) *rand.Rand {
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	return rand.New(rand.NewSource(seed))
}
