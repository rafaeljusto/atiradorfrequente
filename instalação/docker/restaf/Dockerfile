FROM alpine:3.4

RUN mkdir -p /lib64 && ln -s /lib/ld-musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN mkdir -p /usr/local/atiradorfrequente/rest/imagem

COPY rest.af /usr/local/atiradorfrequente/rest
COPY imagem-fonte-source-sans-pro-regular.ttf /usr/local/atiradorfrequente/rest/imagem
COPY imagem-logo.png /usr/local/atiradorfrequente/rest/imagem

ENV AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_FONTE /usr/local/atiradorfrequente/rest/imagem/imagem-fonte-source-sans-pro-regular.ttf 48 300 preto
ENV AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LOGO_IMAGEM /usr/local/atiradorfrequente/rest/imagem/imagem-logo.png
ENV AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LOGO_ESPACAMENTO 100
ENV AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LINHA_FUNDO_COR branco

EXPOSE 80

ENTRYPOINT ["/usr/local/atiradorfrequente/rest/rest.af"]