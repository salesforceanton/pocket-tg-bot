# pocket-tg-bot

## Simple telegram bot to collect your usefull links into Pocket throught Telegram

### Description:

This is example to create telegram bot to collect your usefull links into Pocket throught Telegram.
Bot application contains of 2 parts - telegram bot handler - which listens to messages throught the Telegram Api and
answer after process (saving into your Pocket account or handle any error), and auth server - which listen to users confirmation grant access and processing authorization.

### Stack:

```
1. GO
2. BoltDB
3. Telegram Api golang SDK
4. Pocket API golang SDK (github.com/zhashkevych/go-pocket-sdk)
5. Docker, Makefile
```

### Setting up:

At first we need to create App into Pocket service to take consumer key.

Set on your `.env` file variables according with `.env.example` file
In this example we can run with Docker Container: 

To build image use

```make build-image```

### Running
To start bot application use

```make start-container```