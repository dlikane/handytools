package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	rootDir := flag.String("root", `R:\Dropbox\Apps\my-photo-site`, "Root directory to scan")
	apply := flag.Bool("a", false, "Apply renaming (default is dry-run)")
	flag.Parse()

	// Matches filenames like: name_2005_0001.jpg
	re := regexp.MustCompile(`^([a-zA-Z0-9]+)_(\d{4})_(.+)\.jpg$`)

	err := filepath.Walk(*rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		name := filepath.Base(path)
		ext := strings.ToLower(filepath.Ext(name))
		if ext != ".jpg" {
			return nil
		}

		dir := filepath.Dir(path)

		if m := re.FindStringSubmatch(name); m != nil {
			// Rename name_2005_0001.jpg → 2005_name_0001.jpg
			newName := fmt.Sprintf("%s_%s_%s.jpg", m[2], m[1], m[3])
			newPath := filepath.Join(dir, newName)

			if *apply {
				if err := os.Rename(path, newPath); err != nil {
					fmt.Printf("Failed to rename: %s → %s (%v)\n", path, newPath, err)
				} else {
					fmt.Printf("Renamed: %s → %s\n", path, newPath)
				}
			} else {
				fmt.Printf("[DRY-RUN] Would rename: %s → %s\n", path, newPath)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error walking directory:", err)
	}
}
