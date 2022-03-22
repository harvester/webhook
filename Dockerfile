FROM alpine:3.14
COPY bin/webhook /usr/bin
ENTRYPOINT ["webhook"]
