package interceptador_test

import (
	"testing"

	"strings"

	"fmt"

	"net/http"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
	"github.com/rafaeljusto/atiradorfrequente/rest/interceptador"
	"github.com/rafaeljusto/atiradorfrequente/testes"
	"github.com/rafaeljusto/atiradorfrequente/testes/simulador"
	"github.com/registrobr/gostk/log"
	"github.com/trajber/handy/interceptor"
)

func TestPadronizador_Before(t *testing.T) {
	cenários := []struct {
		descrição          string
		métodoHTTP         string
		campo1             padronizadorCampoSimulado
		campo2             string
		logger             log.Logger
		códigoHTTPEsperado int
		campo1Esperado     padronizadorCampoSimulado
		campo2Esperado     string
		mensagensEsperadas protocolo.Mensagens
	}{
		{
			descrição:  "deve padronizar corretamente o campo",
			métodoHTTP: "GET",
			campo1:     padronizadorCampoSimulado("teste"),
			logger: &simulador.Logger{
				SimulaDebug: func(m ...interface{}) {
					mensagem := fmt.Sprint(m...)
					if mensagem != "Interceptador Antes: Padronizador" {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
			},
			campo1Esperado: padronizadorCampoSimulado("TESTE"),
		},
		{
			descrição:  "deve ignorar quando não for o método HTTP correspondente",
			métodoHTTP: "POST",
			campo1:     padronizadorCampoSimulado("teste"),
			logger: &simulador.Logger{
				SimulaDebug: func(m ...interface{}) {
					mensagem := fmt.Sprint(m...)
					if mensagem != "Interceptador Antes: Padronizador" {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
			},
			campo1Esperado: padronizadorCampoSimulado("teste"),
		},
		{
			descrição:  "deve detectar um campo inválido",
			métodoHTTP: "GET",
			campo1:     padronizadorCampoSimulado("xxx"),
			logger: &simulador.Logger{
				SimulaDebug: func(m ...interface{}) {
					mensagem := fmt.Sprint(m...)
					if mensagem != "Interceptador Antes: Padronizador" {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
			},
			códigoHTTPEsperado: http.StatusBadRequest,
			campo1Esperado:     padronizadorCampoSimulado("XXX"),
			mensagensEsperadas: protocolo.NovasMensagens(
				protocolo.NovaMensagem(protocolo.MensagemCódigoParâmetroInválido),
			),
		},
		{
			descrição:  "deve ignorar quando um campo não suporta padronização",
			métodoHTTP: "PUT",
			campo2:     "teste",
			logger: &simulador.Logger{
				SimulaDebug: func(m ...interface{}) {
					mensagem := fmt.Sprint(m...)
					if mensagem != "Interceptador Antes: Padronizador" {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
			},
			campo2Esperado: "teste",
		},
	}

	for i, cenário := range cenários {
		var handler padronizadorSimulado
		handler.Campo1 = cenário.campo1
		handler.Campo2 = cenário.campo2
		handler.DefineLogger(cenário.logger)
		handler.SimulaRequisição, _ = http.NewRequest(cenário.métodoHTTP, "http://exemplo.com.br", nil)

		estrutura := interceptor.NewIntrospector(&handler)
		if códigoHTTP := estrutura.Before(); códigoHTTP != 0 {
			t.Errorf("Item %d, “%s”: código HTTP %d inesperado",
				i, cenário.descrição, códigoHTTP)
			continue
		}

		padronizador := interceptador.NovoPadronizador(&handler)
		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)

		verificadorResultado.DefinirEsperado(cenário.códigoHTTPEsperado, nil)
		if err := verificadorResultado.VerificaResultado(padronizador.Before(), nil); err != nil {
			t.Error(err)
		}

		verificadorResultado.DefinirEsperado(cenário.campo1Esperado, nil)
		if err := verificadorResultado.VerificaResultado(handler.Campo1, nil); err != nil {
			t.Error(err)
		}

		verificadorResultado.DefinirEsperado(cenário.campo2Esperado, nil)
		if err := verificadorResultado.VerificaResultado(handler.Campo2, nil); err != nil {
			t.Error(err)
		}

		verificadorResultado.DefinirEsperado(cenário.mensagensEsperadas, nil)
		if err := verificadorResultado.VerificaResultado(handler.Mensagens, nil); err != nil {
			t.Error(err)
		}
	}
}

type padronizadorSimulado struct {
	interceptador.LogCompatível
	interceptador.MensagensCompatível
	interceptor.IntrospectorCompliant
	simulador.Handler

	Campo1 padronizadorCampoSimulado `request:"get"`
	Campo2 string                    `request:"put"`
}

type padronizadorCampoSimulado string

func (p *padronizadorCampoSimulado) Normalizar() {
	if p == nil {
		return
	}

	*p = padronizadorCampoSimulado(strings.ToUpper(string(*p)))
}

func (p padronizadorCampoSimulado) Validar() protocolo.Mensagens {
	if p != "TESTE" {
		return protocolo.NovasMensagens(
			protocolo.NovaMensagem(protocolo.MensagemCódigoParâmetroInválido),
		)
	}

	return nil
}
