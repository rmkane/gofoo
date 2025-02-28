package utils

import (
	"log/slog"
	"strings"
)

var (
	LevelNameDebug = "DEBUG"
	LevelNameInfo  = "INFO"
	LevelNameWarn  = "WARN"
	LevelNameError = "ERROR"

	FormatNameText = "text"
	FormatNameJSON = "json"
)

var nameToLevel = map[string]slog.Level{
	LevelNameDebug: slog.LevelDebug,
	LevelNameInfo:  slog.LevelInfo,
	LevelNameWarn:  slog.LevelWarn,
	LevelNameError: slog.LevelError,
}

var supportedLoggingFormats = map[string]bool{
	FormatNameText: true,
	FormatNameJSON: true,
}

func GetLoggingLevelByName(name string) (slog.Level, bool) {
	if level, ok := nameToLevel[strings.ToUpper(name)]; ok {
		return level, true
	}
	return slog.LevelInfo, false
}

func GetLoggingFormatByName(logFormat string) (string, bool) {
	format := strings.ToLower(logFormat)
	if _, ok := supportedLoggingFormats[logFormat]; ok {
		return format, true
	}
	return FormatNameText, false
}
