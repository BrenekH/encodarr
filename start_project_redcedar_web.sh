#!/bin/bash

./stop_project_redcedar_web.sh

echo "Building and running RedCedar Web Edition (this can take a few minutes)"
sudo docker run -d \
--name RedCedarWeb-script-control \
--restart always \
--cpus="1.00" \
--mount type=bind,source=/usr/app/tosearch,target=/usr/app/tosearch \
-p 8123:5000 \
$(sudo docker build -q -t redcedar/web:latest_script .)
echo "RedCedar Web Edition is now running in headless mode"
