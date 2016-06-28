package bd_test

import (
	"testing"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/bd"
	"github.com/rafaeljusto/atiradorfrequente/testes"
)

func TestMarcadores(t *testing.T) {
	cenários := []struct {
		descrição string
		n         int
		esperado  string
	}{
		{
			descrição: "deve gerar os argumentos corretamente",
			n:         3,
			esperado:  "$1,$2,$3",
		},
		{
			descrição: "deve tratar corretamente quando o número de argumentos for negativo",
			n:         -1,
			esperado:  "",
		},
	}

	for i, cenário := range cenários {
		argumentos := bd.MarcadoresPSQL(cenário.n)

		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(cenário.esperado, nil)
		if err := verificadorResultado.VerificaResultado(argumentos, nil); err != nil {
			t.Error(err)
		}
	}
}

func TestMarcadoresComInicio(t *testing.T) {
	cenários := []struct {
		descrição string
		início    int
		n         int
		esperado  string
	}{
		{
			descrição: "deve gerar os argumentos corretamente",
			início:    1,
			n:         3,
			esperado:  "$1,$2,$3",
		},
		{
			descrição: "deve tratar corretamente quando o início for negativo",
			início:    -1,
			n:         3,
			esperado:  "",
		},
		{
			descrição: "deve tratar corretamente quando o número de argumentos for negativo",
			início:    1,
			n:         -1,
			esperado:  "",
		},
		{
			descrição: "deve tratar corretamente quando o início for maior que o número de argumentos",
			início:    4,
			n:         3,
			esperado:  "",
		},
	}

	for i, cenário := range cenários {
		argumentos := bd.MarcadoresPSQLComInício(cenário.início, cenário.n)

		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(cenário.esperado, nil)
		if err := verificadorResultado.VerificaResultado(argumentos, nil); err != nil {
			t.Error(err)
		}
	}
}
