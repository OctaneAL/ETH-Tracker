FROM golang:1.20-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/github.com/OctaneAL/ETH-Tracker
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/ETH-Tracker /go/src/github.com/OctaneAL/ETH-Tracker


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/ETH-Tracker /usr/local/bin/ETH-Tracker
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["ETH-Tracker"]
