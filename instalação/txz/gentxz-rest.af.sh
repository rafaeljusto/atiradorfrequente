#!/usr/bin/env bash
set -e

# informações do pacote
readonly PACOTE_NOME="rest.af"
readonly COMENTARIO="REST - Atirador Frequente"
readonly MANTENEDOR="Rafael Dantas Justo <rafael@justo.net.br>"
readonly URL="https://rafael.net.br"
readonly DESCRICAO="Sistema de controle de frequência para Atiradores com CR em Clubes de Tiro proposto ao Exército Brasileiro"

# informações de instalação
readonly DIRETORIO_TEMPORARIO="/tmp/rest.af/"
readonly DIRETORIO_SCRIPTS="$DIRETORIO_TEMPORARIO/usr/local/etc/rc.d"
readonly DIRETORIO_INSTALACAO="$DIRETORIO_TEMPORARIO/usr/local/atiradorfrequente/rest/"

erro_sair() {
  echo "$1. Abortando" 1>&2
  exit 1
}

preparar() {
  mkdir -p $DIRETORIO_INSTALACAO
  mkdir -p $DIRETORIO_SCRIPTS
}

copiar_arquivos() {
  local diretorio_projeto=`echo $GOPATH | cut -d: -f1`
  diretorio_projeto=$diretorio_projeto/src/github.com/rafaeljusto/atiradorfrequente

  local versao=`echo "$VERSAO" | awk -F "-" '{ print $1 }'`
  local lancamento=`echo "$VERSAO" | awk -F "-" '{ print $2 }'`

  cp $diretorio_projeto/instalação/txz/rest.af.sh $DIRETORIO_SCRIPTS/rest.af.sh || erro_sair "Não foi possível copiar o script de execução ao reiniciar o servidor"

  # calcula tamanho total dos arquivos
  local tamanho_arquivos=0
  for f in `find $DIRETORIO_TEMPORARIO -type f`
  do
    size=`stat --printf="%s" $f`
    tamanho_arquivos=`expr $tamanho_arquivos + $size`
  done

  # calcula hash SHA256 dos arquivos
  local arquivos_manifesto=""
  for file in `find $DIRETORIO_TEMPORARIO -type f`
  do
    hash=`sha256sum $file | awk '{ print $1 }'`
    base=`echo $file | cut -c 12-`
    arquivos_manifesto="$arquivos_manifesto \"$base\":\"1\$${hash}\","
  done

  cat > $DIRETORIO_TEMPORARIO/+COMPACT_MANIFEST <<EOF
{
"name":"$PACOTE_NOME",
"origin":"atiradorfrequente/rest.af",
"version":"$versao,$lancamento",
"comment":"$COMENTARIO",
"maintainer":"$MANTENEDOR",
"www":"$URL",
"abi":"FreeBSD:10:amd64",
"arch":"freebsd:10:x86:64",
"prefix":"/usr/local/atiradorfrequente",
"flatsize":$tamanho_arquivos,
"licenselogic":"single",
"desc":"$DESCRICAO",
"categories":["atiradorfrequente"]
}
EOF

  cat > $DIRETORIO_TEMPORARIO/+MANIFEST <<EOF
{
"name":"$PACOTE_NOME",
"origin":"atiradorfrequente/rest.af",
"version":"$versao,$lancamento",
"comment":"$COMENTARIO",
"maintainer":"$MANTENEDOR",
"www":"$URL",
"abi":"FreeBSD:10:amd64",
"arch":"freebsd:10:x86:64",
"prefix":"/usr/local/atiradorfrequente",
"flatsize":$tamanho_arquivos,
"licenselogic":"single",
"desc":"$DESCRICAO",
"categories":["atiradorfrequente"],
"files":{
$arquivos_manifesto
},
"scripts":{"post-install":"chmod +x /usr/local/etc/rc.d/rest.af.sh;"}
}
EOF
}

compilar() {
  local diretorio_projeto=`echo $GOPATH | cut -d: -f1`
  diretorio_projeto=$diretorio_projeto/src/github.com/rafaeljusto/atiradorfrequente

  cd $diretorio_projeto/rest/rest.af || erro_sair "Cannot change directory"
  env GOOS=freebsd GOARCH=amd64 go build -ldflags "-X github.com/rafaeljusto/atiradorfrequente/rest/config.Version=$version-$release" || erro_sair "Erro de compilação"

  mv $diretorio_projeto/rest/rest.af/rest.af $DIRETORIO_INSTALACAO/ || erro_sair "Erro ao copiar o binário principal"
  cd - 1>/dev/null
}

construir_txz() {
  local versao=`echo "$VERSAO" | awk -F "-" '{ print $1 }'`
  local lancamento=`echo "$VERSAO" | awk -F "-" '{ print $2 }'`
  local diretorio_atual=`pwd`
  local arquivo=rest.af-${versao}-${lancamento}.txz

  cd $DIRETORIO_TEMPORARIO

  find . -type f | cut -c 3- | sort | xargs tar -cJf "$diretorio_atual/$arquivo" --transform 's,^usr,/usr,' --owner=root --group=wheel
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

preparar
compilar
copiar_arquivos
construir_txz
limpar