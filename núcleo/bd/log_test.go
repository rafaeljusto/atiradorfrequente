package bd_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/erikstmartin/go-testdb"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/bd"
	"github.com/rafaeljusto/atiradorfrequente/testes/simulador"
)

func TestNovoSQLogger(t *testing.T) {
	sqlogger := bd.NovoSQLogger(simulador.Tx{
		SimulaExec: func(query string, args ...interface{}) (sql.Result, error) {
			return testdb.NewResult(1, nil, 1, nil), fmt.Errorf("erro na execução")
		},
	})

	if sqlogger == nil {
		t.Errorf("SQLogger não foi inicializado corretamente")
	}

	if _, err := sqlogger.Exec("SELECT"); err.Error() != "erro na execução" {
		t.Errorf("SQLogger não armazenou a transação esperada")
	}
}
