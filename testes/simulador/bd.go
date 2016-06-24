package simulador

import (
	"database/sql"
	"database/sql/driver"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/bd"
)

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

func (b BD) Begin() (bd.Tx, error) {
	return b.SimulaBegin()
}

func (b BD) Close() error {
	return b.SimulaClose()
}

func (b BD) Driver() driver.Driver {
	return b.SimulaDriver()
}

func (b BD) Exec(query string, args ...interface{}) (sql.Result, error) {
	return b.SimulaExec(query, args...)
}

func (b BD) Ping() error {
	return b.SimulaPing()
}

func (b BD) Prepare(query string) (*sql.Stmt, error) {
	return b.SimulaPrepare(query)
}

func (b BD) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return b.SimulaQuery(query, args...)
}

func (b BD) QueryRow(query string, args ...interface{}) *sql.Row {
	return b.SimulaQueryRow(query, args...)
}

func (b BD) SetMaxIdleConns(n int) {
	b.SimulaSetMaxIdleConns(n)
}

func (b BD) SetMaxOpenConns(n int) {
	b.SimulaSetMaxOpenConns(n)
}

type Tx struct {
	SimulaExec     func(query string, args ...interface{}) (sql.Result, error)
	SimulaQuery    func(query string, args ...interface{}) (*sql.Rows, error)
	SimulaQueryRow func(query string, args ...interface{}) *sql.Row
	SimulaPrepare  func(query string) (*sql.Stmt, error)
	SimulaRollback func() error
	SimulaCommit   func() error
}

func (t Tx) Exec(query string, args ...interface{}) (sql.Result, error) {
	return t.SimulaExec(query, args...)
}

func (t Tx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return t.SimulaQuery(query, args...)
}

func (t Tx) QueryRow(query string, args ...interface{}) *sql.Row {
	return t.SimulaQueryRow(query, args...)
}

func (t Tx) Prepare(query string) (*sql.Stmt, error) {
	return t.SimulaPrepare(query)
}

func (t Tx) Rollback() error {
	return t.SimulaRollback()
}

func (t Tx) Commit() error {
	return t.SimulaCommit()
}
