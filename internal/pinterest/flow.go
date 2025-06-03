package pinterest

import (
	"fmt"
	"github.com/disintegration/imaging"
	"handytools/pkg/common"
	"image"
	"image/color"
)

func assembleFlowLayout(images []image.Image, layout string, outputPrefix string) error {
	logger := common.GetLogger()
	cols := getColumnsForLayout(layout)
	if cols <= 0 {
		return fmt.Errorf("unknown layout: %s", layout)
	}
	tileW := (maxWidth - (cols-1)*padding) / cols

	page := 1
	cursor := 0
	for cursor < len(images) {
		canvas := imaging.New(maxWidth, maxHeight, color.NRGBA{255, 255, 255, 255})
		x, y := 0, 0
		rowHeight := 0

		for cursor < len(images) {
			img := images[cursor]
			r := img.Bounds()
			aspect := float64(r.Dx()) / float64(r.Dy())
			th := int(float64(tileW) / aspect)

			if x+tileW > maxWidth {
				x = 0
				y += rowHeight + padding
				rowHeight = 0
			}
			if y+th > maxHeight {
				break
			}

			resized := imaging.Resize(img, tileW, th, imaging.Lanczos)
			canvas = imaging.Paste(canvas, resized, image.Pt(x, y))
			x += tileW + padding
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
