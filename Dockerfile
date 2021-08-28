FROM golang:1.17 AS build
WORKDIR /bot
COPY . .
RUN CGO_ENABLED=0 go build -mod vendor -o unit-bot ./bin/bot

FROM alpine
WORKDIR /bot
COPY --from=build /bot/unit-bot .
ENTRYPOINT ["/bot/unit-bot"]
