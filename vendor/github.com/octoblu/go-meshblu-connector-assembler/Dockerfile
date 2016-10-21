FROM golang:1.6
MAINTAINER Octoblu, Inc. <docker@octoblu.com>

WORKDIR /go/src/github.com/octoblu/meshblu-connector-assembler
COPY . /go/src/github.com/octoblu/meshblu-connector-assembler

RUN env CGO_ENABLED=0 go build -o meshblu-connector-assembler -a -ldflags '-s' .

CMD ["./meshblu-connector-assembler"]
