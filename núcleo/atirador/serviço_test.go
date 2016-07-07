package atirador

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	_ "image/png"
	"testing"
	"time"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/bd"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/config"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/erros"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
	"github.com/rafaeljusto/atiradorfrequente/testes"
	"github.com/registrobr/gostk/errors"
	"golang.org/x/image/font/basicfont"
)

func TestServiço_CadastrarFrequência(t *testing.T) {
	data := time.Now()

	imagemLogoExtraída, err := base64.StdEncoding.DecodeString(imagemLogoPNG)

	if err != nil {
		t.Fatalf("Erro ao extrair a imagem de teste do logo. Detalhes: %s", err)
	}

	imagemLogoBuffer := bytes.NewBuffer(imagemLogoExtraída)
	imagemLogo, _, err := image.Decode(imagemLogoBuffer)

	if err != nil {
		t.Fatalf("Erro ao interpretar imagem. Detalhes: %s", err)
	}

	cenários := []struct {
		descrição                string
		configuração             config.Configuração
		frequênciaPedidoCompleta protocolo.FrequênciaPedidoCompleta
		frequênciaDAO            frequênciaDAO
		imagemNúmeroControleLogo string
		esperado                 protocolo.FrequênciaPendenteResposta
		erroEsperado             error
	}{
		{
			descrição: "deve cadastrar corretamente uma frequência",
			configuração: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.ImagemNúmeroControle.Largura = 3508
				configuração.Atirador.ImagemNúmeroControle.Altura = 2480
				configuração.Atirador.ImagemNúmeroControle.CorFundo.Color = color.RGBA{0xff, 0xff, 0xff, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Fonte.Face = basicfont.Face7x13
				configuração.Atirador.ImagemNúmeroControle.Fonte.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Logo.Imagem.Image = imagemLogo
				configuração.Atirador.ImagemNúmeroControle.Logo.Espaçamento = 100
				configuração.Atirador.ImagemNúmeroControle.Borda.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Cor.Color = color.RGBA{0xee, 0xee, 0xee, 0xff}
				return configuração
			}(),
			frequênciaPedidoCompleta: protocolo.FrequênciaPedidoCompleta{
				CR: "123456789",
				FrequênciaPedido: protocolo.FrequênciaPedido{
					Calibre:           ".380",
					ArmaUtilizada:     "Arma do Clube",
					QuantidadeMunição: 50,
					DataInício:        data,
					DataTérmino:       data.Add(30 * time.Minute),
				},
			},
			frequênciaDAO: simulaFrequênciaDAO{
				simulaCriar: func(frequência *frequência) error {
					if frequência.Controle == 0 {
						t.Errorf("Número aleatório para controle não gerado")
					}

					frequência.ID = 1
					frequência.Controle = 123
					return nil
				},
				simulaAtualizar: func(frequência *frequência) error {
					if frequência.ImagemNúmeroControle == "" {
						t.Errorf("Imagem com o número de controle não gerada")
					}

					return nil
				},
			},
			esperado: protocolo.FrequênciaPendenteResposta{
				NúmeroControle: protocolo.NovoNúmeroControle(1, 123),
				Imagem:         imagemNúmeroControlePNG,
			},
		},
		{
			descrição: "deve detectar um erro ao persistir uma nova frequência",
			configuração: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.ImagemNúmeroControle.Largura = 3508
				configuração.Atirador.ImagemNúmeroControle.Altura = 2480
				configuração.Atirador.ImagemNúmeroControle.CorFundo.Color = color.RGBA{0xff, 0xff, 0xff, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Fonte.Face = basicfont.Face7x13
				configuração.Atirador.ImagemNúmeroControle.Fonte.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Logo.Imagem.Image = imagemLogo
				configuração.Atirador.ImagemNúmeroControle.Logo.Espaçamento = 100
				configuração.Atirador.ImagemNúmeroControle.Borda.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Cor.Color = color.RGBA{0xee, 0xee, 0xee, 0xff}
				return configuração
			}(),
			frequênciaPedidoCompleta: protocolo.FrequênciaPedidoCompleta{
				CR: "123456789",
				FrequênciaPedido: protocolo.FrequênciaPedido{
					Calibre:           ".380",
					ArmaUtilizada:     "Arma do Clube",
					QuantidadeMunição: 50,
					DataInício:        data,
					DataTérmino:       data.Add(30 * time.Minute),
				},
			},
			frequênciaDAO: simulaFrequênciaDAO{
				simulaCriar: func(frequência *frequência) error {
					return errors.Errorf("erro de criação")
				},
			},
			erroEsperado: errors.Errorf("erro de criação"),
		},
		{
			descrição: "deve detectar um erro ao gerar a imagem PNG",
			configuração: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.ImagemNúmeroControle.Largura = 0
				configuração.Atirador.ImagemNúmeroControle.Altura = 0
				configuração.Atirador.ImagemNúmeroControle.CorFundo.Color = color.RGBA{0xff, 0xff, 0xff, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Fonte.Face = basicfont.Face7x13
				configuração.Atirador.ImagemNúmeroControle.Fonte.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Logo.Imagem.Image = imagemLogo
				configuração.Atirador.ImagemNúmeroControle.Logo.Espaçamento = 100
				configuração.Atirador.ImagemNúmeroControle.Borda.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Cor.Color = color.RGBA{0xee, 0xee, 0xee, 0xff}
				return configuração
			}(),
			frequênciaPedidoCompleta: protocolo.FrequênciaPedidoCompleta{
				CR: "123456789",
				FrequênciaPedido: protocolo.FrequênciaPedido{
					Calibre:           ".380",
					ArmaUtilizada:     "Arma do Clube",
					QuantidadeMunição: 50,
					DataInício:        data,
					DataTérmino:       data.Add(30 * time.Minute),
				},
			},
			frequênciaDAO: simulaFrequênciaDAO{
				simulaCriar: func(frequência *frequência) error {
					if frequência.Controle == 0 {
						t.Errorf("Número aleatório para controle não gerado")
					}

					frequência.ID = 1
					frequência.Controle = 123
					return nil
				},
			},
			erroEsperado: errors.Errorf("png: invalid format: invalid image size: 0x0"),
		},
		{
			descrição: "deve detectar um erro ao atualizar uma frequência",
			configuração: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.ImagemNúmeroControle.Largura = 3508
				configuração.Atirador.ImagemNúmeroControle.Altura = 2480
				configuração.Atirador.ImagemNúmeroControle.CorFundo.Color = color.RGBA{0xff, 0xff, 0xff, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Fonte.Face = basicfont.Face7x13
				configuração.Atirador.ImagemNúmeroControle.Fonte.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.Logo.Imagem.Image = imagemLogo
				configuração.Atirador.ImagemNúmeroControle.Logo.Espaçamento = 100
				configuração.Atirador.ImagemNúmeroControle.Borda.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.Borda.Cor.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Largura = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Espaçamento = 50
				configuração.Atirador.ImagemNúmeroControle.LinhaFundo.Cor.Color = color.RGBA{0xee, 0xee, 0xee, 0xff}
				return configuração
			}(),
			frequênciaPedidoCompleta: protocolo.FrequênciaPedidoCompleta{
				CR: "123456789",
				FrequênciaPedido: protocolo.FrequênciaPedido{
					Calibre:           ".380",
					ArmaUtilizada:     "Arma do Clube",
					QuantidadeMunição: 50,
					DataInício:        data,
					DataTérmino:       data.Add(30 * time.Minute),
				},
			},
			frequênciaDAO: simulaFrequênciaDAO{
				simulaCriar: func(frequência *frequência) error {
					frequência.ID = 1
					frequência.Controle = 123
					return nil
				},
				simulaAtualizar: func(frequência *frequência) error {
					return errors.Errorf("erro de atualização")
				},
			},
			erroEsperado: errors.Errorf("erro de atualização"),
		},
	}

	daoOriginal := novaFrequênciaDAO
	defer func() {
		novaFrequênciaDAO = daoOriginal
	}()

	for i, cenário := range cenários {
		novaFrequênciaDAO = func(sqlogger *bd.SQLogger) frequênciaDAO {
			return cenário.frequênciaDAO
		}

		serviço := NovoServiço(nil, cenário.configuração)
		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(cenário.esperado, cenário.erroEsperado)

		if err := verificadorResultado.VerificaResultado(serviço.CadastrarFrequência(cenário.frequênciaPedidoCompleta)); err != nil {
			t.Error(err)
		}
	}
}

