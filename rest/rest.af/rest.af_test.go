package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/rafaeljusto/atiradorfrequente/rest/config"
	"github.com/rafaeljusto/atiradorfrequente/rest/servidor"
	"github.com/rafaeljusto/atiradorfrequente/testes"
)

func Test_main(t *testing.T) {
	cenários := []struct {
		descrição            string
		variáveisAmbiente    map[string]string
		arquivoConfiguração  string
		configuraçãoEsperada *config.Configuração
		saídaPadrãoEsperada  *regexp.Regexp
		saídaErroEsperada    *regexp.Regexp
	}{
		{
			descrição: "deve iniciar o servidor REST carregando o arquivo de configuração",
			arquivoConfiguração: `
binario:
  url: http://localhost:8080/binarios/rest.af
  tempo atualizacao: 1s
servidor:
  endereco: 0.0.0.0:0
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
  tempo maximo cadastro: 11h
  duracao maxima treino: 10h
  chave codigo verificacao: cba321
  imagem numero controle:
    url qrcode: https://exemplo.com.br/frequencia/%s/%s?verificacao=%s
`,
			configuraçãoEsperada: func() *config.Configuração {
				c := new(config.Configuração)
				c.Atirador.PrazoConfirmação = 10 * time.Minute
				c.Atirador.TempoMáximoCadastro = 11 * time.Hour
				c.Atirador.DuraçãoMáximaTreino = 10 * time.Hour
				c.Atirador.ChaveCódigoVerificação = "cba321"
				c.Atirador.ImagemNúmeroControle.URLQRCode = "https://exemplo.com.br/frequencia/%s/%s?verificacao=%s"
				c.Binário.URL = "http://localhost:8080/binarios/rest.af"
				c.Binário.TempoAtualização = 1 * time.Second
				c.Servidor.Endereço = "0.0.0.0:0"
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
			saídaPadrãoEsperada: regexp.MustCompile(`^$`),
			saídaErroEsperada:   regexp.MustCompile(`^$`),
		},
		{
			descrição: "deve detectar um problema no arquivo de configuração",
			arquivoConfiguração: `
- binario:
url: http://localhost:8080/binarios/rest.af
`,
			configuraçãoEsperada: func() *config.Configuração {
				c := new(config.Configuração)
				c.Atirador.PrazoConfirmação = 30 * time.Minute
				c.Atirador.TempoMáximoCadastro = 12 * time.Hour
				c.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				c.Atirador.ImagemNúmeroControle.URLQRCode = "http://localhost/frequencia/%s/%s?verificacao=%s"
				c.Binário.URL = "http://localhost:4000/binarios/rest.af"
				c.Binário.TempoAtualização = 5 * time.Second
				c.Servidor.Endereço = "0.0.0.0:443"
				c.Servidor.TempoEsgotadoLeitura = 5 * time.Second
				c.Syslog.Endereço = "127.0.0.1:514"
				c.Syslog.TempoEsgotadoConexão = 2 * time.Second
				c.BancoDados.Endereço = "127.0.0.1"
				c.BancoDados.Porta = 5432
				c.BancoDados.Nome = "atiradorfrequente"
				c.BancoDados.Usuário = "atiradorfrequente"
				c.BancoDados.TempoEsgotadoConexão = 3 * time.Second
				c.BancoDados.TempoEsgotadoComando = 10 * time.Second
				c.BancoDados.TempoEsgotadoTransação = 3 * time.Second
				c.BancoDados.MáximoNúmeroConexõesInativas = 16
				c.BancoDados.MáximoNúmeroConexõesAbertas = 32
				return c
			}(),
			saídaPadrãoEsperada: regexp.MustCompile(`^$`),
			saídaErroEsperada:   regexp.MustCompile(`^Erro ao carregar o arquivo de configuração. Detalhes: .*yaml: line 2: did not find expected '-' indicator$`),
		},
		{
			descrição: "deve iniciar o servidor REST carregando as configurações de variáveis de ambiente",
			variáveisAmbiente: map[string]string{
				"AF_BINARIO_URL":                                "http://localhost:8080/binarios/rest.af",
				"AF_BINARIO_TEMPO_ATUALIZACAO":                  "1s",
				"AF_SERVIDOR_ENDERECO":                          "0.0.0.0:0",
				"AF_SERVIDOR_TLS_HABILITADO":                    "true",
				"AF_SERVIDOR_TLS_ARQUIVO_CERTIFICADO":           "teste.crt",
				"AF_SERVIDOR_TLS_ARQUIVO_CHAVE":                 "teste.key",
				"AF_SERVIDOR_TEMPO_ESGOTADO_LEITURA":            "5s",
				"AF_SYSLOG_ENDERECO":                            "192.0.2.2:514",
				"AF_SYSLOG_TEMPO_ESGOTADO_CONEXAO":              "5s",
				"AF_BD_ENDERECO":                                "192.0.2.3",
				"AF_BD_PORTA":                                   "5432",
				"AF_BD_NOME":                                    "teste",
				"AF_BD_USUARIO":                                 "usuario_teste",
				"AF_BD_SENHA":                                   "abc123",
				"AF_BD_TEMPO_ESGOTADO_CONEXAO":                  "5s",
				"AF_BD_TEMPO_ESGOTADO_COMANDO":                  "20s",
				"AF_BD_TEMPO_ESGOTADO_TRANSACAO":                "5s",
				"AF_BD_MAXIMO_NUMERO_CONEXOES_INATIVAS":         "10",
				"AF_BD_MAXIMO_NUMERO_CONEXOES_ABERTAS":          "40",
				"AF_PROXIES":                                    "192.0.2.4,192.0.2.5,192.0.2.6",
				"AF_ATIRADOR_PRAZO_CONFIRMACAO":                 "10m",
				"AF_ATIRADOR_TEMPO_MAXIMO_CADASTRO":             "11h",
				"AF_ATIRADOR_DURACAO_MAXIMA_TREINO":             "10h",
				"AF_ATIRADOR_CHAVE_CODIGO_VERIFICACAO":          "cba321",
				"AF_ATIRADOR_IMAGEM_NUMERO_CONTROLE_URL_QRCODE": "https://exemplo.com.br/frequencia/%s/%s?verificacao=%s",
			},
			configuraçãoEsperada: func() *config.Configuração {
				c := new(config.Configuração)
				c.Atirador.PrazoConfirmação = 10 * time.Minute
				c.Atirador.TempoMáximoCadastro = 11 * time.Hour
				c.Atirador.DuraçãoMáximaTreino = 10 * time.Hour
				c.Atirador.ChaveCódigoVerificação = "cba321"
				c.Atirador.ImagemNúmeroControle.URLQRCode = "https://exemplo.com.br/frequencia/%s/%s?verificacao=%s"
				c.Binário.URL = "http://localhost:8080/binarios/rest.af"
				c.Binário.TempoAtualização = 1 * time.Second
				c.Servidor.Endereço = "0.0.0.0:0"
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
			saídaPadrãoEsperada: regexp.MustCompile(`^$`),
			saídaErroEsperada:   regexp.MustCompile(`^$`),
		},
		{
			descrição: "deve detectar um erro ao carregar as variáveis de ambiente",
			variáveisAmbiente: map[string]string{
				"AF_BD_PORTA": "XXXX",
			},
			configuraçãoEsperada: func() *config.Configuração {
				c := new(config.Configuração)
				c.Atirador.PrazoConfirmação = 30 * time.Minute
				c.Atirador.TempoMáximoCadastro = 12 * time.Hour
				c.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				c.Atirador.ImagemNúmeroControle.URLQRCode = "http://localhost/frequencia/%s/%s?verificacao=%s"
				c.Binário.URL = "http://localhost:4000/binarios/rest.af"
				c.Binário.TempoAtualização = 5 * time.Second
				c.Servidor.Endereço = "0.0.0.0:443"
				c.Servidor.TempoEsgotadoLeitura = 5 * time.Second
				c.Syslog.Endereço = "127.0.0.1:514"
				c.Syslog.TempoEsgotadoConexão = 2 * time.Second
				c.BancoDados.Endereço = "127.0.0.1"
				c.BancoDados.Porta = 5432
				c.BancoDados.Nome = "atiradorfrequente"
				c.BancoDados.Usuário = "atiradorfrequente"
				c.BancoDados.TempoEsgotadoConexão = 3 * time.Second
				c.BancoDados.TempoEsgotadoComando = 10 * time.Second
				c.BancoDados.TempoEsgotadoTransação = 3 * time.Second
				c.BancoDados.MáximoNúmeroConexõesInativas = 16
				c.BancoDados.MáximoNúmeroConexõesAbertas = 32
				return c
			}(),
			saídaPadrãoEsperada: regexp.MustCompile(`^$`),
			saídaErroEsperada:   regexp.MustCompile(`^Erro ao carregar as variáveis de ambiente. Detalhes: .*envconfig.Process: assigning AF_BD_PORTA to porta: converting 'XXXX' to type int\. details: strconv.ParseInt: parsing "XXXX": invalid syntax$`),
		},
		{
			descrição: "deve detectar um erro ao escutar em uma interface inválida",
			variáveisAmbiente: map[string]string{
				"AF_BINARIO_URL":               "http://localhost:8080/binarios/rest.af",
				"AF_BINARIO_TEMPO_ATUALIZACAO": "1s",
				"AF_SERVIDOR_ENDERECO":         "X.X.X.X:X",
			},
			configuraçãoEsperada: func() *config.Configuração {
				c := new(config.Configuração)
				c.Atirador.PrazoConfirmação = 30 * time.Minute
				c.Atirador.TempoMáximoCadastro = 12 * time.Hour
				c.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				c.Atirador.ImagemNúmeroControle.URLQRCode = "http://localhost/frequencia/%s/%s?verificacao=%s"
				c.Binário.URL = "http://localhost:8080/binarios/rest.af"
				c.Binário.TempoAtualização = 1 * time.Second
				c.Servidor.Endereço = "X.X.X.X:X"
				c.Servidor.TempoEsgotadoLeitura = 5 * time.Second
				c.Syslog.Endereço = "127.0.0.1:514"
				c.Syslog.TempoEsgotadoConexão = 2 * time.Second
				c.BancoDados.Endereço = "127.0.0.1"
				c.BancoDados.Porta = 5432
				c.BancoDados.Nome = "atiradorfrequente"
				c.BancoDados.Usuário = "atiradorfrequente"
				c.BancoDados.TempoEsgotadoConexão = 3 * time.Second
				c.BancoDados.TempoEsgotadoComando = 10 * time.Second
				c.BancoDados.TempoEsgotadoTransação = 3 * time.Second
				c.BancoDados.MáximoNúmeroConexõesInativas = 16
				c.BancoDados.MáximoNúmeroConexõesAbertas = 32
				return c
			}(),
			saídaPadrãoEsperada: regexp.MustCompile(`^$`),
			saídaErroEsperada:   regexp.MustCompile(`^Erro ao executar a aplicação\. Detalhes: .* Invalid address X\.X\.X\.X:X .*$`),
		},
	}

	teste = true
	defer func() {
		teste = false
	}()

	iniciarOriginal := servidor.Iniciar
	defer func() {
		servidor.Iniciar = iniciarOriginal
	}()

	servidor.Iniciar = func(escuta net.Listener) error {
		time.Sleep(100 * time.Millisecond)
		return nil
	}

	argumentosOriginais := os.Args
	defer func() {
		os.Args = argumentosOriginais
	}()

	for i, cenário := range cenários {
		os.Args = os.Args[:1]
		os.Clearenv()

		for chave, valor := range cenário.variáveisAmbiente {
			os.Setenv(chave, valor)
		}

		if cenário.arquivoConfiguração != "" {
			arquivoConfiguração, err := ioutil.TempFile("", "atirador-frequente-")
			if err != nil {
				t.Fatalf("Item %d, “%s”: erro ao criar o arquivo de configuração. Detalhes: %s",
					i, cenário.descrição, err)
			}

			arquivoConfiguração.WriteString(cenário.arquivoConfiguração)
			arquivoConfiguração.Close()

			os.Args = append(os.Args, []string{"-config", arquivoConfiguração.Name()}...)
		}

		// limpar a configuração a cada execução para tornar os cenários mais reais
		config.AtualizarConfiguração(&config.Configuração{})

		saídaPadrão, saídaErro := capturarSaídas(main)

		// a comparação de fonte e imagem consome muita memória, portanto nos
		// restringimos a uma verificação simples
		c := config.Atual()
		c.Atirador.ImagemNúmeroControle.Fonte.Font = nil
		c.Atirador.ImagemNúmeroControle.ImagemBase.Image = nil
		config.AtualizarConfiguração(c)

		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(cenário.configuraçãoEsperada, nil)
		if err := verificadorResultado.VerificaResultado(config.Atual(), nil); err != nil {
			t.Error(err)
		}

		if !cenário.saídaPadrãoEsperada.MatchString(saídaPadrão) {
			t.Errorf("Item %d, “%s”: saída padrão inesperada. Detalhes: %s",
				i, cenário.descrição, saídaPadrão)
		}

		if !cenário.saídaErroEsperada.MatchString(saídaErro) {
			t.Errorf("Item %d, “%s”: saída de erro inesperada. Detalhes: %s",
				i, cenário.descrição, saídaErro)
		}
	}
}

func capturarSaídas(f func()) (string, string) {
	saídaPadrãoOriginal := os.Stdout
	defer func() {
		os.Stdout = saídaPadrãoOriginal
	}()

	saídaErroOriginal := os.Stderr
	defer func() {
		os.Stderr = saídaErroOriginal
	}()

	leituraPadrão, escritaPadrão, _ := os.Pipe()
	os.Stdout = escritaPadrão

	canalLeituraPadrão := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, leituraPadrão)
		canalLeituraPadrão <- buf.String()
	}()

	leituraErro, escritaErro, _ := os.Pipe()
	os.Stderr = escritaErro

	canalLeituraErro := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, leituraErro)
		canalLeituraErro <- buf.String()
	}()

	f()

	escritaPadrão.Close()
	escritaErro.Close()

	return strings.TrimSpace(<-canalLeituraPadrão), strings.TrimSpace(<-canalLeituraErro)
}
