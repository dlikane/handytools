//go:build windows

package rename

import (
	"os"
	"syscall"
	"time"
)

func creationTime(info os.FileInfo) time.Time {
	if fa, ok := info.Sys().(*syscall.Win32FileAttributeData); ok {
		return time.Unix(0, fa.CreationTime.Nanoseconds()).UTC()
	}
	return info.ModTime().UTC()
}
