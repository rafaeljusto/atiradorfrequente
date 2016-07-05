package handler

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/atirador"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/bd"
	configNúcleo "github.com/rafaeljusto/atiradorfrequente/núcleo/config"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
	configREST "github.com/rafaeljusto/atiradorfrequente/rest/config"
	"github.com/rafaeljusto/atiradorfrequente/testes"
	"github.com/rafaeljusto/atiradorfrequente/testes/simulador"
	"github.com/registrobr/gostk/errors"
	"github.com/registrobr/gostk/log"
)

func TestFrequênciaAtiradorConfirmação_Put(t *testing.T) {
	cenários := []struct {
		descrição                   string
		cr                          string
		númeroControle              protocolo.NúmeroControle
		frequênciaConfirmaçãoPedido protocolo.FrequênciaConfirmaçãoPedido
		logger                      log.Logger
		configuração                *configREST.Configuração
		serviçoAtirador             atirador.Serviço
		códigoHTTPEsperado          int
	}{
		{
			descrição:      "deve confirmar corretamente os dados de frequência do atirador",
			cr:             "123456789",
			númeroControle: protocolo.NovoNúmeroControle(7654, 918273645),
			frequênciaConfirmaçãoPedido: protocolo.FrequênciaConfirmaçãoPedido{
				Imagem: `TWFuIGlzIGRpc3Rpbmd1aXNoZWQsIG5vdCBvbmx5IGJ5IGhpcyByZWFzb24sIGJ1dCBieSB0aGlz
IHNpbmd1bGFyIHBhc3Npb24gZnJvbSBvdGhlciBhbmltYWxzLCB3aGljaCBpcyBhIGx1c3Qgb2Yg
dGhlIG1pbmQsIHRoYXQgYnkgYSBwZXJzZXZlcmFuY2Ugb2YgZGVsaWdodCBpbiB0aGUgY29udGlu
dWVkIGFuZCBpbmRlZmF0aWdhYmxlIGdlbmVyYXRpb24gb2Yga25vd2xlZGdlLCBleGNlZWRzIHRo
ZSBzaG9ydCB2ZWhlbWVuY2Ugb2YgYW55IGNhcm5hbCBwbGVhc3VyZS4=`,
			},
			configuração: func() *configREST.Configuração {
				return new(configREST.Configuração)
			}(),
			serviçoAtirador: simulador.ServiçoAtirador{
				SimulaConfirmarFrequência: func(frequênciaConfirmaçãoPedidoCompleta protocolo.FrequênciaConfirmaçãoPedidoCompleta) error {
					return nil
				},
			},
			códigoHTTPEsperado: http.StatusNoContent,
		},
		{
			descrição:      "deve detectar quando a configuração não foi inicializada",
			cr:             "123456789",
			númeroControle: protocolo.NovoNúmeroControle(7654, 918273645),
			frequênciaConfirmaçãoPedido: protocolo.FrequênciaConfirmaçãoPedido{
				Imagem: `TWFuIGlzIGRpc3Rpbmd1aXNoZWQsIG5vdCBvbmx5IGJ5IGhpcyByZWFzb24sIGJ1dCBieSB0aGlz
IHNpbmd1bGFyIHBhc3Npb24gZnJvbSBvdGhlciBhbmltYWxzLCB3aGljaCBpcyBhIGx1c3Qgb2Yg
dGhlIG1pbmQsIHRoYXQgYnkgYSBwZXJzZXZlcmFuY2Ugb2YgZGVsaWdodCBpbiB0aGUgY29udGlu
dWVkIGFuZCBpbmRlZmF0aWdhYmxlIGdlbmVyYXRpb24gb2Yga25vd2xlZGdlLCBleGNlZWRzIHRo
ZSBzaG9ydCB2ZWhlbWVuY2Ugb2YgYW55IGNhcm5hbCBwbGVhc3VyZS4=`,
			},
			logger: simulador.Logger{
				SimulaCrit: func(m ...interface{}) {
					mensagem := fmt.Sprint(m...)
					if mensagem != "Não existe configuração definida para atender a requisição" {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
				SimulaError: func(e error) {
					if !strings.HasSuffix(e.Error(), "erro de baixo nível") {
						t.Error("não está adicionando o erro correto ao log")
					}
				},
			},
			códigoHTTPEsperado: http.StatusInternalServerError,
		},
		{
			descrição:      "deve detectar um erro na camada de serviço do atirador",
			cr:             "123456789",
			númeroControle: protocolo.NovoNúmeroControle(7654, 918273645),
			frequênciaConfirmaçãoPedido: protocolo.FrequênciaConfirmaçãoPedido{
				Imagem: `TWFuIGlzIGRpc3Rpbmd1aXNoZWQsIG5vdCBvbmx5IGJ5IGhpcyByZWFzb24sIGJ1dCBieSB0aGlz
IHNpbmd1bGFyIHBhc3Npb24gZnJvbSBvdGhlciBhbmltYWxzLCB3aGljaCBpcyBhIGx1c3Qgb2Yg
dGhlIG1pbmQsIHRoYXQgYnkgYSBwZXJzZXZlcmFuY2Ugb2YgZGVsaWdodCBpbiB0aGUgY29udGlu
dWVkIGFuZCBpbmRlZmF0aWdhYmxlIGdlbmVyYXRpb24gb2Yga25vd2xlZGdlLCBleGNlZWRzIHRo
ZSBzaG9ydCB2ZWhlbWVuY2Ugb2YgYW55IGNhcm5hbCBwbGVhc3VyZS4=`,
			},
			logger: simulador.Logger{
				SimulaError: func(e error) {
					if !strings.HasSuffix(e.Error(), "erro de baixo nível") {
						t.Error("não está adicionando o erro correto ao log")
					}
				},
			},
			configuração: func() *configREST.Configuração {
				return new(configREST.Configuração)
			}(),
			serviçoAtirador: simulador.ServiçoAtirador{
				SimulaConfirmarFrequência: func(frequênciaConfirmaçãoPedidoCompleta protocolo.FrequênciaConfirmaçãoPedidoCompleta) error {
					return errors.Errorf("erro de baixo nível")
				},
			},
			códigoHTTPEsperado: http.StatusInternalServerError,
		},
	}

	configuraçãoOriginal := configREST.Atual()
	defer func() {
		configREST.AtualizarConfiguração(configuraçãoOriginal)
	}()

	serviçoAtiradorOriginal := atirador.NovoServiço
	defer func() {
		atirador.NovoServiço = serviçoAtiradorOriginal
	}()

	for i, cenário := range cenários {
		configREST.AtualizarConfiguração(cenário.configuração)

		atirador.NovoServiço = func(s *bd.SQLogger, configuração configNúcleo.Configuração) atirador.Serviço {
			return cenário.serviçoAtirador
		}

		handler := frequênciaAtiradorConfirmação{
			CR:                          cenário.cr,
			NúmeroControle:              cenário.númeroControle,
			FrequênciaConfirmaçãoPedido: cenário.frequênciaConfirmaçãoPedido,
		}
		handler.DefineLogger(cenário.logger)

		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(cenário.códigoHTTPEsperado, nil)
		if err := verificadorResultado.VerificaResultado(handler.Put(), nil); err != nil {
			t.Error(err)
		}
	}
}

func TestFrequênciaAtiradorConfirmação_Interceptors(t *testing.T) {
	esperado := []string{
		"*interceptador.EndereçoRemoto",
		"*interceptador.Log",
		"*interceptor.Introspector",
		"*interceptador.VariáveisEndereço",
		"*interceptador.BD",
	}

	var handler frequênciaAtiradorConfirmação

	verificadorResultado := testes.NovoVerificadorResultados("deve conter os interceptadores corretos", 0)
	verificadorResultado.DefinirEsperado(esperado, nil)
	if err := verificadorResultado.VerificaResultado(testes.TiposDaLista(handler.Interceptors()), nil); err != nil {
		t.Error(err)
	}
}
