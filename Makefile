SHELL           := /bin/bash
ALL_SRC         := $(shell find . -name "*.go" | grep -v -e vendor)
GIT_REMOTE_NAME ?= origin
MASTER_BRANCH   ?= master
RELEASE_BRANCH  ?= master

DOCS_BRANCH     ?= 5.4.0-post

include ./semver.mk

REF := $(shell [ -d .git ] && git rev-parse --short HEAD || echo "none")
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
HOSTNAME := $(shell id -u -n)@$(shell hostname)
RESOLVED_PATH=github.com/confluentinc/cli/cmd/confluent

.PHONY: clean
clean:
	rm -rf $(shell pwd)/dist
	rm -f internal/cmd/local/bindata.go
	rm -f mock/local/shell_runner_mock.go

.PHONY: generate
generate: generate-go mocks

.PHONY: generate-go
generate-go:
	@go generate ./...

.PHONY: deps
deps:
	@GO111MODULE=on go get github.com/goreleaser/goreleaser@v0.106.0
	@GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.21.0
	@GO111MODULE=on go get github.com/mitchellh/golicense@v0.1.1
	@GO111MODULE=on go get github.com/golang/mock/mockgen@v1.2.0
	@GO111MODULE=on go get github.com/kevinburke/go-bindata/...@v3.13.0

build: bindata build-go

ifeq ($(shell uname),Darwin)
GORELEASER_SUFFIX ?= -mac.yml
SHASUM ?= gsha256sum
else ifneq (,$(findstring NT,$(shell uname)))
GORELEASER_SUFFIX ?= -windows.yml
# TODO: I highly doubt this works. Completely untested. The output format is likely very different than expected.
SHASUM ?= CertUtil SHA256 -hashfile
else
GORELEASER_SUFFIX ?= -linux.yml
SHASUM ?= sha256sum
endif

show-args:
	@echo "VERSION: $(VERSION)"

#
# START DEVELOPMENT HELPERS
# Usage: make run-ccloud -- version
#        make run-ccloud -- --version
#

# If the first argument is "run-ccloud"...
ifeq (run-ccloud,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "run-ccloud"
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(RUN_ARGS):;@:)
endif

# If the first argument is "run-confluent"...
ifeq (run-confluent,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "run-confluent"
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(RUN_ARGS):;@:)
endif

.PHONY: run-ccloud
run-ccloud:
	 @go run -ldflags '-X main.cliName=ccloud' cmd/confluent/main.go $(RUN_ARGS)

.PHONY: run-confluent
run-confluent:
	 @go run -ldflags '-X main.cliName=confluent' cmd/confluent/main.go $(RUN_ARGS)

#
# END DEVELOPMENT HELPERS
#

.PHONY: build-go
build-go:
	make build-ccloud
	make build-confluent

.PHONY: build-ccloud
build-ccloud:
	@GO111MODULE=on VERSION=$(VERSION) HOSTNAME=$(HOSTNAME) goreleaser release --snapshot --rm-dist -f .goreleaser-ccloud$(GORELEASER_SUFFIX)

.PHONY: build-confluent
build-confluent:
	@GO111MODULE=on VERSION=$(VERSION) HOSTNAME=$(HOSTNAME) goreleaser release --snapshot --rm-dist -f .goreleaser-confluent$(GORELEASER_SUFFIX)

.PHONY: build-integ
build-integ:
	make build-integ-nonrace
	make build-integ-race

.PHONY: build-integ-nonrace
build-integ-nonrace:
	make build-integ-ccloud-nonrace
	make build-integ-confluent-nonrace

.PHONY: build-integ-ccloud-nonrace
build-integ-ccloud-nonrace:
	binary="ccloud_test" ; \
	[ "$${OS}" = "Windows_NT" ] && binexe=$${binary}.exe || binexe=$${binary} ; \
	GO111MODULE=on go test ./cmd/confluent -ldflags="-s -w -X $(RESOLVED_PATH).cliName=ccloud \
	-X $(RESOLVED_PATH).commit=$(REF) -X $(RESOLVED_PATH).host=$(HOSTNAME) -X $(RESOLVED_PATH).date=$(DATE) \
	-X $(RESOLVED_PATH).version=$(VERSION) -X $(RESOLVED_PATH).isTest=true" -tags testrunmain -coverpkg=./... -c -o $${binexe}

