# build
FROM golang:1.22.6-alpine3.20 as builder

RUN apk update && apk add git

ENV APP_HOME /go/app
ARG ACCESS_TOKEN
ENV ACCESS_TOKEN=$ACCESS_TOKEN
ARG USERNAME
ENV USERNAME=$USERNAME

RUN git config --global url."https://${USERNAME}:${ACCESS_TOKEN}@gitlab.ai-vlab.com".insteadOf "https://gitlab.ai-vlab.com"

WORKDIR "$APP_HOME"

COPY go.* ./
RUN go mod download

COPY . ./

RUN go build -o server ./cmd

# run
FROM surnet/alpine-wkhtmltopdf:3.18.0-0.12.6-small
RUN apk add --no-cache tzdata

ENV APP_HOME /go/app

WORKDIR "$APP_HOME"

COPY config.yml config.yml
COPY i18n i18n
COPY assets assets
COPY --from=builder "$APP_HOME"/server $APP_HOME

# Run the web service on container startup.
ENTRYPOINT ["./server"]
