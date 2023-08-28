FROM golang:latest AS builder

RUN go version

COPY . /github.com/salesforceanton/pocket-tg-bot/
WORKDIR /github.com/salesforceanton/pocket-tg-bot/

RUN go mod download
RUN GOOS=linux go build -o ./.bin/bot ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /github.com/salesforceanton/pocket-tg-bot/.bin/bot .
COPY --from=0 /github.com/salesforceanton/pocket-tg-bot/configs configs/

EXPOSE 80

CMD ["./bot"]