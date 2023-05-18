pwd := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

install:
	dart pub global activate protoc_plugin
	pack i aur.archlinux.org/protoc-gen-go aur.archlinux.org/protoc-gen-go-grpc aur.archlinux.org/buf-bin

.PHONY: gen
gen:
	protoc --proto_path=. --go_out=. --go-grpc_out=require_unimplemented_servers=false:. pack.proto
	protoc --dart_out=grpc:lib/generated -Iproto pack.proto

docker:
	docker build -t fmnx.su/core/repo .
	docker compose up -d

check:
	gofumpt -l -w .
	golangci-lint run
	buf lint

evans:
	evans --proto pack.proto --web --host localhost -p 8080
