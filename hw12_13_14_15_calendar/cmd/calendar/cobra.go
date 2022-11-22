package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "calendar",
	Short: "calendar - a backend for the Calendar app",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.PrintErr()
		cmd.Help()

		return nil
	},
}

var configFile string

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts Calendar server",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return startServer(configFile)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		printVersion()
		return nil
	},
}

func Execute() {
	rootCmd.AddCommand(versionCmd)
	startCmd.Flags().StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
	rootCmd.AddCommand(startCmd)
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "There was an error while executing CLI '%s'", err)
		os.Exit(1)
	}
}
