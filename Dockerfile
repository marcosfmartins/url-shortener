FROM golang:1.24-alpine AS base

WORKDIR /app

COPY . /app

RUN go mod tidy

FROM base AS build

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -installsuffix cgo -o /go/bin/url-shortener /app/cmd/server

FROM alpine:3.22 AS image

COPY --from=build /go/bin/url-shortener /usr/local/bin/url-shortener

ENTRYPOINT ["url-shortener"]
