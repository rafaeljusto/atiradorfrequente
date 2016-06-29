package bd

import (
	"fmt"
	"strings"
	"time"

	"database/sql/driver"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/erros"
)

const (
	// AçãoLogCriação utilizado para identificar a ação de criação de um objeto na
	// base de dados.
	AçãoLogCriação AçãoLog = "CRIACAO"

	// AçãoLogAtualização utilizado para identificar a ação de atualização de um
	// objeto na base de dados.
	AçãoLogAtualização AçãoLog = "ATUALIZACAO"
)

// AçãoLog define a ação realizada sobre um objeto no banco de dados.
type AçãoLog string

// String converte a ação para o formato texto.
func (a AçãoLog) String() string {
	return string(a)
}

// Value converte a ação para um formato aceito no momento de persistência no
// banco de dados.
func (a AçãoLog) Value() (driver.Value, error) {
	return a.String(), nil
}

// Log armazena os dados para rastreamento de todas as modificações do usuário.
type Log struct {
	ID          int64
	DataCriação time.Time

	// TODO(rafaeljusto): Adicionar o IP? Neste caso todos os construtores de DAO
	// receberiam o IP. Futuramente poderíamos até adicionar um identificador da
	// sessão.
}

// SQLogger armazena além dos dados da transação do banco de dados, referências
// para rastrear todas as alterações do usuário nesta transação.
type SQLogger struct {
	sqler
	Log Log
}

// NovoSQLogger gera um novo SQLogger com os dados da transação.
func NovoSQLogger(s sqler) *SQLogger {
	return &SQLogger{
		sqler: s,
	}
}

// Gerar cria uma entrada na tabela log para identificar todas as operações
// feitas na mesma transação.
func (s *SQLogger) Gerar() error {
	if s.Log.ID > 0 {
		return nil
	}

	s.Log.DataCriação = time.Now().UTC()
	resultado, err := s.Exec(logCriaçãoComando, s.Log.DataCriação)

	if err != nil {
		return erros.Novo(err)
	}

	if s.Log.ID, err = resultado.LastInsertId(); err != nil {
		return erros.Novo(err)
	}

	return nil
}

var (
	logTabela = "log"

	logCriaçãoCampos = []string{
		"id",
		"data_criacao",
	}
	logCriaçãoCamposTexto = strings.Join(logCriaçãoCampos, ", ")
	logCriaçãoComando     = fmt.Sprintf(`INSERT INTO %s (%s) VALUES (DEFAULT, %s)`,
		logTabela, logCriaçãoCamposTexto, MarcadoresPSQL(len(logCriaçãoCampos)-1))
)
