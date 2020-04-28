FROM alpine:latest

RUN apk add --no-cache ca-certificates curl git && update-ca-certificates

COPY proaction /proaction

WORKDIR /code

ENTRYPOINT [ "/proaction" ]
