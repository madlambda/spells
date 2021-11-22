all: lint test bench

test:
	go test -race -timeout 10s -coverprofile=coverage.txt -covermode=atomic ./...

test/%:
	go test -race -timeout 10s -coverprofile=coverage.txt -covermode=atomic -run="${*}" ./...

fmt:
	gofmt -s -w .

bench:
	go test -bench=. -benchmem ./...

bench/memory/%:
	@mkdir -p profilling
	go test -bench=. -benchmem -memprofile="profilling/${*}-memory.p" "./${*}"

lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.43.0 run ./...