func TestServiço_ConfirmarFrequência(t *testing.T) {
	data := time.Now()

	cenários := []struct {
		descrição                           string
		configuração                        config.Configuração
		frequênciaConfirmaçãoPedidoCompleta protocolo.FrequênciaConfirmaçãoPedidoCompleta
		frequênciaDAO                       frequênciaDAO
		erroEsperado                        error
	}{
		{
			descrição: "deve confirmar corretamente uma frequência",
			configuração: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.PrazoConfirmação = 20 * time.Minute
				return configuração
			}(),
			frequênciaConfirmaçãoPedidoCompleta: protocolo.FrequênciaConfirmaçãoPedidoCompleta{
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
			frequênciaDAO: simulaFrequênciaDAO{
				simulaAtualizar: func(frequência *frequência) error {
					if frequência.DataConfirmação.Before(data) {
						t.Errorf("Data de confirmação não definida corretamente")
					}

					if frequência.ImagemConfirmação == "" {
						t.Errorf("Imagem de confirmação não definida corretamente")
					}

					return nil
				},
				simulaResgatar: func(id int64) (frequência, error) {
					if id != 7654 {
						t.Errorf("ID %d inesperado", id)
					}

					return frequência{
						ID:                7654,
						Controle:          918273645,
						CR:                "123456789",
						Calibre:           ".380",
						ArmaUtilizada:     "Arma do Clube",
						NúmeroSérie:       "ZA785671",
						GuiaDeTráfego:     "XYZ12345",
						QuantidadeMunição: 50,
						DataInício:        data.Add(-40 * time.Minute),
						DataTérmino:       data.Add(-10 * time.Minute),
						DataCriação:       data.Add(-5 * time.Minute),
						ImagemNúmeroControle: `TWFuIGlzIGRpc3Rpbmd1aXNoZWQsIG5vdCBvbmx5IGJ5IGhpcyByZWFzb24sIGJ1dCBieSB0aGlz
IHNpbmd1bGFyIHBhc3Npb24gZnJvbSBvdGhlciBhbmltYWxzLCB3aGljaCBpcyBhIGx1c3Qgb2Yg
dGhlIG1pbmQsIHRoYXQgYnkgYSBwZXJzZXZlcmFuY2Ugb2YgZGVsaWdodCBpbiB0aGUgY29udGlu
dWVkIGFuZCBpbmRlZmF0aWdhYmxlIGdlbmVyYXRpb24gb2Yga25vd2xlZGdlLCBleGNlZWRzIHRo
ZSBzaG9ydCB2ZWhlbWVuY2Ugb2YgYW55IGNhcm5hbCBwbGVhc3VyZS4=`,
					}, nil
				},
			},
		},
		{
			descrição: "deve detectar um erro ao resgatar a frequência",
			configuração: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.PrazoConfirmação = 20 * time.Minute
				return configuração
			}(),
			frequênciaConfirmaçãoPedidoCompleta: protocolo.FrequênciaConfirmaçãoPedidoCompleta{
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
			frequênciaDAO: simulaFrequênciaDAO{
				simulaResgatar: func(id int64) (frequência, error) {
					return frequência{}, erros.NãoEncontrado
				},
			},
			erroEsperado: erros.NãoEncontrado,
		},
		{
			descrição: "deve detectar quando o CR e o número de controle não conferem",
			configuração: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.PrazoConfirmação = 20 * time.Minute
				return configuração
			}(),
			frequênciaConfirmaçãoPedidoCompleta: protocolo.FrequênciaConfirmaçãoPedidoCompleta{
				CR:             "12345678X",
				NúmeroControle: protocolo.NovoNúmeroControle(7654, 918273640),
				FrequênciaConfirmaçãoPedido: protocolo.FrequênciaConfirmaçãoPedido{
					Imagem: `TWFuIGlzIGRpc3Rpbmd1aXNoZWQsIG5vdCBvbmx5IGJ5IGhpcyByZWFzb24sIGJ1dCBieSB0aGlz
IHNpbmd1bGFyIHBhc3Npb24gZnJvbSBvdGhlciBhbmltYWxzLCB3aGljaCBpcyBhIGx1c3Qgb2Yg
dGhlIG1pbmQsIHRoYXQgYnkgYSBwZXJzZXZlcmFuY2Ugb2YgZGVsaWdodCBpbiB0aGUgY29udGlu
dWVkIGFuZCBpbmRlZmF0aWdhYmxlIGdlbmVyYXRpb24gb2Yga25vd2xlZGdlLCBleGNlZWRzIHRo
ZSBzaG9ydCB2ZWhlbWVuY2Ugb2YgYW55IGNhcm5hbCBwbGVhc3VyZS4=`,
				},
			},
			frequênciaDAO: simulaFrequênciaDAO{
				simulaResgatar: func(id int64) (frequência, error) {
					return frequência{
						ID:                7654,
						Controle:          918273645,
						CR:                "123456789",
						Calibre:           ".380",
						ArmaUtilizada:     "Arma do Clube",
						NúmeroSérie:       "ZA785671",
						GuiaDeTráfego:     "XYZ12345",
						QuantidadeMunição: 50,
						DataInício:        data.Add(-40 * time.Minute),
						DataTérmino:       data.Add(-10 * time.Minute),
						DataCriação:       data.Add(-5 * time.Minute),
						ImagemNúmeroControle: `TWFuIGlzIGRpc3Rpbmd1aXNoZWQsIG5vdCBvbmx5IGJ5IGhpcyByZWFzb24sIGJ1dCBieSB0aGlz
IHNpbmd1bGFyIHBhc3Npb24gZnJvbSBvdGhlciBhbmltYWxzLCB3aGljaCBpcyBhIGx1c3Qgb2Yg
dGhlIG1pbmQsIHRoYXQgYnkgYSBwZXJzZXZlcmFuY2Ugb2YgZGVsaWdodCBpbiB0aGUgY29udGlu
dWVkIGFuZCBpbmRlZmF0aWdhYmxlIGdlbmVyYXRpb24gb2Yga25vd2xlZGdlLCBleGNlZWRzIHRo
ZSBzaG9ydCB2ZWhlbWVuY2Ugb2YgYW55IGNhcm5hbCBwbGVhc3VyZS4=`,
					}, nil
				},
			},
			erroEsperado: protocolo.NovasMensagens(
				protocolo.NovaMensagemComValor(protocolo.MensagemCódigoCRInválido, "12345678X"),
				protocolo.NovaMensagemComValor(protocolo.MensagemCódigoNúmeroControleInválido, "7654-918273640"),
			),
		},
		{
			descrição: "deve detectar quando o prazo de confirmação expirar",
			configuração: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.PrazoConfirmação = 20 * time.Minute
				return configuração
			}(),
			frequênciaConfirmaçãoPedidoCompleta: protocolo.FrequênciaConfirmaçãoPedidoCompleta{
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
			frequênciaDAO: simulaFrequênciaDAO{
				simulaResgatar: func(id int64) (frequência, error) {
					return frequência{
						ID:                7654,
						Controle:          918273645,
						CR:                "123456789",
						Calibre:           ".380",
						ArmaUtilizada:     "Arma do Clube",
						NúmeroSérie:       "ZA785671",
						GuiaDeTráfego:     "XYZ12345",
						QuantidadeMunição: 50,
						DataInício:        data.Add(-60 * time.Minute),
						DataTérmino:       data.Add(-30 * time.Minute),
						DataCriação:       data.Add(-21 * time.Minute),
						ImagemNúmeroControle: `TWFuIGlzIGRpc3Rpbmd1aXNoZWQsIG5vdCBvbmx5IGJ5IGhpcyByZWFzb24sIGJ1dCBieSB0aGlz
IHNpbmd1bGFyIHBhc3Npb24gZnJvbSBvdGhlciBhbmltYWxzLCB3aGljaCBpcyBhIGx1c3Qgb2Yg
dGhlIG1pbmQsIHRoYXQgYnkgYSBwZXJzZXZlcmFuY2Ugb2YgZGVsaWdodCBpbiB0aGUgY29udGlu
dWVkIGFuZCBpbmRlZmF0aWdhYmxlIGdlbmVyYXRpb24gb2Yga25vd2xlZGdlLCBleGNlZWRzIHRo
ZSBzaG9ydCB2ZWhlbWVuY2Ugb2YgYW55IGNhcm5hbCBwbGVhc3VyZS4=`,
					}, nil
				},
			},
			erroEsperado: protocolo.NovasMensagens(
				protocolo.NovaMensagemComValor(protocolo.MensagemCódigoPrazoConfirmaçãoExpirado,
					data.Add(-21*time.Minute).Add(20*time.Minute).Format(time.RFC3339)),
			),
		},
		{
			descrição: "deve detectar um erro ao persistir a frequência existente",
			configuração: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.PrazoConfirmação = 20 * time.Minute
				return configuração
			}(),
			frequênciaConfirmaçãoPedidoCompleta: protocolo.FrequênciaConfirmaçãoPedidoCompleta{
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
			frequênciaDAO: simulaFrequênciaDAO{
				simulaAtualizar: func(*frequência) error {
					return errors.Errorf("erro ao atualizar")
				},
				simulaResgatar: func(id int64) (frequência, error) {
					return frequência{
						ID:                7654,
						Controle:          918273645,
						CR:                "123456789",
						Calibre:           ".380",
						ArmaUtilizada:     "Arma do Clube",
						NúmeroSérie:       "ZA785671",
						GuiaDeTráfego:     "XYZ12345",
						QuantidadeMunição: 50,
						DataInício:        data.Add(-40 * time.Minute),
						DataTérmino:       data.Add(-10 * time.Minute),
						DataCriação:       data.Add(-5 * time.Minute),
						ImagemNúmeroControle: `TWFuIGlzIGRpc3Rpbmd1aXNoZWQsIG5vdCBvbmx5IGJ5IGhpcyByZWFzb24sIGJ1dCBieSB0aGlz
IHNpbmd1bGFyIHBhc3Npb24gZnJvbSBvdGhlciBhbmltYWxzLCB3aGljaCBpcyBhIGx1c3Qgb2Yg
dGhlIG1pbmQsIHRoYXQgYnkgYSBwZXJzZXZlcmFuY2Ugb2YgZGVsaWdodCBpbiB0aGUgY29udGlu
dWVkIGFuZCBpbmRlZmF0aWdhYmxlIGdlbmVyYXRpb24gb2Yga25vd2xlZGdlLCBleGNlZWRzIHRo
ZSBzaG9ydCB2ZWhlbWVuY2Ugb2YgYW55IGNhcm5hbCBwbGVhc3VyZS4=`,
					}, nil
				},
			},
			erroEsperado: errors.Errorf("erro ao atualizar"),
		},
	}

	daoOriginal := novaFrequênciaDAO
	defer func() {
		novaFrequênciaDAO = daoOriginal
	}()

	for i, cenário := range cenários {
		novaFrequênciaDAO = func(sqlogger *bd.SQLogger) frequênciaDAO {
			return cenário.frequênciaDAO
		}

		serviço := NovoServiço(nil, cenário.configuração)
		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(nil, cenário.erroEsperado)

		err := serviço.ConfirmarFrequência(cenário.frequênciaConfirmaçãoPedidoCompleta)
		if err = verificadorResultado.VerificaResultado(nil, err); err != nil {
			t.Error(err)
		}
	}
}

type simulaFrequênciaDAO struct {
	simulaCriar     func(*frequência) error
	simulaAtualizar func(*frequência) error
	simulaResgatar  func(id int64) (frequência, error)
}

func (s simulaFrequênciaDAO) criar(frequência *frequência) error {
	return s.simulaCriar(frequência)
}

func (s simulaFrequênciaDAO) atualizar(frequência *frequência) error {
	return s.simulaAtualizar(frequência)
}

func (s simulaFrequênciaDAO) resgatar(id int64) (frequência, error) {
	return s.simulaResgatar(id)
}

const imagemLogoPNG = `
iVBORw0KGgoAAAANSUhEUgAAAKgAAACoCAMAAABDlVWGAAABI1BMVEX/////////////////////
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
YII=`

const imagemNúmeroControlePNG = `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AAAAAAAAAACwZ6yuro69AlM8oQ0AAAAAAAAAAAAAAACACEEbAAAAAAAAAAAAAAAAABGCNgAAAAAAAAAAAAAAAAAiBG0AAAAAAAAAAAAAAAAARAjaAAAAAAAAAAAAAAAAAIgQtAEAAAAAAAAAAAAAAAAQIWgDAAAAAAAAAAAAAAAAIELQBgAAAAAAAAAAAAAAAECEoA0AAAAAAAAAAAAAAACACEEbAAAAAAAAAAAAAAAAABGCNgAAAAAAAAAAAAAAAAAiBG0AAAAAAAAAAAAAAAAARAjaAAAAAAAAAAAAAAAAAIgQtAEAAAAAAAAAAAAAAAAQIWgDAAAAAAAAAAAAAAAAIELQBgAAAAAAAAAAAAAAAECEoA0AAAAAAAAAAAAAAACACEEbAAAAAAAAAAAAAAAAABGCNgAAAAAAAAAAAAAAAAAiBG0AAAAAAAAAAAAAAAAARAjaAAAAAAAAAAAAAAAAAIgQtAEAAAAAAAAAAAAAAAAQIWgDAAAAAAAAAAAAAAAAIELQBgAAAAAAAAAAAAAAAECEoA0AAAAAAAAAAAAAAACACEEbAAAAAAAAAAAAAAAAABGCNgAAAAAAAAAAAAAAAAAiBG0AAAAAAAAAAAAAAAAARAjaAAAAAAAAAAAAAAAAAIgQtAEAAAAAAAAAAAAAAAAQIWgDAAAAAAAAAAAAAAAAIELQBgAAAAAAAAAAAAAAAECEoA0AAAAAAAAAAAAAAACACEEbAAAAAAAAAAAAAAAAABGCNgAAAAAAAAAAAAAAAAAiBG0AAAAAAAAAAAAAAAAARAjaAAAAAAAAAAAAAAAAAIgQtAEAAAAAAAAAAAAAAAAQIWgDAAAAAAAAAAAAAAAAIELQBgAAAAAAAAAAAAAAAECEoA0AAAAAAAAAAAAAAACACEEbAAAAAAAAAAAAAAAAABGCNgAAAAAAAAAAAAAAAAAiBG0AAAAAAAAAAAAAAAAARAjaAAAAAAAAAAAAAAAAAIgQtAEAAAAAAAAAAAAAAAAQIWgDAAAAAAAAAAAAAAAAIELQBgAAAAAAAAAAAAAAAECEoA0AAAAAAAAAAAAAAACACEEbAAAAAAAAAAAAAAAAABGCNgAAAAAAAAAAAAAAAAAiBG0AAAAAAAAAAAAAAAAARAjaAAAAAAAAAAAAAAAAAIgQtAEAAAAAAAAAAAAAAAAQIWgDAAAAAAAAAAAA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AAAAAAAAIgRtAAAAAAAAAAAAAAAAAEQI2gAAAAAAAAAAAAAAAACIELQBAAAAAAAAAAAAAAAAECFoAwAAAAAAAAAAAAAAACBC0AYAAAAAAAAAAAAAAABAhKANAAAAAAAAAAAAAAAAgAhBGwAAAAAAAAAAAAAAAAARgjYAAAAAAAAAAAAAAAAAIgRtAAAAAAAAAAAAAAAAAEQI2gAAAAAAAAAAAAAAAACIELQBAAAAAAAAAAAAAAAAECFoAwAAAAAAAAAAAAAAACBC0AYAAAAAAAAAAAAAAABAhKANAAAAAAAAAAAAAAAAgAhBGwAAAAAAAAAAAAAAAAARgjYAAAAAAAAAAAAAAAAAIgRtAAAAAAAAAAAAAAAAAEQI2gAAAAAAAAAAAAAAAACIELQBAAAAAAAAAAAAAAAAECFoAwAAAAAAAAAAAAAAACBC0AYAAAAAAAAAAAAAAABAhKANAAAAAAAAAAAAAAAAgAhBGwAAAAAAAAAAAAAAAAARgjYAAAAAAAAAAAAAAAAAIgRtAAAAAAAAAAAAAAAAAEQI2gAAAAAAAAAAAAAAAACIELQBAAAAAAAAAAAAAAAAECFoAwAAAAAAAAAAAAAAACBC0AYAAAAAAAAAAAAAAABAhKANAAAAAAAAAAAAAAAAgAhBGwAAAAAAAAAAAAAAAAARgjYAAAAAAAAAAAAAAAAAIgRtAAAAAAAAAAAAAAAAAEQI2gAAAAAAAAAAAAAAAACIELQBAAAAAAAAAAAAAAAAECFoAwAAAAAAAAAAAAAAACBC0AYAAAAAAAAAAAAAAABAhKANAAAAAAAAAAAAAAAAgAhBGwAAAAAAAAAAAAAAAAARgjYAAAAAAAAAAAAAAAAAIgRtAAAAAAAAAAAAAAAAAEQI2gAAAAAAAAAAAAAAAACIELQBAAAAAAAAAAAAAAAAECFoAwAAAAAAAAAAAAAAACBC0AYAAAAAAAAAAAAAAABAhKANAAAAAAAAAAAAAAAAgAhBGwAAAAAAAAAAAAAAAAARgjYAAAAAAAAAAAAAAAAAIgRtAAAAAAAAAAAAAAAAAEQI2gAAAAAAAAAAAAAAAACIELQBAAAAAAAAAAAAAAAAECFoAwAAAAAAAAAAAAAAACBC0AYAAAAAAAAAAAAAAABAhKANAAAAAAAAAAAAAAAAgAhBGwAAAAAAAAAAAAAAAAARgjYAAAAAAAAAAAAAAAAAIgRtAAAAAAAAAAAAAAAAAEQI2gAAAAAAAAAAAAAAAACIELQBAAAAAAAAAAAAAAAAECFoAwAAAAAAAAAAAAAAACBC0AYAAAAAAAAAAAAAAABAhKANAAAAAAAAAAAAAAAAgAhBGwAAAAAAAAAAAAAAAAARgjYAAAAAAAAAAAAAAAAAIgRtAAAAAAAAAAAAAAAAAEQI2gAAAAAAAAAAAAAAAACIELQBAAAAAAAAAAAAAAAAECFoAwAAAAAAAAAAAAAAACBC0AYAAAAAAAAAAAAAAABAhKANAAAAAAAAAAAAAAAAgAhBGwAAAAAAAAAAAAAAAAARgjYAAAAAAAAAAAAAAAAAIgRtAAAAAAAAAAAAAAAAAEQI2gAAAAAAAAAAAAAAAACIELQBAAAAAAAAAAAAAAAAECFoAwAAAAAAAAAAAAAAACBC0AYAAAAAAAAAAAAAAABAhKANAAAAAAAAAAAAAAAAgAhBGwAAAAAAAAAAAAAAAAARgjYAAAAAAAAAAAAAAAAAIgRtAAAAAAAAAAAAAAAAAEQI2gAAAAAAAAAAAAAAAACIELQBAAAAAAAAAAAAAAAAECFoAwAAAAAAAAAAAAAAACBC0AYAAAAAAAAAAAAAAABAhKANAAAAAAAAAAAAAAAAgAhBGwAAAAAAAAAAAAAAAAARgjYAAAAAAAAAAAAAAAAAIgRtAAAAAAAAAAAAAAAAAEQI2gAAAAAAAAAAAAAAAACIELQBAAAAAAAAAAAAAAAAECFoAwAAAAAAAAAAAAAAACBC0AYAAAAAAAAAAAAAAABAhKANAAAAAAAAAAAAAAAAgAhBGwAAAAAAAAAAAAAAAAARgjYAAAAAAAAAAAAAAAAAIgRtAAAAAAAAAAAAAAAAAEQI2gAAAAAAAAAAAAAAAACIELQBAAAAAAAAAAAAAAAAECFoAwAAAAAAAAAAAAAAACBC0AYAAAAAAAAAAAAAAABAhKANAAAAAAAAAAAAAAAAgAhBGwAAAAAAAAAAAAAAAAARgjYAAAAAAAAAAAAAAAAAIgRtAAAAAAAAAAAAAAAAAEQI2gAAAAAAAAAAAAAAAACIELQBAAAAAAAAAAAAAAAAECFoAwAAAAAAAAAAAAAAACBC0AYAAAAAAAAAAAAAAABAhKANAAAAAAAAAAAAAAAAgAhBGwAAAAAAAAAAAAAAAAARgjYAAAAAAAAAAAAAAAAAIgRtAAAAAAAAAAAAAAAAAEQI2gAAAAAAAAAAAAAAAACIELQBAAAAAAAAAAAAAAAAECFoAwAAAAAAAAAAAAAAACBC0AYAAAAAAAAAAAAAAABAhKANAAAAAAAAAAAAAAAAgAhBGwAAAAAAAAAAAAAAAAARgjYAAAAAAAAAAAAAAAAAIgRtAAAAAAAAAAAAAAAAAEQI2gAAAAAAAAAAAAAAAACIELQBAAAAAAAAAAAAAAAAECFoAwAAAAAAAAAAAAAAACBC0AYAAAAAAAAAAAAAAABAhKANAAAAAAAAAAAAAAAAgAhBGwAAAAAAAAAAAAAAAAARgjYAAAAAAAAAAAAAAAAAIgRtAAAAAAAAAAAAAAAAAEQI2gAAAAAAAAAAAAAAAACIELQBAAAAAAAAAAAAAAAAECFoAwAAAAAAAAAAAAAAACBC0AYAAAAAAAAAAAAAAABAhKANAAAAAAAAAAAAAAAAgAhBGwAAAAAAAAAAAAAAAAARgjYAAAAAAAAAAAAAAAAAIgRtAAAAAAAAAAAAAAAAAEQI2gAAAAAAAAAAAAAAAACIELQBAAAAAAAAAAAAAAAAECFoAwAAAAAAAAAAAAAAACBC0AYAAAAAAAAAAAAAAABAhKANAAAAAAAAAAAAAAAAgAhBGwAAAAAAAAAAAAAAAAARgjYAAAAAAAAAAAAAAAAAIgRtAAAAAAAAAAAAAAAAAEQI2gAAAAAAAAAAAAAAAACIELQBAAAAAAAAAAAAAAAAECFoAwAAAAAAAAAAAAAAACBC0AYAAAAAAAAAAAAAAABAhKANAAAAAAAAAAAAAAAAgAhBGwAAAAAAAAAAAAAAAAARgjYAAAAAAAAAAAAAAAAAIgRtAAAAAAAAAAAAAAAAAEQI2gAAAAAAAAAAAAAAAACIELQBAAAAAAAAAAAAAAAAECFoAwAAAAAAAAAAAAAAACBC0AYAAAAAAAAAAAAAAABAhKANAAAAAAAAAAAAAAAAgAhBGwAAAAAAAAAAAAAAAAARgjYAAAAAAAAAAAAAAAAAIgRtAAAAAAAAAAAAAAAAAEQI2gAAAAAAAAAAAAAAAACIELQBAAAAAAAAAAAAAAAAECFoAwAAAAAAAAAAAAAAACBC0AYAAAAAAAAAAAAAAABAhKANAAAAAAAAAAAAAAAAgAhBGwAAAAAAAAAAAAAAAAARgjYAAAAAAAAAAAAAAAAAIgRtAAAAAAAAAAAAAAAAAEQI2gAAAAAAAAAAAAAAAACIELQBAAAAAAAAAAAAAAAAECFoAwAAAAAAAAAAAAAAACBC0AYAAAAAAAAAAAAAAABAhKANAAAAAAAAAAAAAAAAgAhBGwAAAAAAAAAAAAAAAAARgjYAAAAAAAAAAAAAAAAAIgRtAAAAAAAAAAAAAAAAAEQI2gAAAAAAAAAAAAAAAACIELQBAAAAAAAAAAAAAAAAECFoAwAAAAAAAAAAAAAAACBC0AYAAAAAAAAAAAAAAABAhKANAAAAAAAAAAAAAACAf7JvxwIAAAAAg/ytZ7GrPAJYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAg9u1YAAAAAGCQv/UsdpVHAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAMTeHeSoDYMBFG7QnID7n5ErtGURaRQlI0qjZ8Dft7Rl8++tRwAAAAASgjYAAAAAAAAAAAAAAAAAEoI2AAAAAAAAAAAAAAAAABKCNgAAAAAAAAAAAAAAAAASgjYAAAAAAAAAAAAAAAAAEoI2AAAAAAAAAAAAAAAAABKCNgAAAAAAAAAAAAAAAAASgjYAAAAAAAAAAAAAAAAAEoI2AAAAAAAAAAAAAAAAABKCNgAAAAAAAAAAAAAAAAASgjYAAAAAAAAAAAAAAAAAEoI2AAAAAAAAAAAAAAAAABKCNgAAAAAAAAAAAAAAAAASgjYAAAAAAAAAAAAAAAAAEoI2AAAAAAAAAAAAAAAAABKCNgAAAAAAAAAAAAAAAAASgjYAAAAAAAAAAAAAAAAAEoI2AAAAAAAAAAAAAAAAABKCNgAAAAAAAAAAAAAAAAASgjYAAAAAAAAAAAAAAAAAEoI2AAAAAAAAAAAAAAAAABKCNgAAAAAAAAAAAAAAAAASgjYAAAAAAAAAAAAAAAAAEoI2AAAAAAAAAAAAAAAAABKCNgAAAAAAAAAAAAAAAAASgjYAAAAAAAAAAAAAAAAAEoI2AAAAAAAAAAAAAAAAABKCNgAAAAAAAAAAAAAAAAASgjYAAAAAAAAAAAAAAAAAEoI2AAAAAAAAAAAAAAAAABKCNgAAAAAAAAAAAAAAAAASgjYAAAAAAAAAAAAAAAAAEoI2AAAAAAAAAAAAAAAAABKCNgAAAAAAAAAAAAAAAAASgjYAAAAAAAAAAAAAAAAAEoI2AAAAAAAAAAAAAAAAABKCNgAAAAAAAAAAAAAAAAASgjYAAAAAAAAAAAAAAAAAEoI2AAAAAAAAAAAAAAAAABKCNgAAAAAAAAAAAAAAAAASgjYAAAAAAAAAAAAAAAAAEoI2AAAAAAAAAAAAAAAAABKCNgAAAAAAAAAAAAAAAAASgjYAAAAAAAAAAAAAAAAAEoI2AAAAAAAAAAAAAAAAABKCNgAAAAAAAAAAAAAAAAASgjYAAAAAAAAAAAAAAAAAEoI2AAAAAAAAAAAAAAAAABKCNgAAAAAAAAAAAAAAAAASgjYAAAAAAAAAAAAAAAAAEoI2AAAAAAAAAAAAAAAAABKCNgAAAAAAAAAAAAAAAAASgjYAAAAAAAAAAAAAAAAAEoI2AAAAAAAAAAAAAAAAABKCNgAAAAAAAAAAAAAAAAASgjYAAAAAAAAAAAAAAAAAEoI2AAAAAAAAAAAAAAAAABKCNgAAAAAAAAAAAAAAAAASgjYAAAAAAAAAAAAAAAAAEoI2AAAAAAAAAAAAAAAAABKCNgAAAAAAAAAAAAAAAAASgjYAAAAAAAAAAAAAAAAAEoI2AAAAAAAAAAAAAAAAABKCNgAAAAAAAAAAAAAAAAASgjYAAAAAAAAAAAAAAAAAEoI2AAAAAAAAAAAAAAAAABKCNgAAAAAAAAAAAAAAAAASgjYAAAAAAAAAAAAAAAAAEoI2AAAAAAAAAAAAAAAAABJfowcAAAAAAACAka7X66/b7XbKufvad+v+dn27DwAAAAAAALMQtAEAAAAAADCto9DsmXNHgdt9bS9ce/a3AQAAAAAA4J1dRg8AAAAAAAAAI5z5ZbYRdwAAAAAAAMA7ErQBAAAAAAAwpWeDsqNzj94nZgMAAAAAAGBmX6MHAAAAAAAAgFdzj862/jVC2wvX1ntFbQAAAAAAAMxK0AYAAAAAAAAb/xObrdHa3h3f10RtAAAAAAAAzEjQBgAAAAAAACcRqQEAAAAAAMDPLqMHAAAAAAAAgE/wU8y2frUNAAAAAAAAZrf8/mv0EGdYlmX0CKfwj50AAAAAAACNvcjskbeao3NH0dp653bfuxAAAAAAAEDjU/588EMyMEHbq/FwCQAAAAAAAAAAAAAAAOcRtL2Wy+gBAAAAAAAAAAAAAAAAAJiDoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAAAAAAAAgISgDQAAAAAAAAAAAAAAAICEoA0AAAAAAAAAAAAAAACAhKANAAAAAAAAAPjDvh0LAAAAAAzyt57FrvIIAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAgNi3YwEAAACAQf7Ws9hVHgEAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAMcnV0oAACaLSURBVAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAQ+3YsAAAAADDI33oWu8ojAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAA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`
