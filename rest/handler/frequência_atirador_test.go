package handler

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

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

func TestFrequênciaAtirador_Post(t *testing.T) {
	data := time.Now()

	cenários := []struct {
		descrição          string
		cr                 string
		frequênciaPedido   protocolo.FrequênciaPedido
		logger             log.Logger
		configuração       *configREST.Configuração
		serviçoAtirador    atirador.Serviço
		códigoHTTPEsperado int
		esperado           *protocolo.FrequênciaPendenteResposta
		mensagensEsperadas protocolo.Mensagens
		cabeçalhoEsperado  http.Header
	}{
		{
			descrição: "deve cadastrar corretamente os dados de frequência do atirador",
			cr:        "123456789",
			frequênciaPedido: protocolo.FrequênciaPedido{
				Calibre:           ".380",
				ArmaUtilizada:     "Arma do Clube",
				QuantidadeMunição: 50,
				DataInício:        data,
				DataTérmino:       data.Add(30 * time.Minute),
			},
			configuração: func() *configREST.Configuração {
				return new(configREST.Configuração)
			}(),
			serviçoAtirador: simulador.ServiçoAtirador{
				SimulaCadastrarFrequência: func(frequênciaPedidoCompleta protocolo.FrequênciaPedidoCompleta) (protocolo.FrequênciaPendenteResposta, error) {
					return protocolo.FrequênciaPendenteResposta{
						NúmeroControle: protocolo.NovoNúmeroControle(7654, 918273645),
						Imagem: `TWFuIGlzIGRpc3Rpbmd1aXNoZWQsIG5vdCBvbmx5IGJ5IGhpcyByZWFzb24sIGJ1dCBieSB0aGlz
IHNpbmd1bGFyIHBhc3Npb24gZnJvbSBvdGhlciBhbmltYWxzLCB3aGljaCBpcyBhIGx1c3Qgb2Yg
dGhlIG1pbmQsIHRoYXQgYnkgYSBwZXJzZXZlcmFuY2Ugb2YgZGVsaWdodCBpbiB0aGUgY29udGlu
dWVkIGFuZCBpbmRlZmF0aWdhYmxlIGdlbmVyYXRpb24gb2Yga25vd2xlZGdlLCBleGNlZWRzIHRo
ZSBzaG9ydCB2ZWhlbWVuY2Ugb2YgYW55IGNhcm5hbCBwbGVhc3VyZS4=`,
					}, nil
				},
			},
			códigoHTTPEsperado: http.StatusCreated,
			esperado: &protocolo.FrequênciaPendenteResposta{
				NúmeroControle: protocolo.NovoNúmeroControle(7654, 918273645),
				Imagem: `TWFuIGlzIGRpc3Rpbmd1aXNoZWQsIG5vdCBvbmx5IGJ5IGhpcyByZWFzb24sIGJ1dCBieSB0aGlz
IHNpbmd1bGFyIHBhc3Npb24gZnJvbSBvdGhlciBhbmltYWxzLCB3aGljaCBpcyBhIGx1c3Qgb2Yg
dGhlIG1pbmQsIHRoYXQgYnkgYSBwZXJzZXZlcmFuY2Ugb2YgZGVsaWdodCBpbiB0aGUgY29udGlu
dWVkIGFuZCBpbmRlZmF0aWdhYmxlIGdlbmVyYXRpb24gb2Yga25vd2xlZGdlLCBleGNlZWRzIHRo
ZSBzaG9ydCB2ZWhlbWVuY2Ugb2YgYW55IGNhcm5hbCBwbGVhc3VyZS4=`,
			},
			cabeçalhoEsperado: http.Header{
				"Location": []string{"/frequencia-atirador/123456789/7654-918273645"},
			},
		},
		{
			descrição: "deve detectar quando a configuração não foi inicializada",
			cr:        "123456789",
			frequênciaPedido: protocolo.FrequênciaPedido{
				Calibre:           ".380",
				ArmaUtilizada:     "Arma do Clube",
				QuantidadeMunição: 50,
				DataInício:        data,
				DataTérmino:       data.Add(30 * time.Minute),
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
			descrição: "deve detectar um erro na camada de serviço do atirador",
			cr:        "123456789",
			frequênciaPedido: protocolo.FrequênciaPedido{
				Calibre:           ".380",
				ArmaUtilizada:     "Arma do Clube",
				QuantidadeMunição: 50,
				DataInício:        data,
				DataTérmino:       data.Add(30 * time.Minute),
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
				SimulaCadastrarFrequência: func(frequênciaPedidoCompleta protocolo.FrequênciaPedidoCompleta) (protocolo.FrequênciaPendenteResposta, error) {
					return protocolo.FrequênciaPendenteResposta{}, errors.Errorf("erro de baixo nível")
				},
			},
			códigoHTTPEsperado: http.StatusInternalServerError,
		},
		{
			descrição: "deve detectar mensagens na camada de serviço do atirador",
			cr:        "123456789",
			frequênciaPedido: protocolo.FrequênciaPedido{
				Calibre:           ".380",
				ArmaUtilizada:     "Arma do Clube",
				QuantidadeMunição: 50,
				DataInício:        data,
				DataTérmino:       data.Add(30 * time.Minute),
			},
			logger: simulador.Logger{},
			configuração: func() *configREST.Configuração {
				return new(configREST.Configuração)
			}(),
			serviçoAtirador: simulador.ServiçoAtirador{
				SimulaCadastrarFrequência: func(frequênciaPedidoCompleta protocolo.FrequênciaPedidoCompleta) (protocolo.FrequênciaPendenteResposta, error) {
					return protocolo.FrequênciaPendenteResposta{}, protocolo.NovasMensagens(
						protocolo.NovaMensagem(protocolo.MensagemCódigoCRInválido),
					)
				},
			},
			códigoHTTPEsperado: http.StatusBadRequest,
			mensagensEsperadas: protocolo.NovasMensagens(
				protocolo.NovaMensagem(protocolo.MensagemCódigoCRInválido),
			),
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

		handler := frequênciaAtirador{
			CR:               cenário.cr,
			FrequênciaPedido: cenário.frequênciaPedido,
		}
		handler.DefineLogger(cenário.logger)

		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)

		verificadorResultado.DefinirEsperado(cenário.códigoHTTPEsperado, nil)
		if err := verificadorResultado.VerificaResultado(handler.Post(), nil); err != nil {
			t.Error(err)
		}

		verificadorResultado.DefinirEsperado(cenário.esperado, nil)
		if err := verificadorResultado.VerificaResultado(handler.FrequênciaPendenteResposta, nil); err != nil {
			t.Error(err)
		}

		verificadorResultado.DefinirEsperado(cenário.mensagensEsperadas, nil)
		if err := verificadorResultado.VerificaResultado(handler.Mensagens, nil); err != nil {
			t.Error(err)
		}

		verificadorResultado.DefinirEsperado(cenário.cabeçalhoEsperado, nil)
		if err := verificadorResultado.VerificaResultado(handler.Cabeçalho, nil); err != nil {
			t.Error(err)
		}
	}
}

func TestFrequênciaAtirador_Interceptors(t *testing.T) {
	esperado := []string{
		"*interceptador.EndereçoRemoto",
		"*interceptador.Log",
		"*interceptor.Introspector",
		"*interceptador.VariáveisEndereço",
		"*interceptador.Padronizador",
		"*interceptador.Codificador",
		"*interceptador.BD",
	}

	var handler frequênciaAtirador

	verificadorResultado := testes.NovoVerificadorResultados("deve conter os interceptadores corretos", 0)
	verificadorResultado.DefinirEsperado(esperado, nil)
	if err := verificadorResultado.VerificaResultado(testes.TiposDaLista(handler.Interceptors()), nil); err != nil {
		t.Error(err)
	}
}
