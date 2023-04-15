package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fmnxpkg",
	Short: "📦 Tool for managing personal pacman repository written in go.",
	Long: `📦 Tool for managing personal pacman repository written in go.

	fmnxpkg - Copyright (C) 2023  Dancheg97
	
	This program comes with ABSOLUTELY NO WARRANTY; for details 'use -h flag'.
	This is free software, and you are welcome to redistribute it
	under certain conditions; watch license for details.
`,
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
