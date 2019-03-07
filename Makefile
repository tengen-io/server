.PHONY: test

build:
	go build -o tengen .

gen:
	go run github.com/99designs/gqlgen generate

test:
	go test ./...

docker:
	docker build -t gcr.io/tengen-io/server:latest . -f Dockerfile
