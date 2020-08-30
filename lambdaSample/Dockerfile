FROM golang:1.13.7-buster

WORKDIR /go/src/app

ENV GO111MODULE=on

RUN apt-get update -qq && \
    apt-get install -y git

CMD ["/bin/bash"]