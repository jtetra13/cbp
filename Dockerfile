# Start the Go app build
FROM golang:latest as builder

ENV GO111MODULE=on

COPY . /jbit
WORKDIR /jbit

# Configure Git to use token
# ARG GITHUB_TOKEN=
# RUN echo $(GITHUB_TOKEN)
# RUN git config --global url."https://$(GITHUB_TOKEN):x-oauth-basic@github.com/".insteadOf "https://github.com"
# Get required modules (assumes packages have been added to ./vendor)

# Add support for HTTPS and time zones
#RUN apk update && \
#   apk upgrade && \
#   apk add ca-certificates && \
#   apk add tzdata

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o jbit

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /jbit .

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN chmod 755 jbit
CMD ["./jbit"]