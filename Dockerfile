# syntax=docker/dockerfile:1

FROM golang:1.17.6-buster

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,https://goproxy.io,direct \
    TZ=Asia/Shanghai
    VERSION=1.0.0

#RUN apt-get update && apt-get install -y vim
WORKDIR /tuo-agent

COPY . .

RUN make

#EXPOSE 80

#ENTRYPOINT ["./toc-generator"]
