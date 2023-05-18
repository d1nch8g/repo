// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.su/
// Contact email: help@fmnx.su

package server

import (
	"errors"

	"fmnx.su/core/pack/system"
)

func FormDb(repo string) error {
	repoPath := "/var/cache/pacman/pkg/" + repo + ".db.tar.gz"
	pkgsPath := "/var/cache/pacman/pkg/*.pkg.tar.zst"
	o, err := system.Call("sudo repo-add -n -q " + repoPath + " " + pkgsPath)
	if err != nil {
		return errors.New(o)
	}
	return nil
}
