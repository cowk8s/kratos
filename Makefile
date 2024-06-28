SHELL=/usr/bin/env bash -o pipefail

#  EXECUTABLES = docker-compose docker node npm go
#  K := $(foreach exec,$(EXECUTABLES),\
#          $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH")))

export GO111MODULE        := on
export PATH               := .bin:${PATH}
export PWD                := $(shell pwd)
export BUILD_DATE         := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
export VCS_REF            := $(shell git rev-parse HEAD)
export QUICKSTART_OPTIONS ?= ""
export IMAGE_TAG 					:= $(if $(IMAGE_TAG),$(IMAGE_TAG),latest)

GO_DEPENDENCIES = \
				  github.com/golang/mock/mockgen \
				  github.com/go-swagger/go-swagger/cmd/swagger \
				  golang.org/x/tools/cmd/goimports \
				  github.com/mattn/goveralls \
				  github.com/cortesi/modd/cmd/modd \
				  github.com/mailhog/MailHog

define make-go-dependency
  # go install is responsible for not re-building when the code hasn't changed
  .bin/$(notdir $1): go.mod go.sum
		GOBIN=$(PWD)/.bin/ go install $1
endef
$(foreach dep, $(GO_DEPENDENCIES), $(eval $(call make-go-dependency, $(dep))))
$(call make-lint-dependency)

# Generates the SDK
.PHONY: sdk
sdk: .bin/swagger
	swagger generate spec -m -o spec/swagger.json \
	-c github.com/cowk8s/kratos

	rm -rf internal/httpclient
	mkdir -p internal/httpclient/

node_modules: package-lock.json
	npm ci
	touch node_modules	