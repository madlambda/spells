test:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...

fmt:
	gofmt -s -w .
	
bench:
	- go test ./... -bench .