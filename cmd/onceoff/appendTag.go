package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	dir := flag.String("dir", ".", "Directory to scan (non-recursive)")
	tag := flag.String("tag", "", "Tag to append before .jpg (e.g., _fav)")
	apply := flag.Bool("a", false, "Apply renaming (default is dry-run)")
	flag.Parse()

	if *tag == "" {
		fmt.Println("Missing required --tag value (e.g., --tag _fav)")
		return
	}

	entries, err := os.ReadDir(*dir)
	if err != nil {
		fmt.Println("Failed to read directory:", err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		ext := strings.ToLower(filepath.Ext(name))
		if ext != ".jpg" {
			continue
		}

		base := strings.TrimSuffix(name, ext)
		if strings.HasSuffix(base, *tag) {
			continue // already tagged
		}

		newName := base + *tag + ext
		oldPath := filepath.Join(*dir, name)
		newPath := filepath.Join(*dir, newName)

		if *apply {
			if err := os.Rename(oldPath, newPath); err != nil {
				fmt.Printf("Failed to rename: %s → %s (%v)\n", oldPath, newPath, err)
			} else {
				fmt.Printf("Renamed: %s → %s\n", oldPath, newPath)
			}
		} else {
			fmt.Printf("[DRY-RUN] Would rename: %s → %s\n", oldPath, newPath)
		}
	}
}
