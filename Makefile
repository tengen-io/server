.PHONY: gen

build:
	go build -o tengen .

gen:
	go run github.com/99designs/gqlgen generate

docker:
	docker build -t gcr.io/tengen-io/server:latest . -f Dockerfile
