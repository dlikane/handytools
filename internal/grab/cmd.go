package grab

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var listOnly bool
var excludePatterns []string

var Cmd = &cobra.Command{
	Use:   "grab [files_or_dirs...]",
	Short: "Recursively grab and display file contents from provided files or directories",
	Long: `Recursively grab and display file contents from provided files or directories.

Supports Go-style file selection:
- './.../*.txt' to select all .txt files in all subdirectories
- './*.*' to select all files in a single directory

Examples:
  grab ./.../*.txt  # Grab all .txt files in subdirectories
  grab ./*.go       # Grab all Go files in the current directory
  grab ./docs       # Grab all files inside the 'docs' directory`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		inclusionList := expandWildcards(args)
		exclusionList := expandWildcards(excludePatterns)
		filteredPaths := filterExcluded(inclusionList, exclusionList)
		grabFiles(filteredPaths)
	},
}

func init() {
	Cmd.Flags().BoolVarP(&listOnly, "list", "l", false, "Output list of files only")
	Cmd.Flags().StringSliceVarP(&excludePatterns, "exclude", "e", []string{}, "Exclude files matching the provided wildcard patterns")
}

func expandWildcards(paths []string) []string {
	var expanded []string

	for _, path := range paths {
		path = filepath.ToSlash(filepath.Clean(path))

		if strings.Contains(path, "...") {
			baseDir := strings.Split(path, "...")[0]
			if baseDir == "" {
				baseDir = "."
			}

			_ = filepath.Walk(baseDir, func(fp string, fi os.FileInfo, err error) error {
				if err != nil {
					return nil
				}
				if !fi.IsDir() {
					expanded = append(expanded, fp)
				}
				return nil
			})
		} else if strings.Contains(path, "*") {
			files, _ := filepath.Glob(path)
			expanded = append(expanded, files...)
		} else {
			expanded = append(expanded, path)
		}
	}
	return expanded
}

func filterExcluded(included, excluded []string) []string {
	filtered := []string{}
	excludedSet := make(map[string]bool)
	for _, file := range excluded {
		excludedSet[file] = true
	}

	for _, file := range included {
		if !excludedSet[file] {
			filtered = append(filtered, file)
		}
	}
	return filtered
}

func grabFiles(paths []string) {
	var totalLines int
	fileLineCounts := make(map[string]int)

	for _, file := range paths {
		content, err := ioutil.ReadFile(file)
		if err == nil {
			lines := strings.Count(string(content), "\n") + 1
			totalLines += lines
			fileLineCounts[file] = lines
			if !listOnly {
				fmt.Println(file)
				fmt.Println(string(content))
				fmt.Println()
			}
		}
	}

	fmt.Printf("\nProvided (%d files %d lines):\n", len(paths), totalLines)
	for _, file := range paths {
		fmt.Printf("%s (%d)\n", file, fileLineCounts[file])
	}
}
