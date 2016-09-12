package config_test

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"net"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/rafaeljusto/atiradorfrequente/rest/config"
	"github.com/rafaeljusto/atiradorfrequente/testes"
	"gopkg.in/yaml.v2"
)

func TestConfiguração(t *testing.T) {
	config.AtualizarConfiguração(nil)

	var esperado config.Configuração
	esperado.BancoDados.Nome = "teste"

	config.AtualizarConfiguração(&esperado)
	obtido := config.Atual()

	if !reflect.DeepEqual(&esperado, obtido) {
		t.Errorf("configuração inesperada.\n%v", testes.Diff(esperado, obtido))
	}
}

func TestDefinirValoresPadrão(t *testing.T) {
	config.AtualizarConfiguração(nil)

	esperado := new(config.Configuração)
	esperado.Atirador.PrazoConfirmação = 30 * time.Minute
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
	esperado.Binário.URL = "http://localhost:4000/binarios/rest.af"
	esperado.Binário.TempoAtualização = 5 * time.Second
	esperado.Servidor.Endereço = "0.0.0.0:443"
	esperado.Servidor.TLS.Habilitado = false
	esperado.Servidor.TempoEsgotadoLeitura = 5 * time.Second
	esperado.Syslog.Endereço = "127.0.0.1:514"
	esperado.Syslog.TempoEsgotadoConexão = 2 * time.Second
	esperado.BancoDados.Endereço = "127.0.0.1"
	esperado.BancoDados.Porta = 5432
	esperado.BancoDados.Nome = "atiradorfrequente"
	esperado.BancoDados.Usuário = "atiradorfrequente"
	esperado.BancoDados.TempoEsgotadoConexão = 3 * time.Second
	esperado.BancoDados.TempoEsgotadoComando = 10 * time.Second
	esperado.BancoDados.TempoEsgotadoTransação = 3 * time.Second
	esperado.BancoDados.MáximoNúmeroConexõesInativas = 16
	esperado.BancoDados.MáximoNúmeroConexõesAbertas = 32

	config.DefinirValoresPadrão()

	if !reflect.DeepEqual(config.Atual(), esperado) {
		t.Errorf("Resultados não batem.\n%s", testes.Diff(esperado, config.Atual()))
	}
}

