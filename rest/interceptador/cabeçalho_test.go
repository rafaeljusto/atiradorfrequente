package interceptador_test

import (
	"testing"

	"github.com/rafaeljusto/atiradorfrequente/rest/interceptador"
)

func TestCabeçalhoCompatível_DefinirCabeçalho(t *testing.T) {
	var cabeçalho interceptador.CabeçalhoCompatível
	cabeçalho.DefinirCabeçalho("key", "value1")
	cabeçalho.DefinirCabeçalho("key", "value2")

	if value := cabeçalho.Cabeçalho.Get("key"); value != "value2" {
		t.Errorf("Chave “%s” inexperada", value)
	}
}

func TestCabeçalhoCompatível_AdicionarCabeçalho(t *testing.T) {
	var cabeçalho interceptador.CabeçalhoCompatível
	cabeçalho.AdicionarCabeçalho("key", "value1")
	cabeçalho.AdicionarCabeçalho("key", "value2")

	if value := cabeçalho.Cabeçalho.Get("key"); value != "value1" {
		t.Errorf("Chave “%s” inexperada", value)
	}
}
