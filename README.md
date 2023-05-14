<p align="center">
<img style="align: center; padding-left: 10px; padding-right: 10px; padding-bottom: 10px;" width="238px" height="238px" src="./assets/images/logo.png" />
</p>

<h2 align="center">Pack package repository - repo</h2>

[![Generic badge](https://img.shields.io/badge/LICENSE-GPL-orange.svg)](https://fmnx.su/core/repo/src/branch/main/LICENSE)
[![Generic badge](https://img.shields.io/badge/FMNX-REPO-006db0.svg)](https://fmnx.su/core/repo)
[![Generic badge](https://img.shields.io/badge/CODEBERG-REPO-45a3fb.svg)](https://codeberg.org/fmnx/repo)
[![Generic badge](https://img.shields.io/badge/GITHUB-REPO-white.svg)](https://github.com/fmnx-io/repo)
[![Generic badge](https://img.shields.io/badge/DOCKER-REGISTRY-blue.svg)](https://fmnx.su/core/-/packages/container/repo/latest)
[![Build Status](https://ci.fmnx.su/api/badges/core/repo/status.svg)](https://ci.fmnx.su/core/repo)

Dockerized pack repository with friendly user interface and public API. Project goal is to quickly set up personal repostitory for arch packages. Uses [pack](https://fmnx.su/core/pack) for installing and updating packages (you can install AUR packages aswell).

![](preview.png)

---

## Configurations

Packrepo environment variables/flags:

- `PACKREPO_REPO` - repository name on the web page
- `PACKREPO_PORT` - exposed port, `80` default
- `PACKREPO_API_ADRESS` - adress for backend api calls via `grpc-web`
- `PACKREPO_INIT_PKGS` - initial packages to download on start
- `PACKREPO_WEB_DIR` - directory with flutter web app
- `PACKREPO_LOGINS` - list of logins and passwords separated by '|' symbol

Pack environment variables:

- `PACK_ALLOW_AUR` - automatically try downloading packages from AUR if they are not found in pacman repositories
- `PACK_REMOVE_GIT_REPOS` - remove git repositories after installation
- `PACK_REMOVE_BUILT_PACKAGES` - don't cache built `.pkg.tar.zst` packages
- `PACK_DISABLE_PRETTYPRINT` - disable colors in CLI output
- `PACK_DEBUG_MODE` - watch every system call execution

---

## Deploy

- with `docker`:

```sh
docker run -p 80:80 -e PACKREPO_LOGS_FMT=text fmnx.su/core/repo:latest
```

- with `docker-compose`:

```yml
services:
  pacman:
    image: fmnx.su/core/repo:latest
    environment:
      # Initial packages to install with pack
      PACKREPO_INIT_PKGS: aur.archlinux.org/yay fmnx.su/core/ainst
      # Adress for calls by grpcweb
      PACKREPO_API_ADRESS: http://localhost:80/
      PACKREPO_LOGS_FMT: text
      # Admins logins and passwords
      PACKREPO_LOGINS: user1|pass1|user2|pass2
    volumes:
      # Cache packages between sessions
      - ./pkg:/var/cache/pacman/pkg
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

<!--
Add option to install packages from cache
Add ability to attach shell to execute some scripts from UI when logged
-->
