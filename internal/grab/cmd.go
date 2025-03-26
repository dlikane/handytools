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
	"_*",
}

var Cmd = &cobra.Command{
	Use:   "grab [files_or_dirs...]",
	Short: "Recursively grab and display file contents from provided files or directories",
	Long: `Recursively grab and display file contents from provided files or directories.

⚠️ When using wildcards (* or ...), wrap arguments in quotes to prevent shell expansion.

Go-style file selection examples:

Single Directory (Non-Recursive):
  ".*"             # All files in the current directory
  "./*.go"         # All .go files in the current directory
  "./*.*"          # All files with an extension in the current directory

Recursive (All Subdirectories):
  "./.../*"        # All files in current directory and all subdirectories
  "./.../*.go"     # All .go files in subdirectories
  "src/.../*"      # All files in 'src/' and its subdirectories
  "src/.../*.txt"  # All .txt files in 'src/' and its subdirectories

Specific File Types:
  "./.../*.md"     # All Markdown files
  "./.../*.json"   # All JSON files
  "./.../*.yaml"   # All YAML files

Files with Specific Prefixes/Suffixes:
  "./config*"         # Files starting with 'config' in current dir
  "./.../config*"     # Files starting with 'config' in all subdirs
  "./.../*_test.go"   # All Go test files in all subdirs

To exclude files/directories:
  -e "*.log" -e "node_modules/..." -x

Examples:
  grab "./.../*.txt"
  grab "./*.go"
  grab "./apps/kyc-service/..." -e "node_modules/..." -l
`,
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
