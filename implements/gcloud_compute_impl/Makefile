GO_PACKAGES=.
GO_MODULE_NAME=github.com/mcdonaldseanp/lookout/implements/gcloud_compute_impl
GO_BIN_NAME=gcloud_compute_impl

# Make the build dir, and remove any go bins already there
setup:
	mkdir -p ../../output
	rm -rf ../../output/$(GO_BIN_NAME)
	rm -rf ../../output/$(GO_BIN_NAME).yaml

# Actually build the thing, including copying the impl
# data definition to the output dir
build: setup
	go mod tidy
	go build -o ../../output/ $(GO_MODULE_NAME)
	cp $(GO_BIN_NAME).yaml ../../output/

install:
	go mod tidy
	go install $(GO_MODULE_NAME)

format:
	go fmt $(GO_PACKAGES)