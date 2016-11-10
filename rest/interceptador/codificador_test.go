package interceptador_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/erros"
	"github.com/rafaeljusto/atiradorfrequente/rest/interceptador"
	"github.com/rafaeljusto/atiradorfrequente/testes"
	"github.com/rafaeljusto/atiradorfrequente/testes/simulador"
	"github.com/registrobr/gostk/errors"
	"github.com/registrobr/gostk/log"
	"github.com/trajber/handy/interceptor"
)

func TestCodificador_Before(t *testing.T) {
	cenários := []struct {
		descrição          string
		requisição         *http.Request
		logger             log.Logger
		tipoConteúdo       string
		códigoHTTPEsperado int
		handlerEsperado    codificadorSimulado
	}{
		{
			descrição: "deve preencher corretamente a estrutura de requisição no handler",
			requisição: func() *http.Request {
				requisição, err := http.NewRequest("POST", "https://exemplo.com.br/teste", strings.NewReader(`{
  "campo1": "valor1",
  "campo2": [ 1, 2, 3, 4, 5 ],
  "campo3": "ABCDEFGHIJ1234567890ABCDEFGHIJ1234567890ABCDEFGHIJ1234567890ABCDEFGHIJ1234567890ABCDEFGHIJ1234567890ABCDEFGHIJ1234567890"
}`))

				if err != nil {
					t.Fatalf("Erro ao criar a requisição. Detalhes: %s", err)
				}

				return requisição
			}(),
			logger: &simulador.Logger{
				SimulaDebug: func(m ...interface{}) {
					mensagem := fmt.Sprint(m...)
					if mensagem != "Interceptador Antes: Codificador" {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
				SimulaDebugf: func(m string, a ...interface{}) {
					mensagem := fmt.Sprintf(m, a...)
					if mensagem != `Requisição corpo: “{  "campo1": "valor1",  "campo2": [ 1, 2, 3, 4, 5 ],  "campo3": "ABCDEFGHIJ1234567890ABCDEFGHIJ1234567890ABCDEFGHIJ...1234567890ABCDEFGHIJ1234567890ABCDEFGHIJ1234567890"}”` {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
			},
			tipoConteúdo: "application/json",
			handlerEsperado: codificadorSimulado{
				Requisição: codificadorObjetoSimulada{
					Campo1: "valor1",
					Campo2: []int{1, 2, 3, 4, 5},
					Campo3: "ABCDEFGHIJ1234567890ABCDEFGHIJ1234567890ABCDEFGHIJ1234567890ABCDEFGHIJ1234567890ABCDEFGHIJ1234567890ABCDEFGHIJ1234567890",
				},
			},
		},
		{
			descrição: "deve ignorar quando não houver uma estrutura de requisição correspondente",
			requisição: func() *http.Request {
				requisição, err := http.NewRequest("PUT", "https://exemplo.com.br/teste", strings.NewReader(`{
  "campo1": "valor1",
  "campo2": [ 1, 2, 3, 4, 5 ]
}`))

				if err != nil {
					t.Fatalf("Erro ao criar a requisição. Detalhes: %s", err)
				}

				return requisição
			}(),
			logger: &simulador.Logger{
				SimulaDebug: func(m ...interface{}) {
					mensagem := fmt.Sprint(m...)
					if mensagem != "Interceptador Antes: Codificador" {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
			},
			tipoConteúdo: "application/json",
		},
		{
			descrição: "deve detectar um erro no JSON da requisição",
			requisição: func() *http.Request {
				requisição, err := http.NewRequest("POST", "https://exemplo.com.br/teste", strings.NewReader(`{
  "campo1": "valor1",
  "campo2": [ 1, 2, 3, 4, 5`))

				if err != nil {
					t.Fatalf("Erro ao criar a requisição. Detalhes: %s", err)
				}

				return requisição
			}(),
			logger: &simulador.Logger{
				SimulaError: func(err error) {
					esperado := io.ErrUnexpectedEOF

					if !errors.Equal(err, esperado) {
						t.Errorf("mensagem inesperada: %s", err)
					}
				},
				SimulaDebug: func(m ...interface{}) {
					mensagem := fmt.Sprint(m...)
					if mensagem != "Interceptador Antes: Codificador" {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
			},
			tipoConteúdo:       "application/json",
			códigoHTTPEsperado: http.StatusInternalServerError,
		},
	}

	for i, cenário := range cenários {
		var handler codificadorSimulado
		handler.SimulaRequisição = cenário.requisição
		handler.DefineLogger(cenário.logger)

		estrutura := interceptor.NewIntrospector(&handler)
		if códigoHTTP := estrutura.Before(); códigoHTTP != 0 {
			t.Errorf("Item %d, “%s”: código HTTP %d inesperado",
				i, cenário.descrição, códigoHTTP)
			continue
		}

		codificador := interceptador.NovoCodificador(&handler, cenário.tipoConteúdo)
		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)

		verificadorResultado.DefinirEsperado(cenário.códigoHTTPEsperado, nil)
		if err := verificadorResultado.VerificaResultado(codificador.Before(), nil); err != nil {
			t.Error(err)
		}

		cenário.handlerEsperado.SimulaRequisição = handler.SimulaRequisição
		cenário.handlerEsperado.IntrospectorCompliant = handler.IntrospectorCompliant
		cenário.handlerEsperado.LogCompatível = handler.LogCompatível

		verificadorResultado.DefinirEsperado(cenário.handlerEsperado, nil)
		if err := verificadorResultado.VerificaResultado(handler, nil); err != nil {
			t.Error(err)
		}
	}
}

func TestCodificador_After(t *testing.T) {
	cenários := []struct {
		descrição                  string
		handler                    codificadorSimuladoFlexível
		logger                     log.Logger
		tipoConteúdo               string
		códigoHTTP                 int
		códigoHTTPEsperado         int
		respostaCodificadaEsperada string
		cabeçalhoEsperado          http.Header
	}{
		{
			descrição: "deve escrever corretamente a resposta específica",
			handler: &codificadorSimulado{
				Handler: simulador.Handler{
					SimulaRequisição: func() *http.Request {
						requisição, err := http.NewRequest("GET", "https://exemplo.com.br/teste", nil)

						if err != nil {
							t.Fatalf("Erro ao criar a requisição. Detalhes: %s", err)
						}

						return requisição
					}(),
				},
				Resposta: &codificadorObjetoSimulada{
					Campo1: "valor1",
					Campo2: []int{1, 2, 3, 4, 5},
					Campo3: "ABCDEFGHIJ1234567890ABCDEFGHIJ1234567890ABCDEFGHIJ1234567890ABCDEFGHIJ1234567890ABCDEFGHIJ1234567890ABCDEFGHIJ1234567890",
				},
				CabeçalhoCompatível: interceptador.CabeçalhoCompatível{
					Cabeçalho: http.Header{
						"E-Tag": []string{"ABC123"},
					},
				},
			},
			logger: &simulador.Logger{
				SimulaDebug: func(m ...interface{}) {
					mensagem := fmt.Sprint(m...)
					if mensagem != "Interceptador Depois: Codificador" {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
				SimulaDebugf: func(m string, a ...interface{}) {
					mensagem := fmt.Sprintf(m, a...)
					if mensagem != `Resposta corpo: “{"campo1":"valor1","campo2":[1,2,3,4,5],"campo3":"ABCDEFGHIJ1234567890ABCDEFGHIJ1234567890ABCDEFGHIJ...1234567890ABCDEFGHIJ1234567890ABCDEFGHIJ1234567890"}”` {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
			},
			tipoConteúdo:               "application/json",
			códigoHTTP:                 http.StatusOK,
			códigoHTTPEsperado:         http.StatusOK,
			respostaCodificadaEsperada: `{"campo1":"valor1","campo2":[1,2,3,4,5],"campo3":"ABCDEFGHIJ1234567890ABCDEFGHIJ1234567890ABCDEFGHIJ1234567890ABCDEFGHIJ1234567890ABCDEFGHIJ1234567890ABCDEFGHIJ1234567890"}` + "\n",
			cabeçalhoEsperado: http.Header{
				"Content-Type": []string{"application/json"},
				"E-Tag":        []string{"ABC123"},
			},
		},
		{
			descrição: "deve ignorar quando a resposta é indefinida",
			handler: &codificadorSimulado{
				Handler: simulador.Handler{
					SimulaRequisição: func() *http.Request {
						requisição, err := http.NewRequest("GET", "https://exemplo.com.br/teste", nil)

						if err != nil {
							t.Fatalf("Erro ao criar a requisição. Detalhes: %s", err)
						}

						return requisição
					}(),
				},
			},
			logger: &simulador.Logger{
				SimulaDebug: func(m ...interface{}) {
					mensagem := fmt.Sprint(m...)
					if mensagem != "Interceptador Depois: Codificador" {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
			},
			tipoConteúdo:       "application/json",
			códigoHTTP:         http.StatusOK,
			códigoHTTPEsperado: http.StatusOK,
			cabeçalhoEsperado:  http.Header{},
		},
		{
			descrição: "deve identificar um cabeçalho no formato inválido",
			handler: &codificadorHeaderInválidoSimulado{
				Handler: simulador.Handler{
					SimulaRequisição: func() *http.Request {
						requisição, err := http.NewRequest("GET", "https://exemplo.com.br/teste", nil)

						if err != nil {
							t.Fatalf("Erro ao criar a requisição. Detalhes: %s", err)
						}

						return requisição
					}(),
				},
			},
			logger: &simulador.Logger{
				SimulaErrorf: func(m string, a ...interface{}) {
					mensagem := fmt.Sprintf(m, a...)
					if mensagem != "“Cabeçalho” campo com tipo errado: *int" {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
				SimulaDebug: func(m ...interface{}) {
					mensagem := fmt.Sprint(m...)
					if mensagem != "Interceptador Depois: Codificador" {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
			},
			tipoConteúdo:       "application/json",
			códigoHTTP:         http.StatusOK,
			códigoHTTPEsperado: http.StatusOK,
			cabeçalhoEsperado:  http.Header{},
		},
		{
			descrição: "deve escrever corretamente a resposta genérica",
			handler: &codificadorRespostaGenéricaSimulado{
				Handler: simulador.Handler{
					SimulaRequisição: func() *http.Request {
						requisição, err := http.NewRequest("GET", "https://exemplo.com.br/teste", nil)

						if err != nil {
							t.Fatalf("Erro ao criar a requisição. Detalhes: %s", err)
						}

						return requisição
					}(),
				},
				Resposta: codificadorObjetoGenéricoSimulado{"valor1", "valor2", "valor3"},
				CabeçalhoCompatível: interceptador.CabeçalhoCompatível{
					Cabeçalho: http.Header{
						"E-Tag": []string{"ABC123"},
					},
				},
			},
			logger: &simulador.Logger{
				SimulaDebug: func(m ...interface{}) {
					mensagem := fmt.Sprint(m...)
					if mensagem != "Interceptador Depois: Codificador" {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
				SimulaDebugf: func(m string, a ...interface{}) {
					mensagem := fmt.Sprintf(m, a...)
					if mensagem != `Resposta corpo: “["valor1","valor2","valor3"]”` {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
			},
			tipoConteúdo:               "application/json",
			códigoHTTP:                 http.StatusOK,
			códigoHTTPEsperado:         http.StatusOK,
			respostaCodificadaEsperada: `["valor1","valor2","valor3"]` + "\n",
			cabeçalhoEsperado: http.Header{
				"Content-Type": []string{"application/json"},
				"E-Tag":        []string{"ABC123"},
			},
		},
		{
			descrição: "deve ignorar a resposta genérica quando não definida",
			handler: &codificadorRespostaGenéricaSimulado{
				Handler: simulador.Handler{
					SimulaRequisição: func() *http.Request {
						requisição, err := http.NewRequest("GET", "https://exemplo.com.br/teste", nil)

						if err != nil {
							t.Fatalf("Erro ao criar a requisição. Detalhes: %s", err)
						}

						return requisição
					}(),
				},
			},
			logger: &simulador.Logger{
				SimulaDebug: func(m ...interface{}) {
					mensagem := fmt.Sprint(m...)
					if mensagem != "Interceptador Depois: Codificador" {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
				SimulaDebugf: func(m string, a ...interface{}) {
					mensagem := fmt.Sprintf(m, a...)
					if mensagem != `Resposta corpo: “”` {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
			},
			tipoConteúdo:       "application/json",
			códigoHTTP:         http.StatusInternalServerError,
			códigoHTTPEsperado: http.StatusInternalServerError,
			cabeçalhoEsperado:  http.Header{},
		},
		{
			descrição: "deve detectar um erro ao codificar a resposta",
			handler: &codificadorRespostaInválidaSimulado{
				Handler: simulador.Handler{
					SimulaRequisição: func() *http.Request {
						requisição, err := http.NewRequest("GET", "https://exemplo.com.br/teste", nil)

						if err != nil {
							t.Fatalf("Erro ao criar a requisição. Detalhes: %s", err)
						}

						return requisição
					}(),
				},
				Resposta: &codificadorObjetoInválidoSimulada{
					Campo1: "valor1",
					Campo2: []int{1, 2, 3, 4, 5},
				},
				CabeçalhoCompatível: interceptador.CabeçalhoCompatível{
					Cabeçalho: http.Header{
						"E-Tag": []string{"ABC123"},
					},
				},
			},
			logger: &simulador.Logger{
				SimulaError: func(err error) {
					esperado := erros.Novo(&json.MarshalerError{
						Type: reflect.TypeOf(new(codificadorObjetoInválidoSimulada)),
						Err:  fmt.Errorf("erro de codificação"),
					})

					if !errors.Equal(err, esperado) {
						t.Errorf("mensagem inesperada: %s", err)
					}
				},
				SimulaDebug: func(m ...interface{}) {
					mensagem := fmt.Sprint(m...)
					if mensagem != "Interceptador Depois: Codificador" {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
			},
			tipoConteúdo:       "application/json",
			códigoHTTP:         http.StatusOK,
			códigoHTTPEsperado: http.StatusOK,
			cabeçalhoEsperado: http.Header{
				"Content-Type": []string{"application/json"},
				"E-Tag":        []string{"ABC123"},
			},
		},
	}

	for i, cenário := range cenários {
		gravadorResposta := httptest.NewRecorder()

		cenário.handler.DefineLogger(cenário.logger)
		cenário.handler.DefineResposta(gravadorResposta)

		estrutura := interceptor.NewIntrospector(cenário.handler)
		if códigoHTTP := estrutura.Before(); códigoHTTP != 0 {
			t.Errorf("Item %d, “%s”: código HTTP %d inesperado",
				i, cenário.descrição, códigoHTTP)
			continue
		}

		codificador := interceptador.NovoCodificador(cenário.handler, cenário.tipoConteúdo)
		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)

		verificadorResultado.DefinirEsperado(cenário.códigoHTTPEsperado, nil)
		if err := verificadorResultado.VerificaResultado(codificador.After(cenário.códigoHTTP), nil); err != nil {
			t.Error(err)
		}

		verificadorResultado.DefinirEsperado(cenário.respostaCodificadaEsperada, nil)
		if err := verificadorResultado.VerificaResultado(gravadorResposta.Body.String(), nil); err != nil {
			t.Error(err)
		}

		verificadorResultado.DefinirEsperado(cenário.cabeçalhoEsperado, nil)
		if err := verificadorResultado.VerificaResultado(gravadorResposta.Header(), nil); err != nil {
			t.Error(err)
		}
	}
}

type codificadorSimuladoFlexível interface {
	DefineResposta(http.ResponseWriter)
	SetFields(interceptor.StructFields)
	Field(tag, value string) interface{}
	DefineLogger(log.Logger)
	Logger() log.Logger
	Req() *http.Request
	ResponseWriter() http.ResponseWriter
}

type codificadorSimulado struct {
	interceptador.LogCompatível
	interceptor.IntrospectorCompliant
	interceptador.CabeçalhoCompatível
	simulador.Handler

	Requisição codificadorObjetoSimulada  `request:"post"`
	Resposta   *codificadorObjetoSimulada `response:"get"`
}

func (c *codificadorSimulado) DefineResposta(w http.ResponseWriter) {
	c.SimulaResposta = w
}

type codificadorHeaderInválidoSimulado struct {
	interceptador.LogCompatível
	interceptor.IntrospectorCompliant
	simulador.Handler

	Requisição codificadorObjetoSimulada  `request:"post"`
	Resposta   *codificadorObjetoSimulada `response:"get"`
	Cabeçalho  int                        `response:"header"`
}

func (c *codificadorHeaderInválidoSimulado) DefineResposta(w http.ResponseWriter) {
	c.SimulaResposta = w
}

type codificadorRespostaGenéricaSimulado struct {
	interceptador.LogCompatível
	interceptor.IntrospectorCompliant
	interceptador.CabeçalhoCompatível
	simulador.Handler

	Requisição codificadorObjetoSimulada         `request:"post"`
	Resposta   codificadorObjetoGenéricoSimulado `response:"all"`
}

func (c *codificadorRespostaGenéricaSimulado) DefineResposta(w http.ResponseWriter) {
	c.SimulaResposta = w
}

type codificadorRespostaInválidaSimulado struct {
	interceptador.LogCompatível
	interceptor.IntrospectorCompliant
	interceptador.CabeçalhoCompatível
	simulador.Handler

	Requisição codificadorObjetoSimulada          `request:"post"`
	Resposta   *codificadorObjetoInválidoSimulada `response:"get"`
}

func (c *codificadorRespostaInválidaSimulado) DefineResposta(w http.ResponseWriter) {
	c.SimulaResposta = w
}

type codificadorObjetoSimulada struct {
	Campo1 string `json:"campo1"`
	Campo2 []int  `json:"campo2"`
	Campo3 string `json:"campo3,omitempty"`
}

type codificadorObjetoGenéricoSimulado []string

type codificadorObjetoInválidoSimulada struct {
	Campo1 string `json:"campo1"`
	Campo2 []int  `json:"campo2"`
}

func (c codificadorObjetoInválidoSimulada) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("erro de codificação")
}
