package cmd

import (
	"errors"
	"fmt"
	"strings"

	"fmnx.io/dev/repo/cmd/service"
	"fmnx.io/dev/repo/cmd/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "🚀 Run instance of repo",
	Run:   Run,
}

func Run(cmd *cobra.Command, args []string) {
	var (
		port      = viper.GetInt("port")
		repoName  = viper.GetString("repo")
		webPath   = viper.GetString("web-dir")
		initPkgs  = viper.GetString(`init-pkgs`)
		apiAdress = viper.GetString(`api-adress`)
		logins    = viper.GetString(`logins`)
	)

	fmt.Println("Initial API adress: ", apiAdress)
	fmt.Println("Initial packages: ", initPkgs)

	helper := &utils.OsHelper{}

	err := helper.PrepareInitialPackages()
	CheckErr(err)

	err = helper.ReplaceFileString(
		webPath+`/main.dart.js`,
		`http://localhost:80/`,
		apiAdress,
	)
	CheckErr(err)

	formattedLogins, err := formatLogins(logins)
	CheckErr(err)

	go func() {
		err = helper.Execute("pack get " + initPkgs)
		CheckErr(err)

		err = helper.FormDb(repoName)
		CheckErr(err)
	}()

	err = service.Run(&service.Params{
		Port:     port,
		WebPath:  webPath,
		RepoName: repoName,
		OsHelper: helper,
		Logins:   formattedLogins,
	})
	CheckErr(err)
}

func formatLogins(raw string) (map[string]string, error) {
	splitted := strings.Split(raw, "|")
	if len(splitted)%2 != 0 {
		return nil, errors.New("bad logins/passwords: " + fmt.Sprint(splitted))
	}
	formattedLogins := map[string]string{}
	for i, v := range splitted {
		if i%2 == 1 {
			formattedLogins[splitted[i-1]] = v
		}
	}
	return formattedLogins, nil
}
