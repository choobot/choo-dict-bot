#!/bin/sh

source env.sh
docker-compose up --build -d
ssh -R choo-dict-bot.serveo.net:80:localhost:80 serveo.net
docker-compose stop