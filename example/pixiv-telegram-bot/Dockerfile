# This Dockerfile builds an image for a client_golang example.
#
# Use as (from the root for the client_golang repository):
#    docker build -f examples/$name/Dockerfile -t prometheus/golang-example-$name .

# Builder image, where we build the example.
FROM golang:1.15.6 AS builder
ENV GOPROXY="https://goproxy.io"
WORKDIR /go/src/github.com/abserari/golangbot
COPY . .
# WORKDIR /go/src/github.com/prometheus/client_golang/examples/simple
RUN CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o ./build/pixiv-telegram-bot/main ./example/pixiv-telegram-bot/...

# Final image.
FROM scratch
LABEL maintainer "The Techcats Authors <abserari@google.com>"
COPY --from=builder /go/src/github.com/abserari/golangbot/build/pixiv-telegram-bot .
EXPOSE 9082
ENTRYPOINT ["/main"]