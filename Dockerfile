# Dockerfile References: https://docs.docker.com/engine/reference/builder/

FROM golang:1.15

LABEL maintainer="Jonathan Throne <adramalech707@gmail.com>"

WORKDIR /go/src/snippetbox

COPY . .

RUN go mod download

RUN make build

EXPOSE 80

CMD ["make prod"]
