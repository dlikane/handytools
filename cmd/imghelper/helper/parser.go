package helper

import (
	"errors"
	"flag"
	"path/filepath"
	"strings"
)

type CliParams struct {
	Rows       int
	Cols       int
	Whitespace int
	Images     []string
	Output     string
}

func ParseCliParams() (CliParams, error) {
	// Get the logger instance
	logger := GetLogger()

	var rows, cols, whitespace int
	var output string

	flag.IntVar(&rows, "r", 1, "number of rows in the output image")
	flag.IntVar(&rows, "row", 1, "number of rows in the output image (alternative notation)")
	flag.IntVar(&cols, "c", 1, "number of columns in the output image")
	flag.IntVar(&cols, "col", 1, "number of columns in the output image (alternative notation)")
	flag.IntVar(&whitespace, "w", 0, "pixels of whitespace between images")
	flag.StringVar(&output, "o", "output.jpg", "name of the output file")
	flag.StringVar(&output, "output", "output.jpg", "name of the output file (alternative notation)")

	flag.Parse()

	var images []string
	// Append each file from the file path to the Images slice
	for _, input := range flag.Args() {
		// Check if wildcard is in use
		if strings.Contains(input, "*") {
			matches, err := filepath.Glob(input)
			if err != nil {
				logger.Errorf("An error occurred while expanding the file path: %v", err)
				return CliParams{}, err
			}
			images = append(images, matches...)
		} else {
			images = append(images, input)
		}
	}

	if len(images) == 0 {
		logger.Error("No images provided.")
		return CliParams{}, errors.New("no images provided")
	}

	return CliParams{
		Rows:       rows,
		Cols:       cols,
		Whitespace: whitespace,
		Images:     images,
		Output:     output,
	}, nil
}
