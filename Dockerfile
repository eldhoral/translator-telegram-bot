FROM golang:1.19-alpine as builder

RUN apk update && apk upgrade && \
    apk --no-cache --update add git make gcc \
    libc-dev && \
    mkdir /app

WORKDIR /app

ENV TZ=Asia/Jakarta

# Private Repo
ARG repo_user
ARG repo_token
ARG repo_url
ARG repo_organization
ARG app_name

ENV repo_user=$repo_user
ENV repo_token=$repo_token
ENV repo_url=$repo_url
ENV repo_organization=$repo_organization
ENV app_name=$app_name

RUN git config \
    --global --add \
    url."https://${repo_user}:${repo_token}@${repo_url}/".insteadOf \
    "https://${repo_url}/"

RUN export GOPRIVATE=${repo_url}/${repo_organization}/*

ADD . /app
RUN mkdir params

RUN go mod download

RUN go install -mod=mod github.com/githubnemo/CompileDaemon
RUN go get bitbucket.org/liamstask/goose/cmd/goose

EXPOSE 8080

ENTRYPOINT CompileDaemon --build="go build -o /go/bin/${app_name} ." --command=/go/bin/${app_name} --directory=.
