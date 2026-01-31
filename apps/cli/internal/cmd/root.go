package cmd

import (
	"os"

	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/internal/ui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "gokickstart",
	Short:   "Scaffold a Go Kickstart project",
	Long:    "Go Kickstart scaffolds a customizable monorepo from embedded templates.",
	Version: Version,
}

const Version = "0.1.0"

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		ui.PrintError(err)
		os.Exit(1)
	}
}
