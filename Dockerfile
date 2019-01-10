FROM golang:1.11

WORKDIR /go/src/go_stop
COPY . .

RUN go get -tags 'postgres' -u github.com/golang-migrate/migrate/cmd/migrate
RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 8000

CMD ["go_stop"]
