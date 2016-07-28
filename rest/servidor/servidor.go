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
// Supõe que a configuração já foi carregada. Está função fica bloqueada
// enquanto o servidor estiver executando.
func Iniciar() error {
	if err := iniciarConexãoSyslog(); err != nil {
		log.Critf("Erro ao conectar servidor de log. Detalhes: %s", erros.Novo(err))
		return erros.Novo(err)
	}
	defer func() {
		if err := log.Close(); err != nil {
			log.Errorf("Erro ao fechar a conexão do log. Detalhes: %s", erros.Novo(err))
		}
	}()

	// o sistema não é interrompido caso ocorra um problema de conexão com o banco
	// de dados. Novas tentativas serão feitas a cada tratamento de requisição.
	if err := iniciarConexãoBancoDados(); err != nil {
		log.Critf("Erro ao conectar o banco de dados. Detalhes: %s", erros.Novo(err))
	}
	defer func() {
		if err := bd.Conexão.Close(); err != nil {
			log.Errorf("Erro ao fechar a conexão do banco de dados. Detalhes: %s", erros.Novo(err))
		}
	}()

	if err := iniciarServidor(); err != nil {
		log.Critf("Erro ao iniciar o servidor. Detalhes: %s", erros.Novo(err))
		return erros.Novo(err)
	}

	return nil
}

func iniciarConexãoSyslog() error {
	log.Info("Inicializando conexão com o servidor de log")

	return erros.Novo(log.Dial("tcp", config.Atual().EndereçoSyslog, "atirador-frequente"))
}

func iniciarConexãoBancoDados() error {
	log.Info("Inicializando conexão com o banco de dados")

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
	log.Info("Inicializando servidor")

	// TODO: utilizar bibliotecas para tornar a reinicialização do serviço
	// menos agressiva? Alguns exemplos de bibliotecas:
	//   * github.com/fvbock/endless
	//   * github.com/jpillora/overseer
	//   * github.com/braintree/manners
	//   * github.com/tylerb/graceful
	//   * github.com/facebookgo/httpdown
	//   * github.com/facebookgo/grace
	//
	// Algumas questões a serem levadas em consideração:
	//   * Suporte a múltiplos listeners (?)
	//   * Suporte a arquivo de chave (TLS) com senha
	return nil
}
