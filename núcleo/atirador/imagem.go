package atirador

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/draw"
	"image/png"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/config"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/erros"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
	"github.com/registrobr/gostk/errors"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// gerarImagemNúmeroControle gera uma imagem exibindo uma borda, com um logo no
// canto superior esquerdo e um número de controle centralizado.
func gerarImagemNúmeroControle(númeroControle protocolo.NúmeroControle, configuração config.Configuração) (string, error) {
	var (
		imagemLargura = configuração.Atirador.ImagemNúmeroControle.Largura
		imagemAltura  = configuração.Atirador.ImagemNúmeroControle.Altura
		imagemCor     = configuração.Atirador.ImagemNúmeroControle.CorFundo

		fonte    = configuração.Atirador.ImagemNúmeroControle.Fonte.Face
		fonteCor = configuração.Atirador.ImagemNúmeroControle.Fonte.Cor

		bordaLargura     = configuração.Atirador.ImagemNúmeroControle.Borda.Largura
		bordaEspaçamento = configuração.Atirador.ImagemNúmeroControle.Borda.Espaçamento
		bordaCor         = configuração.Atirador.ImagemNúmeroControle.Borda.Cor

		linhaFundoLargura     = configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Largura
		linhaFundoEspaçamento = configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Espaçamento
		linhaFundoCor         = configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Cor

		logo            = configuração.Atirador.ImagemNúmeroControle.Logo.Imagem
		logoEspaçamento = configuração.Atirador.ImagemNúmeroControle.Logo.Espaçamento
	)

	// não é possível gerar a imagem sem uma fonte definida
	if fonte == nil {
		return "", errors.Errorf("fonte da imagem do número de controle indefinida")
	}

	// define o fundo
	imagem := image.NewRGBA(image.Rect(0, 0, imagemLargura, imagemAltura))
	draw.Draw(imagem, imagem.Bounds(), image.NewUniform(imagemCor), image.ZP, draw.Src)

	// desenha a borda
	if bordaLargura > 0 {
		draw.Draw(imagem, image.Rect(
			bordaEspaçamento,
			bordaEspaçamento,
			imagemLargura-bordaEspaçamento,
			imagemAltura-bordaEspaçamento,
		), image.NewUniform(bordaCor), image.ZP, draw.Src)
		draw.Draw(imagem, image.Rect(
			bordaEspaçamento+bordaLargura,
			bordaEspaçamento+bordaLargura,
			imagemLargura-bordaEspaçamento-bordaLargura,
			imagemAltura-bordaEspaçamento-bordaLargura,
		), image.NewUniform(imagemCor), image.ZP, draw.Src)
	}

	// desenha os riscos no fundo
	if linhaFundoLargura > 0 {
		for y := bordaEspaçamento + bordaLargura; y < imagemAltura-bordaEspaçamento-bordaLargura; y += linhaFundoEspaçamento {
			draw.Draw(imagem, image.Rect(
				bordaEspaçamento+bordaLargura,
				y,
				imagemLargura-bordaEspaçamento-bordaLargura,
				y+linhaFundoLargura,
			), image.NewUniform(linhaFundoCor), image.ZP, draw.Src)
		}
	}

	// desenha o logo no canto superior esquerdo
	if logo.Image != nil {
		draw.Draw(imagem, logo.Bounds().Add(image.Pt(
			bordaLargura+bordaEspaçamento+logoEspaçamento,
			bordaLargura+bordaEspaçamento+logoEspaçamento,
		)), logo, image.ZP, draw.Src)
	}

	desenhador := font.Drawer{
		Dst:  imagem,
		Src:  image.NewUniform(fonteCor),
		Face: fonte,
	}

	desenhador.Dot = fixed.Point26_6{
		X: (fixed.I(imagemLargura) - desenhador.MeasureString(númeroControle.String())) / 2,
		Y: fixed.I(imagemAltura / 2),
	}

	// desenha o número de controle centralizado
	desenhador.DrawString(númeroControle.String())

	var buffer bytes.Buffer
	if err := png.Encode(&buffer, imagem); err != nil {
		return "", erros.Novo(err)
	}
	return base64.StdEncoding.EncodeToString(buffer.Bytes()), nil
}
