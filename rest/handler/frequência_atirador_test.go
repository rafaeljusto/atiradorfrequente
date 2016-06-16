package handler

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/atirador"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
	"github.com/rafaeljusto/atiradorfrequente/testes"
)

func TestFrequênciaAtirador_Post(t *testing.T) {
	data := time.Now()

	cenários := []struct {
		descrição          string
		cr                 string
		frequênciaPedido   protocolo.FrequênciaPedido
		serviçoAtirador    atirador.Serviço
		códigoHTTPEsperado int
		esperado           *protocolo.FrequênciaPendenteResposta
	}{
		{
			descrição: "deve cadastrar corretamente os dados de frequência do atirador",
			cr:        "123456789",
			frequênciaPedido: protocolo.FrequênciaPedido{
				Calibre:           ".380",
				ArmaUtilizada:     "Arma do Clube",
				QuantidadeMunição: 50,
				HorárioInício:     data,
				HorárioTérmino:    data.Add(30 * time.Minute),
			},
			serviçoAtirador: serviçoAtiradorSimulado{
				SimulaCadastrarFrequência: func(frequênciaPedidoCompleta protocolo.FrequênciaPedidoCompleta) (protocolo.FrequênciaPendenteResposta, error) {
					return protocolo.FrequênciaPendenteResposta{
						NúmeroControle: 918273645,
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
				NúmeroControle: 918273645,
				Imagem: `TWFuIGlzIGRpc3Rpbmd1aXNoZWQsIG5vdCBvbmx5IGJ5IGhpcyByZWFzb24sIGJ1dCBieSB0aGlz
IHNpbmd1bGFyIHBhc3Npb24gZnJvbSBvdGhlciBhbmltYWxzLCB3aGljaCBpcyBhIGx1c3Qgb2Yg
dGhlIG1pbmQsIHRoYXQgYnkgYSBwZXJzZXZlcmFuY2Ugb2YgZGVsaWdodCBpbiB0aGUgY29udGlu
dWVkIGFuZCBpbmRlZmF0aWdhYmxlIGdlbmVyYXRpb24gb2Yga25vd2xlZGdlLCBleGNlZWRzIHRo
ZSBzaG9ydCB2ZWhlbWVuY2Ugb2YgYW55IGNhcm5hbCBwbGVhc3VyZS4=`,
			},
		},
		{
			descrição: "deve detectar um erro na camada de serviço do atirador",
			cr:        "123456789",
			frequênciaPedido: protocolo.FrequênciaPedido{
				Calibre:           ".380",
				ArmaUtilizada:     "Arma do Clube",
				QuantidadeMunição: 50,
				HorárioInício:     data,
				HorárioTérmino:    data.Add(30 * time.Minute),
			},
			serviçoAtirador: serviçoAtiradorSimulado{
				SimulaCadastrarFrequência: func(frequênciaPedidoCompleta protocolo.FrequênciaPedidoCompleta) (protocolo.FrequênciaPendenteResposta, error) {
					return protocolo.FrequênciaPendenteResposta{}, fmt.Errorf("erro de baixo nível")
				},
			},
			códigoHTTPEsperado: http.StatusInternalServerError,
		},
	}

	serviçoAtiradorOriginal := atirador.NovoServiço
	defer func() {
		atirador.NovoServiço = serviçoAtiradorOriginal
	}()

	for i, cenário := range cenários {
		atirador.NovoServiço = func() atirador.Serviço {
			return cenário.serviçoAtirador
		}

		handler := frequênciaAtirador{
			CR:               cenário.cr,
			FrequênciaPedido: cenário.frequênciaPedido,
		}

		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)

		verificadorResultado.DefinirEsperado(cenário.códigoHTTPEsperado, nil)
		if err := verificadorResultado.VerificaResultado(handler.Post(), nil); err != nil {
			t.Error(err)
		}

		verificadorResultado.DefinirEsperado(cenário.esperado, nil)
		if err := verificadorResultado.VerificaResultado(handler.FrequênciaPendenteResposta, nil); err != nil {
			t.Error(err)
		}
	}
}

type serviçoAtiradorSimulado struct {
	SimulaCadastrarFrequência func(protocolo.FrequênciaPedidoCompleta) (protocolo.FrequênciaPendenteResposta, error)
	SimulaConfirmarFrequência func(protocolo.FrequênciaConfirmaçãoPedidoCompleta) error
}

func (s serviçoAtiradorSimulado) CadastrarFrequência(frequênciaPedidoCompleta protocolo.FrequênciaPedidoCompleta) (protocolo.FrequênciaPendenteResposta, error) {
	return s.SimulaCadastrarFrequência(frequênciaPedidoCompleta)
}

func (s serviçoAtiradorSimulado) ConfirmarFrequência(frequênciaConfirmaçãoPedidoCompleta protocolo.FrequênciaConfirmaçãoPedidoCompleta) error {
	return s.SimulaConfirmarFrequência(frequênciaConfirmaçãoPedidoCompleta)
}
