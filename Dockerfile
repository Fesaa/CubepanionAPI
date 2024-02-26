FROM golang:latest as go-stage

WORKDIR /app

COPY . ./

RUN go mod download
RUN go build -o /cubepanion_api

FROM debian:stable-slim

WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates libpq-dev
COPY --from=go-stage /cubepanion_api /app/cubepanion_api

EXPOSE 8000

CMD ["/app/cubepanion_api"]