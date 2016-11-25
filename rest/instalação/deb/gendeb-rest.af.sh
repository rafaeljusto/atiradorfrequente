#!/usr/bin/env bash
set -e

# informações do pacote
readonly PACOTE_NOME="rest.af"
readonly FORNECEDOR="Rafael Dantas Justo"
readonly MANTENEDOR="Rafael Dantas Justo <rafael@justo.net.br>"
readonly URL="https://rafael.net.br"
readonly LICENCA="MIT"
readonly DESCRICAO="Sistema de controle de frequência para Atiradores com CR em Clubes de Tiro proposto ao Exército Brasileiro"

# informações de instalação
readonly DIRETORIO_INSTALACAO="/usr/local/atiradorfrequente/rest/"
readonly DIRETORIO_TEMPORARIO="/tmp/rest.af/$DIRETORIO_INSTALACAO"

erro_sair() {
  echo "$1. Abortando" 1>&2
  exit 1
}

preparar() {
  rm -f rest.af*.deb 2>/dev/null

  mkdir -p $DIRETORIO_TEMPORARIO || erro_sair "Não foi possível criar o diretório temporário"
  mkdir -p /tmp/rest.af/usr/share/atiradorfrequente/ || erro_sair "Não foi possível criar o diretório de scripts"
}

copiar_arquivos() {
  local diretorio_projeto=`echo $GOPATH | cut -d: -f1`
  diretorio_projeto=$diretorio_projeto/src/github.com/rafaeljusto/atiradorfrequente

  cp $diretorio_projeto/instalação/deb/rest.af.postscript /tmp/rest.af/usr/share/atiradorfrequente/ || erro_sair "Não foi possível copiar o script de execução ao reiniciar o servidor"
}

compilar() {
  local diretorio_projeto=`echo $GOPATH | cut -d: -f1`
  diretorio_projeto=$diretorio_projeto/src/github.com/rafaeljusto/atiradorfrequente/rest/rest.af

  cd $diretorio_projeto || erro_sair "Não foi possível trocar de diretório"
  go build -ldflags "-X github.com/rafaeljusto/atiradorfrequente/rest/config.Version=$VERSAO" || erro_sair "Erro de compilação"

  mv rest.af $DIRETORIO_TEMPORARIO || erro_sair "Erro ao copiar o binário principal"
  cd - 1>/dev/null
}

construir_deb() {
  local diretorio_projeto=`echo $GOPATH | cut -d: -f1`
  diretorio_projeto=$diretorio_projeto/src/github.com/rafaeljusto/atiradorfrequente

  local versao=`echo "$VERSAO" | awk -F "-" '{ print $1 }'`
  local lancamento=`echo "$VERSAO" | awk -F "-" '{ print $2 }'`

  fpm -s dir -t deb --after-install $diretorio_projeto/instalação/deb/rest.af.postinst \
    --after-upgrade $diretorio_projeto/instalação/deb/rest.af.postinst \
    --exclude=.git -n $PACOTE_NOME -v "$versao" --iteration "$lancamento" --vendor "$FORNECEDOR" \
    --maintainer "$MANTENEDOR" --url $URL --license "$LICENCA" --description "$DESCRICAO" \
    --deb-upstart $diretorio_projeto/instalação/deb/rest.af.upstart \
    --deb-systemd $diretorio_projeto/instalação/deb/rest.af.service \
    --deb-user root --deb-group root \
    --prefix / -C /tmp/rest.af usr/local/atiradorfrequente/rest usr/share/atiradorfrequente
}

limpar() {
  rm -rf /tmp/rest.af
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

limpar
preparar
compilar
copiar_arquivos
construir_deb
limpar