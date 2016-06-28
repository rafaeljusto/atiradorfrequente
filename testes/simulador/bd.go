package simulador

import (
	"database/sql"
	"database/sql/driver"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/bd"
)

// BD estrutura de simulaão de uma conexão com o banco de dados.
type BD struct {
	SimulaBegin           func() (bd.Tx, error)
	SimulaClose           func() error
	SimulaDriver          func() driver.Driver
	SimulaExec            func(query string, args ...interface{}) (sql.Result, error)
	SimulaPing            func() error
	SimulaPrepare         func(query string) (*sql.Stmt, error)
	SimulaQuery           func(query string, args ...interface{}) (*sql.Rows, error)
	SimulaQueryRow        func(query string, args ...interface{}) *sql.Row
	SimulaSetMaxIdleConns func(n int)
	SimulaSetMaxOpenConns func(n int)
}

// Begin inicia uma nova transação.
func (b BD) Begin() (bd.Tx, error) {
	return b.SimulaBegin()
}

// Close encerra a conexão.
func (b BD) Close() error {
	return b.SimulaClose()
}

// Driver retorna a estrutura de baixo nível que implementa a conexão para um
// banco de dados específico.
func (b BD) Driver() driver.Driver {
	return b.SimulaDriver()
}

// Exec executa um comando SQL.
func (b BD) Exec(query string, args ...interface{}) (sql.Result, error) {
	return b.SimulaExec(query, args...)
}

// Ping testa a conexão com o banco de dados.
func (b BD) Ping() error {
	return b.SimulaPing()
}

// Prepare interpreta o comando SQL, substituíndo variáveis quando necessário,
// de maneira a evitar ataques de injeção de SQL.
func (b BD) Prepare(query string) (*sql.Stmt, error) {
	return b.SimulaPrepare(query)
}

// Query executa um comando SQL retornando múltiplus resultados.
func (b BD) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return b.SimulaQuery(query, args...)
}

// QueryRow executa um comando SQL retornando somente um resultado.
func (b BD) QueryRow(query string, args ...interface{}) *sql.Row {
	return b.SimulaQueryRow(query, args...)
}

// SetMaxIdleConns define a quantidade máxima de conexões inativas com o banco
// de dados.
func (b BD) SetMaxIdleConns(n int) {
	b.SimulaSetMaxIdleConns(n)
}

// SetMaxOpenConns define a quantidade máxima de conexões abertas com o banco de
// dados.
func (b BD) SetMaxOpenConns(n int) {
	b.SimulaSetMaxOpenConns(n)
}

// Tx estrutura de simulação de uma transação do banco de dados.
type Tx struct {
	SimulaExec     func(query string, args ...interface{}) (sql.Result, error)
	SimulaQuery    func(query string, args ...interface{}) (*sql.Rows, error)
	SimulaQueryRow func(query string, args ...interface{}) *sql.Row
	SimulaPrepare  func(query string) (*sql.Stmt, error)
	SimulaRollback func() error
	SimulaCommit   func() error
}

// Exec executa um comando SQL.
func (t Tx) Exec(query string, args ...interface{}) (sql.Result, error) {
	return t.SimulaExec(query, args...)
}

// Query executa um comando SQL retornando múltiplus resultados.
func (t Tx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return t.SimulaQuery(query, args...)
}

// QueryRow executa um comando SQL retornando somente um resultado.
func (t Tx) QueryRow(query string, args ...interface{}) *sql.Row {
	return t.SimulaQueryRow(query, args...)
}

// Prepare interpreta o comando SQL, substituíndo variáveis quando necessário,
// de maneira a evitar ataques de injeção de SQL.
func (t Tx) Prepare(query string) (*sql.Stmt, error) {
	return t.SimulaPrepare(query)
}

// Rollback desfaz qualquer alteração feita no banco de dados por esta
// transação.
func (t Tx) Rollback() error {
	return t.SimulaRollback()
}

// Commit confirma qualquer alteração feita no banco de dados por esta
// transação.
func (t Tx) Commit() error {
	return t.SimulaCommit()
}
