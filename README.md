# Project RedCedar

Docker-based webapp for encoding video files to the HEVC \(H.265\) standard.

## Docker Compose

Recommended `docker-compose.yaml` \(Assumes the docker image was built with the tag `zpaw/redcedar:latest`\)

```yaml
version: "2.2"
services:
  redcedar:
    image: zpaw/redcedar:latest
    container_name: RedCedar
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

### Controller

- REDCEDAR_DEBUG (bool)

### Runner

- REDCEDAR_DEBUG (bool)

- REDCEDAR_RUNNER_CONTROLLER_IP (string)

- REDCEDAR_RUNNER_CONTROLLER_PORT (integer)
