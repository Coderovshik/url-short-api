BINARY_NAME = url-short

.PHONY: build
build:
	@go build -o ./bin/${BINARY_NAME} ./cmd/url-short/main.go

.PHONY: run
run: build
	@./bin/${BINARY_NAME}

.PHONY: clean
clean:
	@go clean
	@rm -rf ./bin/*

.PHONY: test
test:
	@go test ./...

.PHONY: test_coverage
test_coverage:
	@go test ./... -coverprofile=coverage.out

.PHONY: dep
dep:
	@go mod download

.PHONY: vet
vet:
	@go vet

.PHONY: lint
lint:
	@golangci-lint run --enable-all

.PHONY: compose-build
compose-build:
	docker compose build

.PHONY: compose-up
compose-up:
	docker compose up

.PHONY: compose-up-build
compose-up-build:
	docker compose up --build

.PHONY: compose-down
compose-down:
	docker compose down