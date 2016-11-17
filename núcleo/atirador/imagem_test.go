package atirador

import (
	"bytes"
	"encoding/base64"
	"image"
	"testing"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/config"
	"golang.org/x/image/font/gofont/goregular"
)

func BenchmarkGerarImagemNúmeroControle(b *testing.B) {
	var configuração config.Configuração
	var err error

	configuração.Atirador.ImagemNúmeroControle.Fonte.Font, err = truetype.Parse(goregular.TTF)
	if err != nil {
		b.Fatalf("Erro ao extrair a fonte de teste. Detalhes: %s", err)
	}

	imagemBase, err := base64.StdEncoding.DecodeString(imagemBasePNG)
	if err != nil {
		b.Fatalf("Erro ao extrair a imagem de teste do logo. Detalhes: %s", err)
	}
	imagemBaseBuffer := bytes.NewBuffer(imagemBase)

	configuração.Atirador.ImagemNúmeroControle.ImagemBase.Image, _, err = image.Decode(imagemBaseBuffer)
	if err != nil {
		b.Fatalf("Erro ao extrair a imagem de teste do logo. Detalhes: %s", err)
	}

	f := frequência{
		ID:                1,
		Controle:          123,
		CR:                123456789,
		Calibre:           ".380",
		ArmaUtilizada:     "Arma do Clube",
		QuantidadeMunição: 50,
		DataInício:        time.Now().Add(-40 * time.Minute),
		DataTérmino:       time.Now().Add(-10 * time.Minute),
	}

	for i := 0; i < b.N; i++ {
		gerarImagemNúmeroControle(f, configuração)
	}
}
