package grab

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Inputs          []string
	ListOnly        bool
	ExcludePatterns []string
	ExtraExclusions bool
}

func RunWorker(config Config) {
	inclusionList := expandWildcards(config.Inputs)
	exclusionList := expandWildcards(config.ExcludePatterns)
	filteredPaths := filterExcluded(inclusionList, exclusionList)
	grabFiles(filteredPaths, config.ListOnly)
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

func grabFiles(paths []string, listOnly bool) {
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
