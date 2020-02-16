
export GO111MODULE=on
export GOPROXY=https://proxy.golang.org

SHELL := /bin/bash -o pipefail
VERSION_PACKAGE = github.com/proactionhq/proaction/pkg/version
VERSION ?=`git describe --tags --dirty`
DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"`

GIT_TREE = $(shell git rev-parse --is-inside-work-tree 2>/dev/null)
ifneq "$(GIT_TREE)" ""
define GIT_UPDATE_INDEX_CMD
git update-index --assume-unchanged
endef
define GIT_SHA
`git rev-parse HEAD`
endef
else
define GIT_UPDATE_INDEX_CMD
echo "Not a git repo, skipping git update-index"
endef
define GIT_SHA
""
endef
endif

define LDFLAGS
-ldflags "\
	-X ${VERSION_PACKAGE}.version=${VERSION} \
	-X ${VERSION_PACKAGE}.gitSHA=${GIT_SHA} \
	-X ${VERSION_PACKAGE}.buildTime=${DATE} \
"
endef

.PHONY: test
test:
	go test ./pkg/... ./cmd/... ./internal/... -coverprofile cover.out

.PHONY: proaction
proaction: fmt vet
	go build ${LDFLAGS} -o bin/proaction github.com/proactionhq/proaction/cmd/proaction

.PHONY: fmt
fmt:
	go fmt ./pkg/... ./internal/... ./cmd/...

.PHONY: vet
vet:
	go vet ./pkg/... ./internal/... ./cmd/... 