func TestCarregarDeArquivo(t *testing.T) {
	cenários := []struct {
		descrição            string
		nomeArquivo          string
		arquivo              string
		configuraçãoEsperada *config.Configuração
		erroEsperado         error
	}{
		{
			descrição: "deve carregar o arquivo de configuração corretamente",
			arquivo: `
binario:
  url: http://localhost:8080/binarios/rest.af
  tempo atualizacao: 1s
servidor:
  endereco: 192.0.2.1:443
  tls:
    habilitado: true
    arquivo certificado: teste.crt
    arquivo chave: teste.key
  tempo esgotado leitura: 5s
syslog:
  endereco: 192.0.2.2:514
  tempo esgotado conexao: 5s
banco de dados:
  endereco: 192.0.2.3
  porta: 5432
  nome: teste
  usuario: usuario_teste
  senha: abc123
  tempo esgotado conexao: 5s
  tempo esgotado comando: 20s
  tempo esgotado transacao: 5s
  maximo numero conexoes inativas: 10
  maximo numero conexoes abertas: 40
proxies:
  - 192.0.2.4
  - 192.0.2.5
  - 192.0.2.6
atirador:
  prazo confirmacao: 10m
  imagem numero controle:
    largura: 3508
    altura: 2480
    cor fundo: branco
    fonte:
      tamanho: 48
      dpi: 300
      cor: preto
    borda:
      largura: 50
      espacamento: 50
      cor: preto
    linha fundo:
      largura: 50
      espacamento: 50
      cor: cinza
`,
			configuraçãoEsperada: func() *config.Configuração {
				c := new(config.Configuração)
				c.Atirador.PrazoConfirmação = 10 * time.Minute
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
				c.Binário.URL = "http://localhost:8080/binarios/rest.af"
				c.Binário.TempoAtualização = 1 * time.Second
				c.Servidor.Endereço = "192.0.2.1:443"
				c.Servidor.TLS.Habilitado = true
				c.Servidor.TLS.ArquivoCertificado = "teste.crt"
				c.Servidor.TLS.ArquivoChave = "teste.key"
				c.Servidor.TempoEsgotadoLeitura = 5 * time.Second
				c.Syslog.Endereço = "192.0.2.2:514"
				c.Syslog.TempoEsgotadoConexão = 5 * time.Second
				c.BancoDados.Endereço = "192.0.2.3"
				c.BancoDados.Porta = 5432
				c.BancoDados.Nome = "teste"
				c.BancoDados.Usuário = "usuario_teste"
				c.BancoDados.Senha = "abc123"
				c.BancoDados.TempoEsgotadoConexão = 5 * time.Second
				c.BancoDados.TempoEsgotadoComando = 20 * time.Second
				c.BancoDados.TempoEsgotadoTransação = 5 * time.Second
				c.BancoDados.MáximoNúmeroConexõesInativas = 10
				c.BancoDados.MáximoNúmeroConexõesAbertas = 40
				c.Proxies = []net.IP{
					net.ParseIP("192.0.2.4"),
					net.ParseIP("192.0.2.5"),
					net.ParseIP("192.0.2.6"),
				}
				return c
			}(),
		},
		{
			descrição:   "deve detectar um erro ao tentar abrir um arquivo inexistente",
			nomeArquivo: "/tmp/atirador-frequente-eu-nao-existo.yaml",
			erroEsperado: &os.PathError{
				Op:   "open",
				Path: "/tmp/atirador-frequente-eu-nao-existo.yaml",
				Err:  fmt.Errorf("no such file or directory"),
			},
		},
		{
			descrição: "deve detectar um formato inválido do arquivo de configuração",
			arquivo: `
servidor:
- endereco: 192.0.2.1:443
`,
			erroEsperado: &yaml.TypeError{
				Errors: []string{
					`line 3: cannot unmarshal !!seq into struct { Endereço string "yaml:\"endereco\" envconfig:\"endereco\""; TLS struct { Habilitado bool "yaml:\"habilitado\" envconfig:\"habilitado\""; ArquivoCertificado string "yaml:\"arquivo certificado\" envconfig:\"arquivo_certificado\""; ArquivoChave string "yaml:\"arquivo chave\" envconfig:\"arquivo_chave\"" } "yaml:\"tls\" envconfig:\"tls\""; TempoEsgotadoLeitura time.Duration "yaml:\"tempo esgotado leitura\" envconfig:\"tempo_esgotado_leitura\"" }`,
				},
			},
		},
	}

	for i, cenário := range cenários {
		config.AtualizarConfiguração(nil)

		arquivo, err := ioutil.TempFile("", "atirador-frequente-")
		if err != nil {
			t.Errorf("Item %d, “%s”: erro ao criar arquivo. Detalhes: %s",
				i, cenário.descrição, err)
		}

		arquivo.WriteString(cenário.arquivo)
		arquivo.Close()

		if cenário.nomeArquivo == "" {
			cenário.nomeArquivo = arquivo.Name()
		}

		err = config.CarregarDeArquivo(cenário.nomeArquivo)

		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(cenário.configuraçãoEsperada, cenário.erroEsperado)
		if err = verificadorResultado.VerificaResultado(config.Atual(), err); err != nil {
			t.Error(err)
		}
	}
}

