package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"gofoo/internal/config"
)

const (
	AppName    = "gofoo"
	ConfigDir  = ".gofoo"
	ConfigName = "config"
)

var logFileHandle *os.File

func main() {
	// Initialize viper config first, so that configuration is available in the cobra commands
	cobra.OnInitialize(initializeConfig)

	// Create and execute the root command
	rootCmd := NewRootCmd()
	rootCmd.Execute()
}

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               AppName,
		PersistentPreRun:  preRun,
		PersistentPostRun: postRun,
	}

	var cfgFile string
	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.gofoo/config.yml)")
	_ = viper.BindPFlag("config", cmd.PersistentFlags().Lookup("config"))

	var verbose bool
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	_ = viper.BindPFlag("verbose", cmd.PersistentFlags().Lookup("verbose"))

	cmd.AddCommand(NewInitCmd())
	cmd.AddCommand(NewShowCmd())

	return cmd
}

func initializeConfig() {
	err := config.InitConfig(ConfigName, ConfigDir)
	if err != nil {
		fmt.Println("Error initializing config:", err)
		os.Exit(1)
	}
}

func preRun(cmd *cobra.Command, args []string) {
	verbose := viper.GetBool("verbose")

	var err error
	logFileHandle, err = config.SetupLogging(AppName, verbose)
	if err != nil {
		fmt.Println("Error setting up logging:", err)
		os.Exit(1)
	}

}

func postRun(cmd *cobra.Command, args []string) {
	if logFileHandle != nil {
		logFileHandle.Close()
	}
}

func NewInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize the configuration file",
		Run: func(cmd *cobra.Command, args []string) {
			format, _ := cmd.Flags().GetString("format")
			force, _ := cmd.Flags().GetBool("force")

			err := config.CreateConfig(ConfigName, ConfigDir, format, force)
			if err != nil {
				fmt.Println("Error creating config:", err)
				os.Exit(1)
			}

			// Print to log
			slog.Debug("Log level DEBUG")
			slog.Info("Log level INFO")
			slog.Warn("Log level WARN")
			slog.Error("Log level ERROR")
		},
	}

	cmd.Flags().StringP("format", "f", "yaml", "config format: json, yaml, toml")
	cmd.Flags().BoolP("force", "", false, "overwrite the configuration file if it exists")

	return cmd
}

func NewShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show the configuration",
		Run: func(cmd *cobra.Command, args []string) {
			config.ShowConfig()

			// Print to log
			slog.Debug("Log level DEBUG")
			slog.Info("Log level INFO")
			slog.Warn("Log level WARN")
			slog.Error("Log level ERROR")
		},
	}
	return cmd
}
