package layout

import "math"

type ImageInfo struct {
	Aspect float64 // width / height
}

type PlacedImage struct {
	Index  int
	X, Y   int
	Width  int
	Height int
}

type Config struct {
	MaxWidth       int     // Total canvas width (e.g. 1080)
	TargetHeight   int     // Desired row height
	Spacing        int     // Space between tiles
	Tolerance      float64 // Acceptable row height deviation (e.g. 0.25 = ±25%)
	MinRowItems    int     // Minimum images per row before finalizing
	MinAspectTotal float64 // Minimum total aspect ratio before forcing row
}

// JustifyWithPageSplits returns layout and vertical split points for paging
func JustifyWithPageSplits(images []ImageInfo, cfg Config, maxPageHeight int) (placed []PlacedImage, canvasHeight int, pageBreaks []int) {
	if cfg.MinRowItems < 1 {
		cfg.MinRowItems = 1
	}
	if cfg.MinAspectTotal == 0 {
		cfg.MinAspectTotal = float64(cfg.MaxWidth-cfg.Spacing*(cfg.MinRowItems-1)) / float64(cfg.TargetHeight)
	}

	var (
		row         []ImageInfo
		rowAspect   float64
		y           = 0
		startIdx    = 0
		lastImageIx = len(images) - 1
	)

	for i, img := range images {
		row = append(row, img)
		rowAspect += img.Aspect

		rowWidth := float64(cfg.MaxWidth - (len(row)-1)*cfg.Spacing)
		rowHeight := rowWidth / rowAspect

		validHeight := rowHeight >= float64(cfg.TargetHeight)*(1-cfg.Tolerance) &&
			rowHeight <= float64(cfg.TargetHeight)*(1+cfg.Tolerance)
		enoughItems := len(row) >= cfg.MinRowItems
		enoughAspect := rowAspect >= cfg.MinAspectTotal
		isLast := i == lastImageIx

		if ((validHeight && enoughItems) || enoughAspect) || isLast {
			scale := rowWidth / rowAspect
			x := 0
			rowTop := y
			for j, r := range row {
				w := int(math.Round(scale * r.Aspect))
				h := int(math.Round(scale))
				placed = append(placed, PlacedImage{
					Index:  startIdx + j,
					X:      x,
					Y:      y,
					Width:  w,
					Height: h,
				})
				x += w + cfg.Spacing
			}
			y += int(math.Round(scale)) + cfg.Spacing
			startIdx = i + 1
			row = nil
			rowAspect = 0

			if y-rowTop > 0 && y > 0 && y%maxPageHeight <= int(scale) {
				pageBreaks = append(pageBreaks, y)
			}
		}
	}

	if len(pageBreaks) == 0 || pageBreaks[len(pageBreaks)-1] != y {
		pageBreaks = append(pageBreaks, y)
	}

	return placed, y, pageBreaks
}
