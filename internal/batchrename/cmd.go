package batchrename

import (
	"fmt"
	"handytools/pkg/common"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ncruces/zenity"
	"github.com/spf13/cobra"
)

type RenameConfig struct {
	Pattern    string
	InputFiles []string
	Apply      bool
}

var (
	logger = common.GetLogger()
	config RenameConfig
)

var Cmd = &cobra.Command{
	Use:   "batchrename",
	Short: "Batch rename files with a given pattern",
	Long:  "Renames files using the provided pattern followed by a counter (e.g., base_0001, base_0002).",
	Run: func(cmd *cobra.Command, args []string) {
		if runtime.GOOS != "windows" {
			logger.Error("This command is only available on Windows.")
			return
		}

		if len(args) == 0 {
			logger.Error("No files selected.")
			return
		}

		// Log arguments for debugging
		logger.Infof("Raw arguments received: %+v", args)

		// Fix for PowerShell passing System.Object[] incorrectly
		var inputFiles []string
		var clipboardText strings.Builder
		for _, arg := range args {
			logger.Infof("Processing argument: %s", arg)
			if arg == "" || strings.Contains(arg, "System.Object[]") {
				logger.Error("Invalid argument detected: skipping.")
				continue
			}
			cleanedPath := strings.TrimSpace(arg)
			if absPath, err := filepath.Abs(cleanedPath); err == nil {
				if _, err := os.Stat(absPath); err == nil {
					inputFiles = append(inputFiles, absPath)
					clipboardText.WriteString(fmt.Sprintf("\"%s\" ", absPath)) // Quote filenames with spaces
				} else {
					logger.WithError(err).Errorf("Skipping invalid file: %s", absPath)
				}
			} else {
				logger.WithError(err).Errorf("Failed to resolve absolute path: %s", cleanedPath)
			}
		}

		if len(inputFiles) == 0 {
			logger.Error("No valid files to rename.")
			return
		}

		// Copy original file list to clipboard
		clipboardData := clipboardText.String()
		if clipboardData != "" {
			common.CopyToClipboard(clipboardData)
			logger.Infof("Original filenames copied to clipboard: %s", clipboardData)
		} else {
			logger.Warn("No valid files found, clipboard not updated.")
		}

		// Prompt for name pattern
		pattern, err := zenity.Entry("Enter base name for renaming files:", zenity.Title("Batch Rename"))
		if err != nil || pattern == "" {
			logger.Info("Operation cancelled.")
			return
		}

		config.Pattern = pattern
		config.InputFiles = inputFiles

		// Prepare preview text
		previewText := generatePreview(config)

		// Ask user for confirmation with preview details
		if err := zenity.Question("Do you want to apply the following renaming?\n\n"+previewText, zenity.Title("Confirm Rename"), zenity.OKLabel("Apply"), zenity.CancelLabel("Cancel")); err == nil {
			applyChanges(config)
		}
	},
}

func generatePreview(cfg RenameConfig) string {
	counter := 1
	var previewText strings.Builder
	for _, filePath := range cfg.InputFiles {
		dir := filepath.Dir(filePath)
		ext := filepath.Ext(filePath)
		newName := fmt.Sprintf("%s_%04d%s", cfg.Pattern, counter, ext)
		newPath := filepath.Join(dir, newName)
		previewText.WriteString(fmt.Sprintf("%s -> %s\n", filepath.Base(filePath), filepath.Base(newPath)))
		counter++
	}
	return previewText.String()
}

func applyChanges(cfg RenameConfig) {
	counter := 1
	for _, filePath := range cfg.InputFiles {
		dir := filepath.Dir(filePath)
		ext := filepath.Ext(filePath)
		newName := fmt.Sprintf("%s_%04d%s", cfg.Pattern, counter, ext)
		newPath := filepath.Join(dir, newName)
		if err := os.Rename(filePath, newPath); err != nil {
			logger.WithError(err).Errorf("Failed to rename: %s", filePath)
		}
		counter++
	}
}
