# Copyright Krzesimir Nowak
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

TOOLS_BIN_DIR := ./tools/bin
TOOLS_SRC_DIR := ./tools

.PHONY: all
all: build lint

.PHONY: build
build:
	go build .
	go test -run xxxxxMatchNothingxxxxx ./... >/dev/null

.PHONY: check
check:
	go test -v .

.PHONY: lint
lint: $(TOOLS_BIN_DIR)/golangci-lint
	$(TOOLS_BIN_DIR)/golangci-lint run --fix
	go mod tidy

.PHONY: lint-tools
lint-tools: $(TOOLS_BIN_DIR)/golangci-lint
	cd '$(TOOLS_SRC_DIR)' && \
	go mod tidy

$(TOOLS_BIN_DIR)/golangci-lint: $(TOOLS_SRC_DIR)/go.mod $(TOOLS_SRC_DIR)/go.sum $(TOOLS_SRC_DIR)/tools.go
	cd '$(TOOLS_SRC_DIR)' && \
	go build -o '$(abspath $(TOOLS_BIN_DIR))/golangci-lint' github.com/golangci/golangci-lint/cmd/golangci-lint
