
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


proto-gen:
	protoc -I/usr/local/include  -I"${GOPATH}"/src   -I./api/  --go_out=plugins=grpc:./api  ./api/job.proto