package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"

	"github.com/rmkane/gofoo/pkg/model"
)

var Formats = [...]string{"json", "yaml", "toml"}

var encoderLookup = map[string]func(f *os.File, configData model.Config) error{
	"json": jsonEncoder,
	"yaml": yamlEncoder,
	"toml": tomlEncoder,
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

	if _, err = os.Stat(configPath); err == nil {
		if !force {
			return fmt.Errorf("config file already exists, use the --force flag to overwrite it")
		}
	}

	if _, err = os.Stat(configHomeDir); os.IsNotExist(err) {
		if err = os.Mkdir(configHomeDir, 0755); err != nil {
			return fmt.Errorf("cannot creating config directory: %v", err)
		}
	}

	f, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("cannot create config file: %v", err)
	}

	err = encodeConfig(f, model.DefaultConfig, format)
	if err != nil {
		return fmt.Errorf("cannot write to config file: %v", err)
	}

	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			slog.Error(err.Error())
		}
	}(f)

	fmt.Println("Successfully created config file:", configPath)

	return nil
}

func encodeConfig(f *os.File, configData model.Config, format string) error {
	encoder := encoderLookup[format]
	if encoder == nil {
		return fmt.Errorf("unsupported format: %s", format)
	}
	return encoder(f, configData)
}

func jsonEncoder(f *os.File, configData model.Config) error {
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(configData)
}

func yamlEncoder(f *os.File, configData model.Config) error {
	enc := yaml.NewEncoder(f)
	enc.SetIndent(2)
	return enc.Encode(configData)
}

func tomlEncoder(f *os.File, configData model.Config) error {
	enc := toml.NewEncoder(f)
	return enc.Encode(configData)
}
