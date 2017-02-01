GO        = go
GOX       = gox
BUILD_DIR = $(CURDIR)/build
GOX_ARGS  = -output="$(BUILD_DIR)/{{.Dir}}_{{.OS}}_{{.Arch}}" -osarch="linux/amd64 freebsd/amd64"


build:
	GOBIN=$(BUILD_DIR) $(GO) install -v

test:
	$(GO) test -v ./...

release-build:
	$(GO) get -u github.com/mitchellh/gox
	$(GOX) $(GOX_ARGS) github.com/thomersch/maillog_exporter
