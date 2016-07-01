package testes

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/aryann/difflib"
	"github.com/davecgh/go-spew/spew"
	"github.com/registrobr/gostk/errors"
)

// Diff facilita a exibição entre as diferenças entre dois objetos.
func Diff(a, b interface{}) []difflib.DiffRecord {
	return difflib.Diff(strings.SplitAfter(spew.Sdump(a), "\n"), strings.SplitAfter(spew.Sdump(b), "\n"))
}

// TiposDaLista retorno todos os tipos de uma lista de interface{} em formato
// texto. Caso o elemento não seja uma lista, nil é retornado.
func TiposDaLista(elemento interface{}) []string {
	valor := reflect.ValueOf(elemento)
	if valor.Kind() != reflect.Slice {
		return nil
	}

	tipos := make([]string, valor.Len())
	for i := 0; i < valor.Len(); i++ {
		tipos[i] = reflect.TypeOf(valor.Index(i).Interface()).String()
	}
	return tipos
}

// VerificadorResultados armazena os dados necessários para comparar dois
// resultados.
type VerificadorResultados struct {
	descrição         string
	índice            int
	resultadoEsperado interface{}
	erroEsperado      error
}

// NovoVerificadorResultados inicializa o tipo VerificadorResultados com o
// cenário analisado.
func NovoVerificadorResultados(descrição string, índice int) *VerificadorResultados {
	return &VerificadorResultados{descrição: descrição, índice: índice}
}

// DefinirEsperado define o resultado esperado para o teste.
func (c *VerificadorResultados) DefinirEsperado(resultado interface{}, err error) {
	c.resultadoEsperado = resultado
	c.erroEsperado = err
}

// VerificaResultado verifica se o resultado apresentado é igual ao resultado
// esperado.
func (c *VerificadorResultados) VerificaResultado(resultado interface{}, err error) error {
	if !errors.Equal(err, c.erroEsperado) {
		return fmt.Errorf("Item %d, “%s”: erros não batem. Esperado “%v”; encontrado “%v”",
			c.índice, c.descrição, c.erroEsperado, err)
	}

	valor := reflect.ValueOf(c.erroEsperado)
	verificaValor := valor.Kind() == reflect.Chan ||
		valor.Kind() == reflect.Func ||
		valor.Kind() == reflect.Interface ||
		valor.Kind() == reflect.Map ||
		valor.Kind() == reflect.Ptr ||
		valor.Kind() == reflect.Slice

	if c.erroEsperado == nil || (verificaValor && valor.IsNil()) {
		if !reflect.DeepEqual(resultado, c.resultadoEsperado) {
			return fmt.Errorf("Item %d, “%s”: resultados não batem.\n%s",
				c.índice, c.descrição, Diff(c.resultadoEsperado, resultado))
		}
	}

	return nil
}
