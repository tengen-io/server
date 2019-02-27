build:
	go build .

test:
	go test ./...

docker:
	docker build -t formomosan/tengen-server:latest . -f Dockerfile
