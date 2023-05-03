package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "repo",
	Short: "📦 Personal repository for arch packages written in go",
}

func Execute() {
	for _, flag := range flags {
		AddFlag(flag)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
