# Build stage
FROM golang:1.16-buster AS builder

WORKDIR /go/src/project-redcedar/controller

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o redcedar .


# Run stage
FROM ubuntu:20.04

ENV TZ=Etc/GMT \
REDCEDAR_SEARCH_DIR="/tosearch" \
REDCEDAR_CONFIG_DIR="/config"

WORKDIR /usr/src/app

RUN chmod 777 /usr/src/app \
 && apt-get update -qq \
 && DEBIAN_FRONTEND="noninteractive" apt-get install -qq -y tzdata mediainfo

COPY --from=builder /go/src/project-redcedar/controller/redcedar ./redcedar

RUN chmod 777 ./redcedar

EXPOSE 8123

CMD ["./redcedar"]
