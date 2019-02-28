FROM golang:1.11 as build

WORKDIR /go/src/server

RUN go get -tags 'postgres' -u github.com/golang-migrate/migrate/cmd/migrate

ENV GO111MODULE on

COPY go.mod go.sum ./
RUN go get -d -v ./...

COPY . .
RUN go build -o /go/bin/tengen -v .

FROM ubuntu:bionic

COPY --from=build /go/bin/tengen /usr/local/bin/tengen
COPY .env .

EXPOSE 8080

CMD ["/usr/local/bin/tengen"]
