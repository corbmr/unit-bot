FROM golang:1.18 AS build
WORKDIR /bot
COPY . .
RUN CGO_ENABLED=0 go build -o unit-bot ./bin/bot

FROM alpine
WORKDIR /bot
COPY --from=build /bot/unit-bot .
COPY channels.txt .
ENTRYPOINT ["/bot/unit-bot"]
