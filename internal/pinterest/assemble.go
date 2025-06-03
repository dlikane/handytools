package pinterest

import (
	"fmt"
	"github.com/disintegration/imaging"
	"handytools/pkg/common"
	"image"
)

const (
	maxWidth  = 1080
	maxHeight = 1920
	padding   = 10
)

func AssembleImages(paths []string, layout string, outputPrefix string) error {
	logger := common.GetLogger()

	var images []image.Image
	for _, path := range paths {
		img, err := imaging.Open(path)
		if err != nil {
			logger.WithError(err).Warnf("Skipping image: %s", path)
			continue
		}
		images = append(images, img)
	}
	if len(images) == 0 {
		return fmt.Errorf("no valid images to assemble")
	}

	if layout == "fit" {
		return assembleTightFit(images, outputPrefix)
	}
	return assembleFlowLayout(images, layout, outputPrefix)
}
