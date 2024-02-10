FROM rust:latest as builder

WORKDIR /app

COPY . ./

RUN cargo build --release

FROM debian:stable-slim

WORKDIR /app

COPY --from=builder /app/target/release/cubepanion_api .

RUN apt-get update && apt-get install -y libssl-dev && apt-get install -y libpq-dev

EXPOSE 8000

CMD ["./cubepanion_api"]