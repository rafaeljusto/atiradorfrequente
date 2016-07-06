package atirador

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/erros"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// gerarImagemNúmeroControle gera uma imagem exibindo uma borda, com um logo no
// canto superior esquerdo e um número de controle centralizado. O logo e a
// fonte estão fixos no código, para alterar o logo calcula o base64 de um
// arquivo PNG, JPEG ou GIF e atualize o arquivo imagem_logo.go, já para alterar
// a fonte obtenha o programa em [1] e execute a linha abaixo com o endereço e
// tamanho da fonte desejada:
//
//   ./genbasicfont -size=200 -pkg=atirador -hinting=full \
//       -var=imagemFonte -fontfile=http://exemplo.com.br/fonte.ttf
//
// O programa gera o arquivo imagemFonte.go, renomeie para imagem_fonte.go para
// manter o padrão.
//
// [1] http://github.com/golang/freetype/example/genbasicfont
func gerarImagemNúmeroControle(númeroControle protocolo.NúmeroControle) (string, error) {
	const (
		imagemLargura         = 3508
		imagemAltura          = 2480
		bordaLargura          = 50
		bordaEspaçamento      = 50
		linhaFundoLargura     = 2
		linhaFundoEspaçamento = 50
		logoEspaçamento       = 100
	)

	var (
		corFundo      = color.White
		corBorda      = color.Black
		corLinhaFundo = color.Gray16{0xeeee}
	)

	imagem := image.NewRGBA(image.Rect(0, 0, imagemLargura, imagemAltura))

	// define o fundo
	draw.Draw(imagem, imagem.Bounds(), image.NewUniform(corFundo), image.ZP, draw.Src)

	// desenha a borda
	draw.Draw(imagem, image.Rect(
		bordaEspaçamento,
		bordaEspaçamento,
		imagemLargura-bordaEspaçamento,
		imagemAltura-bordaEspaçamento,
	), image.NewUniform(corBorda), image.ZP, draw.Src)
	draw.Draw(imagem, image.Rect(
		bordaEspaçamento+bordaLargura,
		bordaEspaçamento+bordaLargura,
		imagemLargura-bordaEspaçamento-bordaLargura,
		imagemAltura-bordaEspaçamento-bordaLargura,
	), image.NewUniform(corFundo), image.ZP, draw.Src)

	// desenha os riscos no fundo
	for y := bordaEspaçamento + bordaLargura; y < imagemAltura-bordaEspaçamento-bordaLargura; y += linhaFundoEspaçamento {
		draw.Draw(imagem, image.Rect(
			bordaEspaçamento+bordaLargura,
			y,
			imagemLargura-bordaEspaçamento-bordaLargura,
			y+linhaFundoLargura,
		), image.NewUniform(corLinhaFundo), image.ZP, draw.Src)
	}

	arquivoLogo, err := base64.StdEncoding.DecodeString(imagemNúmeroControleLogo)
	if err != nil {
		return "", erros.Novo(err)
	}

	logo, _, err := image.Decode(bytes.NewBuffer(arquivoLogo))
	if err != nil {
		return "", erros.Novo(err)
	}

	// desenha o logo no canto superior esquerdo
	draw.Draw(imagem, logo.Bounds().Add(image.Pt(
		bordaLargura+bordaEspaçamento+logoEspaçamento,
		bordaLargura+bordaEspaçamento+logoEspaçamento,
	)), logo, image.ZP, draw.Src)

	desenhador := font.Drawer{
		Dst:  imagem,
		Src:  image.Black,
		Face: &imagemFonte,
	}

	desenhador.Dot = fixed.Point26_6{
		X: (fixed.I(imagemLargura) - desenhador.MeasureString(númeroControle.String())) / 2,
		Y: fixed.I(imagemAltura / 2),
	}

	// desenha o número de controle centralizado
	desenhador.DrawString(númeroControle.String())

	var buffer bytes.Buffer
	if err = png.Encode(&buffer, imagem); err != nil {
		// TODO(rafaeljusto): Não sei ao certo como podemos simular este cenário de
		// erro. Uma possibilidade seria ignorar o erro já que a imagem base sempre
		// será válida.
		return "", erros.Novo(err)
	}
	return base64.StdEncoding.EncodeToString(buffer.Bytes()), nil
}
