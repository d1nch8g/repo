package cmd

import (
	"errors"
	"fmt"
	"strings"

	"fmnx.io/dev/repo/cmd/service"
	"fmnx.io/dev/repo/cmd/utils"
	"github.com/sirupsen/logrus"
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

	setLogFormat()

	logrus.Info("Initial API adress: ", apiAdress)
	logrus.Info("Initial packages: ", initPkgs)

	helper := &utils.OsHelper{}

	err := helper.ReplaceFileString(
		webPath+`/main.dart.js`,
		`http://localhost:8080/`,
		apiAdress,
	)
	checkErr(err)

	formattedLogins, err := formatLogins(logins)
	checkErr(err)

	go func() {
		err = helper.Execute("pack get " + initPkgs)
		checkErr(err)

		err = helper.FormDb(repoName)
		checkErr(err)
	}()

	err = service.Run(&service.Params{
		Port:     port,
		WebPath:  webPath,
		RepoName: repoName,
		OsHelper: helper,
		Logins:   formattedLogins,
	})
	checkErr(err)
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
