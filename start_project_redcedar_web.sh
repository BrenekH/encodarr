#!/bin/bash

echo "Building and running RedCedar Web in headless mode"
sudo docker run -d \
--name RedCedarWeb \
--restart always \
--cpus="1.25" \
--mount source=/media/plex/mnt/PlexMedia,target=/usr/app/tosearch \
-p 8123:5000 \
$(sudo docker build -q -t redcedar/web:latest_script .)
echo "RedCedar Web is running in headless mode"
