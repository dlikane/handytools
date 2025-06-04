package assemble

import (
	"fmt"
	"github.com/disintegration/imaging"
	"handytools/pkg/common"
	"image"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Directory string
	Output    string
	Layout    string
}

var logger = common.GetLogger()

func AssembleFromDirectory(config Config) ([]string, error) {
	var imagePaths []string

	// Read all files from the directory
	entries, err := os.ReadDir(config.Directory)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".jpg") {
			imagePaths = append(imagePaths, filepath.Join(config.Directory, entry.Name()))
		}
	}

	if len(imagePaths) == 0 {
		return nil, fmt.Errorf("no valid images found in directory")
	}

	// Assemble images into output file
	outputPrefix := strings.TrimSuffix(config.Output, ".jpg")
	err = AssembleImages(imagePaths, config.Layout, outputPrefix)
	if err != nil {
		return nil, fmt.Errorf("failed to assemble images: %w", err)
	}

	return imagePaths, nil
}

func AssembleImages(paths []string, layout string, outputPrefix string) error {
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
