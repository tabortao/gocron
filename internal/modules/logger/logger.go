package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gin-gonic/gin"
)

// 日志库

type Level int8

var (
	logger   *slog.Logger
	exitFunc = os.Exit
)

const (
	DEBUG = iota
	INFO
	WARN
	ERROR
	FATAL
)

func InitLogger() {
	logDir := "log"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(err)
	}

	logFile := filepath.Join(logDir, "cron.log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	writer := io.MultiWriter(os.Stdout, file)
	handler := slog.NewTextHandler(writer, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger = slog.New(handler)
}

func Debug(v ...interface{}) {
	if gin.Mode() != gin.DebugMode {
		return
	}
	write(DEBUG, v...)
}

func Debugf(format string, v ...interface{}) {
	if gin.Mode() != gin.DebugMode {
		return
	}
	writef(DEBUG, format, v...)
}

func Info(v ...interface{}) {
	write(INFO, v...)
}

func Infof(format string, v ...interface{}) {
	writef(INFO, format, v...)
}

func Warn(v ...interface{}) {
	write(WARN, v...)
}

func Warnf(format string, v ...interface{}) {
	writef(WARN, format, v...)
}

func Error(v ...interface{}) {
	write(ERROR, v...)
}

func Errorf(format string, v ...interface{}) {
	writef(ERROR, format, v...)
}

func Fatal(v ...interface{}) {
	write(FATAL, v...)
}

func Fatalf(format string, v ...interface{}) {
	writef(FATAL, format, v...)
}

func write(level Level, v ...interface{}) {
	msg := fmt.Sprint(v...)
	args := []any{}

	if gin.Mode() == gin.DebugMode {
		pc, file, line, ok := runtime.Caller(2)
		if ok {
			args = append(args, "file", file, "func", runtime.FuncForPC(pc).Name(), "line", line)
		}
	}

	switch level {
	case DEBUG:
		logger.Debug(msg, args...)
	case INFO:
		logger.Info(msg, args...)
	case WARN:
		logger.Warn(msg, args...)
	case FATAL:
		logger.Error(msg, args...)
		exitFunc(1)
	case ERROR:
		logger.Error(msg, args...)
	}
}

func writef(level Level, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	args := []any{}

	if gin.Mode() == gin.DebugMode {
		pc, file, line, ok := runtime.Caller(2)
		if ok {
			args = append(args, "file", file, "func", runtime.FuncForPC(pc).Name(), "line", line)
		}
	}

	switch level {
	case DEBUG:
		logger.Debug(msg, args...)
	case INFO:
		logger.Info(msg, args...)
	case WARN:
		logger.Warn(msg, args...)
	case FATAL:
		logger.Error(msg, args...)
		exitFunc(1)
	case ERROR:
		logger.Error(msg, args...)
	}
}
