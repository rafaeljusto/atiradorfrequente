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
				DataInício:        data,
				DataTérmino:       data.Add(30 * time.Minute),
			},
			esperado: protocolo.FrequênciaPedidoCompleta{
				CR: "123456789",
				FrequênciaPedido: protocolo.FrequênciaPedido{
					Calibre:           ".380",
					ArmaUtilizada:     "Arma do Clube",
					QuantidadeMunição: 50,
					DataInício:        data,
					DataTérmino:       data.Add(30 * time.Minute),
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
		númeroControle              protocolo.NúmeroControle
		frequênciaConfirmaçãoPedido protocolo.FrequênciaConfirmaçãoPedido
		esperado                    protocolo.FrequênciaConfirmaçãoPedidoCompleta
	}{
		{
			descrição:      "deve inicializar um objeto do tipo FrequênciaPedidoCompleta corretamente",
			cr:             "123456789",
			númeroControle: protocolo.NovoNúmeroControle(7654, 918273645),
			frequênciaConfirmaçãoPedido: protocolo.FrequênciaConfirmaçãoPedido{
				Imagem: `TWFuIGlzIGRpc3Rpbmd1aXNoZWQsIG5vdCBvbmx5IGJ5IGhpcyByZWFzb24sIGJ1dCBieSB0aGlz
IHNpbmd1bGFyIHBhc3Npb24gZnJvbSBvdGhlciBhbmltYWxzLCB3aGljaCBpcyBhIGx1c3Qgb2Yg
dGhlIG1pbmQsIHRoYXQgYnkgYSBwZXJzZXZlcmFuY2Ugb2YgZGVsaWdodCBpbiB0aGUgY29udGlu
dWVkIGFuZCBpbmRlZmF0aWdhYmxlIGdlbmVyYXRpb24gb2Yga25vd2xlZGdlLCBleGNlZWRzIHRo
ZSBzaG9ydCB2ZWhlbWVuY2Ugb2YgYW55IGNhcm5hbCBwbGVhc3VyZS4=`,
			},
			esperado: protocolo.FrequênciaConfirmaçãoPedidoCompleta{
				CR:             "123456789",
				NúmeroControle: protocolo.NovoNúmeroControle(7654, 918273645),
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

func TestNovoNúmeroControle(t *testing.T) {
	cenários := []struct {
		descrição string
		id        int64
		controle  int64
		esperado  protocolo.NúmeroControle
	}{
		{
			descrição: "deve construir corretamente o número de controle",
			id:        123456789,
			controle:  987654321,
			esperado:  protocolo.NúmeroControle("123456789-987654321"),
		},
	}

	for i, cenário := range cenários {
		númeroControle := protocolo.NovoNúmeroControle(cenário.id, cenário.controle)

		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(cenário.esperado, nil)
		if err := verificadorResultado.VerificaResultado(númeroControle, nil); err != nil {
			t.Error(err)
		}
	}
}

func TestNúmeroControle_String(t *testing.T) {
	cenários := []struct {
		descrição      string
		númeroControle protocolo.NúmeroControle
		esperado       string
	}{
		{
			descrição:      "deve construir corretamente o número de controle",
			númeroControle: protocolo.NovoNúmeroControle(123456789, 987654321),
			esperado:       "123456789-987654321",
		},
	}

	for i, cenário := range cenários {
		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(cenário.esperado, nil)
		if err := verificadorResultado.VerificaResultado(cenário.númeroControle.String(), nil); err != nil {
			t.Error(err)
		}
	}
}

func TestNúmeroControle_ID(t *testing.T) {
	cenários := []struct {
		descrição      string
		númeroControle protocolo.NúmeroControle
		idEsperado     int64
	}{
		{
			descrição:      "deve extrair corretamente a identificação do número de controle",
			númeroControle: protocolo.NúmeroControle("123456789-987654321"),
			idEsperado:     123456789,
		},
		{
			descrição:      "deve tratar o caso de quando não existe o número aleatório",
			númeroControle: protocolo.NúmeroControle("123456789"),
			idEsperado:     123456789,
		},
		{
			descrição:      "deve tratar o caso de número de controle indefinido",
			númeroControle: protocolo.NúmeroControle(""),
			idEsperado:     0,
		},
		{
			descrição:      "deve tratar o caso de número de controle inválido",
			númeroControle: protocolo.NúmeroControle("X-X"),
			idEsperado:     0,
		},
	}

	for i, cenário := range cenários {
		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(cenário.idEsperado, nil)
		if err := verificadorResultado.VerificaResultado(cenário.númeroControle.ID(), nil); err != nil {
			t.Error(err)
		}
	}
}

func TestNúmeroControle_Controle(t *testing.T) {
	cenários := []struct {
		descrição        string
		númeroControle   protocolo.NúmeroControle
		controleEsperado int64
	}{
		{
			descrição:        "deve extrair corretamente o número aleatório do número de controle",
			númeroControle:   protocolo.NúmeroControle("123456789-987654321"),
			controleEsperado: 987654321,
		},
		{
			descrição:        "deve detectar quando não existe o número aleatório",
			númeroControle:   protocolo.NúmeroControle("123456789"),
			controleEsperado: 0,
		},
		{
			descrição:        "deve tratar o caso de número de controle indefinido",
			númeroControle:   protocolo.NúmeroControle(""),
			controleEsperado: 0,
		},
		{
			descrição:        "deve tratar o caso de número de controle inválido",
			númeroControle:   protocolo.NúmeroControle("X-X"),
			controleEsperado: 0,
		},
	}

	for i, cenário := range cenários {
		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(cenário.controleEsperado, nil)
		if err := verificadorResultado.VerificaResultado(cenário.númeroControle.Controle(), nil); err != nil {
			t.Error(err)
		}
	}
}

func TestNúmeroControle_Normalizar(t *testing.T) {
	cenários := []struct {
		descrição              string
		númeroControle         *protocolo.NúmeroControle
		númeroControleEsperado *protocolo.NúmeroControle
	}{
		{
			descrição: "deve padronizar corretamente o número de controle",
			númeroControle: func() *protocolo.NúmeroControle {
				n := protocolo.NúmeroControle("  1234567-7654321  ")
				return &n
			}(),
			númeroControleEsperado: func() *protocolo.NúmeroControle {
				n := protocolo.NúmeroControle("1234567-7654321")
				return &n
			}(),
		},
		{
			descrição: "deve tratar corretamente o caso em que o número de controle esta indefinido",
		},
	}

	for i, cenário := range cenários {
		cenário.númeroControle.Normalizar()

		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(cenário.númeroControleEsperado, nil)
		if err := verificadorResultado.VerificaResultado(cenário.númeroControle, nil); err != nil {
			t.Error(err)
		}
	}
}

func TestNúmeroControle_Validar(t *testing.T) {
	cenários := []struct {
		descrição          string
		númeroControle     protocolo.NúmeroControle
		mensagensEsperadas protocolo.Mensagens
	}{
		{
			descrição:      "deve aceitar um número de controle válido",
			númeroControle: protocolo.NúmeroControle("1234567-7654321"),
		},
		{
			descrição:      "deve detectar um número de controle sem hífen",
			númeroControle: protocolo.NúmeroControle("12345677654321"),
			mensagensEsperadas: protocolo.NovasMensagens(
				protocolo.NovaMensagemComValor(protocolo.MensagemCódigoNúmeroControleInválido, "12345677654321"),
			),
		},
	}

	for i, cenário := range cenários {
		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(nil, cenário.mensagensEsperadas)
		if err := verificadorResultado.VerificaResultado(nil, cenário.númeroControle.Validar()); err != nil {
			t.Error(err)
		}
	}
}

func TestNúmeroControle_UnmarshalText(t *testing.T) {
	cenários := []struct {
		descrição              string
		texto                  []byte
		númeroControleEsperado protocolo.NúmeroControle
		mensagensEsperadas     protocolo.Mensagens
	}{
		{
			descrição: "deve aceitar um número de controle válido",
			texto:     []byte("  1234567-7654321  "),
			númeroControleEsperado: protocolo.NúmeroControle("1234567-7654321"),
		},
		{
			descrição: "deve detectar um número de controle inválido",
			texto:     []byte("  1234567 7654321  "),
			mensagensEsperadas: protocolo.NovasMensagens(
				protocolo.NovaMensagemComValor(protocolo.MensagemCódigoNúmeroControleInválido, "1234567 7654321"),
			),
		},
	}

	for i, cenário := range cenários {
		var númeroControle protocolo.NúmeroControle
		mensagens := númeroControle.UnmarshalText(cenário.texto)

		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(cenário.númeroControleEsperado, cenário.mensagensEsperadas)
		if err := verificadorResultado.VerificaResultado(númeroControle, mensagens); err != nil {
			t.Error(err)
		}
	}
}
