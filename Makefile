build:
	go build .

test:
	go test ./...

docker:
	docker build -t formomosan/go_stop:latest . -f Dockerfile
