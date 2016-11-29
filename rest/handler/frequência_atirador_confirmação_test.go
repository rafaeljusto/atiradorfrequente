package handler

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/atirador"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/bd"
	núcleoconfig "github.com/rafaeljusto/atiradorfrequente/núcleo/config"
	núcleolog "github.com/rafaeljusto/atiradorfrequente/núcleo/log"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
	restconfig "github.com/rafaeljusto/atiradorfrequente/rest/config"
	"github.com/rafaeljusto/atiradorfrequente/testes"
	"github.com/rafaeljusto/atiradorfrequente/testes/simulador"
	"github.com/registrobr/gostk/errors"
	gostklog "github.com/registrobr/gostk/log"
)

func TestFrequênciaAtiradorConfirmação_Get(t *testing.T) {
	cenários := []struct {
		descrição          string
		cr                 int
		númeroControle     protocolo.NúmeroControle
		códigoVerificação  string
		logger             gostklog.Logger
		configuração       *restconfig.Configuração
		serviçoAtirador    atirador.Serviço
		códigoHTTPEsperado int
		mensagensEsperadas protocolo.Mensagens
	}{
		{
			descrição:         "deve obter corretamente os dados de frequência do atirador",
			cr:                123456789,
			númeroControle:    protocolo.NovoNúmeroControle(7654, 918273645),
			códigoVerificação: "5JRYo4LFpvhr9gnALUTNJf8v3Z3TwAduwWQy1yxx1c4Q",
			configuração: func() *restconfig.Configuração {
				return new(restconfig.Configuração)
			}(),
			serviçoAtirador: simulador.ServiçoAtirador{
				SimulaObterFrequência: func(cr int, númeroControle protocolo.NúmeroControle, códigoVerificação string) (protocolo.FrequênciaResposta, error) {
					return protocolo.FrequênciaResposta{}, nil
				},
			},
			códigoHTTPEsperado: http.StatusOK,
		},
		{
			descrição:         "deve detectar quando a configuração não foi inicializada",
			cr:                123456789,
			númeroControle:    protocolo.NovoNúmeroControle(7654, 918273645),
			códigoVerificação: "5JRYo4LFpvhr9gnALUTNJf8v3Z3TwAduwWQy1yxx1c4Q",
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
			descrição:         "deve detectar um erro na camada de serviço do atirador",
			cr:                123456789,
			númeroControle:    protocolo.NovoNúmeroControle(7654, 918273645),
			códigoVerificação: "5JRYo4LFpvhr9gnALUTNJf8v3Z3TwAduwWQy1yxx1c4Q",
			logger: simulador.Logger{
				SimulaError: func(e error) {
					if !strings.HasSuffix(e.Error(), "erro de baixo nível") {
						t.Error("não está adicionando o erro correto ao log")
					}
				},
			},
			configuração: func() *restconfig.Configuração {
				return new(restconfig.Configuração)
			}(),
			serviçoAtirador: simulador.ServiçoAtirador{
				SimulaObterFrequência: func(cr int, númeroControle protocolo.NúmeroControle, códigoVerificação string) (protocolo.FrequênciaResposta, error) {
					return protocolo.FrequênciaResposta{}, errors.Errorf("erro de baixo nível")
				},
			},
			códigoHTTPEsperado: http.StatusInternalServerError,
		},
		{
			descrição:         "deve detectar mensagens na camada de serviço do atirador",
			cr:                123456789,
			númeroControle:    protocolo.NovoNúmeroControle(7654, 918273645),
			códigoVerificação: "5JRYo4LFpvhr9gnALUTNJf8v3Z3TwAduwWQy1yxx1c4Q",
			logger:            simulador.Logger{},
			configuração: func() *restconfig.Configuração {
				return new(restconfig.Configuração)
			}(),
			serviçoAtirador: simulador.ServiçoAtirador{
				SimulaObterFrequência: func(cr int, númeroControle protocolo.NúmeroControle, códigoVerificação string) (protocolo.FrequênciaResposta, error) {
					return protocolo.FrequênciaResposta{}, protocolo.NovasMensagens(
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

	configuraçãoOriginal := restconfig.Atual()
	defer func() {
		restconfig.AtualizarConfiguração(configuraçãoOriginal)
	}()

	serviçoAtiradorOriginal := atirador.NovoServiço
	defer func() {
		atirador.NovoServiço = serviçoAtiradorOriginal
	}()

	for i, cenário := range cenários {
		restconfig.AtualizarConfiguração(cenário.configuração)

		atirador.NovoServiço = func(s *bd.SQLogger, l núcleolog.Serviço, configuração núcleoconfig.Configuração) atirador.Serviço {
			return cenário.serviçoAtirador
		}

		handler := frequênciaAtiradorConfirmação{
			CR:                cenário.cr,
			NúmeroControle:    cenário.númeroControle,
			CódigoVerificação: cenário.códigoVerificação,
		}
		handler.DefineLogger(cenário.logger)

		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)

		verificadorResultado.DefinirEsperado(cenário.códigoHTTPEsperado, nil)
		if err := verificadorResultado.VerificaResultado(handler.Get(), nil); err != nil {
			t.Error(err)
		}

		verificadorResultado.DefinirEsperado(cenário.mensagensEsperadas, nil)
		if err := verificadorResultado.VerificaResultado(handler.Mensagens, nil); err != nil {
			t.Error(err)
		}
	}
}

func TestFrequênciaAtiradorConfirmação_Put(t *testing.T) {
	cenários := []struct {
		descrição                   string
		cr                          int
		númeroControle              protocolo.NúmeroControle
		frequênciaConfirmaçãoPedido protocolo.FrequênciaConfirmaçãoPedido
		logger                      gostklog.Logger
		configuração                *restconfig.Configuração
		serviçoAtirador             atirador.Serviço
		códigoHTTPEsperado          int
		mensagensEsperadas          protocolo.Mensagens
	}{
		{
			descrição:      "deve confirmar corretamente os dados de frequência do atirador",
			cr:             123456789,
			númeroControle: protocolo.NovoNúmeroControle(7654, 918273645),
			frequênciaConfirmaçãoPedido: protocolo.FrequênciaConfirmaçãoPedido{
				Imagem: `TWFuIGlzIGRpc3Rpbmd1aXNoZWQsIG5vdCBvbmx5IGJ5IGhpcyByZWFzb24sIGJ1dCBieSB0aGlz
IHNpbmd1bGFyIHBhc3Npb24gZnJvbSBvdGhlciBhbmltYWxzLCB3aGljaCBpcyBhIGx1c3Qgb2Yg
dGhlIG1pbmQsIHRoYXQgYnkgYSBwZXJzZXZlcmFuY2Ugb2YgZGVsaWdodCBpbiB0aGUgY29udGlu
dWVkIGFuZCBpbmRlZmF0aWdhYmxlIGdlbmVyYXRpb24gb2Yga25vd2xlZGdlLCBleGNlZWRzIHRo
ZSBzaG9ydCB2ZWhlbWVuY2Ugb2YgYW55IGNhcm5hbCBwbGVhc3VyZS4=`,
			},
			configuração: func() *restconfig.Configuração {
				return new(restconfig.Configuração)
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
			cr:             123456789,
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
			cr:             123456789,
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
			configuração: func() *restconfig.Configuração {
				return new(restconfig.Configuração)
			}(),
			serviçoAtirador: simulador.ServiçoAtirador{
				SimulaConfirmarFrequência: func(frequênciaConfirmaçãoPedidoCompleta protocolo.FrequênciaConfirmaçãoPedidoCompleta) error {
					return errors.Errorf("erro de baixo nível")
				},
			},
			códigoHTTPEsperado: http.StatusInternalServerError,
		},
		{
			descrição:      "deve detectar mensagens na camada de serviço do atirador",
			cr:             123456789,
			númeroControle: protocolo.NovoNúmeroControle(7654, 918273645),
			frequênciaConfirmaçãoPedido: protocolo.FrequênciaConfirmaçãoPedido{
				Imagem: `TWFuIGlzIGRpc3Rpbmd1aXNoZWQsIG5vdCBvbmx5IGJ5IGhpcyByZWFzb24sIGJ1dCBieSB0aGlz
IHNpbmd1bGFyIHBhc3Npb24gZnJvbSBvdGhlciBhbmltYWxzLCB3aGljaCBpcyBhIGx1c3Qgb2Yg
dGhlIG1pbmQsIHRoYXQgYnkgYSBwZXJzZXZlcmFuY2Ugb2YgZGVsaWdodCBpbiB0aGUgY29udGlu
dWVkIGFuZCBpbmRlZmF0aWdhYmxlIGdlbmVyYXRpb24gb2Yga25vd2xlZGdlLCBleGNlZWRzIHRo
ZSBzaG9ydCB2ZWhlbWVuY2Ugb2YgYW55IGNhcm5hbCBwbGVhc3VyZS4=`,
			},
			logger: simulador.Logger{},
			configuração: func() *restconfig.Configuração {
				return new(restconfig.Configuração)
			}(),
			serviçoAtirador: simulador.ServiçoAtirador{
				SimulaConfirmarFrequência: func(frequênciaConfirmaçãoPedidoCompleta protocolo.FrequênciaConfirmaçãoPedidoCompleta) error {
					return protocolo.NovasMensagens(
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

	configuraçãoOriginal := restconfig.Atual()
	defer func() {
		restconfig.AtualizarConfiguração(configuraçãoOriginal)
	}()

	serviçoAtiradorOriginal := atirador.NovoServiço
	defer func() {
		atirador.NovoServiço = serviçoAtiradorOriginal
	}()

	for i, cenário := range cenários {
		restconfig.AtualizarConfiguração(cenário.configuração)

		atirador.NovoServiço = func(s *bd.SQLogger, l núcleolog.Serviço, configuração núcleoconfig.Configuração) atirador.Serviço {
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

		verificadorResultado.DefinirEsperado(cenário.mensagensEsperadas, nil)
		if err := verificadorResultado.VerificaResultado(handler.Mensagens, nil); err != nil {
			t.Error(err)
		}
	}
}

func TestFrequênciaAtiradorConfirmação_Interceptors(t *testing.T) {
	esperado := []string{
		"*interceptador.EndereçoRemoto",
		"*interceptador.Log",
		"*interceptor.Introspector",
		"*interceptador.Codificador",
		"*interceptador.ParâmetrosConsulta",
		"*interceptador.VariáveisEndereço",
		"*interceptador.Padronizador",
		"*interceptador.BD",
	}

	var handler frequênciaAtiradorConfirmação

	verificadorResultado := testes.NovoVerificadorResultados("deve conter os interceptadores corretos", 0)
	verificadorResultado.DefinirEsperado(esperado, nil)
	if err := verificadorResultado.VerificaResultado(testes.TiposDaLista(handler.Interceptors()), nil); err != nil {
		t.Error(err)
	}
}
