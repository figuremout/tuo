# learn from https://github.com/influxdata/telegraf/blob/master/Makefile

# cmd
GO         = go
GOFMT      = gofmt
GOBUILD    = $(GO) build#-race
GOMOD      = $(GO) mod
GOTEST     = $(GO) test
GORUN      = $(GO) run #-race
DOCKER     = docker


# target platform
GOOS       ?= $(shell go env GOOS) # linux/darwin/windows/freebsd
GOARCH     ?= $(shell go env GOARCH) # amd64/arm/386

# current platform
GOHOSTOS   ?= $(shell go env GOHOSTOS)
GOHOSTARCH ?= $(shell go env GOHOSTARCH)

# protobuf
PB_REL     = https://github.com/protocolbuffers/protobuf/releases# protobuf pre-compiled binaries
# PB_DIR     = $(BINDIR)/pb
PROTOC     = $(BINDIR)/protoc
PROTOC-GEN-GO = $(BINDIR)/protoc-gen-go
PROTOC-GEN-GO-GRPC = $(BINDIR)/protoc-gen-go-grpc

PROTOCFLAGS = --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative
PROTOCIN = $(wildcard ./internal/pkg/grpc/*/*.proto)
#./internal/pkg/grpc/hello/hello.pb.go ./internal/pkg/grpc/hello/hello_grpc.pb.go
PROTOCOUT = $(patsubst %.proto, %.pb.go, $(PROTOCIN)) $(patsubst %.proto, %_grpc.pb.go, $(PROTOCIN))

# proj config
MODULE     = github.com/githubzjm/tuo
GOPROXY    = https://goproxy.cn,https://goproxy.io,direct
GOBIN      = $(abspath $(BINDIR))# GOBIN must be absolute path
BINDIR     = ./bin
AGENT      = $(BINDIR)/agent
SERVER     = $(BINDIR)/server
CLIENT     = $(BINDIR)/client
CMDDIR     = ./cmd
# PATH       := $(PATH):$(GOBIN)
BINS     =  $(AGENT) $(SERVER) $(CLIENT)

VERSION    = 0.0.1
DOCKERFILE = ./Dockerfile
LDFLAGS    = -w -s -X main.version=$(VERSION)
BUILDENV   = CGO_ENABLED=0 GO111MODULE=on GOPROXY=$(GOPROXY)


.PHONY: help
help:
	@echo 'Targets:'
	@echo '    all     - download dependencies and compile'
	@echo '    deps    - download dependencies'
	@echo '    compile - compile'
	@echo '    test    - run unit tests'
	@echo '    fmt     - format source files'
	@echo '    tidy    - tidy go modules'
	@echo '    clean   - delete build artifacts'

# .PHONY: all
# all: clean tidy deps build-linux-amd64 build-windows-amd64
.PHONY: init
init:
	$(GOMOD) init $(MODULE)

.PHONY: deps
deps:
	GOPROXY=$(GOPROXY) $(GOMOD) download -x

.PHONY: run
run:
	$(GORUN) .

.PHONY: tidy
tidy:
	GOPROXY=$(GOPROXY) $(GOMOD) tidy -v
	# GOPROXY=$(GOPROXY) $(GOMOD) verify

.PHONY: clean
clean:
	rm -rf $(BINS)

.PHONY: prepare
prepare:
	echo $(PROTOCOUT)


.PHONY: all
all: clean tidy $(BINS)

.PHONY: agent
agent: tidy
	rm -rf $(AGENT)
	$(MAKE) $(AGENT) --no-print-directory

.PHONY: server
server: tidy
	rm -rf $(SERVER)
	$(MAKE) $(SERVER) --no-print-directory

.PHONY: client
client: tidy
	rm -rf $(CLIENT)
	$(MAKE) $(CLIENT) --no-print-directory

$(BINS): $(PROTOCOUT)
	$(BUILDENV) GOOS=$(GOHOSTOS) GOARCH=$(GOHOSTARCH) $(GOBUILD) -ldflags "$(LDFLAGS)" -o $@ $(CMDDIR)/$(@F)

# grpc
$(PROTOCOUT): $(PROTOC) $(PROTOIN)
	PATH=$(GOBIN):$$PATH $(PROTOC) $(PROTOCFLAGS) $(PROTOCIN)

# Protoc Buffer
$(PROTOC): ./protoc-3.15.8-linux-x86_64.zip $(PROTOC-GEN-GO) $(PROTOC-GEN-GO-GRPC)
	unzip -n -j protoc-3.15.8-linux-x86_64.zip "bin/protoc" -d $(BINDIR)
	# protoc's mtime is old, need to be update so that make will work correctly
	touch -m $(PROTOC)
./protoc-3.15.8-linux-x86_64.zip:
	curl -LO $(PB_REL)/download/v3.15.8/protoc-3.15.8-linux-x86_64.zip
# Go plugins for the protocol compiler
$(PROTOC-GEN-GO):
	GOBIN=$(GOBIN) GOPROXY=$(GOPROXY) $(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
$(PROTOC-GEN-GO-GRPC):
	GOBIN=$(GOBIN) GOPROXY=$(GOPROXY) $(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1


# .PHONY: build
# build:
# 	$(BUILDENV) GOOS=$(GOHOSTOS) GOARCH=$(GOHOSTARCH) $(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BINDIR)/$(TARGET)-$(GOHOSTOS)-$(GOHOSTARCH) $(CMDDIR)

# .PHONY: build-linux-amd64
# build-linux-amd64:
# 	$(BUILDENV) GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BINDIR)/$(TARGET)-linux-amd64 $(CMDDIR)

# .PHONY: build-windows-amd64
# build-windows-amd64:
# 	$(BUILDENV) GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BINDIR)/$(TARGET)-windows-amd64.exe $(CMDDIR)

# .PHONY: deploy
# deploy:
# 	tar -zcf ./$(TARGET)-$(VERSION).tar.gz $(BINDIR)

# .PHONY: init
# init:
# 	$(GOMOD) init github.com/githubzjm/tuo



# .PHONY: test
# test:
# 	$(GOTEST) 



# .PHONY: docker-image
# docker-image:
# 	$(DOCKER) build -f $(DOCKERFILE) -t "tuo-agent/$(VERSION)"

# .PHONY: docker-build
# docker-build:
# 	docker exec
# 	docker cp


