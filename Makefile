project = $(shell basename $(PWD))
server = ./cmd/$(project)
import = github.com/Code-Hex/$(project)
port = 8080
pid = $(PWD)/kamemaru.pid

run:
	@$(GOPATH)/bin/start_server --port=$(port) --pid-file=$(pid) -- ./kamemaru

restart:
	@cat $(pid) | xargs kill -HUP

stop:
	@cat $(pid) | xargs kill -TERM

build-dev: bindata
	@go build -ldflags='-X $(import).DeployMode=develop' $(server)

build-staging: bindata
	@go build -ldflags='-X $(import).DeployMode=staging -X $(import).LogPath=$(PWD)/log' $(server)

bindata:
	@echo 'make bindata.go'
	@go-bindata -pkg $(project) -o bindata.go assets/...