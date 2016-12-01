package atirador

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"strconv"

	"github.com/golang/freetype"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/config"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/erros"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
	"github.com/registrobr/gostk/errors"
	qrcode "github.com/skip2/go-qrcode"
)

// gerarImagemNúmeroControle gera uma imagem com dados da frequência utilizando
// uma imagem base.
func (f *frequência) gerarImagemNúmeroControle(configuração config.Configuração, códigoVerificação string) error {
	fonte := configuração.Atirador.ImagemNúmeroControle.Fonte.Font
	if fonte == nil {
		// não é possível gerar a imagem sem uma fonte definida
		return errors.Errorf("fonte da imagem do número de controle indefinida")
	}

	// define o fundo
	imagem := image.NewRGBA(configuração.Atirador.ImagemNúmeroControle.ImagemBase.Bounds())
	draw.Draw(imagem, imagem.Bounds(), configuração.Atirador.ImagemNúmeroControle.ImagemBase, image.ZP, draw.Src)

	// textos
	textos := []imagemTexto{
		{
			texto:        strconv.Itoa(f.CR),
			fonteCor:     color.RGBA{0x00, 0x00, 0x00, 0xff},
			fonteTamanho: 11,
			posição:      imagemTextoPosição{260, 110},
		},
		{
			texto:        f.Calibre,
			fonteCor:     color.RGBA{0x00, 0x00, 0x00, 0xff},
			fonteTamanho: 11,
			posição:      imagemTextoPosição{260, 126},
		},
		{
			texto:        f.DataInício.Format("02/01/2006 15:04") + " - " + f.DataTérmino.Format("15:04"),
			fonteCor:     color.RGBA{0x00, 0x00, 0x00, 0xff},
			fonteTamanho: 11,
			posição:      imagemTextoPosição{260, 142},
		},
		{
			texto:        f.ArmaUtilizada,
			fonteCor:     color.RGBA{0x00, 0x00, 0x00, 0xff},
			fonteTamanho: 11,
			posição:      imagemTextoPosição{260, 157},
		},
		{
			texto:        strconv.Itoa(f.QuantidadeMunição),
			fonteCor:     color.RGBA{0x00, 0x00, 0x00, 0xff},
			fonteTamanho: 11,
			posição:      imagemTextoPosição{260, 173},
		},
		{
			texto:        string(protocolo.NovoNúmeroControle(f.ID, f.Controle)),
			fonteCor:     color.RGBA{0xff, 0x00, 0x00, 0xff},
			fonteTamanho: 18,
			posição:      imagemTextoPosição{150, 220},
		},
		{
			texto:        códigoVerificação,
			fonteCor:     color.RGBA{0x00, 0x00, 0x00, 0xff},
			fonteTamanho: 11,
			posição:      imagemTextoPosição{32, 486},
		},
	}

	camadaTexto := freetype.NewContext()
	camadaTexto.SetDPI(150)
	camadaTexto.SetClip(imagem.Bounds())
	camadaTexto.SetDst(imagem)
	camadaTexto.SetFont(fonte)

	for _, t := range textos {
		posição := freetype.Pt(int(t.posição.x), int(camadaTexto.PointToFixed(float64(t.fonteTamanho+t.posição.y))>>6))
		camadaTexto.SetSrc(image.NewUniform(t.fonteCor))
		camadaTexto.SetFontSize(t.fonteTamanho)
		camadaTexto.DrawString(t.texto, posição)
	}

	// QR Code
	qrURL := fmt.Sprintf(configuração.Atirador.ImagemNúmeroControle.URLQRCode, strconv.Itoa(f.CR), protocolo.NovoNúmeroControle(f.ID, f.Controle), códigoVerificação)

	qr, err := qrcode.New(qrURL, qrcode.Highest)
	if err != nil {
		return erros.Novo(err)
	}

	qr.BackgroundColor = color.Transparent
	qr.ForegroundColor = color.RGBA{0x00, 0x00, 0x00, 0xff}

	qrPonto := image.Pt(200, 600)
	qrPosição := imagem.Bounds().Min.Sub(qrPonto)
	draw.Draw(imagem, imagem.Bounds(), qr.Image(300), qrPosição, draw.Over)

	// codifica a imagem
	var buffer bytes.Buffer
	if err := png.Encode(&buffer, imagem); err != nil {
		return erros.Novo(err)
	}

	f.ImagemNúmeroControle = base64.StdEncoding.EncodeToString(buffer.Bytes())
	return nil
}

type imagemTexto struct {
	texto        string
	fonteCor     color.RGBA
	fonteTamanho float64
	posição      imagemTextoPosição
}

type imagemTextoPosição struct {
	x float64
	y float64
}
