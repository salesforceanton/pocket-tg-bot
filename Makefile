.PHONY:

build:
	go build -o ./.bin/bot cmd/main.go

run: build
	./.bin/bot

build-image:
	docker build -t salesforceanton/pocket_tg_bot:0.1 .

start-container:
	docker run --env-file .env -p 80:80 salesforceanton/pocket_tg_bot:0.1