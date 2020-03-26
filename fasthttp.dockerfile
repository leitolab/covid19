FROM golang:1.14.1-alpine3.11 AS builder

WORKDIR /ieliot

COPY ./src /ieliot
RUN go mod download
RUN go build -o app -gcflags='-l=4' -ldflags="-s -w" ./server


FROM alpine:3.11.2
EXPOSE 8080

RUN apk update && apk add --no-cache ca-certificates

WORKDIR /home/ieliot
COPY ./src/.env.prod /home/ieliot/.env
COPY ./src/keys /home/ieliot/keys
COPY --from=builder /ieliot/app /home/ieliot/app

CMD ["./app"]
