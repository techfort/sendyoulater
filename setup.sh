#! /bin/bash

CMD=$1
echo "command is ${CMD}"

if [ "$CMD" == "start" ]; then
   echo "setting env vars"
   source .env
   env
   echo "starting docker..."
   docker-compose up -d --remove-orphans
   echo "starting backend"
   cd cmd/api && watcher -watch github.com/techfort/sendyoulater
   cd ../..
   echo "starting frontend"
   ./static/syl-ui/npm run serve
fi

if [ "$CMD" == "stop" ]; then
    echo "stopping docker..."
    docker-compose down -v
fi