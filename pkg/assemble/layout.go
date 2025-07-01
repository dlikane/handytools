package assemble

import (
	"handytools/pkg/layout"
	"image"
	"image/color"

	"github.com/disintegration/imaging"
)

const (
	maxWidth     = 1080
	maxHeight    = 1920
	defaultRowH  = 320
	DefaultSpace = 10
)

func buildContinuousCanvas(images []image.Image, fitOnePage bool) (*image.NRGBA, []int, error) {
	var infos []layout.ImageInfo
	for _, img := range images {
		b := img.Bounds()
		aspect := float64(b.Dx()) / float64(b.Dy())
		infos = append(infos, layout.ImageInfo{Aspect: aspect})
	}

	rowHeight := defaultRowH
	var positions []layout.PlacedImage
	var canvasHeight int
	var pageBreaks []int

	for {
		cfg := layout.Config{
			MaxWidth:       maxWidth,
			TargetHeight:   rowHeight,
			Spacing:        DefaultSpace,
			Tolerance:      0.25,
			MinRowItems:    2,
			MinAspectTotal: 0,
		}

		positions, canvasHeight, pageBreaks = layout.JustifyWithPageSplits(infos, cfg, maxHeight)

		if !fitOnePage || len(pageBreaks) <= 1 || rowHeight <= 100 {
			break
		}

		rowHeight -= 20
	}

	canvas := imaging.New(maxWidth, canvasHeight, color.NRGBA{255, 255, 255, 255})
	for _, pos := range positions {
		if pos.Index >= len(images) {
			continue
		}
		img := images[pos.Index]
		resized := imaging.Resize(img, pos.Width, pos.Height, imaging.Lanczos)
		canvas = imaging.Paste(canvas, resized, image.Pt(pos.X, pos.Y))
	}

	return canvas, pageBreaks, nil
}
