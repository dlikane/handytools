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
		inputPaths := expandWildcards(args)
		grabFiles(inputPaths)
	},
}

func init() {
	Cmd.Flags().BoolVarP(&listOnly, "list", "l", false, "Output list of files only")
	Cmd.Flags().StringSliceVarP(&excludePatterns, "exclude", "e", []string{}, "Exclude files matching the provided wildcard patterns")
}

func expandWildcards(paths []string) []string {
	var expanded []string

	for _, path := range paths {
		if strings.Contains(path, "*") {
			var baseDir, pattern string

			if strings.Contains(path, "...") {
				baseDir = strings.Split(path, "...")[0]
				pattern = strings.TrimPrefix(path, baseDir+"...")
				if pattern == "" {
					pattern = "*"
				}
				filepath.Walk(baseDir, func(fp string, fi os.FileInfo, err error) error {
					if err != nil || fi.IsDir() {
						return nil
					}
					if matched, _ := filepath.Match(pattern, filepath.Base(fp)); matched {
						expanded = append(expanded, fp)
					}
					return nil
				})
			} else {
				baseDir = filepath.Dir(path)
				pattern = filepath.Base(path)
				files, _ := filepath.Glob(filepath.Join(baseDir, pattern))
				expanded = append(expanded, files...)
			}
		} else {
			expanded = append(expanded, path)
		}
	}
	return expanded
}

func matchesExcludePattern(file string) bool {
	for _, pattern := range excludePatterns {
		if matched, _ := filepath.Match(pattern, filepath.Base(file)); matched {
			return true
		}
	}
	return false
}

func grabFiles(paths []string) {
	var collectedFiles []string
	var excludedFiles []string
	var totalLines int
	fileLineCounts := make(map[string]int)

	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			fmt.Printf("Error accessing %s: %v\n", path, err)
			continue
		}

		if info.IsDir() {
			filepath.Walk(path, func(fp string, fi os.FileInfo, err error) error {
				if err != nil || fi.IsDir() {
					return nil
				}
				if matchesExcludePattern(fp) {
					excludedFiles = append(excludedFiles, fp)
					return nil
				}
				collectedFiles = append(collectedFiles, fp)
				return nil
			})
		} else {
			if matchesExcludePattern(path) {
				excludedFiles = append(excludedFiles, path)
				continue
			}
			collectedFiles = append(collectedFiles, path)
		}
	}

	for _, file := range collectedFiles {
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

	fmt.Printf("\nProvided (%d files %d lines):\n", len(collectedFiles), totalLines)
	for _, file := range collectedFiles {
		fmt.Printf("%s (%d)\n", file, fileLineCounts[file])
	}

	if len(excludedFiles) > 0 {
		fmt.Println("Excluded:")
		for _, file := range excludedFiles {
			fmt.Printf("- %s\n", file)
		}
	}
}
