golangci_lint_version=1.21.0

all: analysis test

test:
	go test -race -timeout 10s -coverprofile=coverage.txt -covermode=atomic ./...

test/%:
	go test -race -timeout 10s -coverprofile=coverage.txt -covermode=atomic -run="${*}" ./...

fmt:
	gofmt -s -w .

bench:
	go test ./... -bench .

lint:
	@docker run --rm -v `pwd`:/app -w /app golangci/golangci-lint:v$(golangci_lint_version)  golangci-lint run ./...
