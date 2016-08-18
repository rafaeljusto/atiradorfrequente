package config

import (
	"image"
	"image/color"
	_ "image/gif"  // adiciona o suporte para imagens GIF no image.Decode
	_ "image/jpeg" // adiciona o suporte para imagens JPEG no image.Decode
	_ "image/png"  // adiciona o suporte para imagens PNG no image.Decode
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/erros"
	"github.com/registrobr/gostk/errors"
	"golang.org/x/image/font"
)

// Configuração define os valores configuráveis referentes a regras de negócio e
// políticas nos serviços.
type Configuração struct {
	Atirador struct {
		// PrazoConfirmação define o tempo máximo permitido para confirmar uma
		// frequência a partir do momento de sua criação.
		PrazoConfirmação time.Duration `yaml:"prazo confirmacao" envconfig:"prazo_confirmacao"`

		// ImagemNúmeroControle define as propriedades para geração da imagem que
		// contém o número de controle.
		ImagemNúmeroControle struct {
			// Largura define a largura em pixels da imagem.
			Largura int `yaml:"largura" envconfig:"largura"`
			// Altura define a altura em pixels da imagem.
			Altura int `yaml:"altura" envconfig:"altura"`
			// CorFundo define a cor de fundo da imagem. As possível opções são:
			// "preto", "branco", "verde", "azul", "vermelho", "amarelo" e "cinza"
			CorFundo cor `yaml:"cor fundo" envconfig:"cor_fundo"`
			// Fonte define as propriedades da fonte utilizada na imagem.
			Fonte fonte `yaml:"fonte" envconfig:"fonte"`
			// Logo define o logo a ser exibido na imagem no canto superior esquerdo.
			Logo struct {
				// Imagem caminho para o arquivo que contém a imagem, são suportados os
				// formatos: JPEG, PNG e GIF. Caso o campo esteja indefinido a imagem
				// não será exibida. Esta opção tende a deixar a impressão da imagem
				// mais lenta.
				Imagem imagem `yaml:"imagem" envconfig:"imagem"`
				// Espaçamento define o espaço em pixels entre a borda superior e a
				// borda esquerda da imagem do logo.
				Espaçamento int `yaml:"espacamento" envconfig:"espacamento"`
			} `yaml:"logo" envconfig:"logo"`
			// Borda define as características da borda. Esta opção tende a deixar a
			// impressão da imagem mais lenta.
			Borda struct {
				// Largura define a espessura da borda. Caso o valor seja 0 a borda não
				// será exibida.
				Largura int `yaml:"largura" envconfig:"largura"`
				// Espaçamento define o espeço em pixels entre a borda e as margens da
				// imagem.
				Espaçamento int `yaml:"espacamento" envconfig:"espacamento"`
				// Cor define a cor da borda. As possível opções são:
				// "preto", "branco", "verde", "azul", "vermelho", "amarelo" e "cinza"
				Cor cor `yaml:"cor" envconfig:"cor"`
			} `yaml:"borda" envconfig:"borda"`
			// LinhaFundo define as características da linha que é exibida
			// repetidamente no fundo da imagem. Esta opção tende a deixar a impressão
			// da imagem mais lenta.
			LinhaFundo struct {
				// Largura define a espessura da linha. Caso o valor seja 0 a borda não
				// será exibida.
				Largura int `yaml:"largura" envconfig:"largura"`
				// Espaçamento define o espeço em pixels entre as linhas.
				Espaçamento int `yaml:"espacamento" envconfig:"espacamento"`
				// Cor define a cor da linha. As possível opções são:
				// "preto", "branco", "verde", "azul", "vermelho", "amarelo" e "cinza"
				Cor cor `yaml:"cor" envconfig:"cor"`
			} `yaml:"linha fundo" envconfig:"linha_fundo"`
		} `yaml:"imagem numero controle" envconfig:"imagem_numero_controle"`
	} `yaml:"atirador" envconfig:"atirador"`
}

func DefinirValoresPadrão(c *Configuração) {
	c.Atirador.PrazoConfirmação = 30 * time.Minute
	c.Atirador.ImagemNúmeroControle.Largura = 3508
	c.Atirador.ImagemNúmeroControle.Altura = 2480
	c.Atirador.ImagemNúmeroControle.CorFundo.Color = color.RGBA{0xff, 0xff, 0xff, 0xff}
	c.Atirador.ImagemNúmeroControle.Fonte.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
	c.Atirador.ImagemNúmeroControle.Borda.Largura = 50
	c.Atirador.ImagemNúmeroControle.Borda.Espaçamento = 50
	c.Atirador.ImagemNúmeroControle.Borda.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
	c.Atirador.ImagemNúmeroControle.LinhaFundo.Largura = 50
	c.Atirador.ImagemNúmeroControle.LinhaFundo.Espaçamento = 50
	c.Atirador.ImagemNúmeroControle.LinhaFundo.Cor.Color = color.RGBA{0xee, 0xee, 0xee, 0xff}
}

type imagem struct {
	image.Image
}

// UnmarshalText carrega um arquivo de imagem. Os formatos suportados são JPEG,
// PNG e GIF.
func (i *imagem) UnmarshalText(texto []byte) error {
	arquivoImagem, err := os.Open(string(texto))
	if err != nil {
		return erros.Novo(err)
	}
	defer arquivoImagem.Close()

	if i.Image, _, err = image.Decode(arquivoImagem); err != nil {
		return erros.Novo(err)
	}

	return nil
}

type fonte struct {
	font.Face
	Cor cor
}

// UnmarshalYAML agrupa as propriedades da fonte para gerar o tipo font.Face.
func (f *fonte) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var detalhes fonteDetalhes
	if err := unmarshal(&detalhes); err != nil {
		return erros.Novo(err)
	}

	f.Cor = detalhes.Cor

	if detalhes.Família.Font != nil {
		f.Face = truetype.NewFace(detalhes.Família.Font, &truetype.Options{
			Size:    detalhes.Tamanho,
			DPI:     detalhes.DPI,
			Hinting: font.HintingNone,
		})
	}

	return nil
}

type fonteDetalhes struct {
	// Família caminho para o arquivo que contém a fonte a ser utilizada na
	// imagem. O único formato suportado no momento é o TTF.
	Família fonteFamília `yaml:"familia"`
	// Tamanho define o tamanho da fonte em pixels a ser utilizado.
	Tamanho float64 `yaml:"tamanho"`
	// DPI define a resolução da fonte.
	DPI float64 `yaml:"dpi"`
	// Cor define a cor da fonte.
	Cor cor `yaml:"cor"`
}

type fonteFamília struct {
	*truetype.Font
}

// UnmarshalText carrega um arquivo de fonte TTF.
func (f *fonteFamília) UnmarshalText(texto []byte) error {
	fonte, err := ioutil.ReadFile(string(texto))
	if err != nil {
		return erros.Novo(err)
	}

	if f.Font, err = truetype.Parse(fonte); err != nil {
		return erros.Novo(err)
	}

	return nil
}

type cor struct {
	color.Color
}

// UnmarshalText carrega a cor, convertendo do formato texto para a estrutura
// apropriada. As cores possíveis são: "preto", "branco", "verde", "azul",
// "vermelho", "amarelo" e "cinza"
func (c *cor) UnmarshalText(texto []byte) error {
	cor := string(texto)
	cor = strings.ToLower(cor)

	switch cor {
	case "preto":
		c.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
	case "branco":
		c.Color = color.RGBA{0xff, 0xff, 0xff, 0xff}
	case "verde":
		c.Color = color.RGBA{0x00, 0xff, 0x00, 0xff}
	case "azul":
		c.Color = color.RGBA{0x00, 0x00, 0xff, 0xff}
	case "vermelho":
		c.Color = color.RGBA{0xff, 0x00, 0x00, 0xff}
	case "amarelo":
		c.Color = color.RGBA{0xff, 0xff, 0x00, 0xff}
	case "cinza":
		c.Color = color.RGBA{0xee, 0xee, 0xee, 0xff}
	default:
		return errors.Errorf("cor “%s” desconhecida", cor)
	}

	return nil
}
