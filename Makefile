# SPDX-License-Identifier: Apache-2.0
# Copyright 2021 Authors of KubeArmor

CURDIR     := $(shell pwd)
INSTALLDIR := $(shell go env GOPATH)/bin/

ifeq (, $(shell which govvv))
$(shell go get github.com/ahmetb/govvv@latest)
endif

PKG      := $(shell go list ./version)
GIT_INFO := $(shell govvv -flags -pkg $(PKG))


ifneq "$(findstring cilium, $(MAKECMDGOALS))" ""
BUILD_TAG := -tags cilium
endif

ifneq "$(findstring kubearmor, $(MAKECMDGOALS))" ""
BUILD_TAG := -tags kubearmor
endif

ifneq "$(findstring discovery, $(MAKECMDGOALS))" ""
BUILD_TAG := -tags discovery
endif

.PHONY: build
build:
	cd $(CURDIR); go mod tidy; CGO_ENABLED=0 go build -ldflags "-w -s ${GIT_INFO}" -tags "cilium kubearmor insight" -o karmor

cilium:
	cd $(CURDIR); go mod tidy; CGO_ENABLED=0 go build -ldflags "-w -s ${GIT_INFO}" ${BUILD_TAG} -o karmor

kubearmor:
	cd $(CURDIR); go mod tidy; CGO_ENABLED=0 go build -ldflags "-w -s ${GIT_INFO}" ${BUILD_TAG} -o karmor

insight:
	cd $(CURDIR); go mod tidy; CGO_ENABLED=0 go build -ldflags "-w -s ${GIT_INFO}" ${BUILD_TAG} -o karmor


.PHONY: install
install: build
	install -m 0755 karmor $(DESTDIR)$(INSTALLDIR)

.PHONY: clean
clean:
	cd $(CURDIR); rm -f karmor

.PHONY: protobuf
vm-protobuf:
	cd $(CURDIR)/vm/protobuf; protoc --proto_path=. --go_opt=paths=source_relative --go_out=plugins=grpc:. vm.proto

insight-protobuf:
	cd $(CURDIR)/insight/protobuf; protoc --proto_path=. --go_opt=paths=source_relative --go_out=plugins=grpc:. insight.proto

.PHONY: gofmt
gofmt:
	cd $(CURDIR); gofmt -s -d $(shell find . -type f -name '*.go' -print)

.PHONY: golint
golint:
ifeq (, $(shell which golint))
	@{ \
	set -e ;\
	GOLINT_TMP_DIR=$$(mktemp -d) ;\
	cd $$GOLINT_TMP_DIR ;\
	go mod init tmp ;\
	go get -u golang.org/x/lint/golint ;\
	rm -rf $$GOLINT_TMP_DIR ;\
	}
endif
	cd $(CURDIR); golint ./...

.PHONY: gosec
gosec:
ifeq (, $(shell which gosec))
	@{ \
	set -e ;\
	GOSEC_TMP_DIR=$$(mktemp -d) ;\
	cd $$GOSEC_TMP_DIR ;\
	go mod init tmp ;\
	go get -u github.com/securego/gosec/v2/cmd/gosec ;\
	rm -rf $$GOSEC_TMP_DIR ;\
	}
endif
	cd $(CURDIR); gosec ./...
