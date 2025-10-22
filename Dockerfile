FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN adduser -D -u 1001 -g 1001 gocontentful

COPY gocontentful /usr/bin/

USER gocontentful
WORKDIR /home/gocontentful

ENTRYPOINT ["gocontentful"]
