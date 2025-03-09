package main

import (
	"handytools/cmd/imghelper/helper"
	"os"
)

func main() {
	// Get the logger instance
	logger := helper.GetLogger()

	params, err := helper.ParseCliParams()
	if err != nil {
		logger.Errorf("Failed to parse CLI parameters: %v", err)
		os.Exit(1)
	}

	logger.Infof("Rows: %d, Cols: %d", params.Rows, params.Cols)

	err = helper.MergeImages(params)
	if err != nil {
		logger.Errorf("Failed to merge images: %v", err)
		os.Exit(1)
	}

	logger.Infof("Images merged successfully with Rows: %d, Cols: %d", params.Rows, params.Cols)
}
