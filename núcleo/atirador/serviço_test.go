package atirador

import (
	"bytes"
	"encoding/base64"
	"image"
	_ "image/png"
	"strings"
	"testing"
	"testing/quick"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/bd"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/config"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/erros"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
	"github.com/rafaeljusto/atiradorfrequente/testes"
	"github.com/registrobr/gostk/errors"
	"golang.org/x/image/font/gofont/goregular"
)

func TestServiço_CadastrarFrequência(t *testing.T) {
	data := time.Now()

	imagemBaseExtraída, err := base64.StdEncoding.DecodeString(imagemBasePNG)
	if err != nil {
		t.Fatalf("Erro ao extrair a imagem base de teste. Detalhes: %s", err)
	}

	imagemBaseBuffer := bytes.NewBuffer(imagemBaseExtraída)
	imagemBase, _, err := image.Decode(imagemBaseBuffer)

	if err != nil {
		t.Fatalf("Erro ao interpretar imagem. Detalhes: %s", err)
	}

	imagemBaseInválida := image.NewNRGBA(image.Rect(0, 0, 0, 0))

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
				var configuração config.Configuração
				configuração.Atirador.TempoMáximoCadastro = 12 * time.Hour
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				configuração.Atirador.ImagemNúmeroControle.ImagemBase.Image = imagemBase
				configuração.Atirador.ImagemNúmeroControle.Fonte.Font, err = truetype.Parse(goregular.TTF)

				if err != nil {
					t.Fatalf("Erro ao extrair a fonte de teste. Detalhes: %s", err)
				}

				return configuração
			}(),
			frequênciaPedidoCompleta: protocolo.FrequênciaPedidoCompleta{
				CR: 123456789,
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
				Imagem:         strings.Replace(imagemNúmeroControlePNG, "\n", "", -1),
			},
		},
		{
			descrição: "deve detectar quando o prazo de cadastro do treino já passou",
			configuração: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.TempoMáximoCadastro = 12 * time.Hour
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				return configuração
			}(),
			frequênciaPedidoCompleta: protocolo.FrequênciaPedidoCompleta{
				CR: 1234,
				FrequênciaPedido: protocolo.FrequênciaPedido{
					Calibre:           "",
					ArmaUtilizada:     "Arma do Clube",
					NúmeroSérie:       "XZ23456",
					GuiaDeTráfego:     8734500,
					QuantidadeMunição: 50,
					DataInício:        data.Add(-13 * time.Hour),
					DataTérmino:       data.Add(-12 * time.Hour),
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
			erroEsperado: protocolo.Mensagens{
				protocolo.NovaMensagem(protocolo.MensagemCódigoTempoMáximaCadastroExcedido),
			},
		},
		{
			descrição: "deve detectar quando o tempo de duração máxima do treino é excedida",
			configuração: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.TempoMáximoCadastro = 12 * time.Hour
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				configuração.Atirador.ImagemNúmeroControle.ImagemBase.Image = imagemBase
				configuração.Atirador.ImagemNúmeroControle.Fonte.Font, err = truetype.Parse(goregular.TTF)

				if err != nil {
					t.Fatalf("Erro ao extrair a fonte de teste. Detalhes: %s", err)
				}

				return configuração
			}(),
			frequênciaPedidoCompleta: protocolo.FrequênciaPedidoCompleta{
				CR: 1234,
				FrequênciaPedido: protocolo.FrequênciaPedido{
					Calibre:           "",
					ArmaUtilizada:     "Arma do Clube",
					NúmeroSérie:       "XZ23456",
					GuiaDeTráfego:     8734500,
					QuantidadeMunição: 50,
					DataInício:        data.Add(-13 * time.Hour),
					DataTérmino:       data,
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
			erroEsperado: protocolo.Mensagens{
				protocolo.NovaMensagem(protocolo.MensagemCódigoTreinoMuitoLongo),
			},
		},
		{
			descrição: "deve detectar um erro ao persistir uma nova frequência",
			configuração: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.TempoMáximoCadastro = 12 * time.Hour
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				configuração.Atirador.ImagemNúmeroControle.ImagemBase.Image = imagemBase
				configuração.Atirador.ImagemNúmeroControle.Fonte.Font, err = truetype.Parse(goregular.TTF)

				if err != nil {
					t.Fatalf("Erro ao extrair a fonte de teste. Detalhes: %s", err)
				}

				return configuração
			}(),
			frequênciaPedidoCompleta: protocolo.FrequênciaPedidoCompleta{
				CR: 123456789,
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
				configuração.Atirador.TempoMáximoCadastro = 12 * time.Hour
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				configuração.Atirador.ImagemNúmeroControle.ImagemBase.Image = imagemBaseInválida
				configuração.Atirador.ImagemNúmeroControle.Fonte.Font, err = truetype.Parse(goregular.TTF)

				if err != nil {
					t.Fatalf("Erro ao extrair a fonte de teste. Detalhes: %s", err)
				}

				return configuração
			}(),
			frequênciaPedidoCompleta: protocolo.FrequênciaPedidoCompleta{
				CR: 123456789,
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
			descrição: "deve detectar quando a fonte da imagem não esta definida",
			configuração: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.TempoMáximoCadastro = 12 * time.Hour
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				configuração.Atirador.ImagemNúmeroControle.ImagemBase.Image = imagemBase
				return configuração
			}(),
			frequênciaPedidoCompleta: protocolo.FrequênciaPedidoCompleta{
				CR: 123456789,
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
			erroEsperado: errors.Errorf("fonte da imagem do número de controle indefinida"),
		},
		{
			descrição: "deve detectar um erro ao atualizar uma frequência",
			configuração: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.TempoMáximoCadastro = 12 * time.Hour
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
				configuração.Atirador.ImagemNúmeroControle.ImagemBase.Image = imagemBase
				configuração.Atirador.ImagemNúmeroControle.Fonte.Font, err = truetype.Parse(goregular.TTF)

				if err != nil {
					t.Fatalf("Erro ao extrair a fonte de teste. Detalhes: %s", err)
				}

				return configuração
			}(),
			frequênciaPedidoCompleta: protocolo.FrequênciaPedidoCompleta{
				CR: 123456789,
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

func TestServiço_CadastrarFrequência_valoresAleatórios(t *testing.T) {
	imagemBaseExtraída, err := base64.StdEncoding.DecodeString(imagemBasePNG)
	if err != nil {
		t.Fatalf("Erro ao extrair a imagem base de teste. Detalhes: %s", err)
	}

	imagemBaseBuffer := bytes.NewBuffer(imagemBaseExtraída)
	imagemBase, _, err := image.Decode(imagemBaseBuffer)

	if err != nil {
		t.Fatalf("Erro ao interpretar imagem. Detalhes: %s", err)
	}

	daoOriginal := novaFrequênciaDAO
	defer func() {
		novaFrequênciaDAO = daoOriginal
	}()

	novaFrequênciaDAO = func(sqlogger *bd.SQLogger) frequênciaDAO {
		return simulaFrequênciaDAO{
			simulaCriar: func(frequência *frequência) error {
				frequência.ID = 1
				frequência.Controle = 123
				return nil
			},
			simulaAtualizar: func(frequência *frequência) error {
				return nil
			},
		}
	}

	// não utilizamos diretamente o objeto do protocolo, pois a biblioteca padrão
	// não sabe preencher corretamente o tipo time.Time.
	f := func(cr int, calibre, armaUtilizada, númeroSérie string, guiaDeTráfego, quantidadeMunição int, dataInício, dataTérmino int64) bool {
		var configuração config.Configuração
		configuração.Atirador.TempoMáximoCadastro = 12 * time.Hour
		configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
		configuração.Atirador.ImagemNúmeroControle.ImagemBase.Image = imagemBase
		configuração.Atirador.ImagemNúmeroControle.Fonte.Font, err = truetype.Parse(goregular.TTF)

		if err != nil {
			t.Fatalf("Erro ao extrair a fonte de teste. Detalhes: %s", err)
		}

		frequênciaPedidoCompleta := protocolo.FrequênciaPedidoCompleta{
			CR: cr,
			FrequênciaPedido: protocolo.FrequênciaPedido{
				Calibre:           calibre,
				ArmaUtilizada:     armaUtilizada,
				NúmeroSérie:       númeroSérie,
				GuiaDeTráfego:     guiaDeTráfego,
				QuantidadeMunição: quantidadeMunição,
				DataInício:        time.Unix(dataInício, 0),
				DataTérmino:       time.Unix(dataTérmino, 0),
			},
		}

		serviço := NovoServiço(nil, configuração)
		if _, err := serviço.CadastrarFrequência(frequênciaPedidoCompleta); err != nil {
			if _, ok := err.(protocolo.Mensagens); !ok {
				t.Log(err)
				return false
			}
		}

		return true
	}

	if err := quick.Check(f, &quick.Config{MaxCount: 5}); err != nil {
		t.Error(err)
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
				CR:             123456789,
				NúmeroControle: protocolo.NovoNúmeroControle(7654, 918273645),
				FrequênciaConfirmaçãoPedido: protocolo.FrequênciaConfirmaçãoPedido{
					Imagem: `iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKCAMAAAC67D+PAAAAP1BMVEX///8AezAAzhcIziD//5sA
aygIzos5zoPGpQAArQAArSj/zpsxzkkAWgBCnAAAlBcAvQAApTi1zgApjACMYwCTUqAuAAAAT0lE
QVQImR2MyQ3AMAzDpNjO3bv7z1o1ehEiQABIGv6d/SC3vgu7uTEzZC93HyxRkWK9ozShLJObcMuR
7fZZAWOx4ZMqPIxik+8q19Zk8QFkhgHrQUAyGgAAAABJRU5ErkJggg==`,
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
						CR:                123456789,
						Calibre:           ".380",
						ArmaUtilizada:     "Arma do Clube",
						NúmeroSérie:       "ZA785671",
						GuiaDeTráfego:     762556223,
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
				CR:             123456789,
				NúmeroControle: protocolo.NovoNúmeroControle(7654, 918273645),
				FrequênciaConfirmaçãoPedido: protocolo.FrequênciaConfirmaçãoPedido{
					Imagem: `iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKCAMAAAC67D+PAAAAP1BMVEX///8AezAAzhcIziD//5sA
aygIzos5zoPGpQAArQAArSj/zpsxzkkAWgBCnAAAlBcAvQAApTi1zgApjACMYwCTUqAuAAAAT0lE
QVQImR2MyQ3AMAzDpNjO3bv7z1o1ehEiQABIGv6d/SC3vgu7uTEzZC93HyxRkWK9ozShLJObcMuR
7fZZAWOx4ZMqPIxik+8q19Zk8QFkhgHrQUAyGgAAAABJRU5ErkJggg==`,
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
				CR:             123456781,
				NúmeroControle: protocolo.NovoNúmeroControle(7654, 918273640),
				FrequênciaConfirmaçãoPedido: protocolo.FrequênciaConfirmaçãoPedido{
					Imagem: `iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKCAMAAAC67D+PAAAAP1BMVEX///8AezAAzhcIziD//5sA
aygIzos5zoPGpQAArQAArSj/zpsxzkkAWgBCnAAAlBcAvQAApTi1zgApjACMYwCTUqAuAAAAT0lE
QVQImR2MyQ3AMAzDpNjO3bv7z1o1ehEiQABIGv6d/SC3vgu7uTEzZC93HyxRkWK9ozShLJObcMuR
7fZZAWOx4ZMqPIxik+8q19Zk8QFkhgHrQUAyGgAAAABJRU5ErkJggg==`,
				},
			},
			frequênciaDAO: simulaFrequênciaDAO{
				simulaResgatar: func(id int64) (frequência, error) {
					return frequência{
						ID:                7654,
						Controle:          918273645,
						CR:                123456789,
						Calibre:           ".380",
						ArmaUtilizada:     "Arma do Clube",
						NúmeroSérie:       "ZA785671",
						GuiaDeTráfego:     762556223,
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
				protocolo.NovaMensagemComValor(protocolo.MensagemCódigoCRInválido, "123456781"),
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
				CR:             123456789,
				NúmeroControle: protocolo.NovoNúmeroControle(7654, 918273645),
				FrequênciaConfirmaçãoPedido: protocolo.FrequênciaConfirmaçãoPedido{
					Imagem: `iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKCAMAAAC67D+PAAAAP1BMVEX///8AezAAzhcIziD//5sA
aygIzos5zoPGpQAArQAArSj/zpsxzkkAWgBCnAAAlBcAvQAApTi1zgApjACMYwCTUqAuAAAAT0lE
QVQImR2MyQ3AMAzDpNjO3bv7z1o1ehEiQABIGv6d/SC3vgu7uTEzZC93HyxRkWK9ozShLJObcMuR
7fZZAWOx4ZMqPIxik+8q19Zk8QFkhgHrQUAyGgAAAABJRU5ErkJggg==`,
				},
			},
			frequênciaDAO: simulaFrequênciaDAO{
				simulaResgatar: func(id int64) (frequência, error) {
					return frequência{
						ID:                7654,
						Controle:          918273645,
						CR:                123456789,
						Calibre:           ".380",
						ArmaUtilizada:     "Arma do Clube",
						NúmeroSérie:       "ZA785671",
						GuiaDeTráfego:     762556223,
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
				protocolo.NovaMensagem(protocolo.MensagemCódigoPrazoConfirmaçãoExpirado),
			),
		},
		{
			descrição: "deve detectar quando a imagem de confirmação for igual a imagem do número de controle",
			configuração: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.PrazoConfirmação = 20 * time.Minute
				return configuração
			}(),
			frequênciaConfirmaçãoPedidoCompleta: protocolo.FrequênciaConfirmaçãoPedidoCompleta{
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
			frequênciaDAO: simulaFrequênciaDAO{
				simulaResgatar: func(id int64) (frequência, error) {
					return frequência{
						ID:                7654,
						Controle:          918273645,
						CR:                123456789,
						Calibre:           ".380",
						ArmaUtilizada:     "Arma do Clube",
						NúmeroSérie:       "ZA785671",
						GuiaDeTráfego:     762556223,
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
				protocolo.NovaMensagemComValor(protocolo.MensagemCódigoImagemNãoAceita,
					`TWFuIGlzIGRpc3Rpbmd1aXNoZWQsIG5vdCBvbmx5IGJ5IGhpcyByZWFzb24sIGJ1dCBieSB0aGlz
IHNpbmd1bGFyIHBhc3Npb24gZnJvbSBvdGhlciBhbmltYWxzLCB3aGljaCBpcyBhIGx1c3Qgb2Yg
dGhlIG1pbmQsIHRoYXQgYnkgYSBwZXJzZXZlcmFuY2Ugb2YgZGVsaWdodCBpbiB0aGUgY29udGlu
dWVkIGFuZCBpbmRlZmF0aWdhYmxlIGdlbmVyYXRpb24gb2Yga25vd2xlZGdlLCBleGNlZWRzIHRo
ZSBzaG9ydCB2ZWhlbWVuY2Ugb2YgYW55IGNhcm5hbCBwbGVhc3VyZS4=`),
			),
		},
		{
			descrição: "deve detectar quando a frequência já foi confirmada",
			configuração: func() config.Configuração {
				var configuração config.Configuração
				configuração.Atirador.PrazoConfirmação = 20 * time.Minute
				return configuração
			}(),
			frequênciaConfirmaçãoPedidoCompleta: protocolo.FrequênciaConfirmaçãoPedidoCompleta{
				CR:             123456789,
				NúmeroControle: protocolo.NovoNúmeroControle(7654, 918273645),
				FrequênciaConfirmaçãoPedido: protocolo.FrequênciaConfirmaçãoPedido{
					Imagem: `iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKCAMAAAC67D+PAAAAP1BMVEX///8AezAAzhcIziD//5sA
aygIzos5zoPGpQAArQAArSj/zpsxzkkAWgBCnAAAlBcAvQAApTi1zgApjACMYwCTUqAuAAAAT0lE
QVQImR2MyQ3AMAzDpNjO3bv7z1o1ehEiQABIGv6d/SC3vgu7uTEzZC93HyxRkWK9ozShLJObcMuR
7fZZAWOx4ZMqPIxik+8q19Zk8QFkhgHrQUAyGgAAAABJRU5ErkJggg==`,
				},
			},
			frequênciaDAO: simulaFrequênciaDAO{
				simulaResgatar: func(id int64) (frequência, error) {
					if id != 7654 {
						t.Errorf("ID %d inesperado", id)
					}

					return frequência{
						ID:                7654,
						Controle:          918273645,
						CR:                123456789,
						Calibre:           ".380",
						ArmaUtilizada:     "Arma do Clube",
						NúmeroSérie:       "ZA785671",
						GuiaDeTráfego:     762556223,
						QuantidadeMunição: 50,
						DataInício:        data.Add(-40 * time.Minute),
						DataTérmino:       data.Add(-10 * time.Minute),
						DataCriação:       data.Add(-5 * time.Minute),
						DataConfirmação:   data.Add(-2 * time.Minute),
						ImagemNúmeroControle: `TWFuIGlzIGRpc3Rpbmd1aXNoZWQsIG5vdCBvbmx5IGJ5IGhpcyByZWFzb24sIGJ1dCBieSB0aGlz
IHNpbmd1bGFyIHBhc3Npb24gZnJvbSBvdGhlciBhbmltYWxzLCB3aGljaCBpcyBhIGx1c3Qgb2Yg
dGhlIG1pbmQsIHRoYXQgYnkgYSBwZXJzZXZlcmFuY2Ugb2YgZGVsaWdodCBpbiB0aGUgY29udGlu
dWVkIGFuZCBpbmRlZmF0aWdhYmxlIGdlbmVyYXRpb24gb2Yga25vd2xlZGdlLCBleGNlZWRzIHRo
ZSBzaG9ydCB2ZWhlbWVuY2Ugb2YgYW55IGNhcm5hbCBwbGVhc3VyZS4=`,
						ImagemConfirmação: `iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKCAMAAAC67D+PAAAAP1BMVEX///8AezAAzhcIziD//5sA
aygIzos5zoPGpQAArQAArSj/zpsxzkkAWgBCnAAAlBcAvQAApTi1zgApjACMYwCTUqAuAAAAT0lE
QVQImR2MyQ3AMAzDpNjO3bv7z1o1ehEiQABIGv6d/SC3vgu7uTEzZC93HyxRkWK9ozShLJObcMuR
7fZZAWOx4ZMqPIxik+8q19Zk8QFkhgHrQUAyGgAAAABJRU5ErkJggg==`,
					}, nil
				},
			},
			erroEsperado: protocolo.NovasMensagens(
				protocolo.NovaMensagem(protocolo.MensagemCódigoFrequênciaJáConfirmada),
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
				CR:             123456789,
				NúmeroControle: protocolo.NovoNúmeroControle(7654, 918273645),
				FrequênciaConfirmaçãoPedido: protocolo.FrequênciaConfirmaçãoPedido{
					Imagem: `iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKCAMAAAC67D+PAAAAP1BMVEX///8AezAAzhcIziD//5sA
aygIzos5zoPGpQAArQAArSj/zpsxzkkAWgBCnAAAlBcAvQAApTi1zgApjACMYwCTUqAuAAAAT0lE
QVQImR2MyQ3AMAzDpNjO3bv7z1o1ehEiQABIGv6d/SC3vgu7uTEzZC93HyxRkWK9ozShLJObcMuR
7fZZAWOx4ZMqPIxik+8q19Zk8QFkhgHrQUAyGgAAAABJRU5ErkJggg==`,
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
						CR:                123456789,
						Calibre:           ".380",
						ArmaUtilizada:     "Arma do Clube",
						NúmeroSérie:       "ZA785671",
						GuiaDeTráfego:     762556223,
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

func TestServiço_ConfirmarFrequência_valoresAleatórios(t *testing.T) {
	daoOriginal := novaFrequênciaDAO
	defer func() {
		novaFrequênciaDAO = daoOriginal
	}()

	f := func(frequênciaConfirmaçãoPedidoCompleta protocolo.FrequênciaConfirmaçãoPedidoCompleta) bool {
		novaFrequênciaDAO = func(sqlogger *bd.SQLogger) frequênciaDAO {
			return simulaFrequênciaDAO{
				simulaAtualizar: func(frequência *frequência) error {
					return nil
				},
				simulaResgatar: func(id int64) (frequência, error) {
					return frequência{
						ID:                frequênciaConfirmaçãoPedidoCompleta.NúmeroControle.ID(),
						Controle:          frequênciaConfirmaçãoPedidoCompleta.NúmeroControle.Controle(),
						CR:                frequênciaConfirmaçãoPedidoCompleta.CR,
						Calibre:           ".380",
						ArmaUtilizada:     "Arma do Clube",
						NúmeroSérie:       "ZA785671",
						GuiaDeTráfego:     762556223,
						QuantidadeMunição: 50,
						DataInício:        time.Now().Add(-40 * time.Minute),
						DataTérmino:       time.Now().Add(-10 * time.Minute),
						DataCriação:       time.Now().Add(-5 * time.Minute),
						ImagemNúmeroControle: `TWFuIGlzIGRpc3Rpbmd1aXNoZWQsIG5vdCBvbmx5IGJ5IGhpcyByZWFzb24sIGJ1dCBieSB0aGlz
IHNpbmd1bGFyIHBhc3Npb24gZnJvbSBvdGhlciBhbmltYWxzLCB3aGljaCBpcyBhIGx1c3Qgb2Yg
dGhlIG1pbmQsIHRoYXQgYnkgYSBwZXJzZXZlcmFuY2Ugb2YgZGVsaWdodCBpbiB0aGUgY29udGlu
dWVkIGFuZCBpbmRlZmF0aWdhYmxlIGdlbmVyYXRpb24gb2Yga25vd2xlZGdlLCBleGNlZWRzIHRo
ZSBzaG9ydCB2ZWhlbWVuY2Ugb2YgYW55IGNhcm5hbCBwbGVhc3VyZS4=`,
					}, nil
				},
			}
		}

		var configuração config.Configuração
		configuração.Atirador.PrazoConfirmação = 20 * time.Minute

		serviço := NovoServiço(nil, configuração)
		if err := serviço.ConfirmarFrequência(frequênciaConfirmaçãoPedidoCompleta); err != nil {
			if _, ok := err.(protocolo.Mensagens); !ok {
				t.Log(err)
				return false
			}
		}

		return true
	}

	if err := quick.Check(f, &quick.Config{MaxCount: 5}); err != nil {
		t.Error(err)
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

const imagemBasePNG = `
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

const imagemNúmeroControlePNG = `
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`