.PHONY: build-integ-confluent-nonrace
build-integ-confluent-nonrace:
	binary="confluent_test" ; \
	[ "$${OS}" = "Windows_NT" ] && binexe=$${binary}.exe || binexe=$${binary} ; \
	GO111MODULE=on go test ./cmd/confluent -ldflags="-s -w -X $(RESOLVED_PATH).cliName=confluent \
		    -X $(RESOLVED_PATH).commit=$(REF) -X $(RESOLVED_PATH).host=$(HOSTNAME) -X $(RESOLVED_PATH).date=$(DATE) \
		    -X $(RESOLVED_PATH).version=$(VERSION) -X $(RESOLVED_PATH).isTest=true" -tags testrunmain -coverpkg=./... -c -o $${binexe}

.PHONY: build-integ-race
build-integ-race:
	make build-integ-ccloud-race
	make build-integ-confluent-race

.PHONY: build-integ-ccloud-race
build-integ-ccloud-race:
	binary="ccloud_test_race" ; \
	[ "$${OS}" = "Windows_NT" ] && binexe=$${binary}.exe || binexe=$${binary} ; \
	GO111MODULE=on go test ./cmd/confluent -ldflags="-s -w -X $(RESOLVED_PATH).cliName=ccloud \
	-X $(RESOLVED_PATH).commit=$(REF) -X $(RESOLVED_PATH).host=$(HOSTNAME) -X $(RESOLVED_PATH).date=$(DATE) \
	-X $(RESOLVED_PATH).version=$(VERSION) -X $(RESOLVED_PATH).isTest=true" -tags testrunmain -coverpkg=./... -c -o $${binexe} -race

.PHONY: build-integ-confluent-race
build-integ-confluent-race:
	binary="confluent_test_race" ; \
	[ "$${OS}" = "Windows_NT" ] && binexe=$${binary}.exe || binexe=$${binary} ; \
	GO111MODULE=on go test ./cmd/confluent -ldflags="-s -w -X $(RESOLVED_PATH).cliName=confluent \
		    -X $(RESOLVED_PATH).commit=$(REF) -X $(RESOLVED_PATH).host=$(HOSTNAME) -X $(RESOLVED_PATH).date=$(DATE) \
		    -X $(RESOLVED_PATH).version=$(VERSION) -X $(RESOLVED_PATH).isTest=true" -tags testrunmain -coverpkg=./... -c -o $${binexe} -race

.PHONY: bindata
bindata: internal/cmd/local/bindata.go

internal/cmd/local/bindata.go: cp_cli/* assets/*
	@go-bindata -pkg local -o internal/cmd/local/bindata.go cp_cli/ assets/

.PHONY: authenticate
authenticate:
	# If you setup your laptop following https://github.com/confluentinc/cc-documentation/blob/master/Operations/Laptop%20Setup.md
	# then assuming caas.sh lives here should be fine
	source ~/git/go/src/github.com/confluentinc/cc-dotfiles/caas.sh && caasenv prod

.PHONY: release
release: authenticate get-release-image commit-release tag-release
	@GO111MODULE=on make gorelease
	@GO111MODULE=on VERSION=$(VERSION) make publish
	@GO111MODULE=on VERSION=$(VERSION) make publish-docs
	git checkout go.sum

.PHONY: fakerelease
fakerelease: get-release-image commit-release tag-release
	@GO111MODULE=on make fakegorelease

.PHONY: gorelease
gorelease:
	@GO111MODULE=off go get -u github.com/inconshreveable/mousetrap # dep from cobra -- incompatible with go mod
	@GO111MODULE=on VERSION=$(VERSION) HOSTNAME=$(HOSTNAME) goreleaser release --rm-dist -f .goreleaser-ccloud.yml
	@GO111MODULE=on VERSION=$(VERSION) HOSTNAME=$(HOSTNAME) goreleaser release --rm-dist -f .goreleaser-confluent.yml

.PHONY: fakegorelease
fakegorelease:
	@GO111MODULE=off go get -u github.com/inconshreveable/mousetrap # dep from cobra -- incompatible with go mod
	@GO111MODULE=on VERSION=$(VERSION) HOSTNAME=$(HOSTNAME) goreleaser release --rm-dist -f .goreleaser-ccloud-fake.yml
	@GO111MODULE=on VERSION=$(VERSION) HOSTNAME=$(HOSTNAME) goreleaser release --rm-dist -f .goreleaser-confluent-fake.yml

.PHONY: sign
sign:
	@GO111MODULE=on gon gon_ccloud.hcl
	@GO111MODULE=on gon gon_confluent.hcl
	rm dist/ccloud/darwin_amd64/ccloud_signed.zip || true
	rm dist/confluent/darwin_amd64/confluent_signed.zip || true

.PHONY: download-licenses
download-licenses:
	$(eval token := $(shell (grep github.com ~/.netrc -A 2 | grep password || grep github.com ~/.netrc -A 2 | grep login) | head -1 | awk -F' ' '{ print $$2 }'))
	@# we'd like to use golicense -plain but the exit code is always 0 then so CI won't actually fail on illegal licenses
	@for binary in ccloud confluent; do \
		echo Downloading third-party licenses for $${binary} binary ; \
		GITHUB_TOKEN=$(token) golicense .golicense.hcl ./dist/$${binary}/$(shell go env GOOS)_$(shell go env GOARCH)/$${binary} | GITHUB_TOKEN=$(token) go run cmd/golicense-downloader/main.go -F .golicense-downloader.json -l legal/$${binary}/licenses -n legal/$${binary}/notices ; \
		[ -z "$$(ls -A legal/$${binary}/licenses)" ] && rmdir legal/$${binary}/licenses ; \
		[ -z "$$(ls -A legal/$${binary}/notices)" ] && rmdir legal/$${binary}/notices ; \
	done

.PHONY: dist
dist: download-licenses
	@# unfortunately goreleaser only supports one archive right now (either tar/zip or binaries): https://github.com/goreleaser/goreleaser/issues/705
	@# we had goreleaser upload binaries (they're uncompressed, so goreleaser's parallel uploads will save more time with binaries than archives)
	@for binary in ccloud confluent; do \
		for os in $$(find dist/$${binary} -mindepth 1 -maxdepth 1 -type d | awk -F'/' '{ print $$3 }' | awk -F'_' '{ print $$1 }'); do \
			for arch in $$(find dist/$${binary} -mindepth 1 -maxdepth 1 -iname $${os}_* -type d | awk -F'/' '{ print $$3 }' | awk -F'_' '{ print $$2 }'); do \
				if [ "$${os}" = "darwin" ] && [ "$${arch}" = "386" ] ; then \
					continue ; \
				fi; \
				[ "$${os}" = "windows" ] && binexe=$${binary}.exe || binexe=$${binary} ; \
				rm -rf /tmp/$${binary} && mkdir /tmp/$${binary} ; \
				cp LICENSE /tmp/$${binary} && cp -r legal/$${binary} /tmp/$${binary}/legal ; \
				cp dist/$${binary}/$${os}_$${arch}/$${binexe} /tmp/$${binary} ; \
				suffix="" ; \
				if [ "$${os}" = "windows" ] ; then \
					suffix=zip ; \
					cd /tmp >/dev/null && zip -qr $${binary}.$${suffix} $${binary} && cd - >/dev/null ; \
					mv /tmp/$${binary}.$${suffix} dist/$${binary}/$${binary}_$(VERSION)_$${os}_$${arch}.$${suffix}; \
				else \
					suffix=tar.gz ; \
					tar -czf dist/$${binary}/$${binary}_$(VERSION)_$${os}_$${arch}.$${suffix} -C /tmp $${binary} ; \
				fi ; \
				cp dist/$${binary}/$${binary}_$(VERSION)_$${os}_$${arch}.$${suffix} dist/$${binary}/$${binary}_latest_$${os}_$${arch}.$${suffix} ; \
			done ; \
		done ; \
		cd dist/$${binary}/ ; \
		  $(SHASUM) $${binary}_$(VERSION)_* > $${binary}_$(VERSION)_checksums.txt ; \
		  $(SHASUM) $${binary}_latest_* > $${binary}_latest_checksums.txt ; \
		  cd ../.. ; \
	done

.PHONY: publish
publish: sign dist
	# Note: gorelease target publishes unsigned binaries to the binaries folder in the bucket, we have to overwrite them here after signing
	@for binary in ccloud confluent; do \
		source ~/git/go/src/github.com/confluentinc/cc-dotfiles/caas.sh && caasenv prod && \
		aws s3 cp dist/$${binary}/darwin_amd64/$${binary} s3://confluent.cloud/$${binary}-cli/binaries/$(VERSION:v%=%)/ ; \\
		aws s3 cp dist/$${binary}/ s3://confluent.cloud/$${binary}-cli/archives/$(VERSION:v%=%)/ --recursive --exclude "*" --include "*.tar.gz" --include "*.zip" --include "*_checksums.txt" --exclude "*_latest_*" --acl public-read ; \
		aws s3 cp dist/$${binary}/ s3://confluent.cloud/$${binary}-cli/archives/latest/ --recursive --exclude "*" --include "*.tar.gz" --include "*.zip" --include "*_checksums.txt" --exclude "*_$(VERSION)_*" --acl public-read ; \
	done

.PHONY: publish-installers
## Publish install scripts to S3. You MUST re-run this if/when you update any install script.
publish-installers:
	source ~/git/go/src/github.com/confluentinc/cc-dotfiles/caas.sh && caasenv prod && \
	aws s3 cp install-ccloud.sh s3://confluent.cloud/ccloud-cli/install.sh --acl public-read && \
	aws s3 cp install-confluent.sh s3://confluent.cloud/confluent-cli/install.sh --acl public-read

.PHONY: docs
docs:
#   TODO: we can't enable auto-docs generation for confluent until we migrate go-basher commands into cobra
#	@GO111MODULE=on go run -ldflags '-X main.cliName=confluent' cmd/docs/main.go
	@GO111MODULE=on go run -ldflags '-X main.cliName=ccloud' cmd/docs/main.go

.PHONY: publish-docs
publish-docs: docs
	@TMP_DIR=$$(mktemp -d)/docs || exit 1; \
		git clone git@github.com:confluentinc/docs.git $${TMP_DIR}; \
		cd $${TMP_DIR} || exit 1; \
		git fetch ; \
		git checkout -b cli-$(VERSION) origin/$(DOCS_BRANCH) || exit 1; \
		cd - || exit 1; \
		make publish-docs-internal BASE_DIR=$${TMP_DIR} CLI_NAME=ccloud || exit 1; \
		cd $${TMP_DIR} || exit 1; \
		sed -i '' 's/default "confluent_cli_consumer_[^"]*"/default "confluent_cli_consumer_<uuid>"/' cloud/cli/command-reference/ccloud_kafka_topic_consume.rst || exit 1; \
		git add . || exit 1; \
		git diff --cached --exit-code >/dev/null && echo "nothing to update for docs" && exit 0; \
		git commit -m "chore: updating CLI docs for $(VERSION)" || exit 1; \
		git push origin cli-$(VERSION) || exit 1; \
		hub pull-request -b $(DOCS_BRANCH) -m "chore: updating CLI docs for $(VERSION)" || exit 1; \
		cd - || exit 1; \
		rm -rf $${TMP_DIR}
#   TODO: we can't enable auto-docs generation for confluent until we migrate go-basher commands into cobra
#	    make publish-docs-internal BASE_DIR=$${TMP_DIR} CLI_NAME=confluent || exit 1; \

.PHONY: publish-docs-internal
publish-docs-internal:
ifndef BASE_DIR
	$(error BASE_DIR is not set)
endif
ifeq (ccloud,$(CLI_NAME))
	$(eval DOCS_DIR := cloud/cli/command-reference)
else ifeq (confluent,$(CLI_NAME))
	$(eval DOCS_DIR := cli/command-reference)
else
	$(error CLI_NAME is not set correctly - must be one of "confluent" or "ccloud")
endif
	rm $(BASE_DIR)/$(DOCS_DIR)/*.rst
	cp $(GOPATH)/src/github.com/confluentinc/cli/docs/$(CLI_NAME)/*.rst $(BASE_DIR)/$(DOCS_DIR)

.PHONY: clean-docs
clean-docs:
	rm docs/*/*.rst

