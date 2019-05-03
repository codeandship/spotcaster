FROM golang:1.12-alpine AS builder

RUN apk add --no-cache git

ADD . /build

WORKDIR /build

RUN go build -o spotcaster cmd/spotcaster/main.go

FROM alpine:3.9

RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true

COPY --from=builder /build/spotcaster /spotcaster

ENTRYPOINT ["/spotcaster"]