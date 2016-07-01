package interceptador_test

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"testing"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
	"github.com/rafaeljusto/atiradorfrequente/rest/interceptador"
	"github.com/rafaeljusto/atiradorfrequente/testes"
	"github.com/rafaeljusto/atiradorfrequente/testes/simulador"
	"github.com/registrobr/gostk/log"
	"github.com/trajber/handy"
	"github.com/trajber/handy/interceptor"
)

func TestVariáveisEndereço_Before(t *testing.T) {
	cenários := []struct {
		descrição          string
		variáveisEndereço  handy.URIVars
		logger             log.Logger
		códigoHTTPEsperado int
		handlerEsperado    variáveisEndereçoSimulado
	}{
		{
			descrição: "deve preencher corretamente as variáveis do handler",
			variáveisEndereço: handy.URIVars{
				"campo1":  "valor",
				"campo2":  "true",
				"campo3":  "-1",
				"campo4":  "-2",
				"campo5":  "-3",
				"campo6":  "-4",
				"campo7":  "-5",
				"campo8":  "6",
				"campo9":  "7",
				"campo10": "8",
				"campo11": "9",
				"campo12": "10",
				"campo13": "11.1",
				"campo14": "12.2",
				"campo15": "192.168.1.1",
			},
			logger: &simulador.Logger{
				SimulaDebug: func(m ...interface{}) {
					mensagem := fmt.Sprint(m...)
					if mensagem != "Interceptador Antes: Variáveis Endereço" {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
			},
			handlerEsperado: variáveisEndereçoSimulado{
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
			descrição: "deve detectar quando o campo não existe",
			variáveisEndereço: handy.URIVars{
				"campoX": "valor",
			},
			logger: &simulador.Logger{
				SimulaWarningf: func(m string, a ...interface{}) {
					mensagem := fmt.Sprintf(m, a...)
					if mensagem != "Tentando definir um valor no campo “campoX” que não existe" {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
				SimulaDebug: func(m ...interface{}) {
					mensagem := fmt.Sprint(m...)
					if mensagem != "Interceptador Antes: Variáveis Endereço" {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
			},
		},
		{
			descrição: "deve detectar um número inválido",
			variáveisEndereço: handy.URIVars{
				"campo3": "X",
			},
			logger: &simulador.Logger{
				SimulaError: func(err error) {
					esperado := strconv.NumError{
						Func: "ParseInt",
						Num:  "X",
						Err:  strconv.ErrSyntax,
					}

					if err.Error() != esperado.Error() {
						t.Errorf("mensagem inesperada: %s", err)
					}
				},
				SimulaDebug: func(m ...interface{}) {
					mensagem := fmt.Sprint(m...)
					if mensagem != "Interceptador Antes: Variáveis Endereço" {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
			},
			códigoHTTPEsperado: http.StatusBadRequest,
			handlerEsperado: variáveisEndereçoSimulado{
				MensagensCompatível: interceptador.MensagensCompatível{
					Mensagens: protocolo.NovasMensagens(
						protocolo.NovaMensagemComCampo(protocolo.MensagemCódigoParâmetroInválido, "campo3", "X"),
					),
				},
			},
		},
		{
			descrição: "deve detectar um número sem sinal inválido",
			variáveisEndereço: handy.URIVars{
				"campo8": "X",
			},
			logger: &simulador.Logger{
				SimulaError: func(err error) {
					esperado := strconv.NumError{
						Func: "ParseUint",
						Num:  "X",
						Err:  strconv.ErrSyntax,
					}

					if err.Error() != esperado.Error() {
						t.Errorf("mensagem inesperada: %s", err)
					}
				},
				SimulaDebug: func(m ...interface{}) {
					mensagem := fmt.Sprint(m...)
					if mensagem != "Interceptador Antes: Variáveis Endereço" {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
			},
			códigoHTTPEsperado: http.StatusBadRequest,
			handlerEsperado: variáveisEndereçoSimulado{
				MensagensCompatível: interceptador.MensagensCompatível{
					Mensagens: protocolo.NovasMensagens(
						protocolo.NovaMensagemComCampo(protocolo.MensagemCódigoParâmetroInválido, "campo8", "X"),
					),
				},
			},
		},
		{
			descrição: "deve detectar um número real inválido",
			variáveisEndereço: handy.URIVars{
				"campo13": "X",
			},
			logger: &simulador.Logger{
				SimulaError: func(err error) {
					esperado := strconv.NumError{
						Func: "ParseFloat",
						Num:  "X",
						Err:  strconv.ErrSyntax,
					}

					if err.Error() != esperado.Error() {
						t.Errorf("mensagem inesperada: %s", err)
					}
				},
				SimulaDebug: func(m ...interface{}) {
					mensagem := fmt.Sprint(m...)
					if mensagem != "Interceptador Antes: Variáveis Endereço" {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
			},
			códigoHTTPEsperado: http.StatusBadRequest,
			handlerEsperado: variáveisEndereçoSimulado{
				MensagensCompatível: interceptador.MensagensCompatível{
					Mensagens: protocolo.NovasMensagens(
						protocolo.NovaMensagemComCampo(protocolo.MensagemCódigoParâmetroInválido, "campo13", "X"),
					),
				},
			},
		},
		{
			descrição: "deve detectar um objeto que não consegue interpretar o argumento",
			variáveisEndereço: handy.URIVars{
				"campo16": "X",
			},
			logger: &simulador.Logger{
				SimulaErrorf: func(m string, a ...interface{}) {
					mensagem := fmt.Sprintf(m, a...)
					if mensagem != "tipo de valor não suportado: &struct {}{}" {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
				SimulaDebug: func(m ...interface{}) {
					mensagem := fmt.Sprint(m...)
					if mensagem != "Interceptador Antes: Variáveis Endereço" {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
			},
			códigoHTTPEsperado: http.StatusBadRequest,
			handlerEsperado: variáveisEndereçoSimulado{
				MensagensCompatível: interceptador.MensagensCompatível{
					Mensagens: protocolo.NovasMensagens(
						protocolo.NovaMensagemComCampo(protocolo.MensagemCódigoParâmetroInválido, "campo16", "X"),
					),
				},
			},
		},
		{
			descrição: "deve detectar um objeto inválido",
			variáveisEndereço: handy.URIVars{
				"campo15": "X.X.X.X",
			},
			logger: &simulador.Logger{
				SimulaError: func(err error) {
					esperado := net.ParseError{
						Type: "IP address",
						Text: "X.X.X.X",
					}

					if err.Error() != esperado.Error() {
						t.Errorf("mensagem inesperada: %s", err)
					}
				},
				SimulaDebug: func(m ...interface{}) {
					mensagem := fmt.Sprint(m...)
					if mensagem != "Interceptador Antes: Variáveis Endereço" {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
			},
			códigoHTTPEsperado: http.StatusBadRequest,
			handlerEsperado: variáveisEndereçoSimulado{
				MensagensCompatível: interceptador.MensagensCompatível{
					Mensagens: protocolo.NovasMensagens(
						protocolo.NovaMensagemComCampo(protocolo.MensagemCódigoParâmetroInválido, "campo15", "X.X.X.X"),
					),
				},
			},
		},
		{
			descrição: "deve detectar um objeto inválido com mensagem descritiva",
			variáveisEndereço: handy.URIVars{
				"campo17": "X",
			},
			logger: &simulador.Logger{
				SimulaError: func(err error) {
					esperado := protocolo.NovasMensagens(
						protocolo.NovaMensagem(protocolo.MensagemCódigoParâmetroInválido),
					)

					if err.Error() != esperado.Error() {
						t.Errorf("mensagem inesperada: %s", err)
					}
				},
				SimulaDebug: func(m ...interface{}) {
					mensagem := fmt.Sprint(m...)
					if mensagem != "Interceptador Antes: Variáveis Endereço" {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
			},
			códigoHTTPEsperado: http.StatusBadRequest,
			handlerEsperado: variáveisEndereçoSimulado{
				MensagensCompatível: interceptador.MensagensCompatível{
					Mensagens: protocolo.NovasMensagens(
						protocolo.NovaMensagem(protocolo.MensagemCódigoParâmetroInválido),
					),
				},
			},
		},
	}

	for i, cenário := range cenários {
		var handler variáveisEndereçoSimulado
		handler.SimulaVariáveisEndereço = cenário.variáveisEndereço
		handler.DefineLogger(cenário.logger)

		estrutura := interceptor.NewIntrospector(&handler)
		if códigoHTTP := estrutura.Before(); códigoHTTP != 0 {
			t.Errorf("Item %d, “%s”: código HTTP %d inesperado",
				i, cenário.descrição, códigoHTTP)
			continue
		}

		váriaveisEndereço := interceptador.NovaVariáveisEndereço(&handler)
		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)

		verificadorResultado.DefinirEsperado(cenário.códigoHTTPEsperado, nil)
		if err := verificadorResultado.VerificaResultado(váriaveisEndereço.Before(), nil); err != nil {
			t.Error(err)
		}

		cenário.handlerEsperado.SimulaVariáveisEndereço = handler.SimulaVariáveisEndereço
		cenário.handlerEsperado.IntrospectorCompliant = handler.IntrospectorCompliant
		cenário.handlerEsperado.LogCompatível = handler.LogCompatível

		verificadorResultado.DefinirEsperado(cenário.handlerEsperado, nil)
		if err := verificadorResultado.VerificaResultado(handler, nil); err != nil {
			t.Error(err)
		}
	}
}

type variáveisEndereçoSimulado struct {
	interceptador.LogCompatível
	interceptador.MensagensCompatível
	interceptor.IntrospectorCompliant
	simulador.Handler

	Campo1  string                         `urivar:"campo1"`
	Campo2  bool                           `urivar:"campo2"`
	Campo3  int                            `urivar:"campo3"`
	Campo4  int8                           `urivar:"campo4"`
	Campo5  int16                          `urivar:"campo5"`
	Campo6  int32                          `urivar:"campo6"`
	Campo7  int64                          `urivar:"campo7"`
	Campo8  uint                           `urivar:"campo8"`
	Campo9  uint8                          `urivar:"campo9"`
	Campo10 uint16                         `urivar:"campo10"`
	Campo11 uint32                         `urivar:"campo11"`
	Campo12 uint64                         `urivar:"campo12"`
	Campo13 float32                        `urivar:"campo13"`
	Campo14 float64                        `urivar:"campo14"`
	Campo15 net.IP                         `urivar:"campo15"`
	Campo16 struct{}                       `urivar:"campo16"`
	Campo17 variáveisEndereçoCampoSimulado `urivar:"campo17"`
}

type variáveisEndereçoCampoSimulado struct {
}

func (v *variáveisEndereçoCampoSimulado) UnmarshalText([]byte) error {
	return protocolo.NovasMensagens(
		protocolo.NovaMensagem(protocolo.MensagemCódigoParâmetroInválido),
	)
}
