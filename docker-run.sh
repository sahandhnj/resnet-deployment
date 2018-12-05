#!/bin/bash 

docker run -i --rm --name v3 -v meta:/runtime/meta -p 3002:3002 image_detector:latest