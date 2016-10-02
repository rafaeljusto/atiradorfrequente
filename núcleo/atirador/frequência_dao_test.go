package atirador

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"testing"
	"time"

	"github.com/erikstmartin/go-testdb"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/bd"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/erros"
	"github.com/rafaeljusto/atiradorfrequente/testes"
	"github.com/registrobr/gostk/errors"
)

func TestFrequênciaDAOImpl_criar(t *testing.T) {
	conexão, err := sql.Open("testdb", "")
	if err != nil {
		t.Fatalf("erro ao inicializar a conexão do banco de dados. Detalhes: %s", err)
	}

	data := time.Now()

	cenários := []struct {
		descrição          string
		simulação          func()
		frequência         *frequência
		frequênciaEsperada frequência
		erroEsperado       error
	}{
		{
			descrição: "deve criar corretamente a frequência",
			simulação: func() {
				testdb.StubQuery(frequênciaCriaçãoComando, testdb.RowsFromSlice([]string{"id"}, [][]driver.Value{{1}}))
				testdb.StubExec(frequênciaLogCriaçãoComando, testdb.NewResult(1, nil, 1, nil))

				logCriaçãoComando := `INSERT INTO log (id, data_criacao, endereco_remoto) VALUES (DEFAULT, $1, $2) RETURNING id`
				testdb.StubQuery(logCriaçãoComando, testdb.RowsFromSlice([]string{"id"}, [][]driver.Value{{1}}))
			},
			frequência: &frequência{
				Controle:          98765,
				CR:                1234567890,
				Calibre:           ".380",
				ArmaUtilizada:     "Arma Clube",
				NúmeroSérie:       "ZA785671",
				GuiaDeTráfego:     762556223,
				QuantidadeMunição: 50,
				DataInício:        data.Add(-1 * time.Hour),
				DataTérmino:       data.Add(-10 * time.Minute),
				revisão:           2, // revisão sempre inicia com zero
			},
			frequênciaEsperada: frequência{
				ID:                1,
				Controle:          98765,
				CR:                1234567890,
				Calibre:           ".380",
				ArmaUtilizada:     "Arma Clube",
				NúmeroSérie:       "ZA785671",
				GuiaDeTráfego:     762556223,
				QuantidadeMunição: 50,
				DataInício:        data.Add(-1 * time.Hour),
				DataTérmino:       data.Add(-10 * time.Minute),
				DataCriação:       data,
				revisão:           0,
			},
		},
		{
			descrição:    "deve detectar quando a frequência não está definida",
			erroEsperado: erros.ObjetoIndefinido,
		},
		{
			descrição: "deve detectar um erro ao criar a frequência",
			simulação: func() {
				testdb.StubQueryError(frequênciaCriaçãoComando, fmt.Errorf("erro de execução"))
			},
			frequência: &frequência{
				Controle:          98765,
				CR:                1234567890,
				Calibre:           ".380",
				ArmaUtilizada:     "Arma Clube",
				NúmeroSérie:       "ZA785671",
				GuiaDeTráfego:     762556223,
				QuantidadeMunição: 50,
				DataInício:        data.Add(-1 * time.Hour),
				DataTérmino:       data.Add(-10 * time.Minute),
				revisão:           2, // revisão sempre inicia com zero
			},
			erroEsperado: errors.Errorf("erro de execução"),
		},
		{
			descrição: "deve detectar um erro ao obter a identificação da frequência",
			simulação: func() {
				testdb.StubQuery(frequênciaCriaçãoComando, testdb.RowsFromSlice([]string{"id"}, [][]driver.Value{{"xxx"}}))
			},
			frequência: &frequência{
				Controle:          98765,
				CR:                1234567890,
				Calibre:           ".380",
				ArmaUtilizada:     "Arma Clube",
				NúmeroSérie:       "ZA785671",
				GuiaDeTráfego:     762556223,
				QuantidadeMunição: 50,
				DataInício:        data.Add(-1 * time.Hour),
				DataTérmino:       data.Add(-10 * time.Minute),
				revisão:           2, // revisão sempre inicia com zero
			},
			erroEsperado: errors.Errorf(`sql: Scan error on column index 0: converting driver.Value type string ("xxx") to a int64: invalid syntax`),
		},
		{
			descrição: "deve detectar um erro ao gerar uma identificação de log",
			simulação: func() {
				testdb.StubQuery(frequênciaCriaçãoComando, testdb.RowsFromSlice([]string{"id"}, [][]driver.Value{{1}}))
				testdb.StubExec(frequênciaLogCriaçãoComando, testdb.NewResult(1, nil, 1, nil))

				logCriaçãoComando := `INSERT INTO log (id, data_criacao, endereco_remoto) VALUES (DEFAULT, $1, $2) RETURNING id`
				testdb.StubQueryError(logCriaçãoComando, fmt.Errorf("erro na criação do id log"))
			},
			frequência: &frequência{
				Controle:          98765,
				CR:                1234567890,
				Calibre:           ".380",
				ArmaUtilizada:     "Arma Clube",
				NúmeroSérie:       "ZA785671",
				GuiaDeTráfego:     762556223,
				QuantidadeMunição: 50,
				DataInício:        data.Add(-1 * time.Hour),
				DataTérmino:       data.Add(-10 * time.Minute),
				revisão:           2, // revisão sempre inicia com zero
			},
			erroEsperado: errors.Errorf("erro na criação do id log"),
		},
		{
			descrição: "deve detectar um erro ao gerar uma entrada de log",
			simulação: func() {
				testdb.StubQuery(frequênciaCriaçãoComando, testdb.RowsFromSlice([]string{"id"}, [][]driver.Value{{1}}))
				testdb.StubExecError(frequênciaLogCriaçãoComando, fmt.Errorf("erro na criação do log"))

				logCriaçãoComando := `INSERT INTO log (id, data_criacao, endereco_remoto) VALUES (DEFAULT, $1, $2) RETURNING id`
				testdb.StubQuery(logCriaçãoComando, testdb.RowsFromSlice([]string{"id"}, [][]driver.Value{{1}}))
			},
			frequência: &frequência{
				Controle:          98765,
				CR:                1234567890,
				Calibre:           ".380",
				ArmaUtilizada:     "Arma Clube",
				NúmeroSérie:       "ZA785671",
				GuiaDeTráfego:     762556223,
				QuantidadeMunição: 50,
				DataInício:        data.Add(-1 * time.Hour),
				DataTérmino:       data.Add(-10 * time.Minute),
				revisão:           2, // revisão sempre inicia com zero
			},
			erroEsperado: errors.Errorf("erro na criação do log"),
		},
	}

	for i, cenário := range cenários {
		testdb.Reset()
		if cenário.simulação != nil {
			cenário.simulação()
		}

		dao := novaFrequênciaDAO(bd.NovoSQLogger(conexão, nil))
		err := dao.criar(cenário.frequência)

		if cenário.frequência != nil {
			if cenário.frequência.DataCriação.Before(cenário.frequênciaEsperada.DataCriação) {
				t.Errorf("Item %d, “%s”: data de criação inesperada. Esperava que fosse após “%s”, e foi “%s”",
					i, cenário.descrição, cenário.frequênciaEsperada.DataCriação, cenário.frequência.DataCriação)
			}

			// Após comparar as datas, deixamos elas iguais para comparar os demais
			// campos. Isto é necessário pois não é possível prever a data de criação já
			// que é definida no próprio método.
			cenário.frequênciaEsperada.DataCriação = cenário.frequência.DataCriação
		}

		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(&cenário.frequênciaEsperada, cenário.erroEsperado)
		if err = verificadorResultado.VerificaResultado(cenário.frequência, err); err != nil {
			t.Error(err)
		}
	}
}

func TestFrequênciaDAOImpl_atualizar(t *testing.T) {
	conexão, err := sql.Open("testdb", "")
	if err != nil {
		t.Fatalf("erro ao inicializar a conexão do banco de dados. Detalhes: %s", err)
	}

	data := time.Now()

	cenários := []struct {
		descrição          string
		simulação          func()
		frequência         *frequência
		frequênciaEsperada frequência
		erroEsperado       error
	}{
		{
			descrição: "deve atualizar corretamente a frequência",
			simulação: func() {
				testdb.StubExec(frequênciaAtualizaçãoComando, testdb.NewResult(1, nil, 1, nil))
				testdb.StubExec(frequênciaLogCriaçãoComando, testdb.NewResult(1, nil, 1, nil))

				logCriaçãoComando := `INSERT INTO log (id, data_criacao, endereco_remoto) VALUES (DEFAULT, $1, $2) RETURNING id`
				testdb.StubQuery(logCriaçãoComando, testdb.RowsFromSlice([]string{"id"}, [][]driver.Value{{1}}))
			},
			frequência: &frequência{
				ID:                1,
				Controle:          98765,
				CR:                1234567890,
				Calibre:           ".380",
				ArmaUtilizada:     "Arma Clube",
				NúmeroSérie:       "ZA785671",
				GuiaDeTráfego:     762556223,
				QuantidadeMunição: 50,
				DataInício:        data.Add(-1 * time.Hour),
				DataTérmino:       data.Add(-10 * time.Minute),
				DataCriação:       data.Add(-30 * time.Second),
				DataConfirmação:   data,
				ImagemNúmeroControle: `TWFuIGlzIGRpc3Rpbmd1aXNoZWQsIG5vdCBvbmx5IGJ5IGhpcyByZWFzb24sIGJ1dCBieSB0aGlz
IHNpbmd1bGFyIHBhc3Npb24gZnJvbSBvdGhlciBhbmltYWxzLCB3aGljaCBpcyBhIGx1c3Qgb2Yg
dGhlIG1pbmQsIHRoYXQgYnkgYSBwZXJzZXZlcmFuY2Ugb2YgZGVsaWdodCBpbiB0aGUgY29udGlu
dWVkIGFuZCBpbmRlZmF0aWdhYmxlIGdlbmVyYXRpb24gb2Yga25vd2xlZGdlLCBleGNlZWRzIHRo
ZSBzaG9ydCB2ZWhlbWVuY2Ugb2YgYW55IGNhcm5hbCBwbGVhc3VyZS4=`,
				ImagemConfirmação: `QW5kIGlmIHRoZSBkYXRhIGlzIGEgYml0IGxvbmdlciwgdGhlIGJhc2U2NCBlbmNvZGVkIGRhdGEgd2l
sbCBzcGFuIG11bHRpcGxlIGxpbmVzLgodGhlciBhbmltYWxzLCB3aGljaCBpcyBhIGx1c3Qgb2Yg
dGhlIG1pbmQsIHRoYXQgYnkgYSBwZXJzZXZlcmFuY2Ugb2YgZGVsaWdodCBpbiB0aGUgY29udGlu
dWVkIGFuZCBpbmRlZmF0aWdhYmxlIGdlbmVyYXRpb24gb2Yga25vd2xlZGdlLCBleGNlZWRzIHRo
ZSBzaG9ydCB2ZWhlbWVuY2Ugb2YgYW55IGNhcm5hbCBwbGVhc3VyZS4=`,
				revisão: 0,
			},
			frequênciaEsperada: frequência{
				ID:                1,
				Controle:          98765,
				CR:                1234567890,
				Calibre:           ".380",
				ArmaUtilizada:     "Arma Clube",
				NúmeroSérie:       "ZA785671",
				GuiaDeTráfego:     762556223,
				QuantidadeMunição: 50,
				DataInício:        data.Add(-1 * time.Hour),
				DataTérmino:       data.Add(-10 * time.Minute),
				DataCriação:       data.Add(-30 * time.Second),
				DataAtualização:   data,
				DataConfirmação:   data,
				ImagemNúmeroControle: `TWFuIGlzIGRpc3Rpbmd1aXNoZWQsIG5vdCBvbmx5IGJ5IGhpcyByZWFzb24sIGJ1dCBieSB0aGlz
IHNpbmd1bGFyIHBhc3Npb24gZnJvbSBvdGhlciBhbmltYWxzLCB3aGljaCBpcyBhIGx1c3Qgb2Yg
dGhlIG1pbmQsIHRoYXQgYnkgYSBwZXJzZXZlcmFuY2Ugb2YgZGVsaWdodCBpbiB0aGUgY29udGlu
dWVkIGFuZCBpbmRlZmF0aWdhYmxlIGdlbmVyYXRpb24gb2Yga25vd2xlZGdlLCBleGNlZWRzIHRo
ZSBzaG9ydCB2ZWhlbWVuY2Ugb2YgYW55IGNhcm5hbCBwbGVhc3VyZS4=`,
				ImagemConfirmação: `QW5kIGlmIHRoZSBkYXRhIGlzIGEgYml0IGxvbmdlciwgdGhlIGJhc2U2NCBlbmNvZGVkIGRhdGEgd2l
sbCBzcGFuIG11bHRpcGxlIGxpbmVzLgodGhlciBhbmltYWxzLCB3aGljaCBpcyBhIGx1c3Qgb2Yg
dGhlIG1pbmQsIHRoYXQgYnkgYSBwZXJzZXZlcmFuY2Ugb2YgZGVsaWdodCBpbiB0aGUgY29udGlu
dWVkIGFuZCBpbmRlZmF0aWdhYmxlIGdlbmVyYXRpb24gb2Yga25vd2xlZGdlLCBleGNlZWRzIHRo
ZSBzaG9ydCB2ZWhlbWVuY2Ugb2YgYW55IGNhcm5hbCBwbGVhc3VyZS4=`,
				revisão: 1,
			},
		},
		{
			descrição:    "deve detectar quando a frequência não está definida",
			erroEsperado: erros.ObjetoIndefinido,
		},
		{
			descrição: "deve detectar um erro ao atualizar a frequência",
			simulação: func() {
				testdb.StubExecError(frequênciaAtualizaçãoComando, fmt.Errorf("erro de execução"))
			},
			frequência: &frequência{
				ID:                1,
				Controle:          98765,
				CR:                1234567890,
				Calibre:           ".380",
				ArmaUtilizada:     "Arma Clube",
				NúmeroSérie:       "ZA785671",
				GuiaDeTráfego:     762556223,
				QuantidadeMunição: 50,
				DataInício:        data.Add(-1 * time.Hour),
				DataTérmino:       data.Add(-10 * time.Minute),
				DataCriação:       data.Add(-30 * time.Second),
				DataConfirmação:   data,
				ImagemNúmeroControle: `TWFuIGlzIGRpc3Rpbmd1aXNoZWQsIG5vdCBvbmx5IGJ5IGhpcyByZWFzb24sIGJ1dCBieSB0aGlz
IHNpbmd1bGFyIHBhc3Npb24gZnJvbSBvdGhlciBhbmltYWxzLCB3aGljaCBpcyBhIGx1c3Qgb2Yg
dGhlIG1pbmQsIHRoYXQgYnkgYSBwZXJzZXZlcmFuY2Ugb2YgZGVsaWdodCBpbiB0aGUgY29udGlu
dWVkIGFuZCBpbmRlZmF0aWdhYmxlIGdlbmVyYXRpb24gb2Yga25vd2xlZGdlLCBleGNlZWRzIHRo
ZSBzaG9ydCB2ZWhlbWVuY2Ugb2YgYW55IGNhcm5hbCBwbGVhc3VyZS4=`,
				ImagemConfirmação: `QW5kIGlmIHRoZSBkYXRhIGlzIGEgYml0IGxvbmdlciwgdGhlIGJhc2U2NCBlbmNvZGVkIGRhdGEgd2l
sbCBzcGFuIG11bHRpcGxlIGxpbmVzLgodGhlciBhbmltYWxzLCB3aGljaCBpcyBhIGx1c3Qgb2Yg
dGhlIG1pbmQsIHRoYXQgYnkgYSBwZXJzZXZlcmFuY2Ugb2YgZGVsaWdodCBpbiB0aGUgY29udGlu
dWVkIGFuZCBpbmRlZmF0aWdhYmxlIGdlbmVyYXRpb24gb2Yga25vd2xlZGdlLCBleGNlZWRzIHRo
ZSBzaG9ydCB2ZWhlbWVuY2Ugb2YgYW55IGNhcm5hbCBwbGVhc3VyZS4=`,
				revisão: 0,
			},
			erroEsperado: errors.Errorf("erro de execução"),
		},
		{
			descrição: "deve detectar um erro ao obter as linhas afetadas na base de dados",
			simulação: func() {
				testdb.StubExec(frequênciaAtualizaçãoComando, testdb.NewResult(0, nil, 1, fmt.Errorf("erro com o ID")))
			},
			frequência: &frequência{
				ID:                1,
				Controle:          98765,
				CR:                1234567890,
				Calibre:           ".380",
				ArmaUtilizada:     "Arma Clube",
				NúmeroSérie:       "ZA785671",
				GuiaDeTráfego:     762556223,
				QuantidadeMunição: 50,
				DataInício:        data.Add(-1 * time.Hour),
				DataTérmino:       data.Add(-10 * time.Minute),
				DataCriação:       data.Add(-30 * time.Second),
				DataConfirmação:   data,
				ImagemNúmeroControle: `TWFuIGlzIGRpc3Rpbmd1aXNoZWQsIG5vdCBvbmx5IGJ5IGhpcyByZWFzb24sIGJ1dCBieSB0aGlz
IHNpbmd1bGFyIHBhc3Npb24gZnJvbSBvdGhlciBhbmltYWxzLCB3aGljaCBpcyBhIGx1c3Qgb2Yg
dGhlIG1pbmQsIHRoYXQgYnkgYSBwZXJzZXZlcmFuY2Ugb2YgZGVsaWdodCBpbiB0aGUgY29udGlu
dWVkIGFuZCBpbmRlZmF0aWdhYmxlIGdlbmVyYXRpb24gb2Yga25vd2xlZGdlLCBleGNlZWRzIHRo
ZSBzaG9ydCB2ZWhlbWVuY2Ugb2YgYW55IGNhcm5hbCBwbGVhc3VyZS4=`,
				ImagemConfirmação: `QW5kIGlmIHRoZSBkYXRhIGlzIGEgYml0IGxvbmdlciwgdGhlIGJhc2U2NCBlbmNvZGVkIGRhdGEgd2l
sbCBzcGFuIG11bHRpcGxlIGxpbmVzLgodGhlciBhbmltYWxzLCB3aGljaCBpcyBhIGx1c3Qgb2Yg
dGhlIG1pbmQsIHRoYXQgYnkgYSBwZXJzZXZlcmFuY2Ugb2YgZGVsaWdodCBpbiB0aGUgY29udGlu
dWVkIGFuZCBpbmRlZmF0aWdhYmxlIGdlbmVyYXRpb24gb2Yga25vd2xlZGdlLCBleGNlZWRzIHRo
ZSBzaG9ydCB2ZWhlbWVuY2Ugb2YgYW55IGNhcm5hbCBwbGVhc3VyZS4=`,
				revisão: 0,
			},
			erroEsperado: errors.Errorf("erro com o ID"),
		},
		{
			descrição: "deve detectar quando a atualização não surtiu efeito",
			simulação: func() {
				testdb.StubExec(frequênciaAtualizaçãoComando, testdb.NewResult(0, nil, 0, nil))
			},
			frequência: &frequência{
				ID:                1,
				Controle:          98765,
				CR:                1234567890,
				Calibre:           ".380",
				ArmaUtilizada:     "Arma Clube",
				NúmeroSérie:       "ZA785671",
				GuiaDeTráfego:     762556223,
				QuantidadeMunição: 50,
				DataInício:        data.Add(-1 * time.Hour),
				DataTérmino:       data.Add(-10 * time.Minute),
				DataCriação:       data.Add(-30 * time.Second),
				DataConfirmação:   data,
				ImagemNúmeroControle: `TWFuIGlzIGRpc3Rpbmd1aXNoZWQsIG5vdCBvbmx5IGJ5IGhpcyByZWFzb24sIGJ1dCBieSB0aGlz
IHNpbmd1bGFyIHBhc3Npb24gZnJvbSBvdGhlciBhbmltYWxzLCB3aGljaCBpcyBhIGx1c3Qgb2Yg
dGhlIG1pbmQsIHRoYXQgYnkgYSBwZXJzZXZlcmFuY2Ugb2YgZGVsaWdodCBpbiB0aGUgY29udGlu
dWVkIGFuZCBpbmRlZmF0aWdhYmxlIGdlbmVyYXRpb24gb2Yga25vd2xlZGdlLCBleGNlZWRzIHRo
ZSBzaG9ydCB2ZWhlbWVuY2Ugb2YgYW55IGNhcm5hbCBwbGVhc3VyZS4=`,
				ImagemConfirmação: `QW5kIGlmIHRoZSBkYXRhIGlzIGEgYml0IGxvbmdlciwgdGhlIGJhc2U2NCBlbmNvZGVkIGRhdGEgd2l
sbCBzcGFuIG11bHRpcGxlIGxpbmVzLgodGhlciBhbmltYWxzLCB3aGljaCBpcyBhIGx1c3Qgb2Yg
dGhlIG1pbmQsIHRoYXQgYnkgYSBwZXJzZXZlcmFuY2Ugb2YgZGVsaWdodCBpbiB0aGUgY29udGlu
dWVkIGFuZCBpbmRlZmF0aWdhYmxlIGdlbmVyYXRpb24gb2Yga25vd2xlZGdlLCBleGNlZWRzIHRo
ZSBzaG9ydCB2ZWhlbWVuY2Ugb2YgYW55IGNhcm5hbCBwbGVhc3VyZS4=`,
				revisão: 0,
			},
			erroEsperado: erros.NãoAtualizado,
		},
	}

	for i, cenário := range cenários {
		testdb.Reset()
		if cenário.simulação != nil {
			cenário.simulação()
		}

		dao := novaFrequênciaDAO(bd.NovoSQLogger(conexão, nil))
		err := dao.atualizar(cenário.frequência)

		if cenário.frequência != nil {
			if cenário.frequência.DataAtualização.Before(cenário.frequênciaEsperada.DataAtualização) {
				t.Errorf("Item %d, “%s”: data de atualização inesperada. Esperava que fosse após “%s”, e foi “%s”",
					i, cenário.descrição, cenário.frequênciaEsperada.DataAtualização, cenário.frequência.DataAtualização)
			}

			// Após comparar as datas, deixamos elas iguais para comparar os demais
			// campos. Isto é necessário pois não é possível prever a data de atualização já
			// que é definida no próprio método.
			cenário.frequênciaEsperada.DataAtualização = cenário.frequência.DataAtualização
		}

		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(&cenário.frequênciaEsperada, cenário.erroEsperado)
		if err = verificadorResultado.VerificaResultado(cenário.frequência, err); err != nil {
			t.Error(err)
		}
	}
}

func TestFrequênciaDAOImpl_resgatar(t *testing.T) {
	conexão, err := sql.Open("testdb", "")
	if err != nil {
		t.Fatalf("erro ao inicializar a conexão do banco de dados. Detalhes: %s", err)
	}

	data := time.Now()

	cenários := []struct {
		descrição          string
		simulação          func()
		id                 int64
		frequênciaEsperada frequência
		erroEsperado       error
	}{
		{
			descrição: "deve resgatar corretamente uma frequência",
			simulação: func() {
				testdb.StubQuery(frequênciaResgateComando, testdb.RowsFromSlice(frequênciaResgateCampos, [][]driver.Value{
					{
						1, 98765, 1234567890, ".380", "Arma Clube", "ZA785671", 762556223, 50,
						data.Add(-1 * time.Hour), data.Add(-10 * time.Minute), data, time.Time{}, time.Time{},
						"", "", 0,
					},
				}))
			},
			id: 1,
			frequênciaEsperada: frequência{
				ID:                1,
				Controle:          98765,
				CR:                1234567890,
				Calibre:           ".380",
				ArmaUtilizada:     "Arma Clube",
				NúmeroSérie:       "ZA785671",
				GuiaDeTráfego:     762556223,
				QuantidadeMunição: 50,
				DataInício:        data.Add(-1 * time.Hour),
				DataTérmino:       data.Add(-10 * time.Minute),
				DataCriação:       data,
				revisão:           0,
			},
		},
		{
			descrição: "deve detectar um erro ao resgatar uma frequência",
			simulação: func() {
				testdb.StubQueryError(frequênciaResgateComando, fmt.Errorf("erro de execução"))
			},
			id:           1,
			erroEsperado: errors.Errorf("erro de execução"),
		},
	}

	for i, cenário := range cenários {
		testdb.Reset()
		cenário.simulação()

		dao := novaFrequênciaDAO(bd.NovoSQLogger(conexão, nil))
		f, err := dao.resgatar(cenário.id)

		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(cenário.frequênciaEsperada, cenário.erroEsperado)
		if err = verificadorResultado.VerificaResultado(f, err); err != nil {
			t.Error(err)
		}
	}
}
