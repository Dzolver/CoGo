.PHONY: run-gs
run-gs:
	go run cmd/gameserver/main.go

.PHONY: build-gs
build-gs:
	go build -o bin/gameserver cmd/gameserver/main.go 

.PHONY: test-gs
test-gs:
	go test -v ./...

.PHONY: lint
lint:
	go fmt ./...
	golangci-lint run ./...
