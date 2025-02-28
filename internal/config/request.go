package config

import (
	"fmt"
	"log/slog"

	"github.com/spf13/viper"

	"github.com/rmkane/gofoo/internal/utils"
)

var extensionLookup = map[string]string{
	"json": ".json",
	"yaml": ".yml",
	"toml": ".toml",
}

func GetLoggingLevel() slog.Level {
	levelName := viper.GetString(LoggingLevelKey)
	level, ok := utils.GetLoggingLevelByName(levelName)
	if !ok {
		fmt.Printf("Invalid logging level: %s, using default level: %s\n", levelName, DefaultConfig.Logging.Level)
	}
	return level
}

func GetLoggingDir() string {
	return viper.GetString(LoggingDirKey)
}

func GetLoggingFormat() string {
	logFormat := viper.GetString(LoggingFormatKey)
	format, ok := utils.GetLoggingFormatByName(logFormat)
	if !ok {
		fmt.Printf("Invalid logging format: %s, using default format: %s\n", logFormat, DefaultConfig.Logging.Format)
	}
	return format
}

func GetExtension(format string) string {
	return extensionLookup[format]
}
