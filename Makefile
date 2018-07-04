all: analysis test

test:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...

fmt:
	gofmt -s -w .
	
bench:
	- go test ./... -bench .
	
	
analysis:
	go get golang.org/x/lint/golint
	go get honnef.co/go/tools/cmd/megacheck
	megacheck ./...
	golint ./...