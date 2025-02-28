package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// InitConfig initializes the configuration by checking for a configuration file
// in the following order of precedence:
// 1. If a configuration file is specified using the --config or -c flag, use that.
// 2. If a configuration file exists next to the binary, use that.
// 3. If a configuration file exists in the current working directory, use that.
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

	// If a config file was loaded, print it out.
	configFileUsed := viper.ConfigFileUsed()
	if configFileUsed != "" {
		fmt.Println("Using config file:", configFileUsed)
	}

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

	// Check if the config file exists in the current working directory.
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
