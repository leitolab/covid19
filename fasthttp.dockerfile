FROM golang:1.14.1-alpine3.11 AS builder

WORKDIR /ieliot

COPY ./src /ieliot
RUN go mod download
RUN go build -o app -gcflags='-l=4' -ldflags="-s -w" ./server


FROM alpine:3.11.2

RUN apk update && apk add --no-cache ca-certificates

RUN addgroup -S ieliot && adduser -S ieliot -G ieliot
RUN chown -R ieliot:ieliot /home/ieliot

COPY ./src/.env /home/ieliot/.env
COPY ./src/keys /home/ieliot/keys
COPY --from=builder /ieliot/app /home/ieliot/app
WORKDIR /home/ieliot

CMD ./app
