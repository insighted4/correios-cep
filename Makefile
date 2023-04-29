PROJ=correios-cep
ORG_PATH=github.com/insighted4
REPO_PATH=$(ORG_PATH)/$(PROJ)

DOCKER_IMAGE=insighted4/$(PROJ)

$( shell mkdir -p bin )
$( shell mkdir -p release/bin )
$( shell mkdir -p results )

user=$(shell id -u -n)
group=$(shell id -g -n)

export GOBIN=$(PWD)/bin
# Prefer ./bin instead of system packages for things like protoc, where we want
# to use the version the application uses, not whatever a developer has installed.
export PATH=$(GOBIN):$(shell printenv PATH)

# Version
VERSION ?= $(shell ./scripts/git-version.sh)
COMMIT_HASH ?= $(shell git rev-parse HEAD 2>/dev/null)
BUILD_TIME ?= $(shell date +%FT%T%z)

LD_FLAGS="-s -w -X $(REPO_PATH)/pkg/version.BuildTime=$(BUILD_TIME) -X $(REPO_PATH)/pkg/version.CommitHash=$(COMMIT_HASH) -X $(REPO_PATH)/pkg/version.Version=$(VERSION)"

# Inject .env file
-include .env
export $(shell sed 's/=.*//' .env)

.PHONY: build
build: clean
	@echo "Building: $(REPO_PATH)"
	@go install -v $(REPO_PATH)/cmd/admin

.PHONY: build-ver
build-ver: clean ## build with version number
	@echo "Building: $(REPO_PATH) $(COMMIT_HASH)"
	@go install -v -ldflags $(LD_FLAGS) $(REPO_PATH)/cmd/admin

admin:
	$(MAKE) build-ver

.PHONY: clean
clean:
	@echo "Cleaning binary folders"
	@rm -rf bin/*
	@rm -rf release/*
	@rm -rf results/*

.PHONY: release-binary
release-binary:
	@echo "Releasing binary files: ${COMMIT_HASH}"
	@go build -race -o release/bin/admin -v -ldflags $(LD_FLAGS) $(REPO_PATH)/cmd/admin

.PHONY: docker-image
docker-image: clean
	@echo "Building $(DOCKER_IMAGE) image"
	@docker build -t $(DOCKER_IMAGE) --rm -f Dockerfile .

.PHONY: test
test:
	@echo "Testing"
	@go test -v --short ./...

.PHONY: testcoverage
testcoverage:
	@echo "Testing with coverage"
	@mkdir -p results
	@go test -v $(REPO_PATH)/... | go2xunit -output results/tests.xml
	@gocov test $(REPO_PATH)/... | gocov-xml > results/cobertura-coverage.xml

.PHONY: testrace
testrace:
	@echo "Testing with Race Detection"
	@go test -v --race $(REPO_PATH)/...

.PHONY: vet
vet:
	@echo "Running go tool vet on packages"
	go vet $(REPO_PATH)/...

.PHONY: fmt
fmt:
	@echo "Running gofmt on package sources"
	go fmt $(REPO_PATH)/...

.PHONY: testall
testall: testrace vet fmt testcoverage

