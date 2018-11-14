#!/bin/sh

source env.sh
docker run --rm -it "OXFORD_API_ID=$OXFORD_API_ID" -e "OXFORD_API_KEY=$OXFORD_API_KEY" --name choo-dict-bot -p 80:80 $(docker build -q .) "go test -coverprofile=coverage.out && go tool cover -func=coverage.out"