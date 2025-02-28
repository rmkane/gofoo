package commands

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/rmkane/gofoo/internal/config"
)

func NewInitCmd(configName, configDir string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize the configuration file",
		Run: func(cmd *cobra.Command, args []string) {
			format, _ := cmd.Flags().GetString("format")
			force, _ := cmd.Flags().GetBool("force")

			err := config.CreateConfig(configName, configDir, format, force)
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
