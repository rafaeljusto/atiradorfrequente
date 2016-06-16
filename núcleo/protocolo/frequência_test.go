package protocolo_test

import (
	"testing"
	"time"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
	"github.com/rafaeljusto/atiradorfrequente/testes"
)

func TestNovaFrequênciaPedidoCompleta(t *testing.T) {
	data := time.Now()

	cenários := []struct {
		descrição        string
		cr               string
		frequênciaPedido protocolo.FrequênciaPedido
		esperado         protocolo.FrequênciaPedidoCompleta
	}{
		{
			descrição: "deve inicializar um objeto do tipo FrequênciaPedidoCompleta corretamente",
			cr:        "123456789",
			frequênciaPedido: protocolo.FrequênciaPedido{
				Calibre:           ".380",
				ArmaUtilizada:     "Arma do Clube",
				QuantidadeMunição: 50,
				HorárioInício:     data,
				HorárioTérmino:    data.Add(30 * time.Minute),
			},
			esperado: protocolo.FrequênciaPedidoCompleta{
				CR: "123456789",
				FrequênciaPedido: protocolo.FrequênciaPedido{
					Calibre:           ".380",
					ArmaUtilizada:     "Arma do Clube",
					QuantidadeMunição: 50,
					HorárioInício:     data,
					HorárioTérmino:    data.Add(30 * time.Minute),
				},
			},
		},
	}

	for i, cenário := range cenários {
		freqênciaPedidoCompleta := protocolo.NovaFrequênciaPedidoCompleta(cenário.cr, cenário.frequênciaPedido)
		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(cenário.esperado, nil)

		if err := verificadorResultado.VerificaResultado(freqênciaPedidoCompleta, nil); err != nil {
			t.Error(err)
		}
	}
}

func TestNovaFrequênciaConfirmaçãoPedidoCompleta(t *testing.T) {
	cenários := []struct {
		descrição                   string
		cr                          string
		númeroControle              int64
		frequênciaConfirmaçãoPedido protocolo.FrequênciaConfirmaçãoPedido
		esperado                    protocolo.FrequênciaConfirmaçãoPedidoCompleta
	}{
		{
			descrição:      "deve inicializar um objeto do tipo FrequênciaPedidoCompleta corretamente",
			cr:             "123456789",
			númeroControle: 918273645,
			frequênciaConfirmaçãoPedido: protocolo.FrequênciaConfirmaçãoPedido{
				Imagem: `TWFuIGlzIGRpc3Rpbmd1aXNoZWQsIG5vdCBvbmx5IGJ5IGhpcyByZWFzb24sIGJ1dCBieSB0aGlz
IHNpbmd1bGFyIHBhc3Npb24gZnJvbSBvdGhlciBhbmltYWxzLCB3aGljaCBpcyBhIGx1c3Qgb2Yg
dGhlIG1pbmQsIHRoYXQgYnkgYSBwZXJzZXZlcmFuY2Ugb2YgZGVsaWdodCBpbiB0aGUgY29udGlu
dWVkIGFuZCBpbmRlZmF0aWdhYmxlIGdlbmVyYXRpb24gb2Yga25vd2xlZGdlLCBleGNlZWRzIHRo
ZSBzaG9ydCB2ZWhlbWVuY2Ugb2YgYW55IGNhcm5hbCBwbGVhc3VyZS4=`,
			},
			esperado: protocolo.FrequênciaConfirmaçãoPedidoCompleta{
				CR:             "123456789",
				NúmeroControle: 918273645,
				FrequênciaConfirmaçãoPedido: protocolo.FrequênciaConfirmaçãoPedido{
					Imagem: `TWFuIGlzIGRpc3Rpbmd1aXNoZWQsIG5vdCBvbmx5IGJ5IGhpcyByZWFzb24sIGJ1dCBieSB0aGlz
IHNpbmd1bGFyIHBhc3Npb24gZnJvbSBvdGhlciBhbmltYWxzLCB3aGljaCBpcyBhIGx1c3Qgb2Yg
dGhlIG1pbmQsIHRoYXQgYnkgYSBwZXJzZXZlcmFuY2Ugb2YgZGVsaWdodCBpbiB0aGUgY29udGlu
dWVkIGFuZCBpbmRlZmF0aWdhYmxlIGdlbmVyYXRpb24gb2Yga25vd2xlZGdlLCBleGNlZWRzIHRo
ZSBzaG9ydCB2ZWhlbWVuY2Ugb2YgYW55IGNhcm5hbCBwbGVhc3VyZS4=`,
				},
			},
		},
	}

	for i, cenário := range cenários {
		freqênciaPedidoCompleta := protocolo.NovaFrequênciaConfirmaçãoPedidoCompleta(cenário.cr, cenário.númeroControle, cenário.frequênciaConfirmaçãoPedido)
		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(cenário.esperado, nil)

		if err := verificadorResultado.VerificaResultado(freqênciaPedidoCompleta, nil); err != nil {
			t.Error(err)
		}
	}
}
