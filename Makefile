project = $(shell basename $(PWD))
server = ./cmd/$(project)
import = github.com/Code-Hex/$(project)
port = 8080

build-dev:
	@go build -ldflags='-X $(import).DeployMode=develop -X $(import).Port=$(port)' $(server)

build-staging:
	@go build -ldflags='-X $(import).DeployMode=staging -X $(import).Port=$(port) -X $(import).LogPath=$(PWD)/log' $(server)