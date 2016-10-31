#!/usr/bin/env bash
set -e

executar() {
  if [ "$CENARIO" = "cadastrarfrequencia" ]; then
    go-fuzz-build -o "fuzzer_cadastrarfrequência.zip" -func "FuzzCadastrarFrequência" github.com/rafaeljusto/atiradorfrequente/núcleo/atirador
    go-fuzz -bin=./fuzzer_cadastrarfrequência.zip -dup -workdir=fuzzer/cadastrarfrequência

  elif [ "$CENARIO" = "confirmarfrequencia" ]; then
    go-fuzz-build -o "fuzzer_confirmarfrequência.zip" -func "FuzzConfirmarFrequência" github.com/rafaeljusto/atiradorfrequente/núcleo/atirador
    go-fuzz -bin=./fuzzer_confirmarfrequência.zip -dup -workdir=fuzzer/confirmarfrequência

  else
    echo "Cenário não reconhecido!"
  fi
}

CENARIO=$1

uso() {
  echo "Usage: $1 <cenário>"
}

if [ -z "$CENARIO" ]; then
  echo "CENARIO não definido!"
  uso $0
  exit 1
fi

executar