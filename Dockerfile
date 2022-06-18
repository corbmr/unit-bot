FROM golang:1.18 AS build
WORKDIR /bot
COPY . .
RUN go get
RUN CGO_ENABLED=0 go build -o unit-bot ./bin/bot.go

FROM alpine
WORKDIR /bot
COPY --from=build /bot/unit-bot .
COPY channels.txt .
ENTRYPOINT ["/bot/unit-bot"]
