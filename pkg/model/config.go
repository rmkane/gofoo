package model

var (
	LoggingLevelKey  = "logging.level"
	LoggingDirKey    = "logging.dir"
	LoggingFormatKey = "logging.format"
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
