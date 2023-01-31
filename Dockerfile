FROM debian:bookworm-slim AS base
RUN apt update
RUN apt install -y --reinstall ca-certificates
RUN apt install -y tzdata
RUN mkdir -p /app/bin

FROM base AS build_amd64
ADD ./target/kraken-scheduler-linux-amd64 /app/kraken-scheduler

FROM base AS build_arm64
ADD ./target/kraken-scheduler-linux-arm64 /app/kraken-scheduler

FROM base AS build_armv6
ADD ./target/kraken-scheduler-linux-arm /app/kraken-scheduler

FROM base AS build_armv7
ADD ./target/kraken-scheduler-linux-arm /app/kraken-scheduler

FROM build_${TARGETARCH}${TARGETVARIANT}
ENV PATH="$PATH:/app"