# Copyright The Ratify Authors.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

MODULE         = github.com/ratify-project/ratify-cli/v2
COMMANDS       = ratify
GIT_TAG        = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_COMMIT     = $(shell git rev-parse HEAD)

# if the commit was tagged, BuildMetadata is empty.
ifndef BUILD_METADATA
	BUILD_METADATA := unreleased
endif

ifneq ($(GIT_TAG),)
	BUILD_METADATA := 
endif

# set flags
LDFLAGS := -w \
 -X $(MODULE)/internal/version.GitCommit=$(GIT_COMMIT) \
 -X $(MODULE)/internal/version.BuildMetadata=$(BUILD_METADATA)

GO_BUILD_FLAGS = --ldflags="$(LDFLAGS)"

.PHONY: help
help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}'

.PHONY: all
all: build

.PHONY: FORCE
FORCE:

bin/%: cmd/% FORCE
	go build $(GO_INSTRUMENT_FLAGS) $(GO_BUILD_FLAGS) -o $@ ./$<

.PHONY: download
download: ## download dependencies via go mod
	go mod download

.PHONY: build
build: $(addprefix bin/,$(COMMANDS)) ## builds binaries

.PHONY: test
test: vendor check-line-endings ## run unit tests
	go test -race -v -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: clean
clean:
	git status --ignored --short | grep '^!! ' | sed 's/!! //' | xargs rm -rf

.PHONY: check-line-endings
check-line-endings: ## check line endings
	! find cmd internal -name "*.go" -type f -exec file "{}" ";" | grep CRLF

.PHONY: fix-line-endings
fix-line-endings: ## fix line endings
	find cmd internal -type f -name "*.go" -exec sed -i -e "s/\r//g" {} +

.PHONY: vendor
vendor: ## vendores the go modules
	GO111MODULE=on go mod vendor