package logger

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/gin-gonic/gin"
)

type logEntry struct {
	level slog.Level
	msg   string
}

type recordingHandler struct {
	entries []logEntry
}

func newRecordingHandler() *recordingHandler {
	return &recordingHandler{}
}

func (r *recordingHandler) Enabled(_ context.Context, level slog.Level) bool {
	return true
}

func (r *recordingHandler) Handle(_ context.Context, record slog.Record) error {
	r.entries = append(r.entries, logEntry{level: record.Level, msg: record.Message})
	return nil
}

func (r *recordingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return r
}

func (r *recordingHandler) WithGroup(name string) slog.Handler {
	return r
}

func setupRecordingLogger(t *testing.T) *recordingHandler {
	t.Helper()
	prevLogger := logger
	rec := newRecordingHandler()
	logger = slog.New(rec)
	t.Cleanup(func() { logger = prevLogger })
	return rec
}

func TestDebugLoggingDependsOnGinMode(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	rec := setupRecordingLogger(t)
	Debug("release-mode")
	if hasLevel(rec.entries, slog.LevelDebug) {
		t.Fatalf("expected no debug entry in release mode, got %+v", rec.entries)
	}

	gin.SetMode(gin.DebugMode)
	rec = setupRecordingLogger(t)
	Debug("debug-mode")
	if !hasLevel(rec.entries, slog.LevelDebug) {
		t.Fatalf("expected debug entry in debug mode, got %+v", rec.entries)
	}
}

func TestInfoLogsAndFlushes(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	rec := setupRecordingLogger(t)
	Info("info-message")
	if !hasLevel(rec.entries, slog.LevelInfo) {
		t.Fatalf("expected info entry, got %+v", rec.entries)
	}
}

func TestFatalLogsAndInvokesExit(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	rec := setupRecordingLogger(t)
	prevExit := exitFunc
	exitCalled := 0
	exitCode := 0
	exitFunc = func(code int) {
		exitCalled++
		exitCode = code
	}
	t.Cleanup(func() { exitFunc = prevExit })

	Fatal("fatal-message")
	if !hasLevel(rec.entries, slog.LevelError) {
		t.Fatalf("expected error entry, got %+v", rec.entries)
	}
	if exitCalled != 1 || exitCode != 1 {
		t.Fatalf("expected exitFunc to be called once with code 1, got count=%d code=%d", exitCalled, exitCode)
	}
}

func TestInitLogger(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, nil)
	logger = slog.New(handler)

	Info("test-message")
	if buf.Len() == 0 {
		t.Fatal("expected log output")
	}
}

func hasLevel(entries []logEntry, level slog.Level) bool {
	for _, entry := range entries {
		if entry.level == level {
			return true
		}
	}
	return false
}
