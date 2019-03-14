# Start the Go app build
FROM golang:latest as builder

ENV GO111MODULE=on

COPY . /jbit
WORKDIR /jbit

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o jbit

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /jbit .

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN chmod 755 jbit
CMD ["./jbit"]