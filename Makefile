pwd := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

.PHONY: gen
gen:
	buf format -w
	buf generate
	protoc --dart_out=grpc:lib/generated -Iproto proto/v1/pack.proto

docker:
	docker build -t fmnx.su/core/repo .
	docker compose up -d

check:
	gofumpt -l -w .
	golangci-lint run
	buf lint

evans:
	evans --proto proto/v1/pack.proto --web --host localhost -p 8080
