package interceptador

import (
	"net/http"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/bd"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/erros"
	"github.com/rafaeljusto/atiradorfrequente/rest/config"
	"github.com/registrobr/gostk/db"
	"github.com/registrobr/gostk/log"
)

type sqler interface {
	Logger() log.Logger
	DefineTx(tx *bd.SQLogger)
	Tx() *bd.SQLogger
	Req() *http.Request
}

type BD struct {
	handler sqler
	tx      bd.Tx
}

func NovoBD(h sqler) *BD {
	return &BD{handler: h}
}

func (i *BD) Before() int {
	i.handler.Logger().Debug("Interceptador Antes: BD")

	if bd.Conexão == nil {
		// TODO(rafaeljusto): criptografar a senha
		err := bd.IniciarConexão(db.ConnParams{
			Username:           config.REST.BancoDados.Usuário,
			Password:           config.REST.BancoDados.Senha,
			DatabaseName:       config.REST.BancoDados.Nome,
			Host:               config.REST.BancoDados.Endereço,
			ConnectTimeout:     config.REST.BancoDados.TempoEsgotadoConexão,
			StatementTimeout:   config.REST.BancoDados.TempoEsgotadoComando,
			MaxIdleConnections: config.REST.BancoDados.MáximoNúmeroConexõesInativas,
			MaxOpenConnections: config.REST.BancoDados.MáximoNúmeroConexõesAbertas,
		}, config.REST.BancoDados.TempoEsgotadoTransação)

		if err != nil {
			i.handler.Logger().Critf("Erro ao conectar o banco de dados. Detalhes: %s", err)
			return http.StatusInternalServerError
		}
	}

	var err error
	if i.tx, err = bd.Conexão.Begin(); err != nil {
		i.handler.Logger().Errorf("Erro ao iniciar uma transação no banco de dados. Detalhes: %s", erros.Novo(err))
		return http.StatusInternalServerError
	}

	i.handler.DefineTx(bd.NovoSQLogger(i.tx))
	return 0
}

func (i *BD) After(status int) int {
	i.handler.Logger().Debug("Interceptador Depois: BD")

	if i.tx == nil {
		i.handler.Logger().Errorf("Transação não inicializada detectada")
		return status
	}

	if status >= 200 && status < 400 {
		if err := i.tx.Commit(); err != nil {
			i.handler.Logger().Errorf("Erro ao confirmar uma transação. Detalhes: %s", erros.Novo(err))
			return http.StatusInternalServerError
		}
	} else {
		if err := i.tx.Rollback(); err != nil {
			i.handler.Logger().Errorf("Erro ao desfazer uma transação. Detalhes: %s", erros.Novo(err))
		}
	}

	return status
}

type BDCompatível struct {
	sqlogger *bd.SQLogger
}

func (d *BDCompatível) DefineTx(tx *bd.SQLogger) {
	d.sqlogger = tx
}

func (d BDCompatível) Tx() *bd.SQLogger {
	return d.sqlogger
}
