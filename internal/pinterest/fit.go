package pinterest

import (
	"fmt"
	"github.com/disintegration/imaging"
	"handytools/pkg/common"
	"handytools/pkg/layout"
	"image"
	"image/color"
)

func assembleTightFit(images []image.Image, outputPrefix string) error {
	logger := common.GetLogger()

	var infos []layout.ImageInfo
	for _, img := range images {
		b := img.Bounds()
		aspect := float64(b.Dx()) / float64(b.Dy())
		infos = append(infos, layout.ImageInfo{Aspect: aspect})
	}

	cfg := layout.Config{
		MaxWidth:       1080,
		TargetHeight:   320,
		Spacing:        10,
		Tolerance:      0.25,
		MinRowItems:    2,
		MinAspectTotal: 0, // auto
	}

	positions, canvasHeight := layout.Justify(infos, cfg)
	canvas := imaging.New(cfg.MaxWidth, canvasHeight, color.NRGBA{255, 255, 255, 255})

	for _, pos := range positions {
		if pos.Index >= len(images) {
			continue
		}
		resized := imaging.Resize(images[pos.Index], pos.Width, pos.Height, imaging.Lanczos)
		canvas = imaging.Paste(canvas, resized, image.Pt(pos.X, pos.Y))
	}

	// Crop to 1920 height if needed
	if canvasHeight > 1920 {
		canvas = imaging.Crop(canvas, image.Rect(0, 0, 1080, 1920))
	}

	out := fmt.Sprintf("%s_01.jpg", outputPrefix)
	if err := imaging.Save(canvas, out); err != nil {
		return fmt.Errorf("failed to save fit: %w", err)
	}
	logger.Infof("Saved final justified collage: %s", out)
	return nil
}
