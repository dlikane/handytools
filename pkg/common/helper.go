package common

import "path/filepath"

func ExpandWildcards(patterns []string) []string {
	logger := GetLogger()

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
