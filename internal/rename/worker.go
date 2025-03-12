package rename

import (
	"fmt"
	"handytools/pkg/common"
	"os"
	"path/filepath"
)

func renameFiles(cfg RenameConfig) {
	logger := common.GetLogger()

	if !cfg.Apply {
		common.SetDryRunMode(true)
		logger.Info("Running in DRYRUN mode")
	}

	counter := 1
	for _, filePath := range cfg.InputFiles {
		dir := filepath.Dir(filePath)
		ext := filepath.Ext(filePath)
		newName := fmt.Sprintf("%s_%04d%s", cfg.OutputName, counter, ext)
		newPath := filepath.Join(dir, newName)

		logger.Infof("Would rename: %s -> %s", filePath, newPath)

		if cfg.Apply {
			if err := os.Rename(filePath, newPath); err != nil {
				logger.WithError(err).Errorf("Failed to rename file: %s", filePath)
			}
		}
		counter++
	}
}
