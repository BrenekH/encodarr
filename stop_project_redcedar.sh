#!/bin/bash

echo "Killing and remove docker container RedCedar-script-control"
docker kill RedCedar-script-control && docker rm RedCedar-script-control
