FROM alpine:latest

WORKDIR /
COPY  build/cli /

CMD ["./cli"]