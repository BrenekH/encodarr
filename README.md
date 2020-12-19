# Project RedCedar Web

Docker-based webapp for the handbrake-based [Project RedCedar](http://github.com/zPaw/project-redcedar)

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
