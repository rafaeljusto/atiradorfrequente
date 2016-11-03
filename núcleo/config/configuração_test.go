package config_test

import (
	"encoding/base64"
	"fmt"
	"image/color"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/config"
	"github.com/rafaeljusto/atiradorfrequente/testes"
	"github.com/registrobr/gostk/errors"
	"gopkg.in/yaml.v2"
)

func TestConfiguração_yaml(t *testing.T) {
	arquivoQualquer, err := ioutil.TempFile("", "teste-nucleo-config-")
	if err != nil {
		t.Fatalf("Erro gerar um arquivo qualquer. Detalhes: %s", err)
	}
	arquivoQualquer.Close()

	arquivoFonte, err := ioutil.TempFile("", "teste-nucleo-config-")
	if err != nil {
		t.Fatalf("Erro gerar o arquivo da fonte. Detalhes: %s", err)
	}

	fonteExtraída, err := base64.StdEncoding.DecodeString(fonteTTF)
	if err != nil {
		t.Fatalf("Erro ao extrair a imagem de teste do logo. Detalhes: %s", err)
	}

	arquivoFonte.Write(fonteExtraída)
	arquivoFonte.Close()

	arquivoImagemLogo, err := ioutil.TempFile("", "teste-nucleo-config-")
	if err != nil {
		t.Fatalf("Erro gerar o arquivo do logo. Detalhes: %s", err)
	}

	imagemLogoExtraída, err := base64.StdEncoding.DecodeString(imagemLogoPNG)
	if err != nil {
		t.Fatalf("Erro ao extrair a imagem de teste do logo. Detalhes: %s", err)
	}

	arquivoImagemLogo.Write(imagemLogoExtraída)
	arquivoImagemLogo.Close()

	cenários := []struct {
		descrição            string
		conteúdoArquivo      string
		deveConterFonte      bool
		deveConterImagemLogo bool
		configuraçãoEsperada config.Configuração
		erroEsperado         error
	}{
		{
			descrição: "deve carregar a configuração corretamente",
			conteúdoArquivo: `
atirador:
  prazo confirmacao: 30m
  duracao maxima treino: 12h
  imagem numero controle:
    largura: 3508
    altura: 2480
    cor fundo: branco
    fonte:
      familia: ` + arquivoFonte.Name() + `
      tamanho: 48
      dpi: 300
      cor: preto
    logo:
      imagem: ` + arquivoImagemLogo.Name() + `
      espacamento: 100
    borda:
      largura: 50
      espacamento: 50
      cor: preto
    linha fundo:
      largura: 50
      espacamento: 50
      cor: cinza
`,
			deveConterFonte:      true,
			deveConterImagemLogo: true,
			configuraçãoEsperada: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.PrazoConfirmação = 30 * time.Minute
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				configuração.Atirador.ImagemNúmeroControle.Largura = 3508
				configuração.Atirador.ImagemNúmeroControle.Altura = 2480
				configuração.Atirador.ImagemNúmeroControle.CorFundo.Color = color.RGBA{0xff, 0xff, 0xff, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Fonte.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Logo.Espaçamento = 100
				configuração.Atirador.ImagemNúmeroControle.Borda.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Cor.Color = color.RGBA{0xee, 0xee, 0xee, 0xff}
				return configuração
			}(),
		},
		{
			descrição: "deve ignorar quando a imagem ou a fonte não são informados",
			conteúdoArquivo: `
atirador:
  prazo confirmacao: 30m
  duracao maxima treino: 12h
  imagem numero controle:
    largura: 3508
    altura: 2480
    cor fundo: verde
    fonte:
      familia:
      tamanho: 48
      dpi: 300
      cor: azul
    logo:
      imagem:
      espacamento: 100
    borda:
      largura: 50
      espacamento: 50
      cor: vermelho
    linha fundo:
      largura: 50
      espacamento: 50
      cor: amarelo
`,
			configuraçãoEsperada: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.PrazoConfirmação = 30 * time.Minute
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				configuração.Atirador.ImagemNúmeroControle.Largura = 3508
				configuração.Atirador.ImagemNúmeroControle.Altura = 2480
				configuração.Atirador.ImagemNúmeroControle.CorFundo.Color = color.RGBA{0x00, 0xff, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Fonte.Cor.Color = color.RGBA{0x00, 0x00, 0xff, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Logo.Espaçamento = 100
				configuração.Atirador.ImagemNúmeroControle.Borda.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Cor.Color = color.RGBA{0xff, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Cor.Color = color.RGBA{0xff, 0xff, 0x00, 0xff}
				return configuração
			}(),
		},
		{
			descrição: "deve detectar quando o arquivo da fonte não existe",
			conteúdoArquivo: `
atirador:
  prazo confirmacao: 30m
  duracao maxima treino: 12h
  imagem numero controle:
    largura: 3508
    altura: 2480
    cor fundo: branco
    fonte:
      familia: /tmp/eunaoexisto321.ttf
      tamanho: 48
      dpi: 300
      cor: preto
    logo:
      imagem: ` + arquivoImagemLogo.Name() + `
      espacamento: 100
    borda:
      largura: 50
      espacamento: 50
      cor: preto
    linha fundo:
      largura: 50
      espacamento: 50
      cor: cinza
`,
			configuraçãoEsperada: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.PrazoConfirmação = 30 * time.Minute
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				configuração.Atirador.ImagemNúmeroControle.Largura = 3508
				configuração.Atirador.ImagemNúmeroControle.Altura = 2480
				configuração.Atirador.ImagemNúmeroControle.CorFundo.Color = color.RGBA{0xff, 0xff, 0xff, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Fonte.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Logo.Espaçamento = 100
				configuração.Atirador.ImagemNúmeroControle.Borda.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Cor.Color = color.RGBA{0xee, 0xee, 0xee, 0xff}
				return configuração
			}(),
			erroEsperado: &os.PathError{
				Op:   "open",
				Path: "/tmp/eunaoexisto321.ttf",
				Err:  fmt.Errorf("no such file or directory"),
			},
		},
		{
			descrição: "deve detectar quando o arquivo de fonte esta no formato inválido",
			conteúdoArquivo: `
atirador:
  prazo confirmacao: 30m
  duracao maxima treino: 12h
  imagem numero controle:
    largura: 3508
    altura: 2480
    cor fundo: branco
    fonte:
      familia: ` + arquivoQualquer.Name() + `
      tamanho: 48
      dpi: 300
      cor: preto
    logo:
      imagem: ` + arquivoImagemLogo.Name() + `
      espacamento: 100
    borda:
      largura: 50
      espacamento: 50
      cor: preto
    linha fundo:
      largura: 50
      espacamento: 50
      cor: cinza
`,
			configuraçãoEsperada: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.PrazoConfirmação = 30 * time.Minute
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				configuração.Atirador.ImagemNúmeroControle.Largura = 3508
				configuração.Atirador.ImagemNúmeroControle.Altura = 2480
				configuração.Atirador.ImagemNúmeroControle.CorFundo.Color = color.RGBA{0xff, 0xff, 0xff, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Fonte.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Logo.Espaçamento = 100
				configuração.Atirador.ImagemNúmeroControle.Borda.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Cor.Color = color.RGBA{0xee, 0xee, 0xee, 0xff}
				return configuração
			}(),
			erroEsperado: errors.Errorf("freetype: invalid TrueType format: TTF data is too short"),
		},
		{
			descrição: "deve detectar quando o arquivo de imagem do logo não existe",
			conteúdoArquivo: `
atirador:
  prazo confirmacao: 30m
  duracao maxima treino: 12h
  imagem numero controle:
    largura: 3508
    altura: 2480
    cor fundo: branco
    fonte:
      familia: ` + arquivoFonte.Name() + `
      tamanho: 48
      dpi: 300
      cor: preto
    logo:
      imagem: /tmp/eunaoexisto321.png
      espacamento: 100
    borda:
      largura: 50
      espacamento: 50
      cor: preto
    linha fundo:
      largura: 50
      espacamento: 50
      cor: cinza
`,
			deveConterFonte: true,
			configuraçãoEsperada: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.PrazoConfirmação = 30 * time.Minute
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				configuração.Atirador.ImagemNúmeroControle.Largura = 3508
				configuração.Atirador.ImagemNúmeroControle.Altura = 2480
				configuração.Atirador.ImagemNúmeroControle.CorFundo.Color = color.RGBA{0xff, 0xff, 0xff, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Fonte.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Logo.Espaçamento = 100
				configuração.Atirador.ImagemNúmeroControle.Borda.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Cor.Color = color.RGBA{0xee, 0xee, 0xee, 0xff}
				return configuração
			}(),
			erroEsperado: &os.PathError{
				Op:   "open",
				Path: "/tmp/eunaoexisto321.png",
				Err:  fmt.Errorf("no such file or directory"),
			},
		},
		{
			descrição: "deve detectar quando a imagem do logo esta em um formato inválido",
			conteúdoArquivo: `
atirador:
  prazo confirmacao: 30m
  duracao maxima treino: 12h
  imagem numero controle:
    largura: 3508
    altura: 2480
    cor fundo: branco
    fonte:
      familia: ` + arquivoFonte.Name() + `
      tamanho: 48
      dpi: 300
      cor: preto
    logo:
      imagem: ` + arquivoQualquer.Name() + `
      espacamento: 100
    borda:
      largura: 50
      espacamento: 50
      cor: preto
    linha fundo:
      largura: 50
      espacamento: 50
      cor: cinza
`,
			deveConterFonte: true,
			configuraçãoEsperada: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.PrazoConfirmação = 30 * time.Minute
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				configuração.Atirador.ImagemNúmeroControle.Largura = 3508
				configuração.Atirador.ImagemNúmeroControle.Altura = 2480
				configuração.Atirador.ImagemNúmeroControle.CorFundo.Color = color.RGBA{0xff, 0xff, 0xff, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Fonte.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Logo.Espaçamento = 100
				configuração.Atirador.ImagemNúmeroControle.Borda.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Cor.Color = color.RGBA{0xee, 0xee, 0xee, 0xff}
				return configuração
			}(),
			erroEsperado: errors.Errorf("image: unknown format"),
		},
		{
			descrição: "deve detectar uma cor desconhecida",
			conteúdoArquivo: `
atirador:
  prazo confirmacao: 30m
  duracao maxima treino: 12h
  imagem numero controle:
    largura: 3508
    altura: 2480
    cor fundo: roxo
    fonte:
      familia: ` + arquivoFonte.Name() + `
      tamanho: 48
      dpi: 300
      cor: preto
    logo:
      imagem: ` + arquivoImagemLogo.Name() + `
      espacamento: 100
    borda:
      largura: 50
      espacamento: 50
      cor: preto
    linha fundo:
      largura: 50
      espacamento: 50
      cor: cinza
`,
			configuraçãoEsperada: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.PrazoConfirmação = 30 * time.Minute
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				configuração.Atirador.ImagemNúmeroControle.Largura = 3508
				configuração.Atirador.ImagemNúmeroControle.Altura = 2480
				configuração.Atirador.ImagemNúmeroControle.CorFundo.Color = color.RGBA{0xff, 0xff, 0xff, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Fonte.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Logo.Espaçamento = 100
				configuração.Atirador.ImagemNúmeroControle.Borda.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Cor.Color = color.RGBA{0xee, 0xee, 0xee, 0xff}
				return configuração
			}(),
			erroEsperado: errors.Errorf("cor “roxo” desconhecida"),
		},
	}

	for i, cenário := range cenários {
		var configuração config.Configuração
		err := yaml.Unmarshal([]byte(cenário.conteúdoArquivo), &configuração)

		if cenário.deveConterFonte && configuração.Atirador.ImagemNúmeroControle.Fonte.Face == nil {
			t.Errorf("Item %d, “%s”: fonte não foi carregada corretamente",
				i, cenário.descrição)
		}

		if cenário.deveConterImagemLogo && configuração.Atirador.ImagemNúmeroControle.Logo.Imagem.Image == nil {
			t.Errorf("Item %d, “%s”: imagem do logo não foi carregada corretamente",
				i, cenário.descrição)
		}

		// a comparação de fonte e imagem consome muita memória, portanto nos
		// restringimos a uma verificação simples
		configuração.Atirador.ImagemNúmeroControle.Fonte.Face = nil
		configuração.Atirador.ImagemNúmeroControle.Logo.Imagem.Image = nil

		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(cenário.configuraçãoEsperada, cenário.erroEsperado)
		if err = verificadorResultado.VerificaResultado(configuração, err); err != nil {
			t.Error(err)
		}
	}
}

func TestConfiguração_variáveisDeAmbiente(t *testing.T) {
	arquivoQualquer, err := ioutil.TempFile("", "teste-nucleo-config-")
	if err != nil {
		t.Fatalf("Erro gerar um arquivo qualquer. Detalhes: %s", err)
	}
	arquivoQualquer.Close()

	arquivoFonte, err := ioutil.TempFile("", "teste-nucleo-config-")
	if err != nil {
		t.Fatalf("Erro gerar o arquivo da fonte. Detalhes: %s", err)
	}

	fonteExtraída, err := base64.StdEncoding.DecodeString(fonteTTF)
	if err != nil {
		t.Fatalf("Erro ao extrair a imagem de teste do logo. Detalhes: %s", err)
	}

	arquivoFonte.Write(fonteExtraída)
	arquivoFonte.Close()

	arquivoImagemLogo, err := ioutil.TempFile("", "teste-nucleo-config-")
	if err != nil {
		t.Fatalf("Erro gerar o arquivo do logo. Detalhes: %s", err)
	}

	imagemLogoExtraída, err := base64.StdEncoding.DecodeString(imagemLogoPNG)
	if err != nil {
		t.Fatalf("Erro ao extrair a imagem de teste do logo. Detalhes: %s", err)
	}

	arquivoImagemLogo.Write(imagemLogoExtraída)
	arquivoImagemLogo.Close()

	cenários := []struct {
		descrição            string
		variáveisAmbiente    map[string]string
		deveConterFonte      bool
		deveConterImagemLogo bool
		configuraçãoEsperada config.Configuração
		erroEsperado         error
	}{
		{
			descrição: "deve carregar a configuração corretamente",
			variáveisAmbiente: map[string]string{
				"AF_ATIRADOR_PRAZO_CONFIRMACAO":                              "30m",
				"AF_ATIRADOR_DURACAO_MAXIMA_TREINO":                          "12h",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LARGURA":                 "3508",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_ALTURA":                  "2480",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_COR_FUNDO":               "branco",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_FONTE":                   arquivoFonte.Name() + " 48 300 preto",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LOGO_IMAGEM":             arquivoImagemLogo.Name(),
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LOGO_ESPACAMENTO":        "100",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_BORDA_LARGURA":           "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_BORDA_ESPACAMENTO":       "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_BORDA_COR":               "preto",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LINHA_FUNDO_LARGURA":     "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LINHA_FUNDO_ESPACAMENTO": "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LINHA_FUNDO_COR":         "cinza",
			},
			deveConterFonte:      true,
			deveConterImagemLogo: true,
			configuraçãoEsperada: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.PrazoConfirmação = 30 * time.Minute
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				configuração.Atirador.ImagemNúmeroControle.Largura = 3508
				configuração.Atirador.ImagemNúmeroControle.Altura = 2480
				configuração.Atirador.ImagemNúmeroControle.CorFundo.Color = color.RGBA{0xff, 0xff, 0xff, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Fonte.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Logo.Espaçamento = 100
				configuração.Atirador.ImagemNúmeroControle.Borda.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Cor.Color = color.RGBA{0xee, 0xee, 0xee, 0xff}
				return configuração
			}(),
		},
		{
			descrição: "deve ignorar quando a imagem ou a fonte não são informados",
			variáveisAmbiente: map[string]string{
				"AF_ATIRADOR_PRAZO_CONFIRMACAO":                              "30m",
				"AF_ATIRADOR_DURACAO_MAXIMA_TREINO":                          "12h",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LARGURA":                 "3508",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_ALTURA":                  "2480",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_COR_FUNDO":               "verde",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LOGO_ESPACAMENTO":        "100",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_BORDA_LARGURA":           "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_BORDA_ESPACAMENTO":       "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_BORDA_COR":               "vermelho",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LINHA_FUNDO_LARGURA":     "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LINHA_FUNDO_ESPACAMENTO": "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LINHA_FUNDO_COR":         "amarelo",
			},
			configuraçãoEsperada: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.PrazoConfirmação = 30 * time.Minute
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				configuração.Atirador.ImagemNúmeroControle.Largura = 3508
				configuração.Atirador.ImagemNúmeroControle.Altura = 2480
				configuração.Atirador.ImagemNúmeroControle.CorFundo.Color = color.RGBA{0x00, 0xff, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Logo.Espaçamento = 100
				configuração.Atirador.ImagemNúmeroControle.Borda.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Cor.Color = color.RGBA{0xff, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Cor.Color = color.RGBA{0xff, 0xff, 0x00, 0xff}
				return configuração
			}(),
		},
		{
			descrição: "deve detectar quando a quantidade de argumentos da fonte é menor do que a necessária",
			variáveisAmbiente: map[string]string{
				"AF_ATIRADOR_PRAZO_CONFIRMACAO":                              "30m",
				"AF_ATIRADOR_DURACAO_MAXIMA_TREINO":                          "12h",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LARGURA":                 "3508",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_ALTURA":                  "2480",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_COR_FUNDO":               "branco",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_FONTE":                   "48 300 preto",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LOGO_IMAGEM":             arquivoImagemLogo.Name(),
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LOGO_ESPACAMENTO":        "100",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_BORDA_LARGURA":           "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_BORDA_ESPACAMENTO":       "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_BORDA_COR":               "vermelho",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LINHA_FUNDO_LARGURA":     "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LINHA_FUNDO_ESPACAMENTO": "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LINHA_FUNDO_COR":         "amarelo",
			},
			configuraçãoEsperada: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.PrazoConfirmação = 30 * time.Minute
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				configuração.Atirador.ImagemNúmeroControle.Largura = 3508
				configuração.Atirador.ImagemNúmeroControle.Altura = 2480
				configuração.Atirador.ImagemNúmeroControle.CorFundo.Color = color.RGBA{0x00, 0xff, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Fonte.Cor.Color = color.RGBA{0x00, 0x00, 0xff, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Logo.Espaçamento = 100
				configuração.Atirador.ImagemNúmeroControle.Borda.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Cor.Color = color.RGBA{0xff, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Cor.Color = color.RGBA{0xff, 0xff, 0x00, 0xff}
				return configuração
			}(),
			erroEsperado: errors.Errorf("fonte não contém as informações necessárias"),
		},
		{
			descrição: "deve detectar quando o arquivo da fonte não existe",
			variáveisAmbiente: map[string]string{
				"AF_ATIRADOR_PRAZO_CONFIRMACAO":                              "30m",
				"AF_ATIRADOR_DURACAO_MAXIMA_TREINO":                          "12h",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LARGURA":                 "3508",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_ALTURA":                  "2480",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_COR_FUNDO":               "branco",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_FONTE":                   "/tmp/eunaoexisto321.ttf 48 300 preto",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LOGO_IMAGEM":             arquivoImagemLogo.Name(),
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LOGO_ESPACAMENTO":        "100",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_BORDA_LARGURA":           "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_BORDA_ESPACAMENTO":       "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_BORDA_COR":               "preto",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LINHA_FUNDO_LARGURA":     "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LINHA_FUNDO_ESPACAMENTO": "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LINHA_FUNDO_COR":         "cinza",
			},
			configuraçãoEsperada: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.PrazoConfirmação = 30 * time.Minute
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				configuração.Atirador.ImagemNúmeroControle.Largura = 3508
				configuração.Atirador.ImagemNúmeroControle.Altura = 2480
				configuração.Atirador.ImagemNúmeroControle.CorFundo.Color = color.RGBA{0xff, 0xff, 0xff, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Fonte.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Logo.Espaçamento = 100
				configuração.Atirador.ImagemNúmeroControle.Borda.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Cor.Color = color.RGBA{0xee, 0xee, 0xee, 0xff}
				return configuração
			}(),
			erroEsperado: errors.Errorf("open /tmp/eunaoexisto321.ttf: no such file or directory"),
		},
		{
			descrição: "deve detectar quando o arquivo de fonte esta no formato inválido",
			variáveisAmbiente: map[string]string{
				"AF_ATIRADOR_PRAZO_CONFIRMACAO":                              "30m",
				"AF_ATIRADOR_DURACAO_MAXIMA_TREINO":                          "12h",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LARGURA":                 "3508",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_ALTURA":                  "2480",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_COR_FUNDO":               "branco",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_FONTE":                   arquivoQualquer.Name() + " 48 300 preto",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LOGO_IMAGEM":             arquivoImagemLogo.Name(),
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LOGO_ESPACAMENTO":        "100",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_BORDA_LARGURA":           "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_BORDA_ESPACAMENTO":       "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_BORDA_COR":               "preto",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LINHA_FUNDO_LARGURA":     "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LINHA_FUNDO_ESPACAMENTO": "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LINHA_FUNDO_COR":         "cinza",
			},
			configuraçãoEsperada: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.PrazoConfirmação = 30 * time.Minute
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				configuração.Atirador.ImagemNúmeroControle.Largura = 3508
				configuração.Atirador.ImagemNúmeroControle.Altura = 2480
				configuração.Atirador.ImagemNúmeroControle.CorFundo.Color = color.RGBA{0xff, 0xff, 0xff, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Fonte.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Logo.Espaçamento = 100
				configuração.Atirador.ImagemNúmeroControle.Borda.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Cor.Color = color.RGBA{0xee, 0xee, 0xee, 0xff}
				return configuração
			}(),
			erroEsperado: errors.Errorf("freetype: invalid TrueType format: TTF data is too short"),
		},
		{
			descrição: "deve detectar quando a fonte possui um tamanho inválido",
			variáveisAmbiente: map[string]string{
				"AF_ATIRADOR_PRAZO_CONFIRMACAO":                              "30m",
				"AF_ATIRADOR_DURACAO_MAXIMA_TREINO":                          "12h",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LARGURA":                 "3508",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_ALTURA":                  "2480",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_COR_FUNDO":               "branco",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_FONTE":                   arquivoFonte.Name() + " XX 300 preto",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LOGO_IMAGEM":             arquivoImagemLogo.Name(),
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LOGO_ESPACAMENTO":        "100",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_BORDA_LARGURA":           "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_BORDA_ESPACAMENTO":       "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_BORDA_COR":               "preto",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LINHA_FUNDO_LARGURA":     "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LINHA_FUNDO_ESPACAMENTO": "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LINHA_FUNDO_COR":         "cinza",
			},
			configuraçãoEsperada: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.PrazoConfirmação = 30 * time.Minute
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				configuração.Atirador.ImagemNúmeroControle.Largura = 3508
				configuração.Atirador.ImagemNúmeroControle.Altura = 2480
				configuração.Atirador.ImagemNúmeroControle.CorFundo.Color = color.RGBA{0xff, 0xff, 0xff, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Fonte.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Logo.Espaçamento = 100
				configuração.Atirador.ImagemNúmeroControle.Borda.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Cor.Color = color.RGBA{0xee, 0xee, 0xee, 0xff}
				return configuração
			}(),
			erroEsperado: errors.Errorf(`strconv.ParseFloat: parsing "XX": invalid syntax`),
		},
		{
			descrição: "deve detectar quando a fonte possui um DPI inválido",
			variáveisAmbiente: map[string]string{
				"AF_ATIRADOR_PRAZO_CONFIRMACAO":                              "30m",
				"AF_ATIRADOR_DURACAO_MAXIMA_TREINO":                          "12h",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LARGURA":                 "3508",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_ALTURA":                  "2480",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_COR_FUNDO":               "branco",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_FONTE":                   arquivoFonte.Name() + " 48 XXX preto",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LOGO_IMAGEM":             arquivoImagemLogo.Name(),
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LOGO_ESPACAMENTO":        "100",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_BORDA_LARGURA":           "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_BORDA_ESPACAMENTO":       "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_BORDA_COR":               "preto",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LINHA_FUNDO_LARGURA":     "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LINHA_FUNDO_ESPACAMENTO": "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LINHA_FUNDO_COR":         "cinza",
			},
			configuraçãoEsperada: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.PrazoConfirmação = 30 * time.Minute
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				configuração.Atirador.ImagemNúmeroControle.Largura = 3508
				configuração.Atirador.ImagemNúmeroControle.Altura = 2480
				configuração.Atirador.ImagemNúmeroControle.CorFundo.Color = color.RGBA{0xff, 0xff, 0xff, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Fonte.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Logo.Espaçamento = 100
				configuração.Atirador.ImagemNúmeroControle.Borda.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Cor.Color = color.RGBA{0xee, 0xee, 0xee, 0xff}
				return configuração
			}(),
			erroEsperado: errors.Errorf(`strconv.ParseFloat: parsing "XXX": invalid syntax`),
		},
		{
			descrição: "deve detectar quando a fonte possui uma cor inválida",
			variáveisAmbiente: map[string]string{
				"AF_ATIRADOR_PRAZO_CONFIRMACAO":                              "30m",
				"AF_ATIRADOR_DURACAO_MAXIMA_TREINO":                          "12h",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LARGURA":                 "3508",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_ALTURA":                  "2480",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_COR_FUNDO":               "branco",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_FONTE":                   arquivoFonte.Name() + " 48 300 roxo",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LOGO_IMAGEM":             arquivoImagemLogo.Name(),
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LOGO_ESPACAMENTO":        "100",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_BORDA_LARGURA":           "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_BORDA_ESPACAMENTO":       "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_BORDA_COR":               "preto",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LINHA_FUNDO_LARGURA":     "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LINHA_FUNDO_ESPACAMENTO": "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LINHA_FUNDO_COR":         "cinza",
			},
			configuraçãoEsperada: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.PrazoConfirmação = 30 * time.Minute
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				configuração.Atirador.ImagemNúmeroControle.Largura = 3508
				configuração.Atirador.ImagemNúmeroControle.Altura = 2480
				configuração.Atirador.ImagemNúmeroControle.CorFundo.Color = color.RGBA{0xff, 0xff, 0xff, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Fonte.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Logo.Espaçamento = 100
				configuração.Atirador.ImagemNúmeroControle.Borda.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Cor.Color = color.RGBA{0xee, 0xee, 0xee, 0xff}
				return configuração
			}(),
			erroEsperado: errors.Errorf(`cor “roxo” desconhecida`),
		},
		{
			descrição: "deve detectar quando o arquivo de imagem do logo não existe",
			variáveisAmbiente: map[string]string{
				"AF_ATIRADOR_PRAZO_CONFIRMACAO":                              "30m",
				"AF_ATIRADOR_DURACAO_MAXIMA_TREINO":                          "12h",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LARGURA":                 "3508",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_ALTURA":                  "2480",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_COR_FUNDO":               "branco",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_FONTE":                   arquivoFonte.Name() + " 48 300 preto",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LOGO_IMAGEM":             "/tmp/eunaoexisto321.png",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LOGO_ESPACAMENTO":        "100",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_BORDA_LARGURA":           "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_BORDA_ESPACAMENTO":       "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_BORDA_COR":               "preto",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LINHA_FUNDO_LARGURA":     "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LINHA_FUNDO_ESPACAMENTO": "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LINHA_FUNDO_COR":         "cinza",
			},
			deveConterFonte: true,
			configuraçãoEsperada: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.PrazoConfirmação = 30 * time.Minute
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				configuração.Atirador.ImagemNúmeroControle.Largura = 3508
				configuração.Atirador.ImagemNúmeroControle.Altura = 2480
				configuração.Atirador.ImagemNúmeroControle.CorFundo.Color = color.RGBA{0xff, 0xff, 0xff, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Fonte.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Logo.Espaçamento = 100
				configuração.Atirador.ImagemNúmeroControle.Borda.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Cor.Color = color.RGBA{0xee, 0xee, 0xee, 0xff}
				return configuração
			}(),
			erroEsperado: errors.Errorf("open /tmp/eunaoexisto321.png: no such file or directory"),
		},
		{
			descrição: "deve detectar quando a imagem do logo esta em um formato inválido",
			variáveisAmbiente: map[string]string{
				"AF_ATIRADOR_PRAZO_CONFIRMACAO":                              "30m",
				"AF_ATIRADOR_DURACAO_MAXIMA_TREINO":                          "12h",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LARGURA":                 "3508",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_ALTURA":                  "2480",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_COR_FUNDO":               "branco",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_FONTE":                   arquivoFonte.Name() + " 48 300 preto",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LOGO_IMAGEM":             arquivoQualquer.Name(),
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LOGO_ESPACAMENTO":        "100",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_BORDA_LARGURA":           "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_BORDA_ESPACAMENTO":       "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_BORDA_COR":               "preto",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LINHA_FUNDO_LARGURA":     "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LINHA_FUNDO_ESPACAMENTO": "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LINHA_FUNDO_COR":         "cinza",
			},
			deveConterFonte: true,
			configuraçãoEsperada: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.PrazoConfirmação = 30 * time.Minute
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				configuração.Atirador.ImagemNúmeroControle.Largura = 3508
				configuração.Atirador.ImagemNúmeroControle.Altura = 2480
				configuração.Atirador.ImagemNúmeroControle.CorFundo.Color = color.RGBA{0xff, 0xff, 0xff, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Fonte.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Logo.Espaçamento = 100
				configuração.Atirador.ImagemNúmeroControle.Borda.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Cor.Color = color.RGBA{0xee, 0xee, 0xee, 0xff}
				return configuração
			}(),
			erroEsperado: errors.Errorf("image: unknown format"),
		},
		{
			descrição: "deve detectar uma cor desconhecida",
			variáveisAmbiente: map[string]string{
				"AF_ATIRADOR_PRAZO_CONFIRMACAO":                              "30m",
				"AF_ATIRADOR_DURACAO_MAXIMA_TREINO":                          "12h",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LARGURA":                 "3508",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_ALTURA":                  "2480",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_COR_FUNDO":               "roxo",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_FONTE":                   arquivoFonte.Name() + " 48 300 preto",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LOGO_IMAGEM":             arquivoImagemLogo.Name(),
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LOGO_ESPACAMENTO":        "100",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_BORDA_LARGURA":           "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_BORDA_ESPACAMENTO":       "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_BORDA_COR":               "preto",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LINHA_FUNDO_LARGURA":     "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LINHA_FUNDO_ESPACAMENTO": "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LINHA_FUNDO_COR":         "cinza",
			},
			configuraçãoEsperada: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.PrazoConfirmação = 30 * time.Minute
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				configuração.Atirador.ImagemNúmeroControle.Largura = 3508
				configuração.Atirador.ImagemNúmeroControle.Altura = 2480
				configuração.Atirador.ImagemNúmeroControle.CorFundo.Color = color.RGBA{0xff, 0xff, 0xff, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Fonte.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Logo.Espaçamento = 100
				configuração.Atirador.ImagemNúmeroControle.Borda.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Cor.Color = color.RGBA{0xee, 0xee, 0xee, 0xff}
				return configuração
			}(),
			erroEsperado: errors.Errorf("cor “roxo” desconhecida"),
		},
	}

	for i, cenário := range cenários {
		os.Clearenv()
		for chave, valor := range cenário.variáveisAmbiente {
			os.Setenv(chave, valor)
		}

		var configuração config.Configuração
		err := envconfig.Process("AF", &configuração)

		// para facilitar a comparação dos erros, vamos extrair o erro de baixo
		// nível da biblioteca que interpreta as variáveis de ambiente
		if erroEspecífico, ok := err.(*envconfig.ParseError); ok {
			err = erroEspecífico.Err
		}

		if cenário.deveConterFonte && configuração.Atirador.ImagemNúmeroControle.Fonte.Face == nil {
			t.Errorf("Item %d, “%s”: fonte não foi carregada corretamente",
				i, cenário.descrição)
		}

		if cenário.deveConterImagemLogo && configuração.Atirador.ImagemNúmeroControle.Logo.Imagem.Image == nil {
			t.Errorf("Item %d, “%s”: imagem do logo não foi carregada corretamente",
				i, cenário.descrição)
		}

		// a comparação de fonte e imagem consome muita memória, portanto nos
		// restringimos a uma verificação simples
		configuração.Atirador.ImagemNúmeroControle.Fonte.Face = nil
		configuração.Atirador.ImagemNúmeroControle.Logo.Imagem.Image = nil

		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(cenário.configuraçãoEsperada, cenário.erroEsperado)
		if err = verificadorResultado.VerificaResultado(configuração, err); err != nil {
			t.Error(err)
		}
	}
}

func TestDefinirValoresPadrão(t *testing.T) {
	var esperado config.Configuração
	esperado.Atirador.PrazoConfirmação = 30 * time.Minute
	esperado.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
	esperado.Atirador.ImagemNúmeroControle.Largura = 3508
	esperado.Atirador.ImagemNúmeroControle.Altura = 2480
	esperado.Atirador.ImagemNúmeroControle.CorFundo.Color = color.RGBA{0xff, 0xff, 0xff, 0xff}
	esperado.Atirador.ImagemNúmeroControle.Fonte.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
	esperado.Atirador.ImagemNúmeroControle.Borda.Largura = 50
	esperado.Atirador.ImagemNúmeroControle.Borda.Espaçamento = 50
	esperado.Atirador.ImagemNúmeroControle.Borda.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
	esperado.Atirador.ImagemNúmeroControle.LinhaFundo.Largura = 50
	esperado.Atirador.ImagemNúmeroControle.LinhaFundo.Espaçamento = 50
	esperado.Atirador.ImagemNúmeroControle.LinhaFundo.Cor.Color = color.RGBA{0xee, 0xee, 0xee, 0xff}

	var c config.Configuração
	config.DefinirValoresPadrão(&c)

	if !reflect.DeepEqual(c, esperado) {
		t.Errorf("Resultados não batem.\n%s", testes.Diff(esperado, c))
	}
}

const imagemLogoPNG = `
iVBORw0KGgoAAAANSUhEUgAAAKgAAACoCAMAAABDlVWGAAABI1BMVEX/////////////////////
////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////8yc1n/AAAAYHRS
TlMAAQIDBAUGBwgJCwwNDg8REhQVFxgZGhseISQnKi0wMzY8P0JFSEtOVFdaXWBjZmlsb3J1eHt+
gYSHio2Qk5aZnKKlqKuusbS3ur3Aw8bJzM/S1djb3uHk5+rt8PP2+fz2kcVfAAADv0lEQVR4AdTB
haGEMBBAwQfEsI2n/06/610Dmxl+uTuXoUjJt+PJcrZ8Oosi1p25nQv/hHobFDJ3DfxxJ49SPt38
iLKi1iqRL5egmlx8CGlFtTUFgKV6lPN1AU5BPTmBtqHe1hZCYgIpIAcTOITsmIDLtJUJrI2+MIGl
M17buwsVV7ItDuP/SyAQ0lLkhNvu7u7u7lYd/d7/KYY6rUTGZ/auYX34wn5xrVqKRf83qEENalCD
GtSgBjWoQQ1qUIMa1KAGNahBDWpQg3qcQQ1qUIMadJ+vhv9eqEHLPf7fR1O3wIw8gLaM71y+QH5C
NdsD1uUcmpm/hbfyCdVoFjhLuIamVwp8lqvF6SrBc0ZuoYmZV766qfXATt1DqVduoR13vPWwOt7T
rpptAgtyCx0tEPW82Fbh7/lZpyQNAKdyC50pA+QXUqrokijuJKWfINfiFjpB1F6Tqtq9ubl5gkNJ
28CUnEJ7S0B5Xp8NXF/366PkNXdNUh9wLqfQlhAoDOmrEEJ9tEMuenDd8tW2E2gyIpQG9C0AvTdC
efBt5hq6CDCpOtDMK6vyAtpSAHZVD3rEbcKPN86bQD6oBx2h3CMvoJkCsKh60Bfecg+dB17SdaF4
A70B5lVRNc05NAOUszGAjgOXigF0E1iLA/QUGFFFLQAtfkEfgc6KWUcOINfhFbQABBWzQy56ey/5
7OOd1HixOO4MCpCsmOXJSlm+CvWzHLx6Bm2RWvnqRVFJAK9u+iNOstlTPguHFBUAea8eTJ1vD6aq
MfDo9OlpTBW1HIbhYYsqGgNOnUE3gDX9rtaADWfQMeBKv6srYMwZNChDuUW/o5YyEDiD6gpY/L0f
ra7kDjoDPCf1myVfgGmH0KYCMPP7LlChySFUG0Au0G8U5IB1uYRmC8CefqM9oJB1CtUCwKR+tWmA
ObmFJq6A0qB+paEScJVwDFXLC1D8FelQEXjJyjVUvRGkNKU6TZWAYo/cQzVG1F6gGgUHRI3KB6im
ygD5pbQqSi8XAMqT8gOqkQJR4VKHvtWxFBJVGJYvULXf8NbT2mhXNOgeXXvirZs2+QNVYiqkZuGU
b38sSC8VqCq/lJJvUCmYPC7yreLxZODrXzVSI1tntwC3Z1sjKfvzi0ENalCDGtSgBjWoQQ1qUIMa
1KAGNahBDWpQgxrUoAY1qEENalCD/megsTmFd2xOih6b08zH5sT9sVmF8L+4LJeIxbqOxnScFqDE
ZqVMDJb0NMZt7ZHU2Ozt/TTZ3KhvpX40JuRhicYfqXguO/N/fdwvoyFJPxBTbvQAAAAASUVORK5C
YII=`

const fonteTTF = `
AAEAAAAMAIAAAwBAT1MvMoO1bWYAAAFIAAAATlBDTFQaFycmAABEEAAAADZjbWFwCTUqbwAAAuAA
AAHOZ2x5ZoXl4xEAAAVYAAA0OGhlYWTPbCe3AAAAzAAAADZoaGVhBosC5gAAAQQAAAAkaG10eKzC
Ay0AAAGYAAABSGtlcm7uEfEkAAA5kAAACKBsb2Nh8A/i0gAABLAAAACmbWF4cAEtAJIAAAEoAAAA
IG5hbWUTWiv5AABCMAAAAP5wb3N0MgA4+wAAQzAAAADdAAEAAAABAABBWIJ4Xw889QAAA+gAAAAA
rfq7ggAAAAC/ECjo/+r/KQNiAzwAAQALAAIAAAAAAAAAAQAAAzz/KQAAA6j/6v/CA2IAAQAAAAAA
AAAAAAAAAAAAAFIAAQAAAFIAUAAGAAAAAAACAAgAQAAKAAAAwgAAAAAAAAAAAhcCvAAFAAQCvAKK
AAAAjwK8AooAAAHFADIBAwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABBbHRzACAAICAQAyAA
yAAAAzwA1wAAAfQAPwAAAAABDgAAAQ4AAAE/ACIBMgAbAij/+gJUAA0A4gAoAQsAFADrABsBZgAH
ATwAGwGfADYBGP/6AVwAGwJlABYBYQAPAi3/+gJJ//oChQAHAkkAJAJRABYB8QAHAm0ADwJtAAcB
OAAbARsAAAIj/+0CtgAHAlQADwI+//oCT//6ApIABwJeABYCMQAPAqYAHgEdAA8B/gAPAkEABwHq
AAcDVf/yAiAAHgKwABYCLwAPA6j/8gJcAA0Ctv/yAa7/6gLgABYCXwAHAxb/+gKSAAACSP/6Ah8A
DwK2AAcCVAAPAj7/+gJP//oCkgAHAl4AFgIxAA8CpgAeAR0ADwH+AA8CQQAHAeoABwNV//ICIAAe
ArAAFgIvAA8CWP/yAlwADQK2//IBrv/qAhQAFgJfAAcDFv/6ApIAAAJI//oCHwAPAkcABwAAAAMA
AAAAAAAAHAABAAAAAAB8AAMAAQAAABwABABgAAAAFAAQAAMABAAiACQAKgA7AD8AWgB6AKMgEP//
AAAAIAAkACYALAA/AEEAYQCjIBD////j/+L/4f/g/93/3P/W/67f/QABAAAAAAAAAAAAAAAAAAAA
AAAAAAAABgFSAAAAAACkAAEAAAAAAAAAAAAAAAAAAAABAAMAAAAAAAAAAgAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAABAAAAAAADAAQABQAAAAYAAAAHAAgACQAKAAsAAAAMAA0ADgAPABAA
EQASABMAFAAVABYAFwAYABkAGgAbAAAAAAAAABwAAAAdAB4AHwAgACEAIgAjACQAJQAmACcAKAAp
ACoAKwAsAC0ALgAvADAAMQAyADMANAA1ADYAAAAAAAAAAAAAAAAANwA4ADkAOgA7ADwAPQA+AD8A
QABBAEIAQwBEAEUARgBHAEgASQBKAEsATABNAE4ATwBQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
UQAAAAAAFAAUABQAFABSAIIBAgFsAZABvAHkAgQCLAJOAnACjgLsAyQDjAQEBGIEwgUyBXYF9AZe
Bp4G3gc2B5AH6gg2CJQJDAlkCc4KOgp8CsALIgtiC94MLAx8DNoNOg2uDiYOfA7KDwoPdA/eECgQ
cBDKESQRcBHOEkYSnhMIE3QTthP6FFwUnBUYFWYVthYSFnIW5hdeF7QYAhhCGKwZFhlgGagaHAAA
AAIAPwAAAbYDIAADAAcAADMRIRElMxEjPwF3/sf6+gMg/OA/AqMAAAADACIALwD/AuIADwAaACIA
ABMfAQ8BBhcWBwYnNT8BJjUTFj8BNicmBxYHFhc/ATQnJgcGV40bByEOKQ1DhhQaDhRKKA0iDUMo
Gw0NDRsiGyIvIQ4C4g0bodBKNmsbFFdRV8kHvP6IBy/rShsHNlGnSt4HIjwNDi8iAAAAAwAbAa4B
IQKzAAkAEQAZAAATMxcGFwcjJzYnFxY3NicmBwYXFjc2JyYHBjzDIg4OIsMhDQ0vIQ4HFSEUB6Ei
FAYaIg0HArMhtQ4hIRSvqAcbaw4NL0ohDS5KIgcbbAAABP/6AAACOgKtABwAPwBHAE8AABM2HwEW
FxUGFxYHBgcGIwYnJicmLwE/AiY/AgMGMxYHFz8BNjc2LwImNyY3FxY3JyYnDwEGBx8BFhcWBycT
MhcWByInJhcyFxYHIicm5FcbKVchIQ1XLxqNKRRXBxQbUBsbBxsGQig1UVEGFHgGGygojSINKDZe
NhUOKF4HFCFYIS8obCEUNlcNDi9RPCkNDi8pDQ3JPAcNLzwHDQKmBxsiFFEoISlXlEMUUAcUPQ0b
IT0vISJDrikb/mYvNQ48BkQaWEMoKBQpGygNIQcbLxsoByEOZFE2FCFKFC8BICE2GyI1wiI2GiE2
AAAAAAQADQBKAkEC6QAXADAAOABAAAATNhcWFwcXFhcHFw8DBicmJyY3Jjc2ExY3Fzc1JzcnDwEG
Jz8BNicGBwYXBwYHBhMyFRYHIicmFzIXFgcGJybQoRs8DgcHLxsOLwcavC+AQ1AHDkoNFChlG1d5
GxshISIbKCgHQ1e8gDUUQwc8Ig3WUQcvQwcOB0ovDS9yBg4C4gcUIjyhDgYpp1gaFQ0NBxsvUKg2
oTU9/Z0HLw0bGzxyIg4vG1EvPHlRB14vUDYbUaEB3SJeBiFR62waGw4veQAAAgAoAZMAqAKzAAkA
EQAAEzYXBhcHBic2JxcWNzYnJgcGSkoUDg4iShQODi8iFAcbIhQHAqYNL7UNIg0vFK6nDi9KIQ4v
SgAAAgAU/9gA0AMRAA0AFwAAExYVAyITBgcGAzYnNjcRMzYvATcmDwEUmjY8DkoNFV09Dg4oLxUG
GwYGGgcUAxENFP6w/nsoBxQBSV5eu1j9vimahgcHG0qaAAIAG/+AANACpgAMABMAABMWHwIHAiMm
JxM2AxM/ASYnBxdDPRobGwc1PS8NNg1DZRQhGyEUFAKmB0NKf/L+3wciAU82AUn+MQaArxoGmwAA
AgAHALUBUAHxAAkADgAAEx8BFgcnJicmNxcWNScHrkpRB4CGKA4NKHIoBikB8RReqCIbDURePIAG
FCEGAAAAAAIAG/9YAScA0AALABUAADc2FxYHBicGJzQnJhcyNzYnJgcGHwFymhUGFDYhXg0iFKgo
FAdKNigHFFG8FHKhG0oHBxuNInjWPF4vFEMoFDYAAgA2ALUBQgFQAAkAEQAAExY3FxYHJgcnJhcy
NzUnIwcXZVdXIg0vfy8iDWQvPSlDLw4BUA4OInIHDg4icnkiKBspKAAC//r/xAEMANcABwAQAAA3
NhcWBwYnJhcWNzUmJyYHBmV4KQZQjRsahho9DigbQw28G2yNBhRea40NQyEiFBRDKAAAAAEAG//Y
AUkCrQANAAATNhcHAhUPASY1NxI3NvI8GxSbDRtXG38UBwKmBy9l/jhYFA0UG38BcYAhAAAAAAQA
Fv/UAk8C9wAOAB0AKgA2AAATJBMGFwYHBgcGJyYnEzYTFjc2NyYnJicmDwIXFhM2FRYHFwcXBwYn
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
NxY3JiciBxbYATolF0R3WiwXDzUdLSZ+LQ+k/rAlFQEpHllaBzR/WRcPDw81SlIXDx0OrNk0DwKJ
JtEtLVJDD1klFgctsxcW4AFB/hsIHlIPNEMPNXBSJQdoDwg1WS3Y0QAAAAAE//oAAAJGAvkAEwAm
AC8ANwAAARY3FwcRFg8BJwcGJyY/ATYXNycBBhcyNxc2NREvASIHFwYjJgcGFzYXFiMGNTYnFzY3
NicPARYBjDtSJg8WByZKlv0tFUIsNbMeCP7OD8ItWnc0Dyw8DwcOF2F3PH+GDw80lQ4OWSUXDzQt
Jg8C+Q8PJX/+mjRaJQcWLeC6Wi0tCCVp/nzCPCUWB38BmzQePI4dQzUlPBdDsw80S0qVBy0XQwg7
PAAAAAIABwBLApgDNQAmAEwAABMWByMWBwYfATcFPwE2JwUnJjM2ITYvAQcmJyYzFzc1JyMHJwYn
IjcXNxY3FxYHJgcVFzcXFiMmBwYXMzYHFCMnBicGJzcmNyc3Jjc0PBYHCA8WCBdDNAEFaSUWfv76
JQcWBwEVJQ80nSUIB1mzLS13NGk0UiVSd0tofyUPNJ1oLZwmDzXnHg8twrMWJsmzQ38IDw8XDwce
FwKCPFJhS1lDJg8PCC1KFwglWh4eJSUHDzRhBx0tHg8IFwgsBw8PDyVwDx48JR4HJZUWLQctQ7Nh
Dw8PJTvCfyY7Ui1pHQAAAAIAFv/xAlwC2wAZADQAABMXNxY3FxYjJwcVFzcXFiMnDwEXBycGJxMD
Ez8BJzcXPwEvAQcnJjMWPwE2JyMGJwYHEyMXS2hSeKsmDjT2PDXYJQ80pDQPByWOPB4PB1lpJQcl
pCYWDy3YJg40d5UmHlqzLWg8JQcHBwLUCA8PDyWdDy0sJggllggmJY4lDx48AVcBI/1+DzSOJggX
JSUeByWVDw8eJTQdFgdL/hMsAAAABQAP//kCMALbABEAGgAzADsAPwAAEwUXBhMGByYHJic1Nicm
JzU2EzYnNzYVFgcGExY3NjcRLwEjBycmBwYXNxcWByMnJgcUFxM2Ny8BDwEGEzM1I7MBVyYPDx6d
cEqdDzQHJgclpA8PJpwPNJ0Xd2EtHg8eHi2NpTQHlY4WDzRaSjUlPLMtFg8lNA8ICC0tAtsPJef+
3Y4WDg4Wd0tDD0M0lpz+1kNDJgc0nQcP/rAPSxaOAWYlFx4tDsmGWggefwhLHjw8JQFfBy00Fw8e
NP7sLAACAB7/kAKgAoIAJABCAAATBx8BMz8BJjc2MzYVFgcTDwEGJzcvAQ8CBi8CPwEnNzYzNgM/
AjY3NhcGFxY3ESYnDwEXByYHJzYvAQ8BFwMW/ggXJS00Dw8PD1KVCBcIJllaFg8PLS0lBxdScCUH
CA8HD1JpeGEeCAcmnAgXS1IeCCVoJggmUlIlDw8lSyUPCBcCTrQlFiUllgceDzRSNf4qJQgPNX8s
JggtpEMPBya6s0TYHg/9Sgc8nR0IDzSsHgdSAalLFgc8jiUPDyUmnB4HPLv+uCUAAAACAA8AFgEF
AtsAEQAmAAATFyMXBx8BMz8BJzMmNzYnIwcnFjcXBhcHFwYXByYHJzYnNzUnNictDwgIDxYtNDQP
DwgPFghaJTQIS2ElDw8PDw8PJZUXJQ8PDw8PDwJzu+dSLR4mJcnoHkMsJUMPDyVLhkOkYWEmDw8m
Q5VDWkNwUgAAAAACAA8AQwH7Ay0AEwAmAAABHwEGFwcXBisBJic1NxY3Fj8BCwEWFxY3JjcmJyYH
EwcnJjUmJwcBK6slDg4OBy2GnXceJWklFh4eD9EInNEtDxcPLWEWByU1HRc0JQMtByVElUrvrA+G
liUPDxYHNAFQ/jGrJgfn2ZwtFxZD/iomCAdhJgcWAAAAAgAHAFoCTgM8AB8AOwAAATMXBwYfARYP
AiIvAgYPASYHJxMmNzYXNhcGFzM3ATI3ND8BFxY3Fjc1JzU3JyYPAS8BNycmBxcDFgF1lR4PUggP
cAgef0oIOyYWCCU8OyYIFwglLYYIFy0tS/7zPAclLTw8JQdLaCwHJX8tDwgILWEPCAgPAlwdPEsP
O8IPLRYdaQcHaRYPDyUBzyWHNAgPNGF/NP4qHngdHiWkBxYtLLs0PCUelQgXcJU8D0tp/m0sAAAA
AAIABwAHAigC6gATACIAABM2FQMfATcXFgcjJwYvATYnNgM2EyU3NScHJxMvAQ8BEgcXaJ0HFi2s
JRZDcC32JSYPDx4WHo4BBSU0rCUHDx1pJQ8XJgLbDzT+Zi0eByV3NQ8PDyZDf1IBXzT9SgcmNC0I
JQGbJRcIPP5INDwAAAAAAv/yAGgDYgKRACgATQAAATYXNzYXFhcGFwcmByc3JjcnIwcXBy8BNy8B
Ig8BFgcGIwYnNzUmNzMTPwEnNzYPAR8BFj8BNicmJyYHIyciBwYnIg8BAhc/ASc3NhcGASOONIZ4
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
NxcWHwEWNxM2JyYPAiYnJgcGDwEGLwE0Jw8BEyxXOx0PByQzJTNJJQ4sLIoPLEksM3QeJCUkLIMz
UQ3EMxYkSSUHLCRuFlcIFiwPMywrFjsdLBYrLA8dQjMOQQG2Dw8WMxYHXw4PDx5JD0kdOtu3MwgH
X1AWBwcBrhb+ZwcWOkmECFeSJQc6AQcsHQ8s4gjUJToWDjTFDg6+JSsOHf7NAAACAAAABwK9AssA
IwBDAAATMx8CNzYXFg8CFxYHFCMnIicmDwEGJzUHJyY/ATYnJicmFxUWBwYHFRc2NzYfARYXFj8B
JyYnPwE1JyIHBicmJyZtkiU6HlC2Hg4sdAiSJRZYoBY7Fh0zMzNuHQcdiwcWbR0PZpkVFnUkSQ9Q
JSQsSSUrHg98DyxYHjMdUBYzM0kCyxZQB1cWKx46oUHbSQ8WB18PD1clDwcHFiwr6jMWg1FBQSzF
HVivJR0WHbcWHW4zDw8kLKhCX3wdFiWDBwiDUAAAAAL/+v9YAl4CFQAVACoAABMVEwcWFxY3JjcT
NicmDwEGLwEmJyYnFjcWHwE3NjcfAQcDBxcjJzcnAyZQkgcWHV8WDw+LBx0zHkkOJUEeOg9fblAz
Fh0eHSySHQ+SHQfxBxYkmg0BzCX+1dQzBxZCkl8BHCUWDjOvHRaoOgcPDw8PD2YID2YPDyQ7/txm
viyZbgEzQQACAA8AbQIOAnQAFAArAAATBhczFg8BHwEhNzYjJyY/Ai8BIScFNhcVBgcVNhcPAScF
JzU/ATYnIicmXw8smSUP6QckATMzBxWSJQ+gFhYk/vI7AVA6Hg9QZg8PHVH+oiQOhA4OdQgHAiMr
Dw8k6iQWHTMIDiWvOiUWHQcdO3wsUCQHOmcWDwckUR2DByUsXwADAAf/rgKnApEAFgAqADIAABMX
FhcTFAcmBy8BJgcGBwYvAQYnJjcTAxY/ATYfAhY3Fjc2AyYnJgcDBgEfARYHBic39tksD50leFke
Hi0WLRYmOwhZHgg1lWEeJSUXhiUXShcWPA+kJkNZHqQPAQUlLQ8taA8tApEIHkP92C0PDw8XSjUI
B2EtBwgePDSGAc/9dwgeWiwPFkNLDxZDBwIhJQ8PSv39DwHPD6xDCA81yQAAAAAEAA//8QJkAvIA
EQAkAC0ANQAAExEWFzYXFjc0JyMHIic1JyYHJxc2FwcXMhcWBwYHJwYnNxEmNxM2FxYHBic3NhcW
NzYnIwcGPA8sYUuzQ39SSh4IOxc0FnA7HgcWwi1hNDVv92gXDxYH91IWF0SGBwcWJi0dDywtJg8C
N/5IPA8tNBf3nEslHpU7CB48Dx48aCYtaP6GDw8ePH8BhDRZ/s8WQ6QtDzS7HtEXQ0QsOw8AAAAC
//oAUgJOAq8AFAArAAATJBcWByMnIgcGFzM3Nh8BBgcEJwIBFj8BNicHBic2JzY3MxY3FjcmJyIH
FtgBOiUXRHdaLBcPNR0tJn4tD6T+sCUVASkeWVoHNH9ZFw8PDzVKUhcPHQ6s2TQPAokm0S0tUkMP
WSUWBy2zFxbgAUH+GwgeUg80Qw81cFIlB2gPCDVZLdjRAAAAAAT/+gAAAkYC+QATACYALwA3AAAB
FjcXBxEWDwEnBwYnJj8BNhc3JwEGFzI3FzY1ES8BIgcXBiMmBwYXNhcWIwY1NicXNjc2Jw8BFgGM
O1ImDxYHJkqW/S0VQiw1sx4I/s4Pwi1adzQPLDwPBw4XYXc8f4YPDzSVDg5ZJRcPNC0mDwL5Dw8l
f/6aNFolBxYt4LpaLS0IJWn+fMI8JRYHfwGbNB48jh1DNSU8F0OzDzRLSpUHLRdDCDs8AAAAAgAH
AEsCmAM1ACYATAAAExYHIxYHBh8BNwU/ATYnBScmMzYhNi8BByYnJjMXNzUnIwcnBiciNxc3FjcX
FgcmBxUXNxcWIyYHBhczNgcUIycGJwYnNyY3JzcmNzQ8FgcIDxYIF0M0AQVpJRZ+/volBxYHARUl
DzSdJQgHWbMtLXc0aTRSJVJ3S2h/JQ80nWgtnCYPNeceDy3CsxYmybNDfwgPDxcPBx4XAoI8UmFL
WUMmDw8ILUoXCCVaHh4lJQcPNGEHHS0eDwgXCCwHDw8PJXAPHjwlHgcllRYtBy1Ds2EPDw8lO8J/
JjtSLWkdAAAAAgAW//ECXALbABkANAAAExc3FjcXFiMnBxUXNxcWIycPARcHJwYnEwMTPwEnNxc/
AS8BBycmMxY/ATYnIwYnBgcTIxdLaFJ4qyYONPY8NdglDzSkNA8HJY48Hg8HWWklByWkJhYPLdgm
DjR3lSYeWrMtaDwlBwcHAtQIDw8PJZ0PLSwmCCWWCCYljiUPHjwBVwEj/X4PNI4mCBclJR4HJZUP
Dx4lNB0WB0v+EywAAAAFAA//+QIwAtsAEQAaADMAOwA/AAATBRcGEwYHJgcmJzU2JyYnNTYTNic3
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
Fl9fizMsHQ+oWDoeDw8VmW4rDg4kdSQODlcsHQ4eMyQDIwclHTpJZlj+3DozJG4rDyTpASv9bw86
AfgWMwc6Ol+SXwcVHYMHJQIcByxQUCUHJTpuqA40Oh0PSTsABAANAEoCYwLcABoANAA9AEUAABMS
Bx8BFj8BJzczFwYXFjcvATU2NzYvAgcGNxY3FhcWBxcWDwEXBy8CJgcXByYHJxMDNxcWNxcWBwYn
JhcWNzYnJgcGNg0NFCEOKBsOIoAhBzZlIRtQUBsUPFBs1yFDG+uuDhRDBikOBxUivBsGSjYHIkND
IhUHIcNKKCINL40HDV4NLxU9DS8UAnD+iFAiFAcbL4AhIaEbL3l4Shs2PDY8KRQOFD0ODhVkgEov
ST0bXiENIoZrcn8iDg4iASABGiJ5DQ0ieAcNL3h4Di8NLw4oKQAAAAT/8gBfAsQDOQAFABwAJABJ
AAABBhc3LwIWNxYXFQcVFgcGBycmJyYnJjcnJjc2ExY3NScmBwYXNjc2NzYlBicmNxY3FxY3Nicm
JwYHBgcGFxYzFxYPAS8CBxQBFQ4zMwclUJk7tgckZiQ0p/FYJDQVFV4IXiss8VcIHVgPDlCLOjsk
LP6MSRYPMztXUFEkDzMlmXxfMyUHOkIl2zoWJKhJXxYCeyQWHSQPqA8PLGZJOx1B1IQOBxYdJUlf
HSQ76V/98wgWLA8HFiS3Fh0sUMU7HTpQFg4OQiQ6FjsrHggkJUmKHjoWHVgkB0IHFnwAAAL/6v9Y
AaACBgAcADUAABM2FTIVFxYPAQYHBhUWFzYXFgcGBwYnNycmMzc2Ez8BNicjBi8BJjc2MzYvASYj
BwYHFRcHF6BYFhZmBwgHXx0HJVceBx0eOtssCDsHFlAzUFg6CBYdMxYlBxYOSSUPZg8VJR1JLB4s
Af8HFlcWByVfHQceVxYHJDpmJSQIHXXiWFc7bf2UDkIWBysOJXxfHQ4sLFAOWB0lOr5XAAACABb/
6gIjAhwAFQAuAAATFzYXBxczNzUmNzYXBxcHIwcGJzYDExY3Nhc2NzYvASYHFwcGJzU2LwEPARYH
FjpnOh0HFiwdJTu2Bw4OJHxu4hYPFplfSR1RMwcHFjNBFgcdfAgeHjozHQ4WDwIGDh078CUzmXUI
DjPbviQlHZksATP+UQc6Bx0IXts0JAdX4yQPM4tJMyUPLKAsfAACAAcADwJeAvAADwAjAAAbARYX
FjcTNCcmBwMHJwMmBxY3FgcXFj8CHwELAQcnJjUDNEKDDyuEDpIdMx1QJR1QUVdmSTMHJA8sMyR8
JVhXLLczkgKC/fMlFg5QAeIsFg5Q/owPOgF1ZhYPDywOxjMe6SUIJP6T/vkkBxYzAkErAAL/+v++
Aw0BvQAeAEAAABMWNx8BFjc2FxY3FxY/ATYXBgcGJyMmJyMPAScGAyYTFj8BNDcXFh8BFjcTNicm
DwImJyYHBg8BBi8BNCcPARMsVzsdDwckMyUzSSUOLCyKDyxJLDN0HiQlJCyDM1ENxDMWJEklBywk
bhZXCBYsDzMsKxY7HSwWKywPHUIzDkEBtg8PFjMWB18ODw8eSQ9JHTrbtzMIB19QFgcHAa4W/mcH
FjpJhAhXkiUHOgEHLB0PLOII1CU6Fg40xQ4OviUrDh3+zQAAAgAAAAcCvQLLACMAQwAAEzMfAjc2
FxYPAhcWBxQjJyInJg8BBic1BycmPwE2JyYnJhcVFgcGBxUXNjc2HwEWFxY/AScmJz8BNSciBwYn
JicmbZIlOh5Qth4OLHQIkiUWWKAWOxYdMzMzbh0HHYsHFm0dD2aZFRZ1JEkPUCUkLEklKx4PfA8s
WB4zHVAWMzNJAssWUAdXFiseOqFB20kPFgdfDw9XJQ8HBxYsK+ozFoNRQUEsxR1YryUdFh23Fh1u
Mw8PJCyoQl98HRYlgwcIg1AAAAAC//r/WAJeAhUAFQAqAAATFRMHFhcWNyY3EzYnJg8BBi8BJicm
JxY3Fh8BNzY3HwEHAwcXIyc3JwMmUJIHFh1fFg8PiwcdMx5JDiVBHjoPX25QMxYdHh0skh0Pkh0H
8QcWJJoNAcwl/tXUMwcWQpJfARwlFg4zrx0WqDoHDw8PDw9mCA9mDw8kO/7cZr4smW4BM0EAAgAP
AG0CDgJ0ABQAKwAAEwYXMxYPAR8BITc2IycmPwIvASEnBTYXFQYHFTYXDwEnBSc1PwE2JyInJl8P
LJklD+kHJAEzMwcVkiUPoBYWJP7yOwFQOh4PUGYPDx1R/qIkDoQODnUIBwIjKw8PJOokFh0zCA4l
rzolFh0HHTt8LFAkBzpnFg8HJFEdgwclLF8AAwAH/+UCLQKYABoAPwBGAAATNhcHBicHFw8BBh8B
NxcPASUHJzcnJj8CJhM2FzY3NicHBicmPwEWNSYnJj8BMxcWMz8BJiciBwYPAQYXFgcTFjc2JyYH
6/8hBg5DIQYhShQHV5MiByH+85MiFAZKBjYOB17/LyEOBmRKeQcGGxRkBl4iDiE9FCEpKA0Uf3M8
IQ4oBzYNDaghDgYUIQ4CdyGhXiENKGsiBwYpDRQheSIOByGUFD08Qxvd/Z4NDRQbNQ4ODS8hNg4O
QxsHB40oDlANG3IUPBuoNQdDSl4BUAcUIhQGGgAAAAAAAQAACJwAAQFtBgAACAKOAAUADv+YAAUA
KAAsAAUAN//UAAUAOf+uAAUAO/+4AAUAPgAgAAUAQQAoAAUAQ/6wAAUARP5+AAUARf+XAAUARv6+
AAUAR/+2AAUASP6vAAUASf/XAAUASgA6AAUASwAYAAUATAAsAAUATQAiAAUATgAsAAUATwAnAAUA
UP6vAA4AIgA4AA4ALABPAA4AMP9yAA4AMv9wAA4AM/9nAA4ANf9cAA4APAAzAA4ASAA5AA4ATP97
AA4ATf9uAA4AT/92AB0ALAAqAB0AMP75AB0AMf/VAB0AMv78AB0AM/7XAB0ANf7GAB0AQP9wAB0A
TP9AAB0ATf81AB4AQP9/AB8AQP9AAB8AT/8jAB8AUP/JACAADv+5ACAAHf+bACAAJv+wACAAN/9H
ACAAOQAeACAAOv8nACAAQP90ACAAQwAgACAARf9SACAASf8fACAAUP9dACEAQP9CACIADv3+ACIA
Jv7CACIAOv/EACIAQP9AACIAQ//gACIASf/JACIATf+bACMAMP/qACMANf/PACUAMAAtACcAT/8V
ACcAUP80ACgAH/9mACgAMP7GACgAMv7dACgAM/62ACgAPP/aACgAPf+0ACgASP/kACgATP8YACgA
T/8WACgAUP/VACsAKAApACsAKwAZACsAM/+jACwADv3mACwAHf8rACwAJv7cACwAN/+3ACwAOf+p
ACwAOv+gACwAO/+1ACwAPf+wACwAQP9bACwAQ//rACwARf+TACwAR/+uACwASf+3ACwATP/hACwA
Tf/WACwAT//dACwAUP/fAC0AMv+wAC4AT/9KAC4AUP9CAC8AMP+4AC8AMv/JAC8ANf+nADAADv8O
ADAAHf7xADAAI/+IADAAJv7fADAAKP7RADAAKf9YADAAN/7wADAAOf8DADAAOv7vADAAO/8PADAA
Pf8IADAAQP9dADAAQ/8BADAARP8MADAARf7sADAARv8DADAAR/76ADAASP8AADAASf7yADAAS/8B
ADAATP7zADAATf7pADAATv77ADAAT/7vADAAUP7yADIADv8gADIAHf74ADIAJv8TADIAKAAfADIA
KQAWADIALf+tADIAN/+IADIAOf9/ADIAOv9sADIAO/+LADIAPP9aADIAPf+CADIAPwAaADIAQP9z
ADIAQ/+4ADIARf9pADIASf+EADIATf+wADIAT/+1ADIAUP+5ADMADv7jADMAHf7UADMAJv7jADMA
KAAfADMALf+UADMAN/9qADMAOv9HADMAO/9lADMARf9DADMASP+eADMASf9hADQAMP/EADQAMv/a
ADQASP9HADQAT/8/ADQAUP9OADUADv7VADUAHf7NADUAI/9lADUAJP9gADUAJv7fADUAKP9xADUA
Kf8eADUALf9bADUAN/8OADUAOf7hADUAOv7TADUAO/7tADUAPP/PADUAPf7oADUAQP9eADUAQ/84
ADUARP9DADUARf7LADUARv9FADUAR/7oADUASP9CADUASf7zADUAS/9PADUATP9NADUATf9CADUA
Tv9UADUAT/9IADUAUP9LADYAH/+jADYAIP++ADYAIf+7ADYAIv/GADYAI/+kADYAJP+2ADYAJf+8
ADYAJ//PADYAKP+gADYAKf+tADYAK/+bADYALf+bADYAMP/VADYAMv/lADYAM//gADYANP/iADYA
Nf/GADYANv/cADYAQP9aADYASP9AADYAT/97ADYAUP/gADcABQAoADcALP+yADcAMP6CADcAMv+d
ADcAM/+UADcANf9XADgABf/ZADgAIv8fADgAMP6KADgAMv8YADgAM/8GADgANf73ADgAQP9qADkA
BQAaADkAMP6jADkAM/8yADkANf8VADkAPP9kADkAPf8eADkASP87ADkAT//rADkAUP8NADoASAAz
ADsABQAzADsALP/IADsAMP6EADsAMv+gADsAM/+cADsANf9fADsAOwApADsASAAnADwADv86ADwA
PP+SAD0AMP5yAD0ASAAsAD4AKP7fAD4AUP/RAD8AMAA3AD8AMgApAEEABf/kAEEAIv8vAEEAI/8h
AEEAMP8dAEEAMv8CAEEAM/8TAEEANf8QAEEAOP9LAEEAOv97AEEAPP8FAEEAPf+mAEEAPv8oAEEA
P/9HAEEAQP9iAEEAQf9BAEEAQv81AEEAQwAYAEEARAAjAEEASAAWAEEASf/kAEEASgAcAEEASwAY
AEEATP8GAEEATf9EAEEAT/64AEEAUP6zAEMABf7iAEMAIv81AEMAMP6AAEMASAAqAEMAT//pAEMA
UAAaAEQABQAmAEQAMP58AEUABf/fAEUAIv+1AEUAJwAwAEUALP+cAEUAMP6ZAEUAMv+OAEUAM/97
AEUANf89AEUAQQAYAEUARQBfAEUASAAXAEYABf6MAEYADv7jAEYAMP6SAEYAN/9GAEYAOAArAEYA
OQAoAEYAOgAkAEYAPAAgAEYAPf9aAEYAQP9yAEYAQgAnAEYAQwAmAEYARf9fAEYASf8TAEYAT/+5
AEYAUP9TAEcABQAnAEcAMP6AAEcASAAuAEgABQAcAEgADv79AEgAMP5tAEgAMv/cAEgAM//cAEgA
Ov/JAEgAQ/8BAEgATf/pAEkABf7lAEkALP+JAEkAMP6LAEkANf8HAEoADgBJAEoAMP6CAEoAOv/d
AEoAQwAbAEoATf9fAEsABQAvAEsALAAYAEsAMP55AEsAMv/gAEsANf+XAEsASAAvAEwABf7iAEwA
HQBJAEwAJ/+1AEwAMP5yAEwAM//MAE0ABf7PAE0AHQBTAE0AMP5+AE4ABf/eAE4AMP6UAE8ABf6k
AE8AHQA6AE8AJ/+RAE8AMP5oAE8AMv+7AE8AM/+5AFAABQA2AFAAJf8lAFAAMP5XAFAAPf7wAFAA
SAAXAFAAT/8uAFAAUP84AAAABwBaAAMAAQQJAAAAfgAAAAMAAQQJAAEAEACMAAMAAQQJAAIADgB+
AAMAAQQJAAMAEACMAAMAAQQJAAQAEACMAAMAAQQJAAUACACcAAMAAQQJAAYAEACMAFYAaQBzAGkA
dAAgAGgAdAB0AHAAOgAvAC8AdAB0AGYALgBlAG8AcwBuAGUAdAAuAGMAbwBtAC8AIABmAG8AcgAg
AG0AbwByAGUAIABvAHIAaQBnAGkAbgBhAGwAIABmAG8AbgB0AHMAIABsAGkAawBlACAAdABoAGkA
cwAuAFIAZQBnAHUAbABhAHIAQwBvAG0AaQBjAGEAdABlAC4AdAB0AGYAAAACAAAAAAAA/5wAFAAA
AAAAAAAAAAAAAAAAAAAAAAAAAFIAAAECAQMAAwAEAAUABwAJAAoACwAMAA0ADwAQABEAEgATABQA
FQAWABcAGAAZABoAGwAcAB0AHgAiACQAJQAmACcAKAApACoAKwAsAC0ALgAvADAAMQAyADMANAA1
ADYANwA4ADkAOgA7ADwAPQBEAEUARgBHAEgASQBKAEsATABNAE4ATwBQAFEAUgBTAFQAVQBWAFcA
WABZAFoAWwBcAF0AhQUubnVsbBBub25tYXJraW5ncmV0dXJuAAAAAAEAAAAAAAABDgIzAABvzQMQ
AnVBaXJlb0Rpc3BsYXlCZCAg/////z////9BSVJCMDADAAEAAAA=
`
