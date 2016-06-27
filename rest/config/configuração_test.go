package config_test

import (
	"reflect"
	"testing"

	"github.com/rafaeljusto/atiradorfrequente/rest/config"
	"github.com/rafaeljusto/atiradorfrequente/testes"
)

func TestConfiguração(t *testing.T) {
	var esperado config.Configuração
	esperado.BancoDados.Nome = "teste"

	config.AtualizarConfiguração(&esperado)
	obtido := config.Atual()

	if !reflect.DeepEqual(&esperado, obtido) {
		t.Errorf("configuração inesperada.\n%v", testes.Diff(esperado, obtido))
	}
}