.PHONY: fmt
fmt:
	@gofmt -e -s -l -w $(ALL_SRC)

.PHONY: release-ci
release-ci:
ifneq ($(SEMAPHORE_GIT_PR_BRANCH),)
	true
else ifeq ($(SEMAPHORE_GIT_BRANCH),master)
	make release
else
	true
endif

cmd/lint/en_US.aff:
	@curl -s "https://chromium.googlesource.com/chromium/deps/hunspell_dictionaries/+/master/en_US.aff?format=TEXT" | base64 -D > $@

cmd/lint/en_US.dic:
	@curl -s "https://chromium.googlesource.com/chromium/deps/hunspell_dictionaries/+/master/en_US.dic?format=TEXT" | base64 -D > $@

.PHONY: lint-cli
lint-cli: cmd/lint/en_US.aff cmd/lint/en_US.dic
	@GO111MODULE=on go run cmd/lint/main.go -aff-file $(word 1,$^) -dic-file $(word 2,$^) $(ARGS)

.PHONY: lint-go
lint-go:
	@GO111MODULE=on golangci-lint run --timeout=10m

.PHONY: lint
lint: lint-go lint-cli lint-installers

.PHONY: lint-installers
## Lints the CLI installation scripts
lint-installers:
	@diff install-c* | grep -v -E "^---|^[0-9c0-9]|PROJECT_NAME|BINARY" && echo "diff between install scripts" && exit 1 || exit 0

