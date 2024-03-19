package logger

import (
	"log/slog"
	"os"
	"strings"
)

var lvl slog.LevelVar

func init() {
	lvl.Set(slog.LevelInfo)
	opts := slog.HandlerOptions{
		//AddSource: true,
		Level: &lvl,
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &opts)))
}

func SetLevel(level string) {
	level = strings.ToUpper(level)
	switch level {
	case slog.LevelDebug.String():
		lvl.Set(slog.LevelDebug)
	case slog.LevelInfo.String():
		lvl.Set(slog.LevelInfo)
	case slog.LevelWarn.String():
		lvl.Set(slog.LevelWarn)
	case slog.LevelError.String():
		lvl.Set(slog.LevelError)
	}
}

func Debug(msg string, v ...any) {
	slog.Debug(msg, v...)
}

func Info(msg string, v ...any) {
	slog.Info(msg, v...)
}

func Warn(msg string, v ...any) {
	slog.Warn(msg, v...)
}

func Error(err error, msg string, v ...any) {
	slog.Error(msg, append([]interface{}{"err", err}, v...)...)
}

func FatalWithErr(err error, format string, v ...any) {
	Error(err, format, v...)
	os.Exit(1)
}

func Fatal(format string, v ...any) {
	Error(nil, format, v...)
	os.Exit(1)
}
