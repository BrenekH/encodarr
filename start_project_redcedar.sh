#!/bin/bash

./stop_project_redcedar.sh

echo "Building and running RedCedar (this can take a few minutes)"
docker run -d \
--name RedCedar-script-control \
--restart always \
--cpus="1.00" \
--mount type=bind,source=/usr/app/tosearch,target=/usr/app/tosearch \
-p 8123:5000 \
$(docker build -q -t zpaw/redcedar:latest-script .)
echo "RedCedar is now running in headless mode"
