#!/bin/bash

echo "Killing and remove docker container RedCedarWeb-script-control"
sudo docker kill RedCedarWeb-script-control && sudo docker rm RedCedarWeb-script-control
