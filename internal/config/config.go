package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"

	"github.com/BurntSushi/toml"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var Formats = [...]string{"json", "yaml", "toml"}

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

var encoderLookup = map[string]func(f *os.File, configData Config) error{
	"json": jsonEncoder,
	"yaml": yamlEncoder,
	"toml": tomlEncoder,
}

var extensionLookup = map[string]string{
	"json": ".json",
	"yaml": ".yml",
	"toml": ".toml",
}

func GetLoggingLevel() slog.Level {
	levelName := viper.GetString(LoggingLevelKey)
	level, ok := GetLoggingLevelByName(levelName)
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
	format, ok := GetLoggingFormatByName(logFormat)
	if !ok {
		fmt.Printf("Invalid logging format: %s, using default format: %s\n", logFormat, DefaultConfig.Logging.Format)
	}
	return format
}

func GetExtension(format string) string {
	return extensionLookup[format]
}

func EncodeConfig(f *os.File, configData Config, format string) error {
	encoder := encoderLookup[format]
	if encoder == nil {
		return fmt.Errorf("unsupported format: %s", format)
	}
	return encoder(f, configData)
}

func CreateConfig(configName, configDir, format string, force bool) error {
	if !slices.Contains(Formats[:], format) {
		return fmt.Errorf("unsupported format: %s", format)
	}

	// Create the config in the home directory if it does not exist. If the force flag is used, overwrite it.
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot determine home directory: %v", err)
	}

	ext := GetExtension(format)
	if ext == "" {
		return fmt.Errorf("unsupported format: %s", format)
	}

	configHomeDir := filepath.Join(homeDir, configDir)
	configFile := fmt.Sprintf("%s%s", configName, ext)
	configPath := filepath.Join(configHomeDir, configFile)

	if _, err := os.Stat(configPath); err == nil {
		if !force {
			return fmt.Errorf("config file already exists, use the --force flag to overwrite it")
		}
	}

	if _, err := os.Stat(configHomeDir); os.IsNotExist(err) {
		if err := os.Mkdir(configHomeDir, 0755); err != nil {
			return fmt.Errorf("cannot creating config directory: %v", err)
		}
	}

	f, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("cannot create config file: %v", err)
	}

	err = EncodeConfig(f, DefaultConfig, format)
	if err != nil {
		return fmt.Errorf("cannot write to config file: %v", err)
	}

	defer f.Close()

	fmt.Println("Successfully created config file:", configPath)

	return nil
}

func ShowConfig() {
	fmt.Printf("Showing configuration:\n\n")

	fmt.Printf("- %-15s: %s\n", LoggingDirKey, GetLoggingDir())
	fmt.Printf("- %-15s: %s\n", LoggingFormatKey, GetLoggingFormat())
	fmt.Printf("- %-15s: %s\n", LoggingLevelKey, GetLoggingLevel())
}

// InitConfig initializes the configuration by checking for a configuration file
// in the following order of precedence:
// 1. If a configuration file is specified using the --config or -c flag, use that.
// 2. If a configuration file exists next to the binary, use that.
// 3. If a configuration file exists in the project directory, use that.
// 4. If a configuration file exists in the home directory, use that.
// If no configuration file is found, load the default configuration.
//
// For checks 1-4, the configuration file can have the following extensions:
// .yaml, .yml, .json, .toml.
func InitConfig(configName, configDir string) error {
	verbose := viper.GetBool("verbose")
	cfgFile := viper.GetString("config")

	if cfgFile != "" {
		// Validate the config file.
		if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
			return fmt.Errorf("config file not found: %s", cfgFile)
		}

		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find the config file.
		configFile, err := findConfigFile(configName, configDir)
		if err == nil {
			if verbose {
				fmt.Printf("Found config file: %s\n", configFile)
			}
			return nil
		}
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If no config file is found, use the default configuration.
	fmt.Println("No config file found, using default configuration")
	setDefaultConfig()

	fmt.Println("Using config file:", viper.ConfigFileUsed())

	return nil
}

func findConfigFile(configName, configDir string) (string, error) {
	// Check if the config file exists next to the binary.
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("could not find executable path: %v", err)
	}
	exeDir := filepath.Dir(exePath)
	exeConfigFile, ok := checkConfigFiles(exeDir, configName)
	if ok {
		return exeConfigFile, nil
	}

	// Check if the config file exists in the current directory.
	currDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("could not get current directory: %v", err)
	}
	currDirConfigFile, ok := checkConfigFiles(currDir, configName)
	if ok {
		return currDirConfigFile, nil
	}

	// Find home directory.
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find home directory: %v", err)
	}
	configPath := filepath.Join(homeDir, configDir)
	configFile, ok := checkConfigFiles(configPath, configName)
	if ok {
		return configFile, nil
	}

	return "", fmt.Errorf("no config file found")
}

// checkConfigFiles checks for config files with different extensions in the given directory.
func checkConfigFiles(dir, configName string) (string, bool) {
	extensions := []string{"yaml", "yml", "json", "toml"}
	for _, ext := range extensions {
		configFile := configName + "." + ext
		configPath := filepath.Join(dir, configFile)
		if _, err := os.Stat(configPath); err != nil {
			continue
		}
		if err := loadConfig(configPath, configName); err == nil {
			return configPath, true
		}
	}
	return "", false
}

func loadConfig(configPath string, configName string) error {
	dir := filepath.Dir(configPath)
	viper.AddConfigPath(dir)
	viper.SetConfigName(configName)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	return nil
}

// setDefaultConfig sets the default configuration values.
func setDefaultConfig() {
	viper.SetDefault(LoggingDirKey, DefaultConfig.Logging.Dir)
	viper.SetDefault(LoggingFormatKey, DefaultConfig.Logging.Format)
	viper.SetDefault(LoggingLevelKey, DefaultConfig.Logging.Level)
}

func jsonEncoder(f *os.File, configData Config) error {
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(configData)
}

func yamlEncoder(f *os.File, configData Config) error {
	enc := yaml.NewEncoder(f)
	enc.SetIndent(2)
	return enc.Encode(configData)
}

func tomlEncoder(f *os.File, configData Config) error {
	enc := toml.NewEncoder(f)
	return enc.Encode(configData)
}
