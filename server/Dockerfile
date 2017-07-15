# preliminary stage container for building

FROM golang:1.8 as builder

WORKDIR /go/src/

RUN go-wrapper download -u github.com/labstack/echo/...

RUN mkdir -p /go/src/github.com/william-olson/cmd-server

COPY . /go/src/github.com/william-olson/cmd-server/

RUN go-wrapper install github.com/william-olson/cmd-server



# prod runner container

FROM ubuntu:14.04

WORKDIR /root/

COPY --from=builder /go/bin/cmd-server ./cmd-server

EXPOSE 7447

CMD ["./cmd-server"]

