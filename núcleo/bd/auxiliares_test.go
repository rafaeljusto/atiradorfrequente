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
		argumentos := bd.Marcadores(cenário.n)

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
		inicio    int
		n         int
		esperado  string
	}{
		{
			descrição: "deve gerar os argumentos corretamente",
			inicio:    1,
			n:         3,
			esperado:  "$1,$2,$3",
		},
		{
			descrição: "deve tratar corretamente quando o inicio for negativo",
			inicio:    -1,
			n:         3,
			esperado:  "",
		},
		{
			descrição: "deve tratar corretamente quando o número de argumentos for negativo",
			inicio:    1,
			n:         -1,
			esperado:  "",
		},
		{
			descrição: "deve tratar corretamente quando o inicio for maior que o número de argumentos",
			inicio:    4,
			n:         3,
			esperado:  "",
		},
	}

	for i, cenário := range cenários {
		argumentos := bd.MarcadoresComInicio(cenário.inicio, cenário.n)

		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(cenário.esperado, nil)
		if err := verificadorResultado.VerificaResultado(argumentos, nil); err != nil {
			t.Error(err)
		}
	}
}
