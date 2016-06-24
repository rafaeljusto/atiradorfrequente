package atirador_test

import (
	"testing"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/atirador"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
	"github.com/rafaeljusto/atiradorfrequente/testes"
)

func TestServiço_CadastrarFrequência(t *testing.T) {
	cenários := []struct {
		descrição                string
		frequênciaPedidoCompleta protocolo.FrequênciaPedidoCompleta
		esperado                 protocolo.FrequênciaPendenteResposta
		erroEsperado             error
	}{}

	for i, cenário := range cenários {
		serviço := atirador.NovoServiço(nil)
		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(cenário.esperado, cenário.erroEsperado)

		if err := verificadorResultado.VerificaResultado(serviço.CadastrarFrequência(cenário.frequênciaPedidoCompleta)); err != nil {
			t.Error(err)
		}
	}
}

func TestServiço_ConfirmarFrequência(t *testing.T) {
	cenários := []struct {
		descrição                           string
		frequênciaConfirmaçãoPedidoCompleta protocolo.FrequênciaConfirmaçãoPedidoCompleta
		erroEsperado                        error
	}{}

	for i, cenário := range cenários {
		serviço := atirador.NovoServiço(nil)
		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(nil, cenário.erroEsperado)

		err := serviço.ConfirmarFrequência(cenário.frequênciaConfirmaçãoPedidoCompleta)
		if err = verificadorResultado.VerificaResultado(nil, err); err != nil {
			t.Error(err)
		}
	}
}
