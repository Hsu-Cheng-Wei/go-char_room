FROM golang:latest

RUN apt-get update && apt-get install -y vim

ENV GOPATH=/usr
ENV GO111MODULE=on

COPY . /usr/src/app

WORKDIR /usr/src/app

RUN  go mod download
RUN  go build main.go

EXPOSE 8080

ENTRYPOINT ["./main"]