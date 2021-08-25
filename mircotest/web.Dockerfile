FROM alpine:latest

WORKDIR /
COPY  build/web /

EXPOSE 8000
CMD ["./web"]