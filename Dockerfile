# Dockerfile References: https://docs.docker.com/engine/reference/builder/

FROM golang:1.15.2-alpine3.12

LABEL maintainer="Jonathan Throne <adramalech707@gmail.com>"

WORKDIR /go/src/snippetbox

COPY . .

RUN go mod download

ENV MYSQL_DATABASE_NAME="snippetbox" \
    MYSQL_DATABASE_HOST="127.0.0.1" \
    MYSQL_DATABASE_PORT="3306" \
    MYSQL_USERNAME="web" \
    MYSQL_PASSWORD="password12345!" \
    APP_PORT="80"

RUN apk update && apk add curl \
                          git \
                          bash \
                          make \
                          openssh-client && \
     rm -rf /var/cache/apk/*

RUN make build

EXPOSE 80

CMD [ "make", "prod" ]
