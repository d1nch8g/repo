<p align="center">
<img style="align: center; padding-left: 10px; padding-right: 10px; padding-bottom: 10px;" width="238px" height="238px" src="./assets/images/logo.png" />
</p>

<h2 align="center">Pacman package repository - fleu-pkg</h2>

[![Generic badge](https://img.shields.io/badge/LICENSE-GPLv3-orange.svg)](https://dancheg97.ru/dancheg97/fleupkg/src/branch/main/LICENSE)
[![Generic badge](https://img.shields.io/badge/GITEA-REPO-red.svg)](https://dancheg97.ru/dancheg97/fleupkg)
[![Generic badge](https://img.shields.io/badge/GITHUB-REPO-white.svg)](https://github.com/fleu-io/fleupkg)
[![Generic badge](https://img.shields.io/badge/DOCKER-REGISTRY-blue.svg)](https://dancheg97.ru/dancheg97/-/packages/container/fleupkg/latest)
[![Build Status](https://drone.dancheg97.ru/api/badges/dancheg97/fleupkg/status.svg)](https://drone.dancheg97.ru/dancheg97/fleupkg)

Dockerized pacman repository with friendly user interface and public API. Project goal is to quickly set up personal pacman repostitory without pain and hustle. Base styling is provided to keep in sync with modern gnome apps.

![](preview.gif)

---

## Configurations

Environment variables/flags:

- 📄 - `FLEUPKG_REPO` - `repo` - repository name on the web page
- 😀 - `FLEUPKG_USER` - `user` - user name in system, will be used to eject `yay` packages
- 🌐 - `FLEUPKG_PORT` - `port` - publically exposed port, `8080` default
- 📫 - `FLEUPKG_API_ADRESS` - `api-adress` - adress for backend api calls via `grpc-web`
- 📦 - `FLEUPKG_INIT_PKGS` - `init-pkgs` - initial packages to download on start
- 📒 - `FLEUPKG_LOGS_FORMAT` - `logs-fmt` - format for logs (can be text/json/pretty)
- 📂 - `FLEUPKG_WEB_DIR` - `web-dir` - directory with flutter web app
- 🔐 - `FLEUPKG_LOGINS` - `logins` - list of logins and passwords separated by '|' symbol

---

## Deploy

- with `docker`:

```sh
docker run -p 8080:8080 -e FLEUPKG_LOGS_FMT=text dancheg97.ru/dancheg97/fleupkg:latest
```

- with `docker-compose`:

```yml
services:
  pacman:
    image: dancheg97.ru/dancheg97/fleupkg:latest
    command: run
    environment:
      FLEUPKG_INIT_PKGS: yay
      FLEUPKG_API_ADRESS: http://localhost:8080/
      FLEUPKG_LOGS_FMT: text
      FLEUPKG_LOGINS: user1|pass1|user2|pass2
    ports:
      - 8080:8080
```

## Add to pacman conf

Add those lines to your `/etc/pacman.conf`, to get things to work:

```conf
[localhost]
SigLevel = Optional TrustAll
Server = http://localhost:8080/pkg
```

You can test it with this commands:

```sh
sudo pacman -R yay
sudo pacman -Sy yay
```

## Contribute

For applicaiton development you need to install following software:

- `go`
- `gofumpt`
- `golangci-lint`
- `buf`
- `flutter`
- `flutter webkit`

All frontend dart code is located in `lib` folder, all backend go code is
located in `cmd` folder.
