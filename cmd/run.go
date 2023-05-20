// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.su/
// Contact email: help@fmnx.su

package cmd

import (
	"errors"
	"fmt"
	"strings"

	"fmnx.su/core/pack/system"
	"fmnx.su/core/repo/server"
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

	go Prepare(initPkgs, repoName)

	_, err := system.Callf(
		`sed -i 's|%s|%s|g' %s`,
		`http://localhost:8080/`,
		apiAdress,
		webPath+`/main.dart.js`,
	)
	CheckErr(err)

	formattedLogins, err := formatLogins(logins)
	CheckErr(err)

	err = server.Run(&server.Params{
		Port:     port,
		WebPath:  webPath,
		RepoName: repoName,
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

func Prepare(initPkgs string, repo string) {
	_, err := system.Call("sudo rm -rf /var/lib/pacman/db.lck")
	CheckErr(err)
	_, err = system.Call("sudo rm -rf /tmp/pack.lock")
	CheckErr(err)
	_, err = system.Call("sudo mv -v /var/cache/pacman/initpkg/* /var/cache/pacman/pkg")
	if err != nil {
		fmt.Printf("unable to move initial packages, %+v\n", err)
	}
	_, err = system.Call("pack install " + initPkgs)
	CheckErr(err)
	err = server.FormDb(repo)
	CheckErr(err)
}
