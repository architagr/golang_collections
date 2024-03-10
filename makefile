.PHONY: test test-coverage

test:
	go test ./...

test-coverage:
	go test ./... -coverprofile coverage.txt -covermode=atomic 
	go tool cover -html=coverage.txt -o coverage.html