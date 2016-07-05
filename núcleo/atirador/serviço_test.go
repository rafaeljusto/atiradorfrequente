package atirador

import (
	"testing"
	"time"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/bd"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/config"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/erros"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
	"github.com/rafaeljusto/atiradorfrequente/testes"
	"github.com/registrobr/gostk/errors"
)

func TestServiço_CadastrarFrequência(t *testing.T) {
	data := time.Now()

	cenários := []struct {
		descrição                string
		configuração             config.Configuração
		frequênciaPedidoCompleta protocolo.FrequênciaPedidoCompleta
		frequênciaDAO            frequênciaDAO
		esperado                 protocolo.FrequênciaPendenteResposta
		erroEsperado             error
	}{
		{
			descrição: "deve cadastrar corretamente uma frequência",
			configuração: func() config.Configuração {
				return config.Configuração{}
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
			esperado: protocolo.FrequênciaPendenteResposta{
				NúmeroControle: protocolo.NovoNúmeroControle(1, 123),
			},
		},
		{
			descrição: "deve detectar um erro ao persistir uma nova frequência",
			configuração: func() config.Configuração {
				return config.Configuração{}
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
