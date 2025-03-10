package common

import (
	"github.com/sirupsen/logrus"
	"os"
)

// Global logger instance
var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		DisableColors:   false, // Enable colorized output
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)
}

// GetLogger returns the global logger instance
func GetLogger() *logrus.Logger {
	return log
}

// DryRunHook is a Logrus Hook that adds "DRYRUN: " prefix to log messages
type DryRunHook struct {
	Enabled bool
}

// Fire is called for each log entry and modifies the message
func (h *DryRunHook) Fire(entry *logrus.Entry) error {
	if h.Enabled {
		entry.Message = "DRYRUN: " + entry.Message
	}
	return nil
}

// Levels specifies which log levels this hook applies to
func (h *DryRunHook) Levels() []logrus.Level {
	return logrus.AllLevels // Apply to all log levels
}

// SetDryRunMode enables or disables the DRYRUN prefix
func SetDryRunMode(enabled bool) {
	log.AddHook(&DryRunHook{Enabled: enabled})
}
