FROM debian:bookworm-slim

RUN mkdir -p /app
ADD ./target/kraken-scheduler-linux-amd64 /app
ENV PATH="$PATH:/app"