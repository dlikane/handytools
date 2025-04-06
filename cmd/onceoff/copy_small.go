package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	sourceRoot := `R:\Dropbox\Apps\my-photo-site`
	destDir := `C:\tmp\small`

	if err := os.MkdirAll(destDir, 0755); err != nil {
		fmt.Println("Failed to create destination folder:", err)
		return
	}

	err := filepath.Walk(sourceRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".jpg" {
			return nil
		}

		base := filepath.Base(path)
		if strings.Contains(strings.ToLower(base), "_small") {
			destPath := filepath.Join(destDir, base)

			srcFile, err := os.Open(path)
			if err != nil {
				fmt.Printf("Failed to open source: %s (%v)\n", path, err)
				return nil
			}
			defer srcFile.Close()

			destFile, err := os.Create(destPath)
			if err != nil {
				fmt.Printf("Failed to create dest: %s (%v)\n", destPath, err)
				return nil
			}
			defer destFile.Close()

			if _, err := io.Copy(destFile, srcFile); err != nil {
				fmt.Printf("Failed to copy %s → %s (%v)\n", path, destPath, err)
			} else {
				fmt.Printf("Copied: %s → %s\n", path, destPath)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error during file walk:", err)
	}
}
