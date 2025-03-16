package optimise

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"handytools/pkg/common"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
)

var profiles = map[string]int{
	"x-small": 1080,
	"small":   1440,
	"med":     1920,
	"large":   2560,
	"x-large": 0, // original size, no resizing
}

const jpegQuality = 85

func optimiseImages(cfg Config) {
	logger := common.GetLogger()

	if !cfg.Apply {
		common.SetDryRunMode(true)
		logger.Info("Running in DRYRUN mode")
	}

	maxSize, exists := profiles[cfg.Profile]
	if !exists {
		logger.Errorf("Invalid profile: %s", cfg.Profile)
		return
	}

	for _, filePath := range cfg.InputFiles {
		file, err := os.Open(filePath)
		if err != nil {
			logger.WithError(err).Errorf("Failed to open file: %s", filePath)
			continue
		}

		img, _, err := image.Decode(file)
		if err != nil {
			logger.WithError(err).Errorf("Unsupported format or corrupted image: %s", filePath)
			file.Close()
			continue
		}

		origWidth := img.Bounds().Dx()
		origHeight := img.Bounds().Dy()

		origFileInfo, _ := file.Stat()
		origSize := origFileInfo.Size()
		file.Close()

		if maxSize == 0 || (origWidth <= maxSize && origHeight <= maxSize) {
			logger.Infof("Skipping %s (already within size limits or original size requested)", filePath)
			continue
		}

		newWidth, newHeight := scaleDimensions(origWidth, origHeight, maxSize)
		resizedImg := imaging.Resize(img, newWidth, newHeight, imaging.Lanczos)

		ext := filepath.Ext(filePath)
		tempOutputPath := filePath[:len(filePath)-len(ext)] + "_temp" + ext

		err = imaging.Save(resizedImg, tempOutputPath, imaging.JPEGQuality(jpegQuality))
		if err != nil {
			logger.WithError(err).Errorf("Failed to save resized image: %s", filePath)
			continue
		}

		newFileInfo, _ := os.Stat(tempOutputPath)
		newSize := newFileInfo.Size()

		logger.Infof("Processed: %s", filePath)
		logger.Infof("Original: %dx%d, %.2f MB", origWidth, origHeight, float64(origSize)/(1024*1024))
		logger.Infof("New size: %dx%d, %.2f MB", newWidth, newHeight, float64(newSize)/(1024*1024))

		if cfg.Apply {
			if err := replaceFile(tempOutputPath, filePath, logger); err != nil {
				logger.WithError(err).Errorf("Failed to replace original file: %s", filePath)
				continue
			}
		} else {
			_ = os.Remove(tempOutputPath)
		}
	}
}

func replaceFile(tempPath, originalPath string, logger *logrus.Logger) error {
	if err := os.Remove(originalPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove original file: %w", err)
	}
	if err := os.Rename(tempPath, originalPath); err != nil {
		return fmt.Errorf("failed to rename file: %w", err)
	}
	return nil
}

func scaleDimensions(origWidth, origHeight, maxSize int) (int, int) {
	if origWidth > origHeight {
		return maxSize, int(float64(origHeight) * float64(maxSize) / float64(origWidth))
	}
	return int(float64(origWidth) * (float64(maxSize) / float64(origHeight))), maxSize
}
