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
	"x-large": 0,    // original size, no resizing
	"insta":   1350, // insta optimal 1080 x 1350 4:5 portrait
}

const jpegQuality = 85

type profileResult struct {
	Width, Height int
	SizeBytes     int64
}

type fileStat struct {
	Filename string
	OrigW    int
	OrigH    int
	OrigSize int64
	Profiles map[string]profileResult
}

type statSummary struct {
	Files         []fileStat
	Totals        map[string]int64
	TotalOriginal int64
	TotalResized  int64
}

func optimiseImages(cfg Config) {
	logger := common.GetLogger()

	if !cfg.Apply && !cfg.Stat {
		common.SetDryRunMode(true)
		logger.Info("Running in DRYRUN mode")
	}

	summary := statSummary{
		Totals: make(map[string]int64),
	}

	for _, filePath := range cfg.InputFiles {
		handleImage(filePath, cfg, &summary)
	}

	printSummary(cfg, summary)
}

func handleImage(filePath string, cfg Config, summary *statSummary) {
	logger := common.GetLogger()
	file, err := os.Open(filePath)
	if err != nil {
		logger.WithError(err).Errorf("Failed to open file: %s", filePath)
		return
	}
	img, _, err := image.Decode(file)
	if err != nil {
		logger.WithError(err).Errorf("Unsupported or corrupted image: %s", filePath)
		file.Close()
		return
	}
	origWidth, origHeight := img.Bounds().Dx(), img.Bounds().Dy()
	origFileInfo, _ := file.Stat()
	origSize := origFileInfo.Size()
	file.Close()

	summary.TotalOriginal += origSize
	entry := fileStat{
		Filename: filepath.Base(filePath),
		OrigW:    origWidth,
		OrigH:    origHeight,
		OrigSize: origSize,
		Profiles: map[string]profileResult{},
	}

	if cfg.Stat {
		for name, size := range profiles {
			var resized image.Image
			if size == 0 || (origWidth <= size && origHeight <= size) {
				resized = img
			} else {
				newW, newH := scaleDimensions(origWidth, origHeight, size)
				resized = imaging.Resize(img, newW, newH, imaging.Lanczos)
			}
			tempPath := filepath.Join(os.TempDir(), fmt.Sprintf("%s_%s.jpg", filepath.Base(filePath), name))
			_ = imaging.Save(resized, tempPath, imaging.JPEGQuality(jpegQuality))
			info, _ := os.Stat(tempPath)
			sizeBytes := info.Size()
			os.Remove(tempPath)

			w, h := resized.Bounds().Dx(), resized.Bounds().Dy()
			entry.Profiles[name] = profileResult{w, h, sizeBytes}
			summary.Totals[name] += sizeBytes
		}
		summary.Files = append(summary.Files, entry)
		return
	}

	// Non-stat mode (apply or dry-run)
	maxSize, exists := profiles[cfg.Profile]
	if !exists {
		logger.Errorf("Invalid profile: %s", cfg.Profile)
		return
	}
	if maxSize == 0 || (origWidth <= maxSize && origHeight <= maxSize) {
		logger.Infof("Skipping %s (already within size limits or original size requested)", filePath)
		return
	}

	newWidth, newHeight := scaleDimensions(origWidth, origHeight, maxSize)
	resizedImg := imaging.Resize(img, newWidth, newHeight, imaging.Lanczos)
	ext := filepath.Ext(filePath)
	tempOutputPath := filePath[:len(filePath)-len(ext)] + "_temp" + ext
	err = imaging.Save(resizedImg, tempOutputPath, imaging.JPEGQuality(jpegQuality))
	if err != nil {
		logger.WithError(err).Errorf("Failed to save resized image: %s", filePath)
		return
	}
	newFileInfo, _ := os.Stat(tempOutputPath)
	newSize := newFileInfo.Size()
	summary.TotalResized += newSize

	dry := "dry run"
	if cfg.Apply {
		dry = "applied"
	}
	logger.Infof("%s org: %dx%d %.2f MB %s: %dx%d %.2f MB (%s)",
		filepath.Base(filePath), origWidth, origHeight, float64(origSize)/(1024*1024),
		cfg.Profile, newWidth, newHeight, float64(newSize)/(1024*1024), dry)
	if cfg.Apply {
		if err := replaceFile(tempOutputPath, filePath, logger); err != nil {
			logger.WithError(err).Errorf("Failed to replace original file: %s", filePath)
			return
		}
	} else {
		_ = os.Remove(tempOutputPath)
	}
}

func printSummary(cfg Config, summary statSummary) {
	logger := common.GetLogger()
	if cfg.Stat {
		for _, entry := range summary.Files {
			line := fmt.Sprintf("%s org: %dx%d %.2f MB", entry.Filename, entry.OrigW, entry.OrigH, float64(entry.OrigSize)/(1024*1024))
			for _, key := range []string{"x-large", "large", "med", "small", "x-small"} {
				p := entry.Profiles[key]
				line += fmt.Sprintf(" %s: %dx%d %.2f MB", key, p.Width, p.Height, float64(p.SizeBytes)/(1024*1024))
			}
			logger.Info(line)
		}
		line := fmt.Sprintf("files: %d org: %.2f MB", len(summary.Files), float64(summary.TotalOriginal)/(1024*1024))
		for _, key := range []string{"x-large", "large", "med", "small", "x-small"} {
			line += fmt.Sprintf(" %s: %.2f MB", key, float64(summary.Totals[key])/(1024*1024))
		}
		logger.Info(line)
	} else {
		logger.Infof("files: %d org: %.2f MB %s: %.2f MB", len(summary.Files), float64(summary.TotalOriginal)/(1024*1024), cfg.Profile, float64(summary.TotalResized)/(1024*1024))
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
