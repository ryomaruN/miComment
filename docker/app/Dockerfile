FROM golang:1.18-alpine

RUN apk add --update && apk add --no-cache git && apk add build-base && apk add ffmpeg

RUN mkdir /go/src/app
WORKDIR /go/src/app
ADD ./ /go/src/app/
ENV GOOGLE_APPLICATION_CREDENTIALS=/go/src/app/resources/gcp_token.json