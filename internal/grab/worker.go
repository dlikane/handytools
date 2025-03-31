package grab

import (
	"fmt"
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
	exclusionMap := make(map[string]struct{})
	for _, e := range exclusionList {
		exclusionMap[e] = struct{}{}
	}

	var filteredPaths []string
	for _, incl := range inclusionList {
		skip := false
		for excl := range exclusionMap {
			if strings.HasPrefix(incl, excl+"/") || strings.HasPrefix(incl, excl) || strings.Contains(incl, "/"+excl+"/") {
				skip = true
				break
			}
		}
		if !skip {
			filteredPaths = append(filteredPaths, incl)
		}
	}
	grabFiles(filteredPaths, config.ListOnly)
}

func expandWildcards(paths []string) []string {
	var expanded []string

	for _, path := range paths {
		path = filepath.ToSlash(filepath.Clean(path))

		if strings.Contains(path, "...") {
			parts := strings.Split(path, "...")
			baseDir := filepath.Clean(parts[0])
			if baseDir == "" {
				baseDir = "."
			}
			suffix := ""
			if len(parts) > 1 {
				suffix = parts[1]
			}

			_ = filepath.Walk(baseDir, func(fp string, fi os.FileInfo, err error) error {
				if err != nil {
					return nil
				}
				if fi.IsDir() {
					return nil
				}
				if strings.HasSuffix(filepath.ToSlash(fp), suffix) {
					expanded = append(expanded, filepath.ToSlash(fp))
				}
				return nil
			})
		} else if strings.Contains(path, "*") {
			files, _ := filepath.Glob(path)
			for _, f := range files {
				expanded = append(expanded, filepath.ToSlash(f))
			}
		} else {
			expanded = append(expanded, filepath.ToSlash(path))
		}
	}

	return expanded
}

func grabFiles(paths []string, listOnly bool) {
	var totalLines int
	fileLineCounts := make(map[string]int)

	for _, file := range paths {
		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		if !isText(content) {
			continue
		}

		lines := strings.Count(string(content), "\n") + 1
		totalLines += lines
		fileLineCounts[file] = lines

		if !listOnly {
			fmt.Println(file)
			fmt.Println(string(content))
			fmt.Println()
		}
	}

	fmt.Printf("\nProvided (%d files %d lines):\n", len(fileLineCounts), totalLines)
	for file, count := range fileLineCounts {
		fmt.Printf("%s (%d)\n", file, count)
	}
	fmt.Printf("\nTotal (%d files %d lines):\n", len(fileLineCounts), totalLines)
}

func isText(data []byte) bool {
	for _, b := range data {
		if b == 0 {
			return false // null byte → likely binary
		}
		if b < 0x09 {
			return false
		}
	}
	return true
}
