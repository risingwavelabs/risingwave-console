FROM ubuntu:24.04
WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates

ARG TARGETARCH
COPY ./bin/wavekit-server-${TARGETARCH} /app/wavekit-server

CMD ["./wavekit-server"]
