FROM alpine:latest

WORKDIR /
COPY  build/service /

CMD ["./service"]