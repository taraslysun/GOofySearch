FROM golang

WORKDIR /app

COPY main.go .
COPY go.mod .
COPY go.sum .


