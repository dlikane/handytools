package common

import (
	"bytes"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestCustomFormatter_Format(t *testing.T) {
	t.Parallel()

	f := &CustomFormatter{}
	entry := &logrus.Entry{
		Level:   logrus.InfoLevel,
		Message: "test message",
	}
	out, err := f.Format(entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := string(out)
	if !strings.Contains(got, "info: test message") {
		t.Fatalf("unexpected format: %q", got)
	}
}

func TestGetLogger_ReturnsSameInstance(t *testing.T) {
	t.Parallel()

	l1 := GetLogger()
	l2 := GetLogger()
	if l1 != l2 {
		t.Fatalf("expected same instance for GetLogger()")
	}
}

func TestSetDryRunMode_AddsHook(t *testing.T) {
	t.Parallel()

	logger := GetLogger()
	before := len(logger.Hooks)
	SetDryRunMode(true)
	after := len(logger.Hooks)

	if after <= before {
		t.Fatalf("expected hook count to increase, before=%d after=%d", before, after)
	}
}

func TestDryRunHook_Fire_ModifiesMessageWhenEnabled(t *testing.T) {
	t.Parallel()

	hook := &DryRunHook{Enabled: true}
	entry := &logrus.Entry{Message: "original"}
	if err := hook.Fire(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(entry.Message, "DRYRUN: ") {
		t.Fatalf("expected message to start with DRYRUN:, got %q", entry.Message)
	}
}

func TestDryRunHook_Fire_NoChangeWhenDisabled(t *testing.T) {
	t.Parallel()

	hook := &DryRunHook{Enabled: false}
	entry := &logrus.Entry{Message: "original"}
	if err := hook.Fire(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Message != "original" {
		t.Fatalf("expected unchanged message, got %q", entry.Message)
	}
}

func TestDryRunHook_Levels_AllLevels(t *testing.T) {
	t.Parallel()

	hook := &DryRunHook{}
	got := hook.Levels()
	if len(got) != len(logrus.AllLevels) {
		t.Fatalf("expected all levels, got %v", got)
	}
}

func TestLogger_UsesCustomFormatter(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := GetLogger()
	logger.SetOutput(&buf)

	logger.Info("hello world")
	out := buf.String()
	if !strings.Contains(out, "info: hello world") {
		t.Fatalf("expected formatted output, got %q", out)
	}
}
