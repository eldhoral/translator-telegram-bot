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
RUN go get bitbucket.org/liamstask/goose/cmd/goose

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o "./go-telegram-api" .

EXPOSE 8085

ENTRYPOINT ["./go-telegram-api"]
