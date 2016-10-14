#!/bin/bash

# build API container
docker build -t kbiakov/krot-api .

# create network
docker network create krot

# start Mongo container
docker run -d --net krot -p 27017:27017 mongo --name mymongo

# start NSQ container
docker run -d --net krot -p 4160:4160 -p 4161:4161 nsqio/nsq /nsqlookupd --name lookupd

# start service containers
docker run -d --net krot -p 9010:9010 kbiakov/krot-notification --name krot-notification
docker run -d --net krot -p 9020:9020 kbiakov/krot-subscription --name krot-subscription

# start API container
docker run -d --net krot -p 5000:5000 kbiakov/krot-api --name krot-api
