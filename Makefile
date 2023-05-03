pwd := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

docker:
	docker build -t fmnx.io/core/repo .
	docker run -p 80:80 fmnx.io/core/repo

.PHONY: gen
gen:
	buf format -w
	buf generate
	protoc --dart_out=grpc:lib/generated -Iproto proto/v1/pacman.proto

check:
	gofumpt -l -w .
	golangci-lint run
	buf lint

run:
	docker compose down
	chromium --disalbe-web-security &
	docker compose up --build app

evans:
	evans --proto proto/v1/pacman.proto --web --host localhost -p 8080
