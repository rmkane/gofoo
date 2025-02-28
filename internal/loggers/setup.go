package loggers

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/rmkane/gofoo/internal/config"
	"github.com/rmkane/gofoo/internal/utils"
)

func SetupLogging(prefix string, verbose bool) (*os.File, error) {
	logDir := config.GetLoggingDir()
	logFormat := config.GetLoggingFormat()
	logLevel := config.GetLoggingLevel()

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
	case utils.FormatNameJSON:
		return slog.New(slog.NewJSONHandler(writer, &slog.HandlerOptions{Level: level})), ok
	case utils.FormatNameText:
		return slog.New(slog.NewTextHandler(writer, &slog.HandlerOptions{Level: level})), ok
	}

	return nil, false
}