.PHONY: lint-licenses
## Scan and validate third-party dependeny licenses
lint-licenses: build
	$(eval token := $(shell (grep github.com ~/.netrc -A 2 | grep password || grep github.com ~/.netrc -A 2 | grep login) | head -1 | awk -F' ' '{ print $$2 }'))
	@for binary in ccloud confluent; do \
		echo Licenses for $${binary} binary ; \
		[ -t 0 ] && args="" || args="-plain" ; \
		GITHUB_TOKEN=$(token) golicense $${args} .golicense.hcl ./dist/$${binary}/$(shell go env GOOS)_$(shell go env GOARCH)/$${binary} ; \
		echo ; \
	done

.PHONY: coverage
coverage:
      ifdef CI
	@# Run unit tests with coverage.
	@GO111MODULE=on go test -v -race -coverpkg=$$(go list ./... | grep -v test | grep -v mock | tr '\n' ',' | sed 's/,$$//g') \
		-coverprofile=unit_coverage.txt $$(go list ./... | grep -v vendor | grep -v test)
	@# Run integration tests with coverage.
	@GO111MODULE=on INTEG_COVER=on go test -v $$(go list ./... | grep cli/test) $(TEST_ARGS)
	@echo "mode: atomic" > coverage.txt
	@grep -h -v "mode: atomic" unit_coverage.txt >> coverage.txt
	@grep -h -v "mode: atomic" integ_coverage.txt >> coverage.txt
      else
	@# Run unit tests.
	@GO111MODULE=on go test -race -coverpkg=./... $$(go list ./... | grep -v vendor | grep -v test)
	@# Run integration tests.
	@GO111MODULE=on go test -v -race $$(go list ./... | grep cli/test) $(TEST_ARGS)
      endif

.PHONY: mocks
mocks: mock/local/shell_runner_mock.go

mock/local/shell_runner_mock.go:
	mockgen -source internal/cmd/local/shell_runner.go -destination mock/local/shell_runner_mock.go ShellRunner

.PHONY: test-installers
test-installers:
	@echo Running packaging/installer tests
	@bash test-installers.sh

.PHONY: test
test: bindata mocks lint coverage test-installers

.PHONY: doctoc
doctoc:
	npx doctoc README.md

