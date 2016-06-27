package simulador_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/rafaeljusto/atiradorfrequente/testes/simulador"
	"github.com/trajber/handy"
)

func TestHandler(t *testing.T) {
	var handlerSimulado simulador.Handler
	var métodosSimulados []string

	estruturaBDSimulado := reflect.TypeOf(handlerSimulado)
	for i := 0; i < estruturaBDSimulado.NumField(); i++ {
		// trata somente funções como argumentos, ignorando atributos simples
		if !strings.HasPrefix(estruturaBDSimulado.Field(i).Type.String(), "func (") {
			continue
		}

		métodosSimulados = append(métodosSimulados, estruturaBDSimulado.Field(i).Name)
	}

	visitou := func(métodoSimulado string) {
		for i := len(métodosSimulados) - 1; i >= 0; i-- {
			if métodosSimulados[i] == métodoSimulado {
				métodosSimulados = append(métodosSimulados[:i], métodosSimulados[i+1:]...)
				break
			}
		}
	}

	var err error
	handlerSimulado.SimulaRequisição, err = http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Error(err)
	}

	handlerSimulado.SimulaResposta = httptest.NewRecorder()
	handlerSimulado.SimulaVariáveisEndereço = handy.URIVars{}

	handlerSimulado.SimulaGet = func() int {
		visitou("SimulaGet")
		return 0
	}

	handlerSimulado.SimulaPost = func() int {
		visitou("SimulaPost")
		return 0
	}

	handlerSimulado.SimulaPut = func() int {
		visitou("SimulaPut")
		return 0
	}

	handlerSimulado.SimulaDelete = func() int {
		visitou("SimulaDelete")
		return 0
	}

	handlerSimulado.SimulaPatch = func() int {
		visitou("SimulaPatch")
		return 0
	}

	handlerSimulado.SimulaHead = func() int {
		visitou("SimulaHead")
		return 0
	}

	if requisição := handlerSimulado.Req(); !reflect.DeepEqual(handlerSimulado.SimulaRequisição, requisição) {
		t.Error("requisição inesperada")
	}

	if resposta := handlerSimulado.ResponseWriter(); !reflect.DeepEqual(handlerSimulado.SimulaResposta, resposta) {
		t.Error("resposta inesperada")
	}

	if variáveisEndereço := handlerSimulado.URIVars(); !reflect.DeepEqual(handlerSimulado.URIVars(), variáveisEndereço) {
		t.Errorf("váriaveis de endereço inesperadas")
	}

	handlerSimulado.Get()
	handlerSimulado.Post()
	handlerSimulado.Put()
	handlerSimulado.Delete()
	handlerSimulado.Patch()
	handlerSimulado.Head()

	if len(métodosSimulados) > 0 {
		t.Errorf("métodos %#v não foram chamados", métodosSimulados)
	}
}
