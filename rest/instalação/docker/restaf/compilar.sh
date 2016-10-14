#!/usr/bin/env bash

erro_sair() {
  echo "$1. Abortando" 1>&2
  exit 1
}

compilar() {
  local diretorio_atual=`pwd`
  local diretorio_projeto=`echo $GOPATH | cut -d: -f1`
  diretorio_projeto=$diretorio_projeto/src/github.com/rafaeljusto/atiradorfrequente/rest/rest.af

  cd $diretorio_projeto || erro_sair "Cannot change directory"
  go build -ldflags "-X github.com/rafaeljusto/atiradorfrequente/rest/config.Version=$VERSAO" || erro_sair "Erro de compilação"
  mv rest.af $diretorio_atual || erro_sair "Erro ao copiar binário do servidor REST"
}

VERSAO=$1

uso() {
  echo "Usage: $1 <versão>"
}

if [ -z "$VERSAO" ]; then
  echo "VERSAO não definida!"
  uso $0
  exit 1
fi

compilar