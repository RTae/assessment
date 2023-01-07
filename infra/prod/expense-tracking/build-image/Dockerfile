FROM golang:1.19-alpine as build-base

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY app app

RUN CGO_ENABLED=0 go test -v -tags unit ./...

RUN go build -o ./out/expense-tracking ./app

FROM alpine:3.16.2 as prod
COPY --from=build-base /app/out/expense-tracking /app/expense-tracking

CMD ["/app/expense-tracking"]