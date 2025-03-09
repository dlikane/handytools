package helper

import (
	"github.com/sirupsen/logrus"
	"os"
)

var log = logrus.New()

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.Formatter = new(logrus.JSONFormatter)

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.Out = os.Stdout
}

func GetLogger() *logrus.Logger {
	return log
}
