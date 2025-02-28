package commands

import (
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/rmkane/gofoo/internal/config"
)

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
