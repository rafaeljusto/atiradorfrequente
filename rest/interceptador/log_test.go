package interceptador_test

import (
	"fmt"
	"net"
	"net/http"
	"regexp"
	"testing"

	"github.com/rafaeljusto/atiradorfrequente/rest/interceptador"
	"github.com/rafaeljusto/atiradorfrequente/testes"
	"github.com/rafaeljusto/atiradorfrequente/testes/simulador"
	"github.com/registrobr/gostk/log"
)

func TestLog_Before(t *testing.T) {
	cenários := []struct {
		descrição          string
		endereçoRemoto     net.IP
		ação               func(log.Logger)
		logger             func(id string) log.Logger
		códigoHTTPEsperado int
	}{
		{
			descrição:      "deve inicializar e escrever corretamente no log",
			endereçoRemoto: net.ParseIP("192.168.1.1"),
			ação: func(l log.Logger) {
				l.Debug("Teste")
			},
			logger: func(id string) log.Logger {
				if !regexp.MustCompile(`192\.168\.1\.1 [0-9]+`).MatchString(id) {
					t.Errorf("id do logger incorreto: %s", id)
				}

				return simulador.Logger{
					SimulaInfof: func(m string, a ...interface{}) {
						mensagem := fmt.Sprintf(m, a...)
						if mensagem != "Requisicao GET /teste" {
							t.Errorf("mensagem inesperada: %s", mensagem)
						}
					},
					SimulaDebug: func(m ...interface{}) {
						mensagem := fmt.Sprint(m...)
						if mensagem != "Teste" {
							t.Errorf("mensagem inesperada: %s", mensagem)
						}
					},
				}
			},
		},
	}

	loggerOriginal := log.NewLogger
	defer func() {
		log.NewLogger = loggerOriginal
	}()

	for i, cenário := range cenários {
		log.NewLogger = cenário.logger

		requisição, err := http.NewRequest("GET", "/teste", nil)
		if err != nil {
			t.Fatal(err)
		}

		// executa manualmente o processamento da requisição no servidor
		requisição.RequestURI = requisição.URL.RequestURI()

		handler := &logSimulado{}
		handler.SimulaRequisição = requisição
		handler.DefineEndereçoRemoto(cenário.endereçoRemoto)

		l := interceptador.NovoLog(handler)

		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(cenário.códigoHTTPEsperado, nil)
		if err := verificadorResultado.VerificaResultado(l.Before(), nil); err != nil {
			t.Error(err)
		}

		cenário.ação(handler.Logger())
	}
}

func TestLog_After(t *testing.T) {
	cenários := []struct {
		descrição          string
		códigoHTTP         int
		logger             simulador.Logger
		códigoHTTPEsperado int
	}{
		{
			descrição:  "deve inicializar e escrever corretamente no log",
			códigoHTTP: http.StatusOK,
			logger: simulador.Logger{
				SimulaInfof: func(m string, a ...interface{}) {
					mensagem := fmt.Sprintf(m, a...)
					if mensagem != "Resposta GET /teste 200 OK" {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
			},
			códigoHTTPEsperado: http.StatusOK,
		},
	}

	for i, cenário := range cenários {
		requisição, err := http.NewRequest("GET", "/teste", nil)
		if err != nil {
			t.Fatal(err)
		}

		// executa manualmente o processamento da requisição no servidor
		requisição.RequestURI = requisição.URL.RequestURI()

		handler := &logSimulado{}
		handler.SimulaRequisição = requisição
		handler.DefineLogger(cenário.logger)

		l := interceptador.NovoLog(handler)

		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(cenário.códigoHTTPEsperado, nil)
		if err := verificadorResultado.VerificaResultado(l.After(cenário.códigoHTTP), nil); err != nil {
			t.Error(err)
		}
	}
}

type logSimulado struct {
	interceptador.EndereçoRemotoCompatível
	interceptador.LogCompatível
	simulador.Handler
}
