package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	rootDir := flag.String("root", `R:\Dropbox\Apps\my-photo-site`, "Root directory to search in")
	fromDir := flag.String("from", "b&w", "Source subdirectory (e.g., b&w)")
	toSuffix := flag.String("to", "b&w_fav", "Suffix to add to matched files")
	apply := flag.Bool("a", false, "Apply changes (default is dry-run)")
	flag.Parse()

	sourcePath := filepath.Join(*rootDir, *fromDir)
	sourceFiles := map[string]string{}

	err := filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		sourceFiles[base] = path
		return nil
	})
	if err != nil {
		fmt.Println("Failed to walk source dir:", err)
		return
	}

	found := map[string]string{}
	reverseLookup := map[string]string{} // base -> original file path
	notFound := []string{}

	err = filepath.Walk(*rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || strings.Contains(path, *fromDir) {
			return nil
		}
		base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		if origPath, ok := sourceFiles[base]; ok {
			newName := base + "_" + *toSuffix + filepath.Ext(path)
			newPath := filepath.Join(filepath.Dir(path), newName)
			found[path] = newPath
			reverseLookup[base] = origPath
		}
		return nil
	})
	if err != nil {
		fmt.Println("Failed to walk root dir:", err)
		return
	}

	for base, _ := range sourceFiles {
		if _, ok := reverseLookup[base]; !ok {
			notFound = append(notFound, base)
		}
	}

	for oldPath, newPath := range found {
		base := strings.TrimSuffix(filepath.Base(oldPath), filepath.Ext(oldPath))
		origPath := reverseLookup[base]
		if *apply {
			if err := os.Rename(oldPath, newPath); err != nil {
				fmt.Printf("Failed to rename: %s → %s (%v)\n", oldPath, newPath, err)
			} else {
				fmt.Printf("Renamed: %s → %s\n", oldPath, newPath)
				if err := os.Remove(origPath); err != nil {
					fmt.Printf("  Failed to delete original: %s (%v)\n", origPath, err)
				} else {
					fmt.Printf("  Deleted original: %s\n", origPath)
				}
			}
		} else {
			fmt.Printf("[DRY-RUN] Would rename: %s → %s\n", oldPath, newPath)
			fmt.Printf("[DRY-RUN] Would delete original: %s\n", origPath)
		}
	}

	if len(notFound) > 0 {
		fmt.Println("\n[Not found in other dirs]:")
		for _, base := range notFound {
			fmt.Println("  ", base)
		}
	}
}
