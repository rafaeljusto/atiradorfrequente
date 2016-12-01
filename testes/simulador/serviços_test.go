package simulador_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
	"github.com/rafaeljusto/atiradorfrequente/testes/simulador"
)

func TestServiçoAtirador(t *testing.T) {
	var serviçoAtiradorSimulado simulador.ServiçoAtirador
	var métodosSimulados []string

	estruturaSimulada := reflect.TypeOf(serviçoAtiradorSimulado)
	for i := 0; i < estruturaSimulada.NumField(); i++ {
		// trata somente funções como argumentos, ignorando atributos simples
		if !strings.HasPrefix(estruturaSimulada.Field(i).Type.String(), "func (") {
			continue
		}

		métodosSimulados = append(métodosSimulados, estruturaSimulada.Field(i).Name)
	}

	visitou := func(métodoSimulado string) {
		for i := len(métodosSimulados) - 1; i >= 0; i-- {
			if métodosSimulados[i] == métodoSimulado {
				métodosSimulados = append(métodosSimulados[:i], métodosSimulados[i+1:]...)
				break
			}
		}
	}

	serviçoAtiradorSimulado.SimulaCadastrarFrequência = func(protocolo.FrequênciaPedidoCompleta) (protocolo.FrequênciaPendenteResposta, error) {
		visitou("SimulaCadastrarFrequência")
		return protocolo.FrequênciaPendenteResposta{}, nil
	}

	serviçoAtiradorSimulado.SimulaObterFrequência = func(cr int, númeroControle protocolo.NúmeroControle, códigoVerificação string) (protocolo.FrequênciaResposta, error) {
		visitou("SimulaObterFrequência")
		return protocolo.FrequênciaResposta{}, nil
	}

	serviçoAtiradorSimulado.SimulaConfirmarFrequência = func(protocolo.FrequênciaConfirmaçãoPedidoCompleta) error {
		visitou("SimulaConfirmarFrequência")
		return nil
	}

	serviçoAtiradorSimulado.CadastrarFrequência(protocolo.FrequênciaPedidoCompleta{})
	serviçoAtiradorSimulado.ObterFrequência(0, "", "")
	serviçoAtiradorSimulado.ConfirmarFrequência(protocolo.FrequênciaConfirmaçãoPedidoCompleta{})

	if len(métodosSimulados) > 0 {
		t.Errorf("métodos %#v não foram chamados", métodosSimulados)
	}
}
