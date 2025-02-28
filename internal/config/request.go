package config

import (
	"fmt"
	"log/slog"

	"github.com/spf13/viper"

	"github.com/rmkane/gofoo/internal/utils"
	"github.com/rmkane/gofoo/pkg/model"
)

var extensionLookup = map[string]string{
	"json": ".json",
	"yaml": ".yml",
	"toml": ".toml",
}

func ShowConfig() {
	fmt.Printf("Showing configuration:\n\n")

	fmt.Printf("- %-15s: %s\n", model.LoggingDirKey, GetLoggingDir())
	fmt.Printf("- %-15s: %s\n", model.LoggingFormatKey, GetLoggingFormat())
	fmt.Printf("- %-15s: %s\n", model.LoggingLevelKey, GetLoggingLevel())
}

// setDefaultConfig sets the default configuration values.
func setDefaultConfig() {
	viper.SetDefault(model.LoggingDirKey, model.DefaultConfig.Logging.Dir)
	viper.SetDefault(model.LoggingFormatKey, model.DefaultConfig.Logging.Format)
	viper.SetDefault(model.LoggingLevelKey, model.DefaultConfig.Logging.Level)
}

func GetLoggingLevel() slog.Level {
	levelName := viper.GetString(model.LoggingLevelKey)
	level, ok := utils.GetLoggingLevelByName(levelName)
	if !ok {
		fmt.Printf("Invalid logging level: %s, using default level: %s\n", levelName, model.DefaultConfig.Logging.Level)
	}
	return level
}

func GetLoggingDir() string {
	return viper.GetString(model.LoggingDirKey)
}

func GetLoggingFormat() string {
	logFormat := viper.GetString(model.LoggingFormatKey)
	format, ok := utils.GetLoggingFormatByName(logFormat)
	if !ok {
		fmt.Printf("Invalid logging format: %s, using default format: %s\n", logFormat, model.DefaultConfig.Logging.Format)
	}
	return format
}

func GetExtension(format string) string {
	return extensionLookup[format]
}
