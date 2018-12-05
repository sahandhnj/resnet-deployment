#!/bin/bash 

docker logs -f $(docker run --rm -d --name server2  -v /home/sahand/Projects/Go/src/github.com/sahandhnj/resnet-deployment/meta:/runtime/meta -p 3003:3003 resnet:latest)