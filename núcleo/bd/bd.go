package bd

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/erros"
	"github.com/registrobr/gostk/db"
)

// Conexão armazena a conexão global do sistema com a base de dados. Todas as
// transações com o banco de dados relacional devem ser abertas a partir desta
// variável. Internamente possui um sistema de pool automático, o que garante a
// escalabilidade.
var Conexão BD

// BD representa uma conexão do banco de dados. Util para manipulação em
// cenários de testes de integração.
type BD interface {
	Begin() (Tx, error)
	Close() error
	Driver() driver.Driver
	Exec(query string, args ...interface{}) (sql.Result, error)
	Ping() error
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	SetMaxIdleConns(n int)
	SetMaxOpenConns(n int)
}

// Tx representa uma transação do banco de dados. Util para reutilizar
// transações em cenários de testes de integração.
type Tx interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Prepare(query string) (*sql.Stmt, error)
	Rollback() error
	Commit() error
}

type bd struct {
	*db.DB
}

// Begin sobrescreve o comportamento padrão da biblioteca adicionando a
// capacidade de timeout ao se criar uma nova transação.
func (b *bd) Begin() (Tx, error) {
	return b.DB.Begin()
}

// IniciarConexão conecta-se ao banco de dados com os parâmetros informados.
var IniciarConexão = func(parâmetrosConexão db.ConnParams, txTempoEsgotado time.Duration) error {
	// TODO(rafaeljusto): Adicionar semáforos para evitar concorrência no acesso a
	// variável global Conexão

	if Conexão != nil && Conexão.Ping() == nil {
		return nil
	}

	conexão, err := db.ConnectPostgres(parâmetrosConexão)
	if err != nil {
		return erros.Novo(err)
	}

	Conexão = &bd{
		DB: db.NewDB(conexão, txTempoEsgotado),
	}

	if err := Conexão.Ping(); err != nil {
		return erros.Novo(err)
	}

	return nil
}
