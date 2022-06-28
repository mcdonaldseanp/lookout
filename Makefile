GO_PACKAGES=. ./connection ./local ./localexec ./localfile ./operation ./operparse ./remote ./render ./rgerror ./sanitize ./validator ./version
GO_MODULE_NAME=github.com/mcdonaldseanp/lookout
GO_BIN_NAME=lookout

# Make the build dir, and remove any go bins already there
setup:
	mkdir -p output/
	rm -rf output/$(GO_BIN_NAME)

# Actually build the thing
build-lookout: setup
	go mod tidy
	go build -o output/ $(GO_MODULE_NAME)

build-implements:
	cd implements && \
	for DIR in $$(ls); do \
		cd $$DIR && \
		make build && \
		cd ..; \
	done && \
	cd .. && \
	git checkout -- implements/**/go.mod

build: build-lookout build-implements

install:
	go mod tidy
	go install $(GO_MODULE_NAME)

# Build it before publishing to make sure this publication won't be broken
#
# This also ensures that the clibuild command is available for the version
# command
#
# If NEW_VERSION is set by the user, it will set the new clibuild version
# to that value. Otherwise clibuild will bump the Z version
publish: install format
	NEW_VERSION=$$(lookout update version ./version/version.go "$(NEW_VERSION)") && \
	echo "Tagging and publishing new version $$NEW_VERSION" && \
	git add --all && \
	git commit -m "(release) Update to new version $$NEW_VERSION" && \
	git tag -a $$NEW_VERSION -m "Version $$NEW_VERSION"
	git push
	git push --tags

format:
	go fmt $(lookout_GO_PACKAGES)