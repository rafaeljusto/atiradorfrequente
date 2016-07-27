// Package servidor inicializa o servidor REST e se conecta com os componentes
// necessários.
package servidor

import (
	"github.com/rafaeljusto/atiradorfrequente/núcleo/bd"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/erros"
	"github.com/rafaeljusto/atiradorfrequente/rest/config"
	"github.com/registrobr/gostk/db"
	"github.com/registrobr/gostk/log"
)

// Iniciar realiza todas as inicializações iniciais e sobe o servidor REST.
// Supõe que a configuração já foi carregada.
func Iniciar() error {
	if err := iniciarConexãoSyslog(); err != nil {
		return erros.Novo(err)
	}

	// o sistema não é interrompido caso ocorra um problema de conexão com o banco
	// de dados. Novas tentativas serão feitas a cada tratamento de requisição.
	if err := iniciarConexãoBancoDados(); err != nil {
		log.Critf("Erro ao conectar o banco de dados. Detalhes: %s", erros.Novo(err))
	}

	return erros.Novo(iniciarServidor())
}

func iniciarConexãoSyslog() error {
	log.Info("inicializando conexão com o servidor de log")

	return erros.Novo(log.Dial("tcp", config.Atual().EndereçoSyslog, "atirador-frequente"))
}

func iniciarConexãoBancoDados() error {
	log.Info("inicializando conexão com o banco de dados")

	err := bd.IniciarConexão(db.ConnParams{
		Username:           config.Atual().BancoDados.Usuário,
		Password:           config.Atual().BancoDados.Senha,
		DatabaseName:       config.Atual().BancoDados.Nome,
		Host:               config.Atual().BancoDados.Endereço,
		ConnectTimeout:     config.Atual().BancoDados.TempoEsgotadoConexão,
		StatementTimeout:   config.Atual().BancoDados.TempoEsgotadoComando,
		MaxIdleConnections: config.Atual().BancoDados.MáximoNúmeroConexõesInativas,
		MaxOpenConnections: config.Atual().BancoDados.MáximoNúmeroConexõesAbertas,
	}, config.Atual().BancoDados.TempoEsgotadoTransação)

	return erros.Novo(err)
}

func iniciarServidor() error {
	// TODO: utilizar bibliotecas para tornar a reinicialização do serviço
	// graciosa? Exemplos: https://github.com/fvbock/endless e
	// https://github.com/jpillora/overseer
	return nil
}
