package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/rmkane/gofoo/internal/loggers"
)

var logFileHandle *os.File

func AddEpilog(cmd *cobra.Command, epilog string) {
	cmd.SetHelpTemplate(fmt.Sprintf("%s\n%s\n", cmd.HelpTemplate(), epilog))
}

func NewRootCmd(appName, configName, configDir, version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:               appName,
		Version:           version,
		Long:              "A CLI tool that supports configuration and logging",
		Short:             "A CLI tool that supports configuration and logging",
		PersistentPreRun:  preRunWithAppName(appName),
		PersistentPostRun: postRun,
	}

	AddEpilog(cmd, fmt.Sprintf("Version: %s", version))

	// Determine the location for the default config to be loaded
	defaultPath := filepath.Join("$HOME", configDir, configName+".yml")
	configUsage := fmt.Sprintf("config file (default: %s)", defaultPath)

	var cfgFile string
	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", configUsage)
	_ = viper.BindPFlag("config", cmd.PersistentFlags().Lookup("config"))

	var verbose bool
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	_ = viper.BindPFlag("verbose", cmd.PersistentFlags().Lookup("verbose"))

	cmd.AddCommand(NewInitCmd(configName, configDir))
	cmd.AddCommand(NewShowCmd())

	return cmd
}

func preRunWithAppName(appName string) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		verbose := viper.GetBool("verbose")

		var err error
		logFileHandle, err = loggers.SetupLogging(appName, verbose)
		if err != nil {
			fmt.Println("Error setting up logging:", err)
			os.Exit(1)
		}
	}
}

func postRun(cmd *cobra.Command, args []string) {
	if logFileHandle != nil {
		err := logFileHandle.Close()
		if err != nil {
			panic(err)
		}
	}
}
