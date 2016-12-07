FROM alpine:3.4

RUN mkdir -p /lib64 && ln -s /lib/ld-musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN mkdir -p /usr/local/atiradorfrequente/rest/imagem

COPY rest.af /usr/local/atiradorfrequente/rest
COPY imagem-fonte-source-sans-pro-regular.ttf /usr/local/atiradorfrequente/rest/imagem
COPY imagem-base.png /usr/local/atiradorfrequente/rest/imagem

ENV AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_FONTE /usr/local/atiradorfrequente/rest/imagem/imagem-fonte-source-sans-pro-regular.ttf
ENV AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_IMAGEM_BASE /usr/local/atiradorfrequente/rest/imagem/imagem-base.png

EXPOSE 80

ENTRYPOINT ["/usr/local/atiradorfrequente/rest/rest.af"]