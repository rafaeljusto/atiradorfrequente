package bd_test

import (
	"database/sql/driver"
	"fmt"
	"testing"
	"time"

	"github.com/erikstmartin/go-testdb"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/bd"
	"github.com/rafaeljusto/atiradorfrequente/testes"
	"github.com/rafaeljusto/atiradorfrequente/testes/simulador"
	"github.com/registrobr/gostk/db"
	"github.com/registrobr/gostk/errors"
)

func TestBd_Begin(t *testing.T) {
	cenários := []struct {
		descrição       string
		simulaCriaçãoTx func() (driver.Tx, error)
		erroEsperado    error
	}{
		{
			descrição: "deve criar uma transação do banco de dados corretamente",
			simulaCriaçãoTx: func() (driver.Tx, error) {
				return &testdb.Tx{}, nil
			},
		},
		{
			descrição: "deve detectar um erro ao criar uma transação do banco de dados",
			simulaCriaçãoTx: func() (driver.Tx, error) {
				return nil, fmt.Errorf("erro ao criar")
			},
			erroEsperado: errors.Errorf("erro ao criar"),
		},
	}

	conexãoOriginal := bd.Conexão
	defer func() {
		bd.Conexão = conexãoOriginal
	}()

	driverOriginal := db.PostgresDriver
	defer func() {
		db.PostgresDriver = driverOriginal
	}()

	db.PostgresDriver = "testdb"

	for i, cenário := range cenários {
		testdb.SetBeginFunc(cenário.simulaCriaçãoTx)

		err := bd.IniciarConexão(db.ConnParams{
			Host:               "127.0.0.1",
			DatabaseName:       "teste",
			Username:           "usuario",
			Password:           "senha",
			ConnectTimeout:     2 * time.Second,
			StatementTimeout:   10 * time.Second,
			MaxIdleConnections: 16,
			MaxOpenConnections: 32,
		}, 1*time.Second)

		if err != nil {
			t.Fatalf("Item %d, “%s”: erro ao conectar a base de dados. Detalhes: %s",
				i, cenário.descrição, err)
		}

		_, err = bd.Conexão.Begin()

		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(nil, cenário.erroEsperado)
		if err = verificadorResultado.VerificaResultado(nil, err); err != nil {
			t.Error(err)
		}
	}
}

func TestIniciarConexão(t *testing.T) {
	cenários := []struct {
		descrição         string
		parâmetrosConexão db.ConnParams
		txTempoEsgotado   time.Duration
		driverBancoDados  string
		conexão           bd.BD
		conexãoSimulada   func(dsn string) (driver.Conn, error)
		erroEsperado      error
	}{
		{
			descrição: "deve conectar-se corretamente ao banco de dados",
			parâmetrosConexão: db.ConnParams{
				Host:               "127.0.0.1",
				DatabaseName:       "teste",
				Username:           "usuario",
				Password:           "senha",
				ConnectTimeout:     2 * time.Second,
				StatementTimeout:   10 * time.Second,
				MaxIdleConnections: 16,
				MaxOpenConnections: 32,
			},
			txTempoEsgotado:  1 * time.Second,
			driverBancoDados: "testdb",
			conexãoSimulada: func(dsn string) (driver.Conn, error) {
				return testdb.Conn(), nil
			},
		},
		{
			descrição: "deve detectar um driver inexistente do banco de dados",
			parâmetrosConexão: db.ConnParams{
				Host:               "127.0.0.1",
				DatabaseName:       "teste",
				Username:           "usuario",
				Password:           "senha",
				ConnectTimeout:     2 * time.Second,
				StatementTimeout:   10 * time.Second,
				MaxIdleConnections: 16,
				MaxOpenConnections: 32,
			},
			txTempoEsgotado:  1 * time.Second,
			driverBancoDados: "driver inválido",
			erroEsperado:     errors.Errorf(`sql: unknown driver "driver inválido" (forgotten import?)`),
		},
		{
			descrição: "deve detectar um erro ao conectar-se ao banco de dados",
			parâmetrosConexão: db.ConnParams{
				Host:               "127.0.0.1",
				DatabaseName:       "teste",
				Username:           "usuario",
				Password:           "senha",
				ConnectTimeout:     2 * time.Second,
				StatementTimeout:   10 * time.Second,
				MaxIdleConnections: 16,
				MaxOpenConnections: 32,
			},
			txTempoEsgotado:  1 * time.Second,
			driverBancoDados: "testdb",
			conexãoSimulada: func(dsn string) (driver.Conn, error) {
				return nil, fmt.Errorf("erro de conexão")
			},
			erroEsperado: errors.Errorf("erro de conexão"),
		},
		{
			descrição: "deve ignorar quando o banco de dados já esta conectado",
			parâmetrosConexão: db.ConnParams{
				Host:               "127.0.0.1",
				DatabaseName:       "teste",
				Username:           "usuario",
				Password:           "senha",
				ConnectTimeout:     2 * time.Second,
				StatementTimeout:   10 * time.Second,
				MaxIdleConnections: 16,
				MaxOpenConnections: 32,
			},
			txTempoEsgotado:  1 * time.Second,
			driverBancoDados: "testdb",
			conexão: simulador.BD{
				SimulaPing: func() error {
					return nil
				},
			},
		},
		{
			descrição: "deve conectar-se corretamente ao banco de dados quando a conexão atual não responde",
			parâmetrosConexão: db.ConnParams{
				Host:               "127.0.0.1",
				DatabaseName:       "teste",
				Username:           "usuario",
				Password:           "senha",
				ConnectTimeout:     2 * time.Second,
				StatementTimeout:   10 * time.Second,
				MaxIdleConnections: 16,
				MaxOpenConnections: 32,
			},
			txTempoEsgotado:  1 * time.Second,
			driverBancoDados: "testdb",
			conexão: simulador.BD{
				SimulaPing: func() error {
					return fmt.Errorf("erro de conexão")
				},
			},
			conexãoSimulada: func(dsn string) (driver.Conn, error) {
				return testdb.Conn(), nil
			},
		},
	}

	conexãoOriginal := bd.Conexão
	defer func() {
		bd.Conexão = conexãoOriginal
	}()

	driverOriginal := db.PostgresDriver
	defer func() {
		db.PostgresDriver = driverOriginal
	}()

	for i, cenário := range cenários {
		testdb.SetOpenFunc(cenário.conexãoSimulada)
		db.PostgresDriver = cenário.driverBancoDados
		bd.Conexão = cenário.conexão

		err := bd.IniciarConexão(cenário.parâmetrosConexão, cenário.txTempoEsgotado)

		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(nil, cenário.erroEsperado)
		if err = verificadorResultado.VerificaResultado(nil, err); err != nil {
			t.Error(err)
		}
	}
}
