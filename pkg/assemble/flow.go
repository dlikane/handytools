package assemble

import (
	"fmt"
	"github.com/disintegration/imaging"
	"handytools/pkg/common"
	"image"
	"image/color"
)

const (
	MaxWidth  = 1080
	MaxHeight = 1920
	Padding   = 10
)

func assembleFlowLayout(images []image.Image, layout string, outputPrefix string) error {
	logger := common.GetLogger()
	cols := getColumnsForLayout(layout)
	if cols <= 0 {
		return fmt.Errorf("unknown layout: %s", layout)
	}
	tileW := (MaxWidth - (cols-1)*Padding) / cols

	page := 1
	cursor := 0
	for cursor < len(images) {
		canvas := imaging.New(MaxWidth, MaxHeight, color.NRGBA{255, 255, 255, 255})
		x, y := 0, 0
		rowHeight := 0

		for cursor < len(images) {
			img := images[cursor]
			r := img.Bounds()
			aspect := float64(r.Dx()) / float64(r.Dy())
			th := int(float64(tileW) / aspect)

			if x+tileW > MaxWidth {
				x = 0
				y += rowHeight + Padding
				rowHeight = 0
			}
			if y+th > MaxHeight {
				break
			}

			resized := imaging.Resize(img, tileW, th, imaging.Lanczos)
			canvas = imaging.Paste(canvas, resized, image.Pt(x, y))
			x += tileW + Padding
			if th > rowHeight {
				rowHeight = th
			}
			cursor++
		}

		out := fmt.Sprintf("%s_%02d.jpg", outputPrefix, page)
		if err := imaging.Save(canvas, out); err != nil {
			return fmt.Errorf("failed to save collage page %d: %w", page, err)
		}
		logger.Infof("Saved collage page: %s", out)
		page++
	}
	return nil
}
