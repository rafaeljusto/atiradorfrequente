package config_test

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/config"
	"github.com/rafaeljusto/atiradorfrequente/testes"
	"github.com/registrobr/gostk/errors"
	"golang.org/x/image/font/gofont/goregular"
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

	arquivoFonte.Write(goregular.TTF)
	arquivoFonte.Close()

	arquivoImagemBase, err := ioutil.TempFile("", "teste-nucleo-config-")
	if err != nil {
		t.Fatalf("Erro gerar o arquivo da imagem base. Detalhes: %s", err)
	}

	imagemBaseExtraída, err := base64.StdEncoding.DecodeString(imagemBasePNG)
	if err != nil {
		t.Fatalf("Erro ao extrair a imagem base de teste. Detalhes: %s", err)
	}

	arquivoImagemBase.Write(imagemBaseExtraída)
	arquivoImagemBase.Close()

	cenários := []struct {
		descrição            string
		conteúdoArquivo      string
		deveConterFonte      bool
		deveConterImagemBase bool
		configuraçãoEsperada config.Configuração
		erroEsperado         error
	}{
		{
			descrição: "deve carregar a configuração corretamente",
			conteúdoArquivo: `
atirador:
  prazo confirmacao: 30m
  tempo maximo cadastro: 12h
  duracao maxima treino: 12h
  imagem numero controle:
    fonte: ` + arquivoFonte.Name() + `
    imagem base: ` + arquivoImagemBase.Name() + `
    chave codigo verificacao: abc123
`,
			deveConterFonte:      true,
			deveConterImagemBase: true,
			configuraçãoEsperada: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.PrazoConfirmação = 30 * time.Minute
				configuração.Atirador.TempoMáximoCadastro = 12 * time.Hour
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				configuração.Atirador.ImagemNúmeroControle.ChaveCódigoVerificação = "abc123"
				return configuração
			}(),
		},
		{
			descrição: "deve ignorar quando a imagem ou a fonte não são informados",
			conteúdoArquivo: `
atirador:
  prazo confirmacao: 30m
  tempo maximo cadastro: 12h
  duracao maxima treino: 12h
  imagem numero controle:
    fonte:
    imagem base:
    chave codigo verificacao: abc123
`,
			configuraçãoEsperada: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.PrazoConfirmação = 30 * time.Minute
				configuração.Atirador.TempoMáximoCadastro = 12 * time.Hour
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				configuração.Atirador.ImagemNúmeroControle.ChaveCódigoVerificação = "abc123"
				return configuração
			}(),
		},
		{
			descrição: "deve detectar quando o arquivo da fonte não existe",
			conteúdoArquivo: `
atirador:
  prazo confirmacao: 30m
  tempo maximo cadastro: 12h
  duracao maxima treino: 12h
  imagem numero controle:
    fonte: /tmp/eunaoexisto321.ttf
    imagem base: ` + arquivoImagemBase.Name() + `
    chave codigo verificacao: abc123
`,
			configuraçãoEsperada: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.PrazoConfirmação = 30 * time.Minute
				configuração.Atirador.TempoMáximoCadastro = 12 * time.Hour
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				configuração.Atirador.ImagemNúmeroControle.ChaveCódigoVerificação = "abc123"
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
  tempo maximo cadastro: 12h
  duracao maxima treino: 12h
  imagem numero controle:
    fonte: ` + arquivoQualquer.Name() + `
    imagem base: ` + arquivoImagemBase.Name() + `
    chave codigo verificacao: abc123
`,
			configuraçãoEsperada: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.PrazoConfirmação = 30 * time.Minute
				configuração.Atirador.TempoMáximoCadastro = 12 * time.Hour
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				configuração.Atirador.ImagemNúmeroControle.ChaveCódigoVerificação = "abc123"
				return configuração
			}(),
			erroEsperado: errors.Errorf("freetype: invalid TrueType format: TTF data is too short"),
		},
		{
			descrição: "deve detectar quando o arquivo de imagem base não existe",
			conteúdoArquivo: `
atirador:
  prazo confirmacao: 30m
  tempo maximo cadastro: 12h
  duracao maxima treino: 12h
  imagem numero controle:
    fonte: ` + arquivoFonte.Name() + `
    imagem base: /tmp/eunaoexisto321.png
    chave codigo verificacao: abc123
`,
			deveConterFonte: true,
			configuraçãoEsperada: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.PrazoConfirmação = 30 * time.Minute
				configuração.Atirador.TempoMáximoCadastro = 12 * time.Hour
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				configuração.Atirador.ImagemNúmeroControle.ChaveCódigoVerificação = "abc123"
				return configuração
			}(),
			erroEsperado: &os.PathError{
				Op:   "open",
				Path: "/tmp/eunaoexisto321.png",
				Err:  fmt.Errorf("no such file or directory"),
			},
		},
		{
			descrição: "deve detectar quando a imagem base esta em um formato inválido",
			conteúdoArquivo: `
atirador:
  prazo confirmacao: 30m
  tempo maximo cadastro: 12h
  duracao maxima treino: 12h
  imagem numero controle:
    fonte: ` + arquivoFonte.Name() + `
    imagem base: ` + arquivoQualquer.Name() + `
    chave codigo verificacao: abc123
`,
			deveConterFonte: true,
			configuraçãoEsperada: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.PrazoConfirmação = 30 * time.Minute
				configuração.Atirador.TempoMáximoCadastro = 12 * time.Hour
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				configuração.Atirador.ImagemNúmeroControle.ChaveCódigoVerificação = "abc123"
				return configuração
			}(),
			erroEsperado: errors.Errorf("image: unknown format"),
		},
	}

	for i, cenário := range cenários {
		var configuração config.Configuração
		err := yaml.Unmarshal([]byte(cenário.conteúdoArquivo), &configuração)

		if cenário.deveConterFonte && configuração.Atirador.ImagemNúmeroControle.Fonte.Font == nil {
			t.Errorf("Item %d, “%s”: fonte não foi carregada corretamente",
				i, cenário.descrição)
		}

		if cenário.deveConterImagemBase && configuração.Atirador.ImagemNúmeroControle.ImagemBase.Image == nil {
			t.Errorf("Item %d, “%s”: imagem base não foi carregada corretamente",
				i, cenário.descrição)
		}

		// a comparação de fonte e imagem consome muita memória, portanto nos
		// restringimos a uma verificação simples
		configuração.Atirador.ImagemNúmeroControle.Fonte.Font = nil
		configuração.Atirador.ImagemNúmeroControle.ImagemBase.Image = nil

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

	arquivoFonte.Write(goregular.TTF)
	arquivoFonte.Close()

	arquivoImagemBase, err := ioutil.TempFile("", "teste-nucleo-config-")
	if err != nil {
		t.Fatalf("Erro gerar o arquivo da imagem base. Detalhes: %s", err)
	}

	imagemBaseExtraída, err := base64.StdEncoding.DecodeString(imagemBasePNG)
	if err != nil {
		t.Fatalf("Erro ao extrair a imagem base de teste. Detalhes: %s", err)
	}

	arquivoImagemBase.Write(imagemBaseExtraída)
	arquivoImagemBase.Close()

	cenários := []struct {
		descrição            string
		variáveisAmbiente    map[string]string
		deveConterFonte      bool
		deveConterImagemBase bool
		configuraçãoEsperada config.Configuração
		erroEsperado         error
	}{
		{
			descrição: "deve carregar a configuração corretamente",
			variáveisAmbiente: map[string]string{
				"AF_ATIRADOR_PRAZO_CONFIRMACAO":                               "30m",
				"AF_ATIRADOR_TEMPO_MAXIMO_CADASTRO":                           "12h",
				"AF_ATIRADOR_DURACAO_MAXIMA_TREINO":                           "12h",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_FONTE":                    arquivoFonte.Name(),
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_IMAGEM_BASE":              arquivoImagemBase.Name(),
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_CHAVE_CODIGO_VERIFICACAO": "abc123",
			},
			deveConterFonte:      true,
			deveConterImagemBase: true,
			configuraçãoEsperada: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.PrazoConfirmação = 30 * time.Minute
				configuração.Atirador.TempoMáximoCadastro = 12 * time.Hour
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				configuração.Atirador.ImagemNúmeroControle.ChaveCódigoVerificação = "abc123"
				return configuração
			}(),
		},
		{
			descrição: "deve ignorar quando a imagem ou a fonte não são informados",
			variáveisAmbiente: map[string]string{
				"AF_ATIRADOR_PRAZO_CONFIRMACAO":                               "30m",
				"AF_ATIRADOR_TEMPO_MAXIMO_CADASTRO":                           "12h",
				"AF_ATIRADOR_DURACAO_MAXIMA_TREINO":                           "12h",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_CHAVE_CODIGO_VERIFICACAO": "abc123",
			},
			configuraçãoEsperada: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.PrazoConfirmação = 30 * time.Minute
				configuração.Atirador.TempoMáximoCadastro = 12 * time.Hour
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				configuração.Atirador.ImagemNúmeroControle.ChaveCódigoVerificação = "abc123"
				return configuração
			}(),
		},
		{
			descrição: "deve detectar quando o arquivo da fonte não existe",
			variáveisAmbiente: map[string]string{
				"AF_ATIRADOR_PRAZO_CONFIRMACAO":                               "30m",
				"AF_ATIRADOR_TEMPO_MAXIMO_CADASTRO":                           "12h",
				"AF_ATIRADOR_DURACAO_MAXIMA_TREINO":                           "12h",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_FONTE":                    "/tmp/eunaoexisto321.ttf",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_IMAGEM_BASE":              arquivoImagemBase.Name(),
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_CHAVE_CODIGO_VERIFICACAO": "abc123",
			},
			configuraçãoEsperada: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.PrazoConfirmação = 30 * time.Minute
				configuração.Atirador.TempoMáximoCadastro = 12 * time.Hour
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				configuração.Atirador.ImagemNúmeroControle.ChaveCódigoVerificação = "abc123"
				return configuração
			}(),
			erroEsperado: errors.Errorf("open /tmp/eunaoexisto321.ttf: no such file or directory"),
		},
		{
			descrição: "deve detectar quando o arquivo de fonte esta no formato inválido",
			variáveisAmbiente: map[string]string{
				"AF_ATIRADOR_PRAZO_CONFIRMACAO":                               "30m",
				"AF_ATIRADOR_TEMPO_MAXIMO_CADASTRO":                           "12h",
				"AF_ATIRADOR_DURACAO_MAXIMA_TREINO":                           "12h",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_FONTE":                    arquivoQualquer.Name(),
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_IMAGEM_BASE":              arquivoImagemBase.Name(),
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_CHAVE_CODIGO_VERIFICACAO": "abc123",
			},
			configuraçãoEsperada: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.PrazoConfirmação = 30 * time.Minute
				configuração.Atirador.TempoMáximoCadastro = 12 * time.Hour
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				configuração.Atirador.ImagemNúmeroControle.ChaveCódigoVerificação = "abc123"
				return configuração
			}(),
			erroEsperado: errors.Errorf("freetype: invalid TrueType format: TTF data is too short"),
		},
		{
			descrição: "deve detectar quando o arquivo de imagem base não existe",
			variáveisAmbiente: map[string]string{
				"AF_ATIRADOR_PRAZO_CONFIRMACAO":                               "30m",
				"AF_ATIRADOR_TEMPO_MAXIMO_CADASTRO":                           "12h",
				"AF_ATIRADOR_DURACAO_MAXIMA_TREINO":                           "12h",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_FONTE":                    arquivoFonte.Name(),
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_IMAGEM_BASE":              "/tmp/eunaoexisto321.png",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_CHAVE_CODIGO_VERIFICACAO": "abc123",
			},
			deveConterFonte: true,
			configuraçãoEsperada: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.PrazoConfirmação = 30 * time.Minute
				configuração.Atirador.TempoMáximoCadastro = 12 * time.Hour
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				configuração.Atirador.ImagemNúmeroControle.ChaveCódigoVerificação = "abc123"
				return configuração
			}(),
			erroEsperado: errors.Errorf("open /tmp/eunaoexisto321.png: no such file or directory"),
		},
		{
			descrição: "deve detectar quando a imagem base esta em um formato inválido",
			variáveisAmbiente: map[string]string{
				"AF_ATIRADOR_PRAZO_CONFIRMACAO":                               "30m",
				"AF_ATIRADOR_TEMPO_MAXIMO_CADASTRO":                           "12h",
				"AF_ATIRADOR_DURACAO_MAXIMA_TREINO":                           "12h",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_FONTE":                    arquivoFonte.Name(),
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_IMAGEM_BASE":              arquivoQualquer.Name(),
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_CHAVE_CODIGO_VERIFICACAO": "abc123",
			},
			deveConterFonte: true,
			configuraçãoEsperada: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.PrazoConfirmação = 30 * time.Minute
				configuração.Atirador.TempoMáximoCadastro = 12 * time.Hour
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				configuração.Atirador.ImagemNúmeroControle.ChaveCódigoVerificação = "abc123"
				return configuração
			}(),
			erroEsperado: errors.Errorf("image: unknown format"),
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

		if cenário.deveConterFonte && configuração.Atirador.ImagemNúmeroControle.Fonte.Font == nil {
			t.Errorf("Item %d, “%s”: fonte não foi carregada corretamente",
				i, cenário.descrição)
		}

		if cenário.deveConterImagemBase && configuração.Atirador.ImagemNúmeroControle.ImagemBase.Image == nil {
			t.Errorf("Item %d, “%s”: imagem base não foi carregada corretamente",
				i, cenário.descrição)
		}

		// a comparação de fonte e imagem consome muita memória, portanto nos
		// restringimos a uma verificação simples
		configuração.Atirador.ImagemNúmeroControle.Fonte.Font = nil
		configuração.Atirador.ImagemNúmeroControle.ImagemBase.Image = nil

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
	esperado.Atirador.TempoMáximoCadastro = 12 * time.Hour
	esperado.Atirador.DuraçãoMáximaTreino = 12 * time.Hour

	var c config.Configuração
	config.DefinirValoresPadrão(&c)

	// a comparação de fonte e imagem consome muita memória, portanto nos
	// restringimos a uma verificação simples
	c.Atirador.ImagemNúmeroControle.Fonte.Font = nil

	if !reflect.DeepEqual(c, esperado) {
		t.Errorf("Resultados não batem.\n%s", testes.Diff(esperado, c))
	}
}

const imagemBasePNG = `
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
