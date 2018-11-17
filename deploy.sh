#!/bin/sh

source env.sh

heroku container:login

heroku config:set OXFORD_API_ID=$OXFORD_API_ID OXFORD_API_KEY=$OXFORD_API_KEY LINE_BOT_SECRET=$LINE_BOT_SECRET LINE_BOT_TOKEN=$LINE_BOT_TOKEN --app=$HEROKU_APP

heroku container:push web --app=$HEROKU_APP
heroku container:release web --app=$HEROKU_APP