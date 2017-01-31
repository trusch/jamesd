#!/bin/bash

docker run -d --name jamesd-db mongo

docker run \
  -e DATABASE=mongodb://db/jamesd \
  --link jamesd-db:db \
  -p 8080:80 \
  trusch/jamesd \
    jamesd serve

docker kill jamesd-db
docker rm jamesd-db

exit 0
