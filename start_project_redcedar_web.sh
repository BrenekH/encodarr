#!/bin/bash

cd ./redcedar-worker
chmod +x start_redcedar_worker.sh && ./start_redcedar_worker.sh

cd ../webapp
chmod +x start_redcedar_webapp.sh && ./start_redcedar_webapp.sh

cd ..
