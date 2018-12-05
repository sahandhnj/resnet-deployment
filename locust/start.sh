#!/bin/bash

docker run --name standalone --hostname standalone -e ATTACKED_HOST=http://localhost:3200 -p 8089:8089 -d mylocust