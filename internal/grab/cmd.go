package grab

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var Cmd = &cobra.Command{
	Use:   "grab [files_or_dirs...]",
	Short: "Recursively grab and display file contents from provided files or directories",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		inputPaths := expandWildcards(args)
		grabFiles(inputPaths)
	},
}

func expandWildcards(paths []string) []string {
	var expanded []string

	for _, path := range paths {
		if strings.Contains(path, "*") {
			root := filepath.Dir(path)

			pattern := filepath.Base(path)
			filepath.Walk(root, func(fp string, fi os.FileInfo, err error) error {
				if err != nil || fi.IsDir() {
					return nil
				}
				if matched, _ := filepath.Match(pattern, fi.Name()); matched {
					expanded = append(expanded, fp)
				}
				return nil
			})
		} else {
			expanded = append(expanded, path)
		}
	}
	return expanded
}

func grabFiles(paths []string) {
	var collectedFiles []string

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
				fmt.Println(fp)
				content, err := ioutil.ReadFile(fp)
				if err == nil {
					fmt.Println(string(content))
				}
				collectedFiles = append(collectedFiles, fp)
				return nil
			})
		} else {
			fmt.Println(path)
			content, err := ioutil.ReadFile(path)
			if err == nil {
				fmt.Println(string(content))
				fmt.Println()
				fmt.Println()
			}
			collectedFiles = append(collectedFiles, path)
		}
	}

	fmt.Println("\nProvided:")
	for _, file := range collectedFiles {
		fmt.Println(file)
	}
}
