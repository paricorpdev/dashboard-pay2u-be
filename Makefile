ifeq (,$(wildcard ./config.yml))
    $(shell cp config.yml.example config.yml)
endif


## for linux .env

NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
GOLANGCI_CMD := $(shell command -v golangci-lint 2> /dev/null)
PKGS := $(shell go list ./... | grep -v /vendor/ | grep -v /internal/mock)
ALL_PACKAGES := $(shell go list ./... | grep -v /vendor/ | grep -v /internal/mock)

run-check-golangci:
ifndef GOLANGCI_CMD
	$(error "Please install golangci linters from https://golangci-lint.run/usage/install/")
endif


## running application
docker-dev-up:
	cp config.yml.example config.yml
	docker-compose -f docker-compose-dev.yml up --build -d

docker-dev-down:
	docker-compose -f docker-compose-dev.yml down

run-lint: check-golangci fmt
	@echo -e "$(OK_COLOR)==> linting projects$(NO_COLOR)..."
	@golangci-lint run --fix

run-packages-fmt:
	@go fmt $(ALL_PACKAGES)

run-test:
	@go test -coverprofile=test_coverage.out $(PKGS)
	go tool cover -html=test_coverage.out -o test_coverage.html
	rm test_coverage.out
	@echo -e "$(OK_COLOR)==> Open test_coverage.html file on your web browser for detailed coverage$(NO_COLOR)..."

run-deps:
	go mod tidy && go mod vendor

run-api-service:
	go run cmd/api/main.go
