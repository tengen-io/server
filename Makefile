build:
	go build .

docker:
	docker build -t formomosan/go_stop:latest . -f Dockerfile
