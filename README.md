# choo-dict-bot

## Note
- There is Load testing with 10,000 concurrent requests in Unit Testing but Oxford API limit is 60 requests per minute, So Live Testing will support only 60 requests per minute
- Repo: https://github.com/choobot/choo-dict-bot/

## Live Testing
- The webhook URL for LINE Messaging API will be https://choo-dict-bot.herokuapp.com/callback
- Bot ID is [@lzp7933a](http://line.me/ti/p/~@lzp7933a)

## Prerequisites for Development
- Mac or Linux which can run shell script
- Docker
- Heroku CLI (for deployment only)

## Local Running and Expose to the internet
- Config environment variables in env.sh
- $ ./run.sh
- $ ./expose.sh (in the new Terminal)
- The webhook URL for LINE Messaging API will be https://choo-dict-bot.serveo.net/callback

## Unit Testing
- Config environment variables in env.sh
- $ ./test.sh

## Deployment
- Config environment variables in env.sh
- Config webhook URL for LINE Messaging API
- $ ./deploy.sh

## Tech Stack
- Go
- Oxford Dictionaries API
- LINE Messaging API
- Docker
- Heroku