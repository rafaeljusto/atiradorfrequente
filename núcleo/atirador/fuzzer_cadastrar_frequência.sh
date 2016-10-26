#!/bin/sh
set -e
go-fuzz-build -o "fuzzer_cadastrar_frequência.zip" -func "FuzzCadastrarFrequência" github.com/rafaeljusto/atiradorfrequente/núcleo/atirador
go-fuzz -bin=./fuzzer_cadastrar_frequência.zip -dup -workdir=fuzzer
