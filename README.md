# choo-dict-bot

## Live Testing
- The webhook URL for Line Messaging API will be https://choo-dict-bot.herokuapp.com/callback
- Bot ID is [@lzp7933a](http://line.me/ti/p/~@lzp7933a)

![alt text](https://qr-official.line.me/M/fdAHRY11Na.png "QR Code")

## Prerequisites for Development
- Mac or Linux which can run shell script
- Docker
- Heroku CLI (for deployment only)

## Local Running and Expose to the internet
- Config API Id, Key, Secret, Token in env.sh
- $ ./run.sh
- $ ./expose.sh (in the new Terminal)
- The webhook URL for Line Messaging API will be https://choo-dict-bot.serveo.net/callback

## Unit Testing
- Config API Id, Key, Secret, Token in env.sh
- $ ./test.sh

## Deployment
- Config Heroku App name in deploy.sh
- $ ./deploy.sh

## Tech Stack
- Go
- Oxford Dictionaries API
- Line Messaging API
- Docker
- Heroku