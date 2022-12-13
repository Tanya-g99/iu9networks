FROM golang:1.13 as BASE

WORKDIR /code/

COPY . /code/

RUN sh ./build.sh

ENV LOGXI *
ENV LOGXI_FORMAT pretty,happy