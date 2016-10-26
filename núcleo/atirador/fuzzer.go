// +build gofuzz

package atirador

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"

	"github.com/erikstmartin/go-testdb"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/bd"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/config"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
)

// FuzzCadastrarFrequência é utilizado pela ferramenta go-fuzz, responsável por
// testar o cadastro de frequência com dados aleatórios.
func FuzzCadastrarFrequência(dados []byte) int {
	var frequênciaPedidoCompleta protocolo.FrequênciaPedidoCompleta
	if err := json.Unmarshal(dados, frequênciaPedidoCompleta); err != nil {
		return -1
	}

	conexão, err := sql.Open("testdb", "")
	if err != nil {
		panic(err)
	}

	testdb.StubQuery(frequênciaCriaçãoComando, testdb.RowsFromSlice([]string{"id"}, [][]driver.Value{{1}}))
	testdb.StubExec(frequênciaLogCriaçãoComando, testdb.NewResult(1, nil, 1, nil))

	logCriaçãoComando := `INSERT INTO log (id, data_criacao, endereco_remoto) VALUES (DEFAULT, $1, $2) RETURNING id`
	testdb.StubQuery(logCriaçãoComando, testdb.RowsFromSlice([]string{"id"}, [][]driver.Value{{1}}))

	var configuração config.Configuração
	config.DefinirValoresPadrão(&configuração)

	serviço := NovoServiço(bd.NovoSQLogger(conexão, nil), configuração)
	if _, err := serviço.CadastrarFrequência(frequênciaPedidoCompleta); err != nil {
		if _, ok := err.(protocolo.Mensagens); !ok {
			panic(err)
		}

		return 0
	}

	return 1
}
