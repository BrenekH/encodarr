# Encodarr

Docker-based webapp for encoding video files to the HEVC \(H.265\) standard.

## Docker Compose

Recommended `docker-compose.yaml` \(Assumes the docker image was built with the tag `brenekh/encodarr:latest`\)

```yaml
version: "2.2"
services:
  encodarr:
    image: brenekh/encodarr:latest
    container_name: Encodarr
    volumes:
      - /config:/config:rw
      - /media/folder/to/search:/usr/app/tosearch:rw
    ports:
      - 5000:5000
    restart: unless-stopped
    cpus: 2.00
    stop_signal: SIGINT
```

## Environment Variables

### Common

- ENCODARR_DEBUG (bool)

- ENCODARR_LOG_FILE (string)

### Runner

- ENCODARR_RUNNER_NAME (string)

- ENCODARR_RUNNER_CONTROLLER_IP (string)

- ENCODARR_RUNNER_CONTROLLER_PORT (integer)

## Attributions

`controller/controller/mediainfo.go` was modified from [pascoej/go-mediainfo](https://github.com/pascoej/go-mediainfo/blob/509f5adb9998a8fe497be4eed69c73d75161709e/mediainfo.go)
