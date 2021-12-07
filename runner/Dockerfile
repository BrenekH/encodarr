# Build stage
FROM --platform=$BUILDPLATFORM golang:1.16-buster AS builder

ARG TARGETOS
ARG TARGETARCH
ARG LDFLAGS_VERSION=development

WORKDIR /go/src/encodarr/runner

COPY . .

# Disable CGO so that we have a static binary
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -tags timetzdata -o runner -ldflags="-X 'github.com/BrenekH/encodarr/runner/options.Version=${LDFLAGS_VERSION}'" cmd/EncodarrRunner/main.go

# Acquire FFmpeg stage
FROM --platform=$BUILDPLATFORM alpine as ffmpeg-extracter

ARG TARGETARCH

WORKDIR /usr/bin

RUN apk add --no-cache wget
RUN wget -q https://johnvansickle.com/ffmpeg/builds/ffmpeg-git-${TARGETARCH}-static.tar.xz -O ffmpeg.tar.xz
RUN mkdir /usr/bin/ffmpeg
RUN tar -xf ffmpeg.tar.xz -C /usr/bin/ffmpeg --strip-components 1

# Run stage
FROM scratch

ENV TZ=Etc/GMT \
	PATH="/bin" \
	ENCODARR_CONFIG_DIR="/config"

WORKDIR /usr/src/app

COPY --from=builder /go/src/encodarr/runner/runner ./runner

COPY --from=ffmpeg-extracter /usr/bin/ffmpeg/ffmpeg /bin/ffmpeg

CMD ["./runner"]
