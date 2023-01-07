FROM golang:1.19-alpine as build-base

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY app app