FROM golang:1.11 as build

WORKDIR /go/src/server
COPY . .

RUN go get -tags 'postgres' -u github.com/golang-migrate/migrate/cmd/migrate

ENV GO111MODULE on

RUN go get -d -v ./...
RUN go install -v ./...

FROM ubuntu:bionic

COPY --from=build /go/bin/server /usr/local/bin/server

EXPOSE 8000

CMD ["/usr/local/bin/server"]
