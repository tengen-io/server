build:
	go build -o tengen .

gen:
	go run scripts/gqlgen.go generate

test:
	go test ./...

docker:
	docker build -t gcr.io/tengen-io/server:latest . -f Dockerfile
