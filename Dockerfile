# Dockerfile References: https://docs.docker.com/engine/reference/builder/

FROM golang:1.15

LABEL maintainer="Jonathan Throne <adramalech707@gmail.com>"

WORKDIR /go/src/snippetbox

COPY . .

RUN go mod download

RUN make build

EXPOSE 8081

ENV MYSQL_DATABASE_NAME="snippetbox" \
    MYSQL_DATABASE_HOST="127.0.0.1" \
    MYSQL_DATABASE_PORT="8080" \
    MYSQL_USERNAME="web" \
    MYSQL_PASSWORD="password12345!" \
    APP_PORT="8081"

CMD ["make prod"]
