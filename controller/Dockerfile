# Build stage
FROM golang:1.16-buster AS builder

WORKDIR /go/src/encodarr/controller

COPY . .

# We disable CGo via an env var, but reenable it using an installsuffix so that a statically linked binary is produced.
RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o encodarr .


# Run stage
FROM ubuntu:20.04

ENV TZ=Etc/GMT \
ENCODARR_SEARCH_DIR="/tosearch" \
ENCODARR_CONFIG_DIR="/config"

WORKDIR /usr/src/app

RUN chmod 777 /usr/src/app \
 && apt-get update -qq \
 && DEBIAN_FRONTEND="noninteractive" apt-get install -qq -y tzdata mediainfo

COPY --from=builder /go/src/encodarr/controller/encodarr ./encodarr

RUN chmod 777 ./encodarr

EXPOSE 8123

CMD ["./encodarr"]
