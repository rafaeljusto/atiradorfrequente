package protocolo_test

import (
	_ "image/png"
	"testing"
	"time"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
	"github.com/rafaeljusto/atiradorfrequente/testes"
)

func TestFrequênciaPedido_Normalizar(t *testing.T) {
	data := time.Now()

	cenários := []struct {
		descrição        string
		frequênciaPedido protocolo.FrequênciaPedido
		esperado         protocolo.FrequênciaPedido
	}{
		{
			descrição: "deve normalizar os campos corretamente",
			frequênciaPedido: protocolo.FrequênciaPedido{
				Calibre:           "  calibre .380  ",
				ArmaUtilizada:     "  arma do clube  ",
				NúmeroSérie:       "  za785671  ",
				GuiaDeTráfego:     762556223,
				QuantidadeMunição: 50,
				DataInício:        data,
				DataTérmino:       data.Add(30 * time.Minute),
			},
			esperado: protocolo.FrequênciaPedido{
				Calibre:           "CALIBRE .380",
				ArmaUtilizada:     "ARMA DO CLUBE",
				NúmeroSérie:       "ZA785671",
				GuiaDeTráfego:     762556223,
				QuantidadeMunição: 50,
				DataInício:        data,
				DataTérmino:       data.Add(30 * time.Minute),
			},
		},
	}

	for i, cenário := range cenários {
		cenário.frequênciaPedido.Normalizar()

		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(cenário.esperado, nil)
		if err := verificadorResultado.VerificaResultado(cenário.frequênciaPedido, nil); err != nil {
			t.Error(err)
		}
	}
}

func TestFrequênciaPedido_Validar(t *testing.T) {
	data := time.Now()

	cenários := []struct {
		descrição        string
		frequênciaPedido protocolo.FrequênciaPedido
		esperado         protocolo.Mensagens
	}{
		{
			descrição: "deve aceitar um pedido válido",
			frequênciaPedido: protocolo.FrequênciaPedido{
				Calibre:           "380",
				ArmaUtilizada:     "Arma do Clube",
				NúmeroSérie:       "HG785671",
				QuantidadeMunição: 100,
				DataInício:        data.Add(-30 * time.Minute),
				DataTérmino:       data.Add(-10 * time.Minute),
			},
		},
		{
			descrição: "deve detectar erros de validação na maioria dos campos",
			frequênciaPedido: protocolo.FrequênciaPedido{
				NúmeroSérie: "785671",
				DataInício:  data,
				DataTérmino: data.Add(-10 * time.Minute),
			},
			esperado: protocolo.Mensagens{
				protocolo.NovaMensagemComCampo(protocolo.MensagemCódigoCampoNãoPreenchido, "calibre", ""),
				protocolo.NovaMensagemComCampo(protocolo.MensagemCódigoCampoNãoPreenchido, "armaUtilizada", ""),
				protocolo.NovaMensagemComValor(protocolo.MensagemCódigoNúmeroSérieInválido, "785671"),
				protocolo.NovaMensagemComCampo(protocolo.MensagemCódigoCampoNãoPreenchido, "quantidadeMunicao", "0"),
				protocolo.NovaMensagem(protocolo.MensagemCódigoDatasPeríodoIncorreto),
			},
		},
		{
			descrição: "deve detectar quando a data de término está depois do momento atual",
			frequênciaPedido: protocolo.FrequênciaPedido{
				Calibre:           "380",
				ArmaUtilizada:     "Arma do Clube",
				NúmeroSérie:       "HG785671",
				QuantidadeMunição: 100,
				DataInício:        data.Add(-30 * time.Minute),
				DataTérmino:       data.Add(10 * time.Minute),
			},
			esperado: protocolo.Mensagens{
				protocolo.NovaMensagem(protocolo.MensagemCódigoDatasPeríodoIncorreto),
			},
		},
	}

	for i, cenário := range cenários {
		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(cenário.esperado, nil)
		if err := verificadorResultado.VerificaResultado(cenário.frequênciaPedido.Validar(), nil); err != nil {
			t.Error(err)
		}
	}
}

func TestNovaFrequênciaPedidoCompleta(t *testing.T) {
	data := time.Now()

	cenários := []struct {
		descrição        string
		cr               int
		frequênciaPedido protocolo.FrequênciaPedido
		esperado         protocolo.FrequênciaPedidoCompleta
	}{
		{
			descrição: "deve inicializar um objeto do tipo FrequênciaPedidoCompleta corretamente",
			cr:        123456789,
			frequênciaPedido: protocolo.FrequênciaPedido{
				Calibre:           ".380",
				ArmaUtilizada:     "Arma do Clube",
				QuantidadeMunição: 50,
				DataInício:        data,
				DataTérmino:       data.Add(30 * time.Minute),
			},
			esperado: protocolo.FrequênciaPedidoCompleta{
				CR: 123456789,
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
		cr                          int
		númeroControle              protocolo.NúmeroControle
		frequênciaConfirmaçãoPedido protocolo.FrequênciaConfirmaçãoPedido
		esperado                    protocolo.FrequênciaConfirmaçãoPedidoCompleta
	}{
		{
			descrição:      "deve inicializar um objeto do tipo FrequênciaPedidoCompleta corretamente",
			cr:             123456789,
			númeroControle: protocolo.NovoNúmeroControle(7654, 918273645),
			frequênciaConfirmaçãoPedido: protocolo.FrequênciaConfirmaçãoPedido{
				Imagem: `TWFuIGlzIGRpc3Rpbmd1aXNoZWQsIG5vdCBvbmx5IGJ5IGhpcyByZWFzb24sIGJ1dCBieSB0aGlz
IHNpbmd1bGFyIHBhc3Npb24gZnJvbSBvdGhlciBhbmltYWxzLCB3aGljaCBpcyBhIGx1c3Qgb2Yg
dGhlIG1pbmQsIHRoYXQgYnkgYSBwZXJzZXZlcmFuY2Ugb2YgZGVsaWdodCBpbiB0aGUgY29udGlu
dWVkIGFuZCBpbmRlZmF0aWdhYmxlIGdlbmVyYXRpb24gb2Yga25vd2xlZGdlLCBleGNlZWRzIHRo
ZSBzaG9ydCB2ZWhlbWVuY2Ugb2YgYW55IGNhcm5hbCBwbGVhc3VyZS4=`,
			},
			esperado: protocolo.FrequênciaConfirmaçãoPedidoCompleta{
				CR:             123456789,
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

func TestFrequênciaConfirmaçãoPedido_Normalizar(t *testing.T) {
	cenários := []struct {
		descrição                   string
		frequênciaConfirmaçãoPedido protocolo.FrequênciaConfirmaçãoPedido
		esperado                    protocolo.FrequênciaConfirmaçãoPedido
	}{
		{
			descrição: "deve normalizar os campos corretamente",
			frequênciaConfirmaçãoPedido: protocolo.FrequênciaConfirmaçãoPedido{
				Imagem: `     TWFuIGlzIGRpc3Rpbmd1aXNoZWQsIG5vdCBvbmx5IGJ5IGhpcyByZWFzb24sIGJ1dCBieSB0aGlz
IHNpbmd1bGFyIHBhc3Npb24gZnJvbSBvdGhlciBhbmltYWxzLCB3aGljaCBpcyBhIGx1c3Qgb2Yg
dGhlIG1pbmQsIHRoYXQgYnkgYSBwZXJzZXZlcmFuY2Ugb2YgZGVsaWdodCBpbiB0aGUgY29udGlu
dWVkIGFuZCBpbmRlZmF0aWdhYmxlIGdlbmVyYXRpb24gb2Yga25vd2xlZGdlLCBleGNlZWRzIHRo
ZSBzaG9ydCB2ZWhlbWVuY2Ugb2YgYW55IGNhcm5hbCBwbGVhc3VyZS4=     `,
			},
			esperado: protocolo.FrequênciaConfirmaçãoPedido{
				Imagem: `TWFuIGlzIGRpc3Rpbmd1aXNoZWQsIG5vdCBvbmx5IGJ5IGhpcyByZWFzb24sIGJ1dCBieSB0aGlz
IHNpbmd1bGFyIHBhc3Npb24gZnJvbSBvdGhlciBhbmltYWxzLCB3aGljaCBpcyBhIGx1c3Qgb2Yg
dGhlIG1pbmQsIHRoYXQgYnkgYSBwZXJzZXZlcmFuY2Ugb2YgZGVsaWdodCBpbiB0aGUgY29udGlu
dWVkIGFuZCBpbmRlZmF0aWdhYmxlIGdlbmVyYXRpb24gb2Yga25vd2xlZGdlLCBleGNlZWRzIHRo
ZSBzaG9ydCB2ZWhlbWVuY2Ugb2YgYW55IGNhcm5hbCBwbGVhc3VyZS4=`,
			},
		},
	}

	for i, cenário := range cenários {
		cenário.frequênciaConfirmaçãoPedido.Normalizar()

		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(cenário.esperado, nil)
		if err := verificadorResultado.VerificaResultado(cenário.frequênciaConfirmaçãoPedido, nil); err != nil {
			t.Error(err)
		}
	}
}

func TestFrequênciaConfirmaçãoPedido_Validar(t *testing.T) {
	cenários := []struct {
		descrição                   string
		frequênciaConfirmaçãoPedido protocolo.FrequênciaConfirmaçãoPedido
		esperado                    protocolo.Mensagens
	}{
		{
			descrição: "deve aceitar uma imagem válida",
			frequênciaConfirmaçãoPedido: protocolo.FrequênciaConfirmaçãoPedido{
				Imagem: `iVBORw0KGgoAAAANSUhEUgAAAKgAAACoCAMAAABDlVWGAAABI1BMVEX/////////////////////
////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////8yc1n/AAAAYHRS
TlMAAQIDBAUGBwgJCwwNDg8REhQVFxgZGhseISQnKi0wMzY8P0JFSEtOVFdaXWBjZmlsb3J1eHt+
gYSHio2Qk5aZnKKlqKuusbS3ur3Aw8bJzM/S1djb3uHk5+rt8PP2+fz2kcVfAAADv0lEQVR4AdTB
haGEMBBAwQfEsI2n/06/610Dmxl+uTuXoUjJt+PJcrZ8Oosi1p25nQv/hHobFDJ3DfxxJ49SPt38
iLKi1iqRL5egmlx8CGlFtTUFgKV6lPN1AU5BPTmBtqHe1hZCYgIpIAcTOITsmIDLtJUJrI2+MIGl
M17buwsVV7ItDuP/SyAQ0lLkhNvu7u7u7lYd/d7/KYY6rUTGZ/auYX34wn5xrVqKRf83qEENalCD
GtSgBjWoQQ1qUIMa1KAGNahBDWpQg3qcQQ1qUIMadJ+vhv9eqEHLPf7fR1O3wIw8gLaM71y+QH5C
NdsD1uUcmpm/hbfyCdVoFjhLuIamVwp8lqvF6SrBc0ZuoYmZV766qfXATt1DqVduoR13vPWwOt7T
rpptAgtyCx0tEPW82Fbh7/lZpyQNAKdyC50pA+QXUqrokijuJKWfINfiFjpB1F6Tqtq9ubl5gkNJ
28CUnEJ7S0B5Xp8NXF/366PkNXdNUh9wLqfQlhAoDOmrEEJ9tEMuenDd8tW2E2gyIpQG9C0AvTdC
efBt5hq6CDCpOtDMK6vyAtpSAHZVD3rEbcKPN86bQD6oBx2h3CMvoJkCsKh60Bfecg+dB17SdaF4
A70B5lVRNc05NAOUszGAjgOXigF0E1iLA/QUGFFFLQAtfkEfgc6KWUcOINfhFbQABBWzQy56ey/5
7OOd1HixOO4MCpCsmOXJSlm+CvWzHLx6Bm2RWvnqRVFJAK9u+iNOstlTPguHFBUAea8eTJ1vD6aq
MfDo9OlpTBW1HIbhYYsqGgNOnUE3gDX9rtaADWfQMeBKv6srYMwZNChDuUW/o5YyEDiD6gpY/L0f
ra7kDjoDPCf1myVfgGmH0KYCMPP7LlChySFUG0Au0G8U5IB1uYRmC8CefqM9oJB1CtUCwKR+tWmA
ObmFJq6A0qB+paEScJVwDFXLC1D8FelQEXjJyjVUvRGkNKU6TZWAYo/cQzVG1F6gGgUHRI3KB6im
ygD5pbQqSi8XAMqT8gOqkQJR4VKHvtWxFBJVGJYvULXf8NbT2mhXNOgeXXvirZs2+QNVYiqkZuGU
b38sSC8VqCq/lJJvUCmYPC7yreLxZODrXzVSI1tntwC3Z1sjKfvzi0ENalCDGtSgBjWoQQ1qUIMa
1KAGNahBDWpQgxrUoAY1qEENalCD/megsTmFd2xOih6b08zH5sT9sVmF8L+4LJeIxbqOxnScFqDE
ZqVMDJb0NMZt7ZHU2Ozt/TTZ3KhvpX40JuRhicYfqXguO/N/fdwvoyFJPxBTbvQAAAAASUVORK5C
YII=`,
			},
		},
		{
			descrição: "deve detectar uma imagem com base64 incorreto",
			frequênciaConfirmaçãoPedido: protocolo.FrequênciaConfirmaçãoPedido{
				Imagem: `XXXXXXXXXXXXXXXXXXXXXXXXXX=`,
			},
			esperado: protocolo.Mensagens{
				protocolo.NovaMensagemComCampo(protocolo.MensagemCódigoImagemBase64Inválido, "imagem", `XXXXXXXXXXXXXXXXXXXXXXXXXX=`),
			},
		},
		{
			descrição: "deve detectar uma imagem não suportada ou inválida",
			frequênciaConfirmaçãoPedido: protocolo.FrequênciaConfirmaçãoPedido{
				Imagem: `VGVzdGUK`,
			},
			esperado: protocolo.Mensagens{
				protocolo.NovaMensagemComCampo(protocolo.MensagemCódigoImagemFormatoInválido, "imagem", `VGVzdGUK`),
			},
		},
	}

	for i, cenário := range cenários {
		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(cenário.esperado, nil)
		if err := verificadorResultado.VerificaResultado(cenário.frequênciaConfirmaçãoPedido.Validar(), nil); err != nil {
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
		erroEsperado           error
	}{
		{
			descrição: "deve aceitar um número de controle válido",
			texto:     []byte("  1234567-7654321  "),
			númeroControleEsperado: protocolo.NúmeroControle("1234567-7654321"),
		},
		{
			descrição: "deve detectar um número de controle inválido",
			texto:     []byte("  1234567 7654321  "),
			erroEsperado: protocolo.NovasMensagens(
				protocolo.NovaMensagemComValor(protocolo.MensagemCódigoNúmeroControleInválido, "1234567 7654321"),
			),
		},
	}

	for i, cenário := range cenários {
		var númeroControle protocolo.NúmeroControle
		mensagens := númeroControle.UnmarshalText(cenário.texto)

		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(cenário.númeroControleEsperado, cenário.erroEsperado)
		if err := verificadorResultado.VerificaResultado(númeroControle, mensagens); err != nil {
			t.Error(err)
		}
	}
}
