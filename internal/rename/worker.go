package rename

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	"handytools/pkg/common"
)

type fileItem struct {
	Path       string
	Name       string
	CreateTime time.Time
	ModTime    time.Time
}

// Cross-platform creation time.
// Windows: uses CreationTime
// Others: falls back to ModTime()
func creationTime(info os.FileInfo) time.Time {
	if fa, ok := info.Sys().(*syscall.Win32FileAttributeData); ok {
		ns := fa.CreationTime.Nanoseconds()
		return time.Unix(0, ns).UTC()
	}
	return info.ModTime().UTC()
}

func sortFiles(paths []string, sortBy, order string) ([]fileItem, error) {
	items := make([]fileItem, 0, len(paths))
	for _, p := range paths {
		info, err := os.Stat(p)
		if err != nil {
			common.GetLogger().WithError(err).Warnf("Skipping (stat failed): %s", p)
			continue
		}
		items = append(items, fileItem{
			Path:       p,
			Name:       strings.ToLower(filepath.Base(p)),
			CreateTime: creationTime(info),
			ModTime:    info.ModTime().UTC(),
		})
	}

	less := func(i, j int) bool {
		switch sortBy {
		case "name":
			if items[i].Name == items[j].Name {
				if order == "asc" {
					return items[i].CreateTime.Before(items[j].CreateTime)
				}
				return items[i].CreateTime.After(items[j].CreateTime)
			}
			if order == "asc" {
				return items[i].Name < items[j].Name
			}
			return items[i].Name > items[j].Name

		case "modified":
			if items[i].ModTime.Equal(items[j].ModTime) {
				if order == "asc" {
					return items[i].Name < items[j].Name
				}
				return items[i].Name > items[j].Name
			}
			if order == "asc" {
				return items[i].ModTime.Before(items[j].ModTime)
			}
			return items[i].ModTime.After(items[j].ModTime)

		case "created":
			if items[i].CreateTime.Equal(items[j].CreateTime) {
				if order == "asc" {
					return items[i].Name < items[j].Name
				}
				return items[i].Name > items[j].Name
			}
			if order == "asc" {
				return items[i].CreateTime.Before(items[j].CreateTime)
			}
			return items[i].CreateTime.After(items[j].CreateTime)
		}

		// default stable fallback
		return items[i].Path < items[j].Path
	}

	sort.Slice(items, less)
	return items, nil
}

func renameFiles(cfg RenameConfig) {
	logger := common.GetLogger()

	if !cfg.Apply {
		common.SetDryRunMode(true)
		logger.Info("Running in DRYRUN mode")
	}

	// Sort input files
	items, _ := sortFiles(cfg.InputFiles, cfg.SortBy, cfg.Order)

	counter := 1
	for _, it := range items {
		dir := filepath.Dir(it.Path)
		ext := filepath.Ext(it.Path)
		newName := fmt.Sprintf("%s_%04d%s", cfg.OutputName, counter, ext)
		newPath := filepath.Join(dir, newName)

		logger.Infof("Rename: %s -> %s", it.Path, newPath)

		if cfg.Apply {
			if err := os.Rename(it.Path, newPath); err != nil {
				logger.WithError(err).Errorf("Failed to rename file: %s", it.Path)
			}
		}
		counter++
	}
}
