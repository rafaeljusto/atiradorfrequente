package config

import (
	"image"
	_ "image/gif"  // adiciona o suporte para imagens GIF no image.Decode
	_ "image/jpeg" // adiciona o suporte para imagens JPEG no image.Decode
	_ "image/png"  // adiciona o suporte para imagens PNG no image.Decode
	"io/ioutil"
	"os"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/erros"
	"golang.org/x/image/font/gofont/goregular"
)

// Configuração define os valores configuráveis referentes a regras de negócio e
// políticas nos serviços.
type Configuração struct {
	Atirador struct {
		// PrazoConfirmação define o tempo máximo permitido para confirmar uma
		// frequência a partir do momento de sua criação.
		PrazoConfirmação time.Duration `yaml:"prazo confirmacao" envconfig:"prazo_confirmacao"`

		// DuraçãoMáximaTreino tempo máximo permitido de duração de um treino.
		DuraçãoMáximaTreino time.Duration `yaml:"duracao maxima treino" envconfig:"duracao_maxima_treino"`

		// ImagemNúmeroControle define as propriedades para geração da imagem que
		// contém o número de controle.
		ImagemNúmeroControle struct {
			// Fonte define as propriedades da fonte utilizada na imagem.
			Fonte fonteFamília `yaml:"fonte" envconfig:"fonte"`
			// ImagemBase caminho para o arquivo que contém a imagem, são suportados
			// os formatos: JPEG, PNG e GIF.
			ImagemBase imagem `yaml:"imagem base" envconfig:"imagem_base"`
		} `yaml:"imagem numero controle" envconfig:"imagem_numero_controle"`
	} `yaml:"atirador" envconfig:"atirador"`
}

// DefinirValoresPadrão utiliza valores padrão em todos os campos da
// configuração caso o usuário não informe. O usuário também tem a opção de
// sobrescrever somente alguns valores, mantendo os demais com valores padrão.
func DefinirValoresPadrão(c *Configuração) {
	c.Atirador.PrazoConfirmação = 30 * time.Minute
	c.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
	c.Atirador.ImagemNúmeroControle.Fonte.Font, _ = truetype.Parse(goregular.TTF)
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
