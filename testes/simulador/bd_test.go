package simulador_test

import (
	"database/sql"
	"database/sql/driver"
	"reflect"
	"strings"
	"testing"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/bd"
	"github.com/rafaeljusto/atiradorfrequente/testes/simulador"
)

func TestBD(t *testing.T) {
	var bdSimulado simulador.BD
	var métodosSimulados []string

	estruturaBDSimulado := reflect.TypeOf(bdSimulado)
	for i := 0; i < estruturaBDSimulado.NumField(); i++ {
		// trata somente funções como argumentos, ignorando atributos simples
		if !strings.HasPrefix(estruturaBDSimulado.Field(i).Type.String(), "func (") {
			continue
		}

		métodosSimulados = append(métodosSimulados, estruturaBDSimulado.Field(i).Name)
	}

	visitou := func(métodoSimulado string) {
		for i := len(métodosSimulados) - 1; i >= 0; i-- {
			if métodosSimulados[i] == métodoSimulado {
				métodosSimulados = append(métodosSimulados[:i], métodosSimulados[i+1:]...)
				break
			}
		}
	}

	bdSimulado.SimulaBegin = func() (bd.Tx, error) {
		visitou("SimulaBegin")
		return nil, nil
	}

	bdSimulado.SimulaClose = func() error {
		visitou("SimulaClose")
		return nil
	}

	bdSimulado.SimulaDriver = func() driver.Driver {
		visitou("SimulaDriver")
		return nil
	}

	bdSimulado.SimulaExec = func(query string, args ...interface{}) (sql.Result, error) {
		visitou("SimulaExec")
		return nil, nil
	}

	bdSimulado.SimulaPing = func() error {
		visitou("SimulaPing")
		return nil
	}

	bdSimulado.SimulaPrepare = func(query string) (*sql.Stmt, error) {
		visitou("SimulaPrepare")
		return nil, nil
	}

	bdSimulado.SimulaQuery = func(query string, args ...interface{}) (*sql.Rows, error) {
		visitou("SimulaQuery")
		return nil, nil
	}

	bdSimulado.SimulaQueryRow = func(query string, args ...interface{}) *sql.Row {
		visitou("SimulaQueryRow")
		return nil
	}

	bdSimulado.SimulaSetMaxIdleConns = func(n int) {
		visitou("SimulaSetMaxIdleConns")
	}

	bdSimulado.SimulaSetMaxOpenConns = func(n int) {
		visitou("SimulaSetMaxOpenConns")
	}

	bdSimulado.Begin()
	bdSimulado.Close()
	bdSimulado.Driver()
	bdSimulado.Exec("")
	bdSimulado.Ping()
	bdSimulado.Prepare("")
	bdSimulado.Query("")
	bdSimulado.QueryRow("")
	bdSimulado.SetMaxIdleConns(0)
	bdSimulado.SetMaxOpenConns(0)

	if len(métodosSimulados) > 0 {
		t.Errorf("métodos %#v não foram chamados", métodosSimulados)
	}
}

func TestTx(t *testing.T) {
	var txSimulado simulador.Tx
	var métodosSimulados []string

	estruturaBDSimulado := reflect.TypeOf(txSimulado)
	for i := 0; i < estruturaBDSimulado.NumField(); i++ {
		// trata somente funções como argumentos, ignorando atributos simples
		if !strings.HasPrefix(estruturaBDSimulado.Field(i).Type.String(), "func (") {
			continue
		}

		métodosSimulados = append(métodosSimulados, estruturaBDSimulado.Field(i).Name)
	}

	visitou := func(métodoSimulado string) {
		for i := len(métodosSimulados) - 1; i >= 0; i-- {
			if métodosSimulados[i] == métodoSimulado {
				métodosSimulados = append(métodosSimulados[:i], métodosSimulados[i+1:]...)
				break
			}
		}
	}

	txSimulado.SimulaExec = func(query string, args ...interface{}) (sql.Result, error) {
		visitou("SimulaExec")
		return nil, nil
	}

	txSimulado.SimulaQuery = func(query string, args ...interface{}) (*sql.Rows, error) {
		visitou("SimulaQuery")
		return nil, nil
	}

	txSimulado.SimulaQueryRow = func(query string, args ...interface{}) *sql.Row {
		visitou("SimulaQueryRow")
		return nil
	}

	txSimulado.SimulaPrepare = func(query string) (*sql.Stmt, error) {
		visitou("SimulaPrepare")
		return nil, nil
	}

	txSimulado.SimulaRollback = func() error {
		visitou("SimulaRollback")
		return nil
	}

	txSimulado.SimulaCommit = func() error {
		visitou("SimulaCommit")
		return nil
	}

	txSimulado.Exec("")
	txSimulado.Query("")
	txSimulado.QueryRow("")
	txSimulado.Prepare("")
	txSimulado.Rollback()
	txSimulado.Commit()

	if len(métodosSimulados) > 0 {
		t.Errorf("métodos %#v não foram chamados", métodosSimulados)
	}
}
