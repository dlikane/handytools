package collage

import (
	"handytools/pkg/common"
	"path/filepath"

	"github.com/spf13/cobra"
)

type Config struct {
	InputFiles  []string
	OutputFile  string
	Rows        int
	Columns     int
	AspectRatio string // "free" or "4x5"
}

var (
	logger = common.GetLogger()
	config Config
)

var Cmd = &cobra.Command{
	Use:   "collage",
	Short: "Create an image collage",
	Long:  "Combines multiple images into a collage with the specified rows and columns.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			logger.Error("No images provided. Use wildcard or file list.")
			return
		}
		config.InputFiles = expandWildcards(args)
		if config.AspectRatio != "free" && config.AspectRatio != "4x5" {
			logger.Error("Invalid aspect ratio. Use 'free' or '4x5'")
			return
		}
		logger.Infof("Running collage with config: %+v", config)
		createCollage(config)
	},
}

func init() {
	Cmd.Flags().IntVarP(&config.Rows, "rows", "r", 1, "Number of rows")
	Cmd.Flags().IntVarP(&config.Columns, "columns", "c", 1, "Number of columns")
	Cmd.Flags().StringVarP(&config.OutputFile, "output", "o", "collage.jpg", "Output file")
	Cmd.Flags().StringVarP(&config.AspectRatio, "aspect", "a", "free", "Output aspect ratio: 'free' or '4x5'")
}

func expandWildcards(patterns []string) []string {
	var files []string
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			logger.WithField("pattern", pattern).Error("Error processing wildcard")
			continue
		}
		files = append(files, matches...)
	}
	return files
}
