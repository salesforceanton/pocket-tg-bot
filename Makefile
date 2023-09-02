.PHONY:

build-image:
	docker build -t salesforceanton/pocket_tg_bot:0.1 .

start-container:
	docker run --name pocket_tg_bot --env-file .env -p 80:80 salesforceanton/pocket_tg_bot:0.1