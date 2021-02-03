# Build stage
FROM golang:1.15.7-alpine AS builder

WORKDIR /go/src/project-redcedar/controller

COPY . .

RUN GOOS=linux go build -o redcedar .


# Run stage
FROM ubuntu:20.04

ENV TZ=America/Boise

RUN apt-get update -qq && DEBIAN_FRONTEND="noninteractive" apt-get install -qq -y tzdata mediainfo

WORKDIR /usr/src/app

COPY --from=builder /go/src/project-redcedar/controller .

ADD webfiles/ ./webfiles/

CMD ["./redcedar"]
