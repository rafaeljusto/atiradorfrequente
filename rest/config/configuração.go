// Package config armazena as configurações relacionadas ao servidor REST. Para
// facilitar a usabilidade uma variável global estará disponível para acessar os
// dados de configuração.
package config

import (
	"net"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/config"
)

var configuração unsafe.Pointer

// Configuração estrutura que representa todas as possíveis configurações do
// servidor REST.
type Configuração struct {
	config.Configuração

	BancoDados struct {
		Endereço                     string        `yaml:"endereco"`
		Nome                         string        `yaml:"nome"`
		Usuário                      string        `yaml:"usuario"`
		Senha                        string        `yaml:"senha"`
		TempoEsgotadoConexão         time.Duration `yaml:"tempo esgotado conexao"`
		TempoEsgotadoComando         time.Duration `yaml:"tempo esgotado comando"`
		TempoEsgotadoTransação       time.Duration `yaml:"tempo esgotado transacao"`
		MáximoNúmeroConexõesInativas int           `yaml:"maximo numero conexoes inativas"`
		MáximoNúmeroConexõesAbertas  int           `yaml:"maximo numero conexoes abertas"`
	} `yaml:"banco de dados"`

	// Proxies define a lista de endereços IPs que podem informar os cabeçalhos
	// HTTP X-Forwarded-For ou X-Real-IP para identificar os clientes finais.
	Proxies []net.IP `yaml:"proxies"`
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
