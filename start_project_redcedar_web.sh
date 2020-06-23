#!/bin/bash

echo "Building RedCedar Web"
sudo docker run -d \
--name RedCedarWeb \
--restart always \
--cpus="1.00" \
--mount type=bind,source=/media/plex/mnt/PlexMedia,target=/usr/app/tosearch \
-p 8123:5000 \
$(sudo docker build -q -t redcedar/web:latest_script .)
echo "RedCedar Web is now running in headless mode"
