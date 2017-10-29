FROM golang:1.9.2-alpine3.6

RUN apk update
RUN apk add --virtual build-dependencies build-base git
RUN go get github.com/yoheimuta/go-warmcache/example
RUN apk del build-dependencies && rm -rf /var/cache/apk/*
