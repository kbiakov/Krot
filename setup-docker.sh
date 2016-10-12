#!/bin/bash

# build API container
docker build -t kbiakov/krot-api .

# create network
docker network create krot

# start NSQ container
docker run -d --net krot -p 4160:4160 -p 4161:4161 nsqio/nsq /nsqlookupd --name lookupd 

# start API container
docker run -d --net krot -p 5000:5000 --name krot-api kbiakov/krot-api
