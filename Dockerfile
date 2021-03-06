# preliminary stage container for building

FROM golang:1.8 as builder

WORKDIR /go/src/

RUN go-wrapper download -u github.com/labstack/echo/... && \
    go-wrapper download -u database/sql/driver && \
    go-wrapper download -u upper.io/db.v3 && \
    go-wrapper download -u github.com/lib/pq && \
    go-wrapper download -u gopkg.in/matryer/try.v1 && \
    go-wrapper download -u github.com/go-shadow/moment &&\
    go-wrapper download -u github.com/kpango/glg && \
    go-wrapper download -u github.com/getsentry/raven-go

RUN mkdir -p /go/src/github.com/william-olson/cmd-server

COPY . /go/src/github.com/william-olson/cmd-server/

RUN go-wrapper install github.com/william-olson/cmd-server/cmddeps && \
    go-wrapper install github.com/william-olson/cmd-server/cmddb && \
    go-wrapper install github.com/william-olson/cmd-server/cmdutils && \
    go-wrapper install github.com/william-olson/cmd-server/cmdserver && \
    go-wrapper install github.com/william-olson/cmd-server



# prod runner container

FROM ubuntu:14.04

WORKDIR /root/

RUN apt-get update && apt-get install -y ca-certificates

COPY --from=builder /go/bin/cmd-server ./cmd-server
COPY ./cmdversions/version.json ./version.json

EXPOSE 7447

CMD ["./cmd-server"]

