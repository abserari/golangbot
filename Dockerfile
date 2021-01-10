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
RUN CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o main 


# Final image.
FROM scratch
LABEL maintainer "The Techcats Authors <abserari@google.com>"
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/github.com/abserari/golangbot .
EXPOSE 8080
ENTRYPOINT ["/main"]