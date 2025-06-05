FROM ubuntu:24.04
WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates

ARG TARGETARCH
COPY ./bin/risingwave-console-server-${TARGETARCH} /app/risingwave-console-server

CMD ["./risingwave-console-server"]
