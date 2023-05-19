pwd := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

.PHONY: install
install:
	dart pub global activate protoc_plugin
	pack i aur.archlinux.org/protoc-gen-go aur.archlinux.org/protoc-gen-go-grpc

.PHONY: gen
gen:
	protoc --proto_path=. --go_out=. --go-grpc_out=require_unimplemented_servers=false:. pack.proto
	protoc --proto_path=. --dart_out=grpc:lib/generated -Iproto pack.proto

evans:
	evans --proto pack.proto --web --host localhost -p 80

run:
	flutter build web
	go build . 
	./repo -w build/web -l 'user|password' run
