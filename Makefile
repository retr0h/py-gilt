VENDOR := vendor
PKGS := $(shell go list ./... | grep -v /$(VENDOR)/)
$(if $(PKGS), , $(error "go list failed"))
PKGS_DELIM := $(shell echo $(PKGS) | sed -e 's/ /,/g')
GITCOMMIT ?= $(shell git rev-parse --short HEAD)
$(if $(GITCOMMIT), , $(error "git rev-parse failed"))
GITUNTRACKEDCHANGES := $(shell git status --porcelain --untracked-files=no)
ifneq ($(GITUNTRACKEDCHANGES),)
GITCOMMIT := $(GITCOMMIT)-dirty
endif
VERSION ?= $(shell git describe --tags --always)
$(if $(VERSION), , $(error "git describe failed"))
BUILDDATE := $(shell date '+%Y/%m/%d %H:%M:%S')
LDFLAGS := -s \
		-w \
		-X 'main.version=$(VERSION)' \
		-X 'main.buildHash=$(GITCOMMIT)' \
		-X 'main.buildDate=$(BUILDDATE)'
BUILDDIR := .build
export GOCACHE = off

test: fmt lint vet bats
	@echo "+ $@"
	@go test -parallel 5 -covermode=count ./...

bats:
	@echo "+ $@"
	@./test/integration/vendor/bats/bin/bats test/integration

cover:
	@echo "+ $@"
	$(shell [ -e coverage.out ] && rm coverage.out)
	@echo "mode: count" > coverage-all.out
	$(foreach pkg,$(PKGS),\
		go test -coverprofile=coverage.out -covermode=count $(pkg);\
		tail -n +2 coverage.out >> coverage-all.out;)
	@go tool cover -html=coverage-all.out -o=coverage-all.html

fmt:
	@echo "+ $@"
	@gofmt -s -l . | grep -v $(VENDOR) | tee /dev/stderr

lint:
	@echo "+ $@"
	@golint ./... | grep -v $(VENDOR) | tee /dev/stderr

vet:
	@echo "+ $@"
	@go vet $(shell go list ./... | grep -v $(VENDOR))

clean:
	@echo "+ $@"
	@rm -rf $(BUILDDIR)

build: clean
	@echo "+ $@"
	@docker run --rm -it \
		-v $(CURDIR):/gopath/src/github.com/retr0h/go-gilt \
		-w /gopath/src/github.com/retr0h/go-gilt \
		tcnksm/gox:1.10.3 \
		gox \
			-ldflags="$(LDFLAGS)" \
			-osarch="linux/amd64 darwin/amd64" \
			-output="$(BUILDDIR)/{{.Dir}}_{{.OS}}_{{.Arch}}"
