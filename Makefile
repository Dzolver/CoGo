.PHONY: run-gs
run-gs:
	go run cmd/gameserver/main.go

.PHONY: build
build: test
	go build -o bin/gameserver cmd/gameserver/main.go 
	go build -o bin/validation cmd/validation/main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint:
	go fmt ./...
	golangci-lint run ./...

.PHONY: build-docker
build-docker:
	docker build -t validation .

