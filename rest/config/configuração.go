// Package config armazena as configurações relacionadas ao servidor REST. Para
// facilitar a usabilidade uma variável global estará disponível para acessar os
// dados de configuração.
package config

import (
	"io/ioutil"
	"net"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/kelseyhightower/envconfig"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/config"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/erros"
	"gopkg.in/yaml.v2"
)

// prefixo define o prefixo utilizado nas variáveis de ambiente para definir
// valores no arquivo de configuração. Como exemplo, se o prefixo for "AF",
// variáveis de ambiente com o prefixo "AF_" serão analisadas.
const prefixo = "AF"

var configuração unsafe.Pointer

// Configuração estrutura que representa todas as possíveis configurações do
// relacionadas ao sistema REST.
type Configuração struct {
	config.Configuração `yaml:",inline"`

	// Binário define as informações necessárias para se obter um novo binário e
	// troca-lo com o atual sem que o sistema pare de funcionar.
	Binário struct {
		// URL endereço para se obter o novo binário.
		URL string `yaml:"url" envconfig:"url"`

		// TempoAtualização intervalo de tempo que o sistema irá verificar se existe
		// um binário novo.
		TempoAtualização time.Duration `yaml:"tempo atualizacao" envconfig:"tempo_atualizacao"`
	} `yaml:"binario" envconfig:"binario"`

	Servidor struct {
		// Endereço interface que o servidor irá escutar, formado pelo IP com porta
		// (exemplo 127.0.0.1:443).
		Endereço string `yaml:"endereco" envconfig:"endereco"`

		TLS struct {
			// Habilitado define se o TLS será utilizado ou não.
			Habilitado bool `yaml:"habilitado" envconfig:"habilitado"`

			// ArquivoCertificado caminho para o arquivo em formato PEM que contém o
			// certificado.
			ArquivoCertificado string `yaml:"arquivo certificado" envconfig:"arquivo_certificado"`

			// ArquivoCertificado caminho para o arquivo em formato PEM que contém a
			// chave privada referente ao certificado.
			ArquivoChave string `yaml:"arquivo chave" envconfig:"arquivo_chave"`
		} `yaml:"tls" envconfig:"tls"`

		// TempoEsgotadoLeitura define o tempo em que o servidor irá aguardar após
		// um cliente se conectar para que alguma requisição seja recebida.
		TempoEsgotadoLeitura time.Duration `yaml:"tempo esgotado leitura" envconfig:"tempo_esgotado_leitura"`
	} `yaml:"servidor" envconfig:"servidor"`

	Syslog struct {
		// EndereçoSyslog endereço IP com porta (exemplo  127.0.0.1:514) para
		// conexão com o servidor de log central. A conexão será feita utilizando o
		// protocolo TCP para minimizar a perda de mensagens.
		Endereço             string        `yaml:"endereco" envconfig:"endereco"`
		TempoEsgotadoConexão time.Duration `yaml:"tempo esgotado conexao" envconfig:"tempo_esgotado_conexao"`
	} `yaml:"syslog" envconfig:"syslog"`

	BancoDados struct {
		Endereço                     string        `yaml:"endereco" envconfig:"endereco"`
		Porta                        int           `yaml:"porta" envconfig:"porta"`
		Nome                         string        `yaml:"nome" envconfig:"nome"`
		Usuário                      string        `yaml:"usuario" envconfig:"usuario"`
		Senha                        string        `yaml:"senha" envconfig:"senha"`
		TempoEsgotadoConexão         time.Duration `yaml:"tempo esgotado conexao" envconfig:"tempo_esgotado_Conexao"`
		TempoEsgotadoComando         time.Duration `yaml:"tempo esgotado comando" envconfig:"tempo_esgotado_comando"`
		TempoEsgotadoTransação       time.Duration `yaml:"tempo esgotado transacao" envconfig:"tempo_esgotado_transacao"`
		MáximoNúmeroConexõesInativas int           `yaml:"maximo numero conexoes inativas" envconfig:"maximo_numero_conexoes_inativas"`
		MáximoNúmeroConexõesAbertas  int           `yaml:"maximo numero conexoes abertas" envconfig:"maximo_numero_conexoes_abertas"`
	} `yaml:"banco de dados" envconfig:"bd"`

	// Proxies define a lista de endereços IPs que podem informar os cabeçalhos
	// HTTP X-Forwarded-For ou X-Real-IP para identificar os clientes finais.
	Proxies []net.IP `yaml:"proxies" envconfig:"proxies"`
}

// Atual retorna a configuração atual do sistema, armazenada internamente em uma
// variável global.
func Atual() *Configuração {
	return (*Configuração)(atomic.LoadPointer(&configuração))
}

// AtualizarConfiguração modifica a atual configuração do sistema de maneira
// segura. Muito útil para cenários de testes que precisam simular uma
// configuração específica.
func AtualizarConfiguração(c *Configuração) {
	atomic.StorePointer(&configuração, unsafe.Pointer(c))
}

// DefinirValoresPadrão utiliza valores padrão em todos os campos da
// configuração caso o usuário não informe. O usuário também tem a opção de
// sobrescrever somente alguns valores, mantendo os demais com valores padrão.
func DefinirValoresPadrão() {
	c := Atual()
	if c == nil {
		c = new(Configuração)
	}

	config.DefinirValoresPadrão(&c.Configuração)
	c.Binário.URL = "http://localhost:4000/binarios/rest.af"
	c.Binário.TempoAtualização = 5 * time.Second
	c.Servidor.Endereço = "0.0.0.0:443"
	c.Servidor.TLS.Habilitado = false
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

	AtualizarConfiguração(c)
}

// CarregarDeArquivo interpreta um arquivo de configuração em formato YAML e
// preenche as variáveis da configuração global do servidor REST. O que não for
// informado no arquivo de configuração YAML não será sobrescrito.
func CarregarDeArquivo(arquivo string) error {
	conteúdo, err := ioutil.ReadFile(arquivo)
	if err != nil {
		return erros.Novo(err)
	}

	c := Atual()
	if c == nil {
		c = new(Configuração)
	}

	if err = yaml.Unmarshal(conteúdo, c); err != nil {
		return erros.Novo(err)
	}

	AtualizarConfiguração(c)
	return nil
}

// CarregarDeVariávelAmbiente analisa as variáveis de ambiente e preenche as
// variáveis da configuração global do servidor REST. O que não for informado
// nas variáveis de ambiente não será sobrescrito.
func CarregarDeVariávelAmbiente() error {
	c := Atual()
	if c == nil {
		c = new(Configuração)
	}

	if err := envconfig.Process(prefixo, c); err != nil {
		return erros.Novo(err)
	}

	AtualizarConfiguração(c)
	return nil
}
