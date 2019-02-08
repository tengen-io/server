FROM golang:1.11 as build

WORKDIR /go/src/go_stop
COPY . .

RUN go get -tags 'postgres' -u github.com/golang-migrate/migrate/cmd/migrate

ENV GO111MODULE on

RUN go get -d -v ./...
RUN go install -v ./...

FROM ubuntu:bionic

COPY --from=build /go/bin/go_stop /usr/local/bin/go_stop

EXPOSE 8000

CMD ["/usr/local/bin/go_stop"]
