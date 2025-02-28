package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/rmkane/gofoo/cmd/gofoo/commands"
	"github.com/rmkane/gofoo/internal/config"
)

const (
	AppName    = "gofoo"
	ConfigDir  = ".gofoo"
	ConfigName = "config"
)

var Version = "dev"

func main() {
	// Set the version

	// Initialize viper config first, so that configuration is available in the cobra commands
	cobra.OnInitialize(initializeConfig)

	// Create and execute the root command
	rootCmd := commands.NewRootCmd(AppName, ConfigName, ConfigDir, Version)
	rootCmd.Execute()
}

func initializeConfig() {
	err := config.InitConfig(ConfigName, ConfigDir)
	if err != nil {
		fmt.Println("Error initializing config:", err)
		os.Exit(1)
	}
}
