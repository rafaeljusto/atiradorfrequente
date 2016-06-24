package config

import "time"

var REST configuração

type configuração struct {
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
