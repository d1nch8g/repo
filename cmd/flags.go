// Copyright 2023 FMNX Linux team.
// This code is covered by GPL license, which can be found in LICENSE file.
// Additional information could be found on official web page: https://fmnx.io/
// Email: help@fmnx.io
package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var flags = []Flag{
	{
		Cmd:         rootCmd,
		Name:        "repo",
		ShortName:   "r",
		Env:         "PACKREPO_REPO",
		Value:       "localhost",
		Description: "📄 repository name on the web page",
	},
	{
		Cmd:         rootCmd,
		Name:        "port",
		ShortName:   "p",
		Env:         "PACKREPO_PORT",
		Value:       "80",
		Type:        "int",
		Description: "🌐 publically exposed port for both HTTP and gRPC calls",
	},
	{
		Cmd:         rootCmd,
		Name:        "init-pkgs",
		ShortName:   "i",
		Env:         "PACKREPO_INIT_PKGS",
		Description: "📦 initial packages for download",
	},
	{
		Cmd:         rootCmd,
		Name:        "web-dir",
		ShortName:   "w",
		Env:         "PACKREPO_WEB_DIR",
		Value:       "/web",
		Description: "📂 directory with flutter web app",
	},
	{
		Cmd:         rootCmd,
		Name:        "api-adress",
		ShortName:   "a",
		Env:         "PACKREPO_API_ADRESS",
		Value:       "http://localhost:80/",
		Description: "📫 adress for backend api calls via grpc-web",
	},
	{
		Cmd:         rootCmd,
		Name:        "logins",
		ShortName:   "l",
		Env:         "PACKREPO_LOGINS",
		Value:       "user|password",
		Description: "🔐 list of logins and passwords separated by '|' symbol",
	},
}

// Short description of contents for command.
type Flag struct {
	// Cobra command that we will bound our cmd to
	Cmd *cobra.Command
	// Name of command in CLI
	Name string
	// Optional short name for command, leave empty to skip short name
	ShortName string
	// Environment variable to read from
	Env string
	// Regular name for the flag
	Value string
	// Wether this value should be provided by user
	IsRequired bool
	// Leave empty if type is string: ["", "strarr", "bool"]
	Type string
	// Description for flag
	Description string
}

// Function to add new command to CLI tool.
func AddFlag(cmd Flag) {
	if cmd.Type == "" {
		cmd.Cmd.PersistentFlags().StringP(cmd.Name, cmd.ShortName, cmd.Value, cmd.Description)
		err := viper.BindPFlag(cmd.Name, cmd.Cmd.PersistentFlags().Lookup(cmd.Name))
		CheckErr(err)
	}

	if cmd.Type == "strarr" {
		cmd.Cmd.PersistentFlags().StringArrayP(cmd.Name, cmd.ShortName, nil, cmd.Description)
		err := viper.BindPFlag(cmd.Name, cmd.Cmd.PersistentFlags().Lookup(cmd.Name))
		CheckErr(err)
	}

	if cmd.Type == "bool" {
		cmd.Cmd.PersistentFlags().BoolP(cmd.Name, cmd.ShortName, false, cmd.Description)
		err := viper.BindPFlag(cmd.Name, cmd.Cmd.PersistentFlags().Lookup(cmd.Name))
		CheckErr(err)
	}

	if cmd.Type == "int" {
		if cmd.Value != "" {
			i, err := strconv.Atoi(cmd.Value)
			if err != nil {
				err = fmt.Errorf("value for flag "+cmd.Name+" should be int: %w", err)
				CheckErr(err)
			}
			cmd.Cmd.PersistentFlags().IntP(cmd.Name, cmd.ShortName, i, cmd.Description)
			err = viper.BindPFlag(cmd.Name, cmd.Cmd.PersistentFlags().Lookup(cmd.Name))
			CheckErr(err)
			return
		}
		cmd.Cmd.PersistentFlags().IntP(cmd.Name, cmd.ShortName, 0, cmd.Description)
		err := viper.BindPFlag(cmd.Name, cmd.Cmd.PersistentFlags().Lookup(cmd.Name))
		CheckErr(err)
	}

	if cmd.Env != `` {
		err := viper.BindEnv(cmd.Name, cmd.Env)
		CheckErr(err)
	}

	if cmd.IsRequired {
		err := cmd.Cmd.MarkFlagRequired(cmd.Name)
		CheckErr(err)
	}
}

func CheckErr(err error) {
	if err != nil {
		fmt.Println("Error occured: ", fmt.Sprintf("%+v", err))
		os.Exit(1)
	}
}
