//go:build !windows

package rename

import (
	"os"
	"time"
)

func creationTime(info os.FileInfo) time.Time {
	return info.ModTime().UTC()
}
