package atirador

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"testing"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/config"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
	"golang.org/x/image/font/basicfont"
)

func BenchmarkGerarImagemNúmeroControle(b *testing.B) {
	númeroControle := protocolo.NovoNúmeroControle(1, 123)

	var configuração config.Configuração
	configuração.Atirador.ImagemNúmeroControle.Largura = 0
	configuração.Atirador.ImagemNúmeroControle.Altura = 0
	configuração.Atirador.ImagemNúmeroControle.CorFundo.Color = color.RGBA{0xff, 0xff, 0xff, 0xff}
	configuração.Atirador.ImagemNúmeroControle.Fonte.Face = basicfont.Face7x13
	configuração.Atirador.ImagemNúmeroControle.Fonte.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}

	imagemLogo, err := base64.StdEncoding.DecodeString(imagemLogoPNG)
	if err != nil {
		b.Fatalf("Erro ao extrair a imagem de teste do logo. Detalhes: %s", err)
	}
	imagemLogoBuffer := bytes.NewBuffer(imagemLogo)

	configuração.Atirador.ImagemNúmeroControle.Logo.Imagem.Image, _, err = image.Decode(imagemLogoBuffer)
	if err != nil {
		b.Fatalf("Erro ao extrair a imagem de teste do logo. Detalhes: %s", err)
	}

	configuração.Atirador.ImagemNúmeroControle.Logo.Espaçamento = 100
	configuração.Atirador.ImagemNúmeroControle.Borda.Largura = 50
	configuração.Atirador.ImagemNúmeroControle.Borda.Espaçamento = 50
	configuração.Atirador.ImagemNúmeroControle.Borda.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
	configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Largura = 50
	configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Espaçamento = 50
	configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Cor.Color = color.RGBA{0xee, 0xee, 0xee, 0xff}

	for i := 0; i < b.N; i++ {
		gerarImagemNúmeroControle(númeroControle, configuração)
	}
}
