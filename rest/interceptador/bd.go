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

// BD disponibiliza uma transação do banco de dados para o handler.
type BD struct {
	handler sqler
	tx      bd.Tx
}

// NovoBD cria um novo interceptador BD.
func NovoBD(h sqler) *BD {
	return &BD{handler: h}
}

// Before inicia uma conexão com o banco de dados caso não existe, e começa uma
// transação que ficará disponível para o handler. Ao iniciarmos a conexão com o
// banco de dados neste ponto garantimos que falhas temporárias de comunicação
// com o banco de dados serão resolvidas assim que o nível de rede voltar ao
// normal.
func (i *BD) Before() int {
	i.handler.Logger().Debug("Interceptador Antes: BD")

	if bd.Conexão == nil {
		// TODO(rafaeljusto): criptografar a senha
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

// After responsável por confirmar a transação (commit) ou desfazer as alterções
// (rollback). A confirmação somente é feita se o handler ou outros
// interceptadores retornarem um código HTTP de sucesso.
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

// BDCompatível implementa os métodos que serão utilizados pelo handler para
// acessar a transação criada por este interceptador.
type BDCompatível struct {
	sqlogger *bd.SQLogger
}

// DefineTx define a transação do banco de dados que será utilizada pelo
// handler.
func (d *BDCompatível) DefineTx(tx *bd.SQLogger) {
	d.sqlogger = tx
}

// Tx retorna a transação do banco de dados que será utilizada pelo handler.
func (d BDCompatível) Tx() *bd.SQLogger {
	return d.sqlogger
}
