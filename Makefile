
test:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...
	
bench:
	- go test ./... -bench .