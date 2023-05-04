pwd := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

.PHONY: gen
gen:
	buf format -w
	buf generate
	protoc --dart_out=grpc:lib/generated -Iproto proto/v1/pack.proto

docker:
	docker build -t fmnx.io/core/repo .
	docker compose up -d
	chromium --disalbe-web-security http://localhost/ &

check:
	gofumpt -l -w .
	golangci-lint run
	buf lint

test:
	docker compose down
	docker run -p 80:80 -e PACKREPO_LOGS_FORMAT=text -e PACKREPO_INIT_PKGS='xdg-user-dirs-gtk aur.archlinux.org/papirus-icon-theme aur.archlinux.org/adw-gtk-theme aur.archlinux.org/gnome-browser-connector aur.archlinux.org/gnome-shell-extension-dash-to-dock aur.archlinux.org/zsh-theme-powerlevel10k-bin-git aur.archlinux.org/zsh-autosuggestions aur.archlinux.org/zsh-syntax-highlighting aur.archlinux.org/onlyoffice-bin aur.archlinux.org/telegram-desktop aur.archlinux.org/visual-studio-code-bin docker docker-compose aur.archlinux.org/flutter aur.archlinux.org/buf-bin aur.archlinux.org/golangci-lint-bin aur.archlinux.org/protoc-gen-go-grpc aur.archlinux.org/bluez-utils aur.archlinux.org/gnome-shell-extension-dash-to-dock aur.archlinux.org/gnome-shell-extension-gtile cmake clang aur.archlinux.org/neovim-git aur.archlinux.org/vlang qemu-desktop edk2-ovmf archiso archinstall' fmnx.io/core/repo

evans:
	evans --proto proto/v1/pack.proto --web --host localhost -p 8080
