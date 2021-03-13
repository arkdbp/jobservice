
lint:
	golint ./...

test:
	go test ./...

test-race:
	go test -race ./...

test-race-log:
	go test -race -v ./...

clean-cache:
	go clean -testcache