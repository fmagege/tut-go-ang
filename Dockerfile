FROM golang:alpine AS base
MAINTAINER FH Magege <fmagege@mulax.io>
LABEL vendor=Mulax\ Crypto\ Exchange \
      io.mulax.server.is-production="false" \
      io.mulax.server.version="0.0.1-alpha" \
      io.mulax.server.release-date="2019-03-05"

RUN apk --no-cache add git ca-certificates tzdata \
    && update-ca-certificates

ENV GO111MODULE on
ENV BUILD_PATH /go/src/golang-angular
RUN mkdir -p $BUILD_PATH

RUN adduser -S -D -H -h $BUILD_PATH builder && chown builder $BUILD_PATH
USER builder
WORKDIR $BUILD_PATH

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . $BUILD_PATH

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -a -o server .

FROM scratch

ENV BUILD_PATH /go/src/golang-angular
ENV BIN_PATH /go/bin
ENV AUTH0_API_IDENTIFIER https://golang-angular-api
ENV AUTH0_DOMAIN mulax.auth0.com

COPY --from=base $BUILD_PATH/server $BIN_PATH/server
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo

EXPOSE 5005

CMD ["/go/bin/server"]
