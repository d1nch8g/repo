<p align="center">
<img style="align: center; padding-left: 10px; padding-right: 10px; padding-bottom: 10px;" width="238px" height="238px" src="./assets/images/logo.png" />
</p>

<h2 align="center">Pack package repository - repo</h2>

[![Generic badge](https://img.shields.io/badge/LICENSE-GPL-orange.svg)](https://fmnx.io/core/repo/src/branch/main/LICENSE)
[![Generic badge](https://img.shields.io/badge/GITEA-REPO-yellow.svg)](https://fmnx.io/core/repo)
[![Generic badge](https://img.shields.io/badge/GITHUB-REPO-white.svg)](https://github.com/fmnx-io/repo)
[![Generic badge](https://img.shields.io/badge/DOCKER-REGISTRY-blue.svg)](https://fmnx.io/core/-/packages/container/repo/latest)
[![Build Status](https://ci.fmnx.io/api/badges/core/repo/status.svg)](https://ci.fmnx.io/core/repo)

Dockerized pack repository with friendly user interface and public API. Project goal is to quickly set up personal repostitory for arch packages.

![](preview.png)

---

## Configurations

Environment variables/flags:

- 📄 - `PACKREPO_REPO` - `repo` - repository name on the web page
- 😀 - `PACKREPO_USER` - `user` - user name in system, will be used to eject `pack` packages, by default it's 'pack'
- 🌐 - `PACKREPO_PORT` - `port` - publically exposed port, `80` default
- 📫 - `PACKREPO_API_ADRESS` - `api-adress` - adress for backend api calls via `grpc-web`
- 📦 - `PACKREPO_INIT_PKGS` - `init-pkgs` - initial packages to download on start
- 📥 - `PACKREPO_INIT_PKGS_LINKS` - `init-pkgs-links` - initial packages to download using links, separated with space
- 📒 - `PACKREPO_LOGS_FORMAT` - `logs-fmt` - format for logs (can be text/json/pretty)
- 📂 - `PACKREPO_WEB_DIR` - `web-dir` - directory with flutter web app
- 🔐 - `PACKREPO_LOGINS` - `logins` - list of logins and passwords separated by '|' symbol

---

## Deploy

- with `docker`:

```sh
docker run -p 80:80 -e PACKREPO_LOGS_FMT=text fmnx.io/core/repo:latest
```

- with `docker-compose`:

```yml
services:
  pacman:
    image: fmnx.io/core/repo:latest
    command: run
    environment:
      PACKREPO_INIT_PKGS: aur.archlinux.org/yamux fmnx.io/core/ainst
      PACKREPO_API_ADRESS: http://localhost:80/
      PACKREPO_LOGS_FMT: text
      PACKREPO_LOGINS: user1|pass1|user2|pass2
    ports:
      - 80:80
```

## Add to pacman conf

Add those lines to your `/etc/pacman.conf`, to get things to work:

```conf
[localhost]
SigLevel = Optional TrustAll
Server = http://localhost/pkg
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
Добавить файл конфигурации который будет автоматически изменяться с поступающими командами для возможности бэкапа в слуаче проблем
Добавить OAuth через сторонние приложения акки гити
Добавить tree view для просмотра зависимостей пакетов
add del lockfile on startup
add sanitizer on each request
-->
