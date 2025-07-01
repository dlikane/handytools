package assemble

import (
	"fmt"
	"handytools/pkg/common"
	"image"

	"github.com/disintegration/imaging"
)

var logger = common.GetLogger()

func AssembleImagesWithMax(paths []string, outputPrefix string, fitOnePage bool) error {
	logger.Infof("Assemble images: %d", len(paths))
	var images []image.Image
	for _, path := range paths {
		img, err := imaging.Open(path)
		if err != nil {
			logger.WithError(err).Warn("Skipping image: ", path)
			continue
		}
		images = append(images, img)
		if len(images)%5 == 0 {
			logger.Infof("Assembled images: %d", len(images))
		}
	}
	if len(images) == 0 {
		return fmt.Errorf("no valid images to assemble")
	}

	canvas, pageBreaks, err := buildContinuousCanvas(images, fitOnePage)
	if err != nil {
		return fmt.Errorf("failed to build canvas: %w", err)
	}

	out := fmt.Sprintf("%s.jpg", outputPrefix)
	for i, y := range pageBreaks {
		var top int
		if i > 0 {
			top = pageBreaks[i-1]
		}
		cropped := imaging.Crop(canvas, image.Rect(0, top, maxWidth, y-DefaultSpace))
		if len(pageBreaks) > 1 {
			out = fmt.Sprintf("%s_%02d.jpg", outputPrefix, i+1)
		}
		if err := imaging.Save(cropped, out); err != nil {
			return fmt.Errorf("failed to save page %d: %w", i+1, err)
		}
		logger.Infof("Saved page: %s", out)
	}

	return nil
}