func TestCarregarDeVariávelAmbiente(t *testing.T) {
	cenários := []struct {
		descrição            string
		variáveisAmbiente    map[string]string
		configuraçãoEsperada *config.Configuração
		erroEsperado         error
	}{
		{
			descrição: "deve carregar corretamente das variáveis de ambiente",
			variáveisAmbiente: map[string]string{
				"AF_BINARIO_URL":                                             "http://localhost:8080/binarios/rest.af",
				"AF_BINARIO_TEMPO_ATUALIZACAO":                               "1s",
				"AF_SERVIDOR_ENDERECO":                                       "192.0.2.1:443",
				"AF_SERVIDOR_TLS_HABILITADO":                                 "true",
				"AF_SERVIDOR_TLS_ARQUIVO_CERTIFICADO":                        "teste.crt",
				"AF_SERVIDOR_TLS_ARQUIVO_CHAVE":                              "teste.key",
				"AF_SERVIDOR_TEMPO_ESGOTADO_LEITURA":                         "5s",
				"AF_SYSLOG_ENDERECO":                                         "192.0.2.2:514",
				"AF_SYSLOG_TEMPO_ESGOTADO_CONEXAO":                           "5s",
				"AF_BD_ENDERECO":                                             "192.0.2.3",
				"AF_BD_PORTA":                                                "5432",
				"AF_BD_NOME":                                                 "teste",
				"AF_BD_USUARIO":                                              "usuario_teste",
				"AF_BD_SENHA":                                                "abc123",
				"AF_BD_TEMPO_ESGOTADO_CONEXAO":                               "5s",
				"AF_BD_TEMPO_ESGOTADO_COMANDO":                               "20s",
				"AF_BD_TEMPO_ESGOTADO_TRANSACAO":                             "5s",
				"AF_BD_MAXIMO_NUMERO_CONEXOES_INATIVAS":                      "10",
				"AF_BD_MAXIMO_NUMERO_CONEXOES_ABERTAS":                       "40",
				"AF_PROXIES":                                                 "192.0.2.4,192.0.2.5,192.0.2.6",
				"AF_ATIRADOR_PRAZO_CONFIRMACAO":                              "10m",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LARGURA":                 "3508",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_ALTURA":                  "2480",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_COR_FUNDO":               "branco",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_BORDA_LARGURA":           "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_BORDA_ESPACAMENTO":       "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_BORDA_COR":               "preto",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LINHA_FUNDO_LARGURA":     "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LINHA_FUNDO_ESPACAMENTO": "50",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_LINHA_FUNDO_COR":         "cinza",
			},
			configuraçãoEsperada: func() *config.Configuração {
				c := new(config.Configuração)
				c.Atirador.PrazoConfirmação = 10 * time.Minute
				c.Atirador.ImagemNúmeroControle.Largura = 3508
				c.Atirador.ImagemNúmeroControle.Altura = 2480
				c.Atirador.ImagemNúmeroControle.CorFundo.Color = color.RGBA{0xff, 0xff, 0xff, 0xff}
				c.Atirador.ImagemNúmeroControle.Borda.Largura = 50
				c.Atirador.ImagemNúmeroControle.Borda.Espaçamento = 50
				c.Atirador.ImagemNúmeroControle.Borda.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				c.Atirador.ImagemNúmeroControle.LinhaFundo.Largura = 50
				c.Atirador.ImagemNúmeroControle.LinhaFundo.Espaçamento = 50
				c.Atirador.ImagemNúmeroControle.LinhaFundo.Cor.Color = color.RGBA{0xee, 0xee, 0xee, 0xff}
				c.Binário.URL = "http://localhost:8080/binarios/rest.af"
				c.Binário.TempoAtualização = 1 * time.Second
				c.Servidor.Endereço = "192.0.2.1:443"
				c.Servidor.TLS.Habilitado = true
				c.Servidor.TLS.ArquivoCertificado = "teste.crt"
				c.Servidor.TLS.ArquivoChave = "teste.key"
				c.Servidor.TempoEsgotadoLeitura = 5 * time.Second
				c.Syslog.Endereço = "192.0.2.2:514"
				c.Syslog.TempoEsgotadoConexão = 5 * time.Second
				c.BancoDados.Endereço = "192.0.2.3"
				c.BancoDados.Porta = 5432
				c.BancoDados.Nome = "teste"
				c.BancoDados.Usuário = "usuario_teste"
				c.BancoDados.Senha = "abc123"
				c.BancoDados.TempoEsgotadoConexão = 5 * time.Second
				c.BancoDados.TempoEsgotadoComando = 20 * time.Second
				c.BancoDados.TempoEsgotadoTransação = 5 * time.Second
				c.BancoDados.MáximoNúmeroConexõesInativas = 10
				c.BancoDados.MáximoNúmeroConexõesAbertas = 40
				c.Proxies = []net.IP{
					net.ParseIP("192.0.2.4"),
					net.ParseIP("192.0.2.5"),
					net.ParseIP("192.0.2.6"),
				}
				return c
			}(),
		},
		{
			descrição: "deve detectar um problema nas variáveis de ambiente",
			variáveisAmbiente: map[string]string{
				"AF_BD_PORTA": "XXX",
			},
			erroEsperado: &envconfig.ParseError{
				KeyName:   "AF_BD_PORTA",
				FieldName: "porta",
				TypeName:  "int",
				Value:     "XXX",
			},
		},
	}

	for i, cenário := range cenários {
		config.AtualizarConfiguração(nil)

		os.Clearenv()
		for chave, valor := range cenário.variáveisAmbiente {
			os.Setenv(chave, valor)
		}

		err := config.CarregarDeVariávelAmbiente()

		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(cenário.configuraçãoEsperada, cenário.erroEsperado)
		if err = verificadorResultado.VerificaResultado(config.Atual(), err); err != nil {
			t.Error(err)
		}
	}
}
