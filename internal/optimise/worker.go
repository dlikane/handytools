package optimise

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"handytools/pkg/common"
	"image"
	_ "image/jpeg" // Support JPEG decoding
	_ "image/png"  // Support PNG decoding
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
)

const midSize = 3000
const smallSize = 1000

func optimiseImages(cfg Config) {
	logger := common.GetLogger()

	if !cfg.Apply {
		common.SetDryRunMode(true)
		logger.Info("Running in DRYRUN mode")
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

		maxSize := midSize
		if cfg.Small {
			maxSize = smallSize
		}
		if origWidth <= maxSize && origHeight <= maxSize {
			logger.Infof("Skipping %s (already within size limits)", filePath)
			continue
		}

		newWidth, newHeight := scaleDimensions(origWidth, origHeight, maxSize)
		resizedImg := imaging.Resize(img, newWidth, newHeight, imaging.Lanczos)

		ext := filepath.Ext(filePath)
		tempOutputPath := filePath[:len(filePath)-len(ext)] + "_temp" + ext

		// Save using correct format
		err = imaging.Save(resizedImg, tempOutputPath)
		if err != nil {
			logger.WithError(err).Errorf("Failed to save resized image: %s", filePath)
			continue
		}

		newFileInfo, _ := os.Stat(tempOutputPath)
		newSize := newFileInfo.Size()

		logger.Infof("Processed: %s", filePath)
		logger.Infof("Original: %dx%d, %.2f MB", origWidth, origHeight, float64(origSize)/1024/1024)
		logger.Infof("Resized : %dx%d, %.2f MB", newWidth, newHeight, float64(newSize)/1024/1024)

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
	// Ensure the original file is writable
	if err := os.Chmod(originalPath, 0666); err != nil {
		logger.WithError(err).Warnf("Failed to modify permissions for: %s", originalPath)
	}

	// Remove the original file before renaming
	if err := os.Remove(originalPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove original file: %w", err)
	}

	// Rename temp file to original filename
	if err := os.Rename(tempPath, originalPath); err != nil {
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	logger.Infof("Successfully replaced original file: %s", originalPath)
	return nil
}

func scaleDimensions(width, height, maxSize int) (int, int) {
	if width > height {
		newWidth := maxSize
		newHeight := int(float64(height) * (float64(maxSize) / float64(width)))
		return newWidth, newHeight
	} else {
		newHeight := maxSize
		newWidth := int(float64(width) * (float64(maxSize) / float64(height)))
		return newWidth, newHeight
	}
}
