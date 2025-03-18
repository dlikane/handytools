package grab

import (
	"github.com/spf13/cobra"
)

var config = Config{
	Inputs:          []string{},
	ListOnly:        false,
	ExcludePatterns: []string{},
	ExtraExclusions: false,
}

var defaultExclusions = []string{
	".git/...",
	".idea/...",
	"dist/...",
	"build/...",
	"package-lock.json",
	"pnpm-lock.yaml",
	"node_modules/...",
	"go.sum",
	".*",
}

var Cmd = &cobra.Command{
	Use:   "grab [files_or_dirs...]",
	Short: "Recursively grab and display file contents from provided files or directories",
	Long: `Recursively grab and display file contents from provided files or directories.

Supports doublestar-style file selection:
- './**/*.txt' to select all .txt files in all subdirectories
- './*.*' to select all files in a single directory

Examples:
  grab ./**/*.txt  # Grab all .txt files in subdirectories
  grab ./*.go      # Grab all Go files in the current directory
  grab ./docs      # Grab all files inside the 'docs' directory`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config.Inputs = args
		if config.ExtraExclusions {
			config.ExcludePatterns = append(config.ExcludePatterns, defaultExclusions...)
		}
		RunWorker(config)
	},
}

func init() {
	Cmd.Flags().BoolVarP(&config.ListOnly, "list", "l", false, "Output list of files only")
	Cmd.Flags().StringSliceVarP(&config.ExcludePatterns, "exclude", "e", []string{}, "Exclude files matching the provided wildcard patterns")
	Cmd.Flags().BoolVarP(&config.ExtraExclusions, "exclude-defaults", "x", false, "Exclude common directories and files such as .git, node_modules, dist, and build")
}
