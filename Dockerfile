FROM golang:1.6
MAINTAINER Octoblu, Inc. <docker@octoblu.com>

WORKDIR /go/src/github.com/octoblu/meshblu-connector-installer
COPY . /go/src/github.com/octoblu/meshblu-connector-installer

RUN env CGO_ENABLED=0 go build -o meshblu-connector-installer -a -ldflags '-s' .

CMD ["./meshblu-connector-installer"]
