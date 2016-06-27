package config

import (
	"sync/atomic"
	"time"
	"unsafe"
)

var configuração unsafe.Pointer

// Configuração estrutura que representa todas as possíveis configurações do
// servidor REST.
type Configuração struct {
	BancoDados struct {
		Endereço                     string
		Nome                         string
		Usuário                      string
		Senha                        string
		TempoEsgotadoConexão         time.Duration `yaml:"tempo esgotado conexao"`
		TempoEsgotadoComando         time.Duration `yaml:"tempo esgotado comando"`
		TempoEsgotadoTransação       time.Duration `yaml:"tempo esgotado transacao"`
		MáximoNúmeroConexõesInativas int           `yaml:"maximo numero conexoes inativas"`
		MáximoNúmeroConexõesAbertas  int           `yaml:"maximo numero conexoes abertas"`
	} `yaml:"banco de dados"`
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
