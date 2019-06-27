# Build stage
FROM golang:1.12.5-alpine3.9 as builder

RUN apk update && apk add --no-cache git

WORKDIR /go/src/github.com/bcmendoza/xds-explorer
ADD . /go/src/github.com/bcmendoza/xds-explorer

# Compile necessary binaries for final image
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build --mod=vendor -a -installsuffix cgo -ldflags="-w -s"

# Run stage
FROM alpine:3.9

WORKDIR /app

# Copy over binary
COPY --from=builder /go/src/github.com/bcmendoza/xds-explorer/xds-explorer /app/

RUN chown -R 0:0 /app
RUN chmod -R g=u /app

USER 1000

CMD ./xds-explorer
