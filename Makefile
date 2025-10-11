.PHONY: build run test clean deploy

build:
	go build -o bin/uuid-inspector ./cmd/uuid-inspector

run:
	go run ./cmd/uuid-inspector

test:
	go test -v -cover ./...

clean:
	rm -rf bin/

deploy:
	flyctl deploy -c deploy/fly.toml