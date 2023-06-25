
BASE_VERSION = 1.0.0
EXTRA_VERSION ?= $(shell git rev-parse  HEAD)
WORKSPACE=$(shell pwd)
PROJECT_NAME=controller
BASE_DOCKER_LABEL=task_tracker
DOCKER_NS=task_tracker
BASE_DOCKER_NS=task_tracker


# Tool commands (overridable)
GO_CMD             ?= go
GO_DEP_CMD         ?= dep
DOCKER_CMD         ?= docker
DOCKER_COMPOSE_CMD ?= docker-compose
IS_RELEASE=true

# defined in version/version.go
METADATA_VAR = Version=$(BASE_VERSION)
METADATA_VAR += CommitSHA=$(EXTRA_VERSION)
#METADATA_VAR += BaseVersion=$(BASEIMAGE_RELEASE)
METADATA_VAR += BaseDockerLabel=$(BASE_DOCKER_LABEL)
METADATA_VAR += DockerNamespace=$(DOCKER_NS)
METADATA_VAR += BaseDockerNamespace=$(BASE_DOCKER_NS)
METADATA_VAR += BuildDate="$(shell date "+%Y年%m月%d日%H:%M:%S" )"

DOCKER_REPO = harbor.mty.wang/mtoss
GO_LDFLAGS = $(patsubst %,-X $(PROJECT_NAME)/version.%,$(METADATA_VAR))
ARCH=$(shell go env GOARCH)
ifneq ($(IS_RELEASE),true)
PROJECT_VERSION=$(BASE_VERSION)-$(EXTRA_VERSION)
APP_TAG ?= v0
else
PROJECT_VERSION=$(BASE_VERSION)
#APP_TAG ?= $(ARCH)-$(BASE_VERSION)
APP_TAG ?= v0
endif
export GO_LDFLAGS

BUILD_DIR ?= .build
SHELL := /bin/bash

#cp  -r ./conf/  .build/

# build chunker
.PHONY: all
all:  format build docker
	@echo "copy file"
	@mkdir -p $(BUILD_DIR)
	cp  -r  conf  .build/
	cp  -r  $(DOCKER_NS) .build/
	cp  -r  docker-compose.yaml .build/
	cp  -r  Dockerfile   .build/

.PHONY: build
build:
	@echo "building ... "

	go build -o $(DOCKER_NS)  -ldflags "$(GO_LDFLAGS) -X $(PROJECT_NAME)/version.ProgramName=task_tracker"  main.go


.PHONY: docker
docker:
	@echo "Building  docker images"
	docker build -t task_tracker:v0 .
	docker save -o task_tracker.imagev0 task_tracker:v0
	#scp -r task_tracker.imagev0  lftadmin@10.80.115.41:/mnt/data/kepler/controller/task_tracker/


format:
	@echo "go fmt"
	@go fmt ./...
	@echo "gofmt finished"

# download go module
.PHONY: test
test:
	go test -v ./...

vet:
	@echo "go vet"
	@go vet ./...
	@echo "ok"

# download go module
.PHONY: mod
mod:
	go mod download


# static code check
.PHONY: check
check:
	@echo $(WORKSPACE)
	docker run --env CODECLIMATE_CODE="$WORKSPACE" -v /tmp/cc:/tmp/cc -v $WORKSPACE:/code -v /var/run/docker.sock:/var/run/docker.sock codeclimate/codeclimate analyze -f text

# clean
.PHONY: clean
clean:
	@rm -rf $(BUILD_DIR)
