FROM alpine:3.4

RUN mkdir -p /lib64 && ln -s /lib/ld-musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

RUN mkdir -p /usr/local/atiradorfrequente/rest
COPY rest.af /usr/local/atiradorfrequente/rest

EXPOSE 80

ENTRYPOINT ["/usr/local/atiradorfrequente/rest/rest.af"]