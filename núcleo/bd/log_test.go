package bd_test

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/erikstmartin/go-testdb"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/bd"
	"github.com/rafaeljusto/atiradorfrequente/testes"
	"github.com/rafaeljusto/atiradorfrequente/testes/simulador"
	"github.com/registrobr/gostk/errors"
)

func TestAçãoLog_String(t *testing.T) {
	cenários := []struct {
		ação     bd.AçãoLog
		esperado string
	}{
		{ação: bd.AçãoLogCriação, esperado: string(bd.AçãoLogCriação)},
		{ação: bd.AçãoLogAtualização, esperado: string(bd.AçãoLogAtualização)},
	}

	for _, cenário := range cenários {
		if texto := cenário.ação.String(); texto != cenário.esperado {
			t.Errorf("Ação não se converteu corretamente para o formato texto. Esperado “%s”, mas obteve “%s”",
				cenário.esperado, texto)
		}
	}
}

func TestAçãoLog_Value(t *testing.T) {
	cenários := []struct {
		ação          bd.AçãoLog
		valorEsperado driver.Value
		erroEsperado  error
	}{
		{ação: bd.AçãoLogCriação, valorEsperado: string(bd.AçãoLogCriação)},
		{ação: bd.AçãoLogAtualização, valorEsperado: string(bd.AçãoLogAtualização)},
	}

	for _, cenário := range cenários {
		valor, err := cenário.ação.Value()

		if !errors.Equal(cenário.erroEsperado, err) {
			t.Errorf("Erros não batem. Esperado “%v”; encontrado “%v”",
				cenário.erroEsperado, err)
		}

		if valor != cenário.valorEsperado {
			t.Errorf("Ação não se converteu corretamente para o formato texto. Esperado “%s”, mas obteve “%s”",
				cenário.valorEsperado, valor)
		}
	}
}

func TestNovoSQLogger(t *testing.T) {
	sqlogger := bd.NovoSQLogger(simulador.Tx{
		SimulaExec: func(query string, args ...interface{}) (sql.Result, error) {
			return testdb.NewResult(1, nil, 1, nil), fmt.Errorf("erro na execução")
		},
	}, net.ParseIP("192.168.1.1"))

	if sqlogger == nil {
		t.Errorf("SQLogger não foi inicializado corretamente")
	}

	if _, err := sqlogger.Exec("SELECT"); err.Error() != "erro na execução" {
		t.Errorf("SQLogger não armazenou a transação esperada")
	}
}

func TestSQLogger_Gerar(t *testing.T) {
	conexão, err := sql.Open("testdb", "")
	if err != nil {
		t.Fatalf("erro ao inicializar a conexão do banco de dados. Detalhes: %s", err)
	}

	data := time.Now()

	cenários := []struct {
		descrição    string
		simulação    func()
		log          *bd.Log
		logEsperado  bd.Log
		erroEsperado error
	}{
		{
			descrição: "deve gerar um log corretamente",
			simulação: func() {
				logCriaçãoComando := `INSERT INTO log (id, data_criacao, endereco_remoto) VALUES (DEFAULT, $1, $2)`
				testdb.StubExec(logCriaçãoComando, testdb.NewResult(1, nil, 1, nil))
			},
			logEsperado: bd.Log{
				ID:             1,
				DataCriação:    data,
				EndereçoRemoto: net.ParseIP("192.168.1.1"),
			},
		},
		{
			descrição: "deve ignorar se já existir um log gerado",
			simulação: func() {
				logCriaçãoComando := `INSERT INTO log (id, data_criacao, endereco_remoto) VALUES (DEFAULT, $1, $2)`
				testdb.StubExec(logCriaçãoComando, testdb.NewResult(1, nil, 1, nil))
			},
			log: &bd.Log{
				ID:             1,
				DataCriação:    data,
				EndereçoRemoto: net.ParseIP("192.168.1.1"),
			},
			logEsperado: bd.Log{
				ID:             1,
				DataCriação:    data,
				EndereçoRemoto: net.ParseIP("192.168.1.1"),
			},
		},
		{
			descrição: "deve detectar um erro ao criar um log",
			simulação: func() {
				logCriaçãoComando := `INSERT INTO log (id, data_criacao, endereco_remoto) VALUES (DEFAULT, $1, $2)`
				testdb.StubExecError(logCriaçãoComando, fmt.Errorf("erro ao gerar o log"))
			},
			erroEsperado: errors.Errorf("erro ao gerar o log"),
		},
		{
			descrição: "deve detectar um erro ao obter o número de identificação do log",
			simulação: func() {
				logCriaçãoComando := `INSERT INTO log (id, data_criacao, endereco_remoto) VALUES (DEFAULT, $1, $2)`
				testdb.StubExec(logCriaçãoComando, testdb.NewResult(1, fmt.Errorf("erro ao obter id"), 1, nil))
			},
			erroEsperado: errors.Errorf("erro ao obter id"),
		},
	}

	for i, cenário := range cenários {
		testdb.Reset()
		if cenário.simulação != nil {
			cenário.simulação()
		}

		sqlogger := bd.NovoSQLogger(conexão, net.ParseIP("192.168.1.1"))
		if cenário.log != nil {
			sqlogger.Log = *cenário.log
		}

		err := sqlogger.Gerar()

		if sqlogger.Log.DataCriação.Before(cenário.logEsperado.DataCriação) {
			t.Errorf("Item %d, “%s”: data de criação inesperada. Esperava que fosse após “%s”, e foi “%s”",
				i, cenário.descrição, cenário.logEsperado.DataCriação, sqlogger.Log.DataCriação)
		}

		// Após comparar as datas, deixamos elas iguais para comparar os demais
		// campos. Isto é necessário pois não é possível prever a data de criação já
		// que é definida no próprio método.
		cenário.logEsperado.DataCriação = sqlogger.Log.DataCriação

		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(cenário.logEsperado, cenário.erroEsperado)
		if err = verificadorResultado.VerificaResultado(sqlogger.Log, err); err != nil {
			t.Error(err)
		}
	}
}
