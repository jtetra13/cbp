# Start the Go app build
FROM golang:latest

#Copy source
WORKDIR src/github.com/jtetra13/cbp
COPY . .

# Configure Git to use token
# ARG GITHUB_TOKEN=
# RUN echo $(GITHUB_TOKEN)
# RUN git config --global url."https://$(GITHUB_TOKEN):x-oauth-basic@github.com/".insteadOf "https://github.com"
# Get required modules (assumes packages have been added to ./vendor)

# Build a statically-linked Go binary for Linux
# RUN go build

# New build phase -- create binary-only image
# FROM alpine:latest

# Add support for HTTPS and time zones
#RUN apk update && \
#   apk upgrade && \
#   apk add ca-certificates && \
#   apk add tzdata

RUN go get -d -v ./...
RUN pwd && find .


# Start the application
# CMD ["github.com/go-coinbasepro/main"]
CMD ["go", "run", "main.go"]