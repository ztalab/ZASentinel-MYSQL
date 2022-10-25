# Copyright 2022-present The Ztalab Authors.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#    http://www.apache.org/licenses/LICENSE-2.0
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

PROG=bin/za-mysql

SRCS=cmd/proxy/main.go

CFLAGS = -ldflags "-s -w "

INSTALL_PREFIX=/usr/local

build:
	if [ ! -d "./bin/" ]; then \
		mkdir bin; \
	fi
	go build $(CFLAGS) -o $(PROG) $(SRCS)

install:
	cp $(PROG) $(INSTALL_PREFIX)/bin

race:
	if [ ! -d "./bin/" ]; then \
    	mkdir bin; \
    fi
	go build $(CFLAGS) -race -o $(PROG) $(SRCS)

clean:
	rm -rf ./bin

run:
	go run --race cmd/main.go -c config/config.yaml

################################################################################
# Target: format                                                              #
################################################################################
.PHONY: format
format: modtidy
	gofumpt -l -w . && goimports -local github.com/ztalab/ -w $(shell find ./ -type f -name '*.go' -not -path "./test")

.PHONY: modtidy
modtidy:
	go mod tidy