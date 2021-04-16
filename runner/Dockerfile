# Build stage
FROM golang:1.16-buster AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /go/src/encodarr/runner

COPY . .

# Disable CGO so that we have a static binary
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o runner .


# Run stage
FROM ubuntu:20.04

ENV TZ=Etc/GMT \
ENCODARR_CONFIG_DIR="/config"

WORKDIR /usr/src/app

RUN chmod 777 /usr/src/app \
 && apt-get update -qq \
 && DEBIAN_FRONTEND="noninteractive" apt-get install -qq -y tzdata ffmpeg

COPY --from=builder /go/src/encodarr/runner/runner ./runner

RUN chmod 777 ./runner

EXPOSE 8123

CMD ["./runner"]
