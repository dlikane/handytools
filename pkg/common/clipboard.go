package common

import (
	"github.com/atotto/clipboard"
	"github.com/sirupsen/logrus"
)

// CopyToClipboard copies the provided text to the system clipboard.
func CopyToClipboard(text string) {
	if err := clipboard.WriteAll(text); err != nil {
		logrus.WithError(err).Error("Failed to copy to clipboard")
	} else {
		logrus.Info("Copied to clipboard successfully.")
	}
}
