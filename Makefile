build:
	go build .

test:
	go test ./models ./resolvers ./rules

docker:
	docker build -t formomosan/go_stop:latest . -f Dockerfile
