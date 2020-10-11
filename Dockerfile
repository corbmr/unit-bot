FROM golang AS build
RUN mkdir /bot
COPY  . /bot
WORKDIR /bot
RUN go mod download
RUN CGO_ENABLED=0 go build -o unit-bot

FROM alpine
WORKDIR /unit-bot
COPY --from=build /bot/unit-bot ./bot
ENTRYPOINT ["/unit-bot/bot"]
