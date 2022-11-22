package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "calendar",
	Short: "calendar - a backend for the Calendar app",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var configFile string

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts Calendar server",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		startServer(configFile)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

func Execute() {
	rootCmd.AddCommand(versionCmd)
	startCmd.Flags().StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
	rootCmd.AddCommand(startCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
