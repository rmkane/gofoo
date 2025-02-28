package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type LoggingConfig struct {
	Level  string `yaml:"level" json:"level" toml:"level"`
	Dir    string `yaml:"dir" json:"dir" toml:"dir"`
	Format string `yaml:"format" json:"format" toml:"format"`
}

type Config struct {
	Logging LoggingConfig `yaml:"logging" json:"logging" toml:"logging"`
}

var DefaultConfig = Config{
	Logging: LoggingConfig{
		Level:  "INFO",
		Dir:    "./logs",
		Format: "text",
	},
}

var LoggingLevelKey = "logging.level"

var LoggingDirKey = "logging.dir"

var LoggingFormatKey = "logging.format"

func ShowConfig() {
	fmt.Printf("Showing configuration:\n\n")

	fmt.Printf("- %-15s: %s\n", LoggingDirKey, GetLoggingDir())
	fmt.Printf("- %-15s: %s\n", LoggingFormatKey, GetLoggingFormat())
	fmt.Printf("- %-15s: %s\n", LoggingLevelKey, GetLoggingLevel())
}

// setDefaultConfig sets the default configuration values.
func setDefaultConfig() {
	viper.SetDefault(LoggingDirKey, DefaultConfig.Logging.Dir)
	viper.SetDefault(LoggingFormatKey, DefaultConfig.Logging.Format)
	viper.SetDefault(LoggingLevelKey, DefaultConfig.Logging.Level)
}
