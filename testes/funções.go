package testes

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/aryann/difflib"
	"github.com/davecgh/go-spew/spew"
)

// Diff facilita a exibição entre as diferenças entre dois objetos.
func Diff(a, b interface{}) []difflib.DiffRecord {
	return difflib.Diff(strings.SplitAfter(spew.Sdump(a), "\n"), strings.SplitAfter(spew.Sdump(b), "\n"))
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
	if ((err == nil || c.erroEsperado == nil) && err != c.erroEsperado) ||
		(err != nil && c.erroEsperado != nil && err.Error() != c.erroEsperado.Error()) {

		return fmt.Errorf("Item %d, “%s”: erros não batem. Esperado ‘%v’; encontrado ‘%v’",
			c.índice, c.descrição, c.erroEsperado, err)
	}

	value := reflect.ValueOf(c.erroEsperado)
	checkValue := value.Kind() == reflect.Chan ||
		value.Kind() == reflect.Func ||
		value.Kind() == reflect.Interface ||
		value.Kind() == reflect.Map ||
		value.Kind() == reflect.Ptr ||
		value.Kind() == reflect.Slice

	if c.erroEsperado == nil || (checkValue && value.IsNil()) {
		if !reflect.DeepEqual(resultado, c.resultadoEsperado) {
			return fmt.Errorf("Item %d, “%s”: resultados não batem.\n%s",
				c.índice, c.descrição, Diff(c.resultadoEsperado, resultado))
		}
	}

	return nil
}
