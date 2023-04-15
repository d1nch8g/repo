<p align="center">
<img style="align: center; padding-left: 10px; padding-right: 10px; padding-bottom: 10px;" width="238px" height="238px" src="./assets/images/logo.png" />
</p>

<h2 align="center">Pacman package repository - fmnxpkg</h2>

[![Generic badge](https://img.shields.io/badge/LICENSE-GPLv3-orange.svg)](https://fmnx.ru/dancheg97/fmnxpkg/src/branch/main/LICENSE)
[![Generic badge](https://img.shields.io/badge/GITEA-REPO-red.svg)](https://fmnx.ru/dancheg97/fmnxpkg)
[![Generic badge](https://img.shields.io/badge/GITHUB-REPO-white.svg)](https://github.com/fmnx-ru/fmnxpkg)
[![Generic badge](https://img.shields.io/badge/DOCKER-REGISTRY-blue.svg)](https://fmnx.ru/dancheg97/-/packages/container/fmnxpkg/latest)
[![Build Status](https://ci.fmnx.ru/api/badges/dancheg97/fmnxpkg/status.svg)](https://ci.fmnx.ru/dancheg97/fmnxpkg)

Dockerized pacman repository with friendly user interface and public API. Project goal is to quickly set up personal pacman repostitory without pain and hustle. Base styling is provided to keep in sync with modern gnome apps.

![](preview.gif)

---

## Configurations

Environment variables/flags:

- 📄 - `FMNXPKG_REPO` - `repo` - repository name on the web page
- 😀 - `FMNXPKG_USER` - `user` - user name in system, will be used to eject `yay` packages
- 🌐 - `FMNXPKG_PORT` - `port` - publically exposed port, `8080` default
- 📫 - `FMNXPKG_API_ADRESS` - `api-adress` - adress for backend api calls via `grpc-web`
- 📦 - `FMNXPKG_INIT_PKGS` - `init-pkgs` - initial packages to download on start
- 📥 - `FMNXPKG_INIT_PKGS_LINKS` - `init-pkgs-links` - initial packages to download using links, separated with space
- 📒 - `FMNXPKG_LOGS_FORMAT` - `logs-fmt` - format for logs (can be text/json/pretty)
- 📂 - `FMNXPKG_WEB_DIR` - `web-dir` - directory with flutter web app
- 🔐 - `FMNXPKG_LOGINS` - `logins` - list of logins and passwords separated by '|' symbol

---

## Deploy

- with `docker`:

```sh
docker run -p 8080:8080 -e FMNXPKG_LOGS_FMT=text fmnx.ru/dancheg97/fmnxpkg:latest
```

- with `docker-compose`:

```yml
services:
  pacman:
    image: fmnx.ru/dancheg97/fmnxpkg:latest
    command: run
    environment:
      FMNXPKG_INIT_PKGS: yay
      FMNXPKG_API_ADRESS: http://localhost:8080/
      FMNXPKG_LOGS_FMT: text
      FMNXPKG_LOGINS: user1|pass1|user2|pass2
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

<!--
Добавить установку пакетов загруженных из ссылок
Добавить ui для загрузки пакетов через ссылки
Добавить файл конфигурации который будет автоматически изменяться с поступающими командами для возможности бэкапа в слуаче проблем
Добавить OAuth через сторонние приложения акки гити
 -->
