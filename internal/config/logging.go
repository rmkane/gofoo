package config

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"
)

const (
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

func SetupLogging(prefix string, verbose bool) (*os.File, error) {
	logDir := GetLoggingDir()
	logFormat := GetLoggingFormat()
	logLevel := GetLoggingLevel()

	// Ensure log directory exists
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.MkdirAll(logDir, 0755)
	}

	now := time.Now().Format("2006-01-02")
	logFileName := fmt.Sprintf("%s/%s-%s_%s.log", logDir, prefix, now, logFormat)
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	// Set up slog handler based on format
	logger, ok := getLoggerForFormat(logFormat, logLevel, logFile)
	if logger == nil {
		return nil, fmt.Errorf("unsupported log format: %s", logFormat)
	}

	if verbose && ok {
		fmt.Println("Logging to: ", logFileName)
	}

	slog.SetDefault(logger)

	return logFile, nil
}

func getLoggerForFormat(format string, level slog.Level, logFile *os.File) (*slog.Logger, bool) {
	ok := false
	var writer io.Writer
	if logFile != nil {
		writer = logFile
		ok = true
	} else {
		writer = os.Stdout
	}

	switch format {
	case FormatNameJSON:
		return slog.New(slog.NewJSONHandler(writer, &slog.HandlerOptions{Level: level})), ok
	case FormatNameText:
		return slog.New(slog.NewTextHandler(writer, &slog.HandlerOptions{Level: level})), ok
	}

	return nil, false
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
