build:
	go build -o tengen .

test:
	go test ./...

docker:
	docker build -t gcr.io/tengen-io/server:latest . -f Dockerfile
