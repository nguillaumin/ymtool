GOCMD=go
GO_TEST=$(GOCMD) test
GO_BUILD=$(GOCMD) build
ZIP=zip

GITHUB_REF_NAME ?= master
SRC = cmd/ymtool/ymtool.go

all: test build_linux build_windows build_macosx zip

test:
	$(GO_TEST) -v ./...

build_linux: $(SRC)
	GOOS=linux GOARCH=amd64 $(GO_BUILD) -v -o ymtool-linux_amd64-$(GITHUB_REF_NAME) $<

build_windows: $(SRC)
	GOOS=windows GOARCH=amd64 $(GO_BUILD) -v -o ymtool-windows_amd64-$(GITHUB_REF_NAME) $<

build_macosx: $(SRC)
	GOOS=darwin GOARCH=amd64 $(GO_BUILD) -v -o ymtool-darwin_amd64-$(GITHUB_REF_NAME) $<

zip:
	$(ZIP) ymtool-linux_amd64.zip ymtool-linux_amd64-*
	$(ZIP) ymtool-windows_amd64.zip ymtool-windows_amd64-*
	$(ZIP) ymtool-darwin_amd64.zip ymtool-darwin_amd64-*

clean:
	rm -f ymtool-*
	 
