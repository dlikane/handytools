package collage

import (
	"handytools/pkg/common"
	"path/filepath"

	"github.com/spf13/cobra"
)

type Config struct {
	InputFiles []string
	OutputFile string
	Rows       int
	Columns    int
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
		logger.Info("Running collage with config: %+v\n", config)
		createCollage(config)
	},
}

func init() {
	Cmd.Flags().IntVarP(&config.Rows, "rows", "r", 1, "Number of rows")
	Cmd.Flags().IntVarP(&config.Columns, "columns", "c", 1, "Number of columns")
	Cmd.Flags().StringVarP(&config.OutputFile, "output", "o", "collage.jpg", "Output file")
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
