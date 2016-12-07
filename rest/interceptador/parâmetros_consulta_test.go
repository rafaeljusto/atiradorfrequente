package interceptador_test

import (
	"fmt"
	"net"
	"net/http"
	"testing"

	"net/url"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
	"github.com/rafaeljusto/atiradorfrequente/rest/interceptador"
	"github.com/rafaeljusto/atiradorfrequente/testes"
	"github.com/rafaeljusto/atiradorfrequente/testes/simulador"
	"github.com/registrobr/gostk/log"
	"github.com/trajber/handy/interceptor"
)

func TestParâmetrosConsulta_Before(t *testing.T) {
	cenários := []struct {
		descrição          string
		requisição         *http.Request
		logger             log.Logger
		códigoHTTPEsperado int
		handlerEsperado    parâmetrosConsultaSimulado
	}{
		{
			descrição: "deve preencher corretamente as variáveis do handler",
			requisição: func() *http.Request {
				req, err := http.NewRequest("GET", "/test", nil)
				if err != nil {
					t.Fatal(err)
				}

				req.Form = make(url.Values)
				req.Form.Set("campo1", "valor")
				req.Form.Set("campo2", "true")
				req.Form.Set("campo3", "-1")
				req.Form.Set("campo4", "-2")
				req.Form.Set("campo5", "-3")
				req.Form.Set("campo6", "-4")
				req.Form.Set("campo7", "-5")
				req.Form.Set("campo8", "6")
				req.Form.Set("campo9", "7")
				req.Form.Set("campo10", "8")
				req.Form.Set("campo11", "9")
				req.Form.Set("campo12", "10")
				req.Form.Set("campo13", "11.1")
				req.Form.Set("campo14", "12.2")
				req.Form.Set("campo15", "192.168.1.1")
				req.Form["campo16"] = []string{}
				req.Form["campoinexistente"] = []string{"teste"}
				return req
			}(),
			logger: &simulador.Logger{
				SimulaDebug: func(m ...interface{}) {
					mensagem := fmt.Sprint(m...)
					if mensagem != "Interceptador Antes: Parâmetros Consulta" {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
			},
			handlerEsperado: parâmetrosConsultaSimulado{
				Campo1:  "valor",
				Campo2:  true,
				Campo3:  -1,
				Campo4:  -2,
				Campo5:  -3,
				Campo6:  -4,
				Campo7:  -5,
				Campo8:  6,
				Campo9:  7,
				Campo10: 8,
				Campo11: 9,
				Campo12: 10,
				Campo13: 11.1,
				Campo14: 12.2,
				Campo15: net.ParseIP("192.168.1.1"),
			},
		},
		{
			descrição: "deve detectar um objeto que não consegue interpretar o argumento",
			requisição: func() *http.Request {
				req, err := http.NewRequest("GET", "/test", nil)
				if err != nil {
					t.Fatal(err)
				}

				req.Form = make(url.Values)
				req.Form.Set("campo16", "X")
				return req
			}(),
			logger: &simulador.Logger{
				SimulaErrorf: func(m string, a ...interface{}) {
					mensagem := fmt.Sprintf(m, a...)
					if mensagem != "tipo de valor não suportado: &struct {}{}" {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
				SimulaDebug: func(m ...interface{}) {
					mensagem := fmt.Sprint(m...)
					if mensagem != "Interceptador Antes: Parâmetros Consulta" {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
			},
			códigoHTTPEsperado: http.StatusBadRequest,
			handlerEsperado: parâmetrosConsultaSimulado{
				MensagensCompatível: interceptador.MensagensCompatível{
					Mensagens: protocolo.NovasMensagens(
						protocolo.NovaMensagemComCampo(protocolo.MensagemCódigoParâmetroInválido, "campo16", "X"),
					),
				},
			},
		},
		{
			descrição: "deve detectar quando o formulário ainda não foi interpretado",
			requisição: func() *http.Request {
				req, err := http.NewRequest("GET", "/test?campo1=valor&campo2=true&campo3=-1", nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			}(),
			logger: &simulador.Logger{
				SimulaDebug: func(m ...interface{}) {
					mensagem := fmt.Sprint(m...)
					if mensagem != "Interceptador Antes: Parâmetros Consulta" {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
			},
			handlerEsperado: parâmetrosConsultaSimulado{
				Campo1: "valor",
				Campo2: true,
				Campo3: -1,
			},
		},
	}

	for i, cenário := range cenários {
		var handler parâmetrosConsultaSimulado
		handler.SimulaRequisição = cenário.requisição
		handler.DefineLogger(cenário.logger)

		estrutura := interceptor.NewIntrospector(&handler)
		if códigoHTTP := estrutura.Before(); códigoHTTP != 0 {
			t.Errorf("Item %d, “%s”: código HTTP %d inesperado",
				i, cenário.descrição, códigoHTTP)
			continue
		}

		parâmetrosConsulta := interceptador.NovoParâmetrosConsulta(&handler)
		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)

		verificadorResultado.DefinirEsperado(cenário.códigoHTTPEsperado, nil)
		if err := verificadorResultado.VerificaResultado(parâmetrosConsulta.Before(), nil); err != nil {
			t.Error(err)
		}

		cenário.handlerEsperado.SimulaRequisição = handler.Req()
		cenário.handlerEsperado.IntrospectorCompliant = handler.IntrospectorCompliant
		cenário.handlerEsperado.LogCompatível = handler.LogCompatível

		verificadorResultado.DefinirEsperado(cenário.handlerEsperado, nil)
		if err := verificadorResultado.VerificaResultado(handler, nil); err != nil {
			t.Error(err)
		}
	}
}

type parâmetrosConsultaSimulado struct {
	interceptador.LogCompatível
	interceptador.MensagensCompatível
	interceptor.IntrospectorCompliant
	simulador.Handler

	Campo1  string   `query:"campo1"`
	Campo2  bool     `query:"campo2"`
	Campo3  int      `query:"campo3"`
	Campo4  int8     `query:"campo4"`
	Campo5  int16    `query:"campo5"`
	Campo6  int32    `query:"campo6"`
	Campo7  int64    `query:"campo7"`
	Campo8  uint     `query:"campo8"`
	Campo9  uint8    `query:"campo9"`
	Campo10 uint16   `query:"campo10"`
	Campo11 uint32   `query:"campo11"`
	Campo12 uint64   `query:"campo12"`
	Campo13 float32  `query:"campo13"`
	Campo14 float64  `query:"campo14"`
	Campo15 net.IP   `query:"campo15"`
	Campo16 struct{} `query:"campo16"`
}
