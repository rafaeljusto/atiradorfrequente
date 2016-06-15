package interceptador_test

import (
	"testing"

	"github.com/rafaeljusto/cr/rest/interceptador"
)

func TestCabeçalhoCompatível_DefinirCabeçalho(t *testing.T) {
	var header interceptador.CabeçalhoCompatível
	header.DefinirCabeçalho("key", "value1")
	header.DefinirCabeçalho("key", "value2")

	if value := header.Cabeçalho.Get("key"); value != "value2" {
		t.Errorf("Unexpected key “%s”", value)
	}
}

func TestCabeçalhoCompatível_AdicionarCabeçalho(t *testing.T) {
	var header interceptador.CabeçalhoCompatível
	header.AdicionarCabeçalho("key", "value1")
	header.AdicionarCabeçalho("key", "value2")

	if value := header.Cabeçalho.Get("key"); value != "value1" {
		t.Errorf("Unexpected key “%s”", value)
	}
}
