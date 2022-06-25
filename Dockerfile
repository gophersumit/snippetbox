# syntax=docker/dockerfile:1

FROM golang:1.18-alpine

WORKDIR /

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY . ./

RUN go build -o /snippetbox

EXPOSE 8080

CMD [ "/snippetbox" ]