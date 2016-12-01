package handler

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/erikstmartin/go-testdb"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/bd"
	"github.com/rafaeljusto/atiradorfrequente/testes"
	"github.com/rafaeljusto/atiradorfrequente/testes/simulador"
	"github.com/registrobr/gostk/log"
)

func TestPing_Get(t *testing.T) {
	conexão, err := sql.Open("testdb", "")
	if err != nil {
		t.Fatalf("erro ao inicializar a conexão do banco de dados. Detalhes: %s", err)
	}

	cenários := []struct {
		descrição          string
		simulação          func()
		logger             log.Logger
		códigoHTTPEsperado int
	}{
		{
			descrição: "deve reportar o funcionamento correto do sistema",
			simulação: func() {
				testdb.StubQuery("SELECT NOW() AT TIME ZONE 'UTC'", testdb.RowsFromSlice([]string{"NOW()"}, [][]driver.Value{
					{time.Now().UTC()},
				}))
			},
			códigoHTTPEsperado: http.StatusNoContent,
		},
		{
			descrição: "deve detectar uma falha no banco de dados",
			simulação: func() {
				testdb.StubQueryError("SELECT NOW() AT TIME ZONE 'UTC'", fmt.Errorf("erro de conexão com o banco de dados"))
			},
			logger: simulador.Logger{
				SimulaError: func(e error) {
					if !strings.HasSuffix(e.Error(), "erro de conexão com o banco de dados") {
						t.Error("não está adicionando o erro correto ao log")
					}
				},
			},
			códigoHTTPEsperado: http.StatusInternalServerError,
		},
	}

	for i, cenário := range cenários {
		testdb.Reset()
		if cenário.simulação != nil {
			cenário.simulação()
		}

		var handler ping
		handler.DefineLogger(cenário.logger)
		handler.DefineTx(bd.NovoSQLogger(conexão, net.ParseIP("127.0.0.1")))

		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(cenário.códigoHTTPEsperado, nil)
		if err := verificadorResultado.VerificaResultado(handler.Get(), nil); err != nil {
			t.Error(err)
		}
	}
}

func TestPing_Interceptors(t *testing.T) {
	esperado := []string{
		"*interceptador.EndereçoRemoto",
		"*interceptador.Log",
		"*interceptor.Introspector",
		"*interceptador.Codificador",
		"*interceptador.ParâmetrosConsulta",
		"*interceptador.VariáveisEndereço",
		"*interceptador.Padronizador",
		"*interceptador.BD",
	}

	var handler ping

	verificadorResultado := testes.NovoVerificadorResultados("deve conter os interceptadores corretos", 0)
	verificadorResultado.DefinirEsperado(esperado, nil)
	if err := verificadorResultado.VerificaResultado(testes.TiposDaLista(handler.Interceptors()), nil); err != nil {
		t.Error(err)
	}
}
