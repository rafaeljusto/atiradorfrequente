package handler_test

import (
	"testing"

	"github.com/rafaeljusto/atiradorfrequente/rest/handler"
)

func TestRotas(t *testing.T) {
	if handler.Rotas == nil {
		t.Errorf("As rotas do servidor REST estão vazias")
	}

	if h, ok := handler.Rotas["/ping"]; !ok {
		t.Error("Handler ping não encontrado")
	} else if h() == nil {
		t.Error("Handler ping corrompido")
	}

	if h, ok := handler.Rotas["/frequencia/{cr}"]; !ok {
		t.Error("Handler de cadastro da frequência do atirador não encontrado")
	} else if h() == nil {
		t.Error("Handler de cadastro da frequência do atirador corrompido")
	}

	if h, ok := handler.Rotas["/frequencia/{cr}/{numeroControle}"]; !ok {
		t.Error("Handler de confirmação da frequência do atirador não encontrado")
	} else if h() == nil {
		t.Error("Handler de confirmação da frequência do atirador corrompido")
	}
}
