build:
	go build -o tengen .

test:
	go test ./...

docker:
	docker build -t formomosan/tengen-server:latest . -f Dockerfile
