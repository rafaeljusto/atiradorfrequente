package testes_test

import (
	"reflect"
	"testing"

	"fmt"

	"github.com/aryann/difflib"
	"github.com/rafaeljusto/atiradorfrequente/testes"
)

func TestDiff(t *testing.T) {
	a := struct {
		valor1 int
		valor2 string
	}{valor1: 123, valor2: "abc"}

	b := struct {
		valor1 int
		valor2 string
	}{valor1: 321, valor2: "cba"}

	esperado := []difflib.DiffRecord{
		{Payload: "(struct { valor1 int; valor2 string }) {\n", Delta: difflib.Common},
		{Payload: " valor1: (int) 123,\n", Delta: difflib.LeftOnly},
		{Payload: " valor2: (string) (len=3) \"abc\"\n", Delta: difflib.LeftOnly},
		{Payload: " valor1: (int) 321,\n", Delta: difflib.RightOnly},
		{Payload: " valor2: (string) (len=3) \"cba\"\n", Delta: difflib.RightOnly},
		{Payload: "}\n", Delta: difflib.Common},
		{Payload: "", Delta: difflib.Common},
	}

	if diferenças := testes.Diff(a, b); !reflect.DeepEqual(esperado, diferenças) {
		t.Errorf("Existem diferenças não esperadas. Esperado: %#v\nObtido:%#v", esperado, diferenças)
	}

	esperado = []difflib.DiffRecord{}

	if diferenças := testes.Diff(a, a); reflect.DeepEqual(esperado, diferenças) {
		t.Errorf("Existem diferenças não esperadas. Esperado: %#v\nObtido:%#v", esperado, diferenças)
	}
}

func TestTiposDaLista(t *testing.T) {
	cenários := []struct {
		descrição string
		lista     interface{}
		esperado  []string
	}{
		{
			descrição: "deve detectar os tipos da lista corretamente",
			lista:     []interface{}{test1{}, test2{}},
			esperado:  []string{"testes_test.test1", "testes_test.test2"},
		},
		{
			descrição: "deve detectar quando não é uma lista",
			lista:     test1{},
		},
		{
			descrição: "deve detectar quando a lista não esta definida",
		},
		{
			descrição: "deve detectar uma lista vazia",
			lista:     []interface{}{},
			esperado:  []string{},
		},
	}

	for i, cenário := range cenários {
		verificadorResultados := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultados.DefinirEsperado(cenário.esperado, nil)
		if err := verificadorResultados.VerificaResultado(testes.TiposDaLista(cenário.lista), nil); err != nil {
			t.Error(err)
		}
	}
}

func TestVerificadorResultados_VerificaResultado(t *testing.T) {
	cenários := []struct {
		descrição         string
		esperado          interface{}
		erroEsperado      error
		obtido            interface{}
		erroObtido        error
		resultadoEsperado error
	}{
		{
			descrição:         "deve identificar quando os erros não batem (1)",
			esperado:          1,
			erroEsperado:      fmt.Errorf("erro 1"),
			obtido:            1,
			erroObtido:        fmt.Errorf("erro 2"),
			resultadoEsperado: fmt.Errorf("Item 0, “deve identificar quando os erros não batem (1)”: erros não batem. Esperado “erro 1”; encontrado “erro 2”"),
		},
		{
			descrição:         "deve identificar quando os erros não batem (2)",
			esperado:          1,
			erroEsperado:      nil,
			obtido:            1,
			erroObtido:        fmt.Errorf("erro 2"),
			resultadoEsperado: fmt.Errorf("Item 1, “deve identificar quando os erros não batem (2)”: erros não batem. Esperado “<nil>”; encontrado “erro 2”"),
		},
		{
			descrição:         "deve identificar quando os erros não batem (3)",
			esperado:          1,
			erroEsperado:      fmt.Errorf("erro 1"),
			obtido:            1,
			erroObtido:        nil,
			resultadoEsperado: fmt.Errorf("Item 2, “deve identificar quando os erros não batem (3)”: erros não batem. Esperado “erro 1”; encontrado “<nil>”"),
		},
		{
			descrição:         "deve identificar quando os resultados não batem",
			esperado:          1,
			erroEsperado:      nil,
			obtido:            2,
			erroObtido:        nil,
			resultadoEsperado: fmt.Errorf("Item 3, “deve identificar quando os resultados não batem”: resultados não batem.\n[- (int) 1\n + (int) 2\n   ]"),
		},
		{
			descrição:    "deve identificar quando os resultados são esperados",
			esperado:     1,
			erroEsperado: fmt.Errorf("erro 1"),
			obtido:       1,
			erroObtido:   fmt.Errorf("erro 1"),
		},
	}

	for i, cenário := range cenários {
		verificadorResultados := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultados.DefinirEsperado(cenário.esperado, cenário.erroEsperado)
		err := verificadorResultados.VerificaResultado(cenário.obtido, cenário.erroObtido)

		if ((err == nil || cenário.resultadoEsperado == nil) && err != cenário.resultadoEsperado) ||
			(err != nil && cenário.resultadoEsperado != nil && err.Error() != cenário.resultadoEsperado.Error()) {
			t.Errorf("Item %d, “%s”: resultados não batem. Esperado: %#v\nObtido: %#v",
				i, cenário.descrição, cenário.resultadoEsperado, err)
		}
	}
}

type test1 struct{}
type test2 struct{}
