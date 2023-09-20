FROM golang:1.19-alpine as builder

RUN apk update && apk upgrade && \
    apk --no-cache --update add git make gcc \
    libc-dev && \
    mkdir /app

WORKDIR /app

ENV TZ=Asia/Jakarta

ADD . /app
RUN mkdir -p audio \
    && chown -R $(id -u $(whoami)):0 audio \
    && chmod -R g+w audio

RUN go mod download

RUN go install -mod=mod github.com/githubnemo/CompileDaemon
RUN go get bitbucket.org/liamstask/goose/cmd/goose

EXPOSE 8085

ENTRYPOINT ["/app/nobuapiloan"]
