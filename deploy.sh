#!/bin/sh

heroku container:login
heroku container:push web --app=choo-dict-bot
heroku container:release web --app=choo-dict-bot
heroku open --app=choo-dict-bot