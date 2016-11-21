package atirador

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	_ "image/png"
	"strings"
	"testing"
	"testing/quick"
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
				configuração.Atirador.TempoMáximoCadastro = 12 * time.Hour
				configuração.Atirador.DuraçãoMáximaTreino = 12 * time.Hour
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
				configuração.Atirador.ImagemNúmeroControle.Largura = 0
				configuração.Atirador.ImagemNúmeroControle.Altura = 0
				configuração.Atirador.ImagemNúmeroControle.CorFundo.Color = color.RGBA{0xff, 0xff, 0xff, 0xff}
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
	imagemLogoExtraída, err := base64.StdEncoding.DecodeString(imagemLogoPNG)

	if err != nil {
		t.Fatalf("Erro ao extrair a imagem de teste do logo. Detalhes: %s", err)
	}

	imagemLogoBuffer := bytes.NewBuffer(imagemLogoExtraída)
	imagemLogo, _, err := image.Decode(imagemLogoBuffer)

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

const imagemNúmeroControlePNG = `
iVBORw0KGgoAAAANSUhEUgAADbQAAAmwCAYAAAD4UEzwAACAAElEQVR4nOzdT6jl91nH8c8zmZlO
m7HBgZHKrFwVhEpAyUIKQsCVUFEsSMt0Y1AiFqHSgii6EotFN1MaAi0K0lLRjeCiUOhGQRQLgUAg
EBSCgUDkkokTMn/zyMETPTmduXPvzL3Pb86d1wsu9/7+nPN9At/M7s33dHd3AAAAAAAAAAAAAAAA
AOCYnVp6AAAAAAAAAAAAAAAAAAAeD4I2AAAAAAAAAAAAAAAAAEYI2gAAAAAAAAAAAAAAAAAYIWgD
AAAAAAAAAAAAAAAAYISgDQAAAAAAAAAAAAAAAIARgjYAAAAAAAAAAAAAAAAARgjaAAAAAAAAAAAA
AAAAABghaAMAAAAAAAAAAAAAAABghKANAAAAAAAAAAAAAAAAgBGCNgAAAAAAAAAAAAAAAABGCNoA
AAAAAAAAAAAAAAAAGCFoAwAAAAAAAAAAAAAAAGCEoA0AAAAAAAAAAAAAAACAEYI2AAAAAAAAAAAA
AAAAAEYI2gAAAAAAAAAAAAAAAAAYIWgDAAAAAAAAAAAAAAAAYISgDQAAAAAAAAAAAAAAAIARgjYA
AAAAAAAAAAAAAAAARgjaAAAAAAAAAAAAAAAAABghaAMAAAAAAAAAAAAAAABghKANAAAAAAAAAAAA
AAAAgBGCNgAAAAAAAAAAAAAAAABGCNoAAAAAAAAAAAAAAAAAGCFoAwAAAAAAAAAAAAAAAGCEoA0A
AAAAAAAAAAAAAACAEYI2AAAAAAAAAAAAAAAAAEYI2gAAAAAAAAAAAAAAAAAYIWgDAAAAAAAAAAAA
AAAAYISgDQAAAAAAAAAAAAAAAIARgjYAAAAAAAAAAAAAAAAARgjaAAAAAAAAAAAAAAAAABghaAMA
AAAAAAAAAAAAAABghKANAAAAAAAAAAAAAAAAgBGCNgAAAAAAAAAAAAAAAABGCNoAAAAAAAAAAAAA
AAAAGCFoAwAAAAAAAAAAAAAAAGCEoA0AAAAAAAAAAAAAAACAEYI2AAAAAAAAAAAAAAAAAEYI2gAA
AAAAAAAAAAAAAAAYIWgDAAAAAAAAAAAAAAAAYISgDQAAAAAAAAAAAAAAAIARgjYAAAAAAAAAAAAA
AAAARgjaAAAAAAAAAAAAAAAAABghaAMAAAAAAAAAAAAAAABghKANAAAAAAAAAAAAAAAAgBGCNgAA
AAAAAAAAAAAAAABGCNoAAAAAAAAAAAAAAAAAGCFoAwAAAAAAAAAAAAAAAGCEoA0AAAAAAAAAAAAA
AACAEaeXHuCoVNXSIwAAAAAAAAAAAAAAAAAci+5eeoQj4YQ2AAAAAAAAAAAAAAAAAEYI2gAAAAAA
AAAAAAAAAAAYIWgDAAAAAAAAAAAAAAAAYISgDQAAAAAAAAAAAAAAAIARgjYAAAAAAAAAAAAAAAAA
RgjaAAAAAAAAAAAAAAAAABghaAMAAAAAAAAAAAAAAABghKANAAAAAAAAAAAAAAAAgBGCNgAAAAAA
AAAAAAAAAABGCNoAAAAAAAAAAAAAAAAAGCFoAwAAAAAAAAAAAAAAAGCEoA0AAAAAAAAAAAAAAACA
EYI2AAAAAAAAAAAAAAAAAEYI2gAAAAAAAAAAAAAAAAAYIWgDAAAAAAAAAAAAAAAAYISgDQAAAAAA
AAAAAAAAAIARgjYAAAAAAAAAAAAAAAAARgjaAAAAAAAAAAAAAAAAABghaAMAAAAAAAAAAAAAAABg
hKANAAAAAAAAAAAAAAAAgBGCNgAAAAAAAAAAAAAAAABGCNoAAAAAAAAAAAAAAAAAGCFoAwAAAAAA
AAAAAAAAAGCEoA0AAAAAAAAAAAAAAACAEYI2AAAAAAAAAAAAAAAAAEYI2gAAAAAAAAAAAAAAAAAY
IWgDAAAAAAAAAAAAAAAAYISgDQAAAAAAAAAAAAAAAIARgjYAAAAAAAAAAAAAAAAARgjaAAAAAAAA
AAAAAAAAABghaAMAAAAAAAAAAAAAAABghKANAAAAAAAAAAAAAAAAgBGCNgAAAAAAAAAAAAAAAABG
CNoAAAAAAAAAAAAAAAAAGCFoAwAAAAAAAAAAAAAAAGCEoA0AAAAAAAAAAAAAAACAEYI2AAAAAAAA
AAAAAAAAAEYI2gAAAAAAAAAAAAAAAAAYIWgDAAAAAAAAAAAAAAAAYISgDQAAAAAAAAAAAAAAAIAR
gjYAAAAAAAAAAAAAAAAARgjaAAAAAAAAAAAAAAAAABghaAMAAAAAAAAAAAAAAABghKANAAAAAAAA
AAAAAAAAgBGCNgAAAAAAAAAAAAAAAABGCNoAAAAAAAAAAAAAAAAAGCFoAwAAAAAAAAAAAAAAAGCE
oA0AAAAAAAAAAAAAAACAEYI2AAAAAAAAAAAAAAAAAEacXnoAPmxvb2/pEQAAAAAAAAAAAAAAAODE
uHDhwtIjsMEJbQAAAAAAAAAAAAAAAACMELQBAAAAAAAAAAAAAAAAMELQBgAAAAAAAAAAAAAAAMAI
QRsAAAAAAAAAAAAAAAAAIwRtAAAAAAAAAAAAAAAAAIwQtAEAAAAAAAAAAAAAAAAwQtAGAAAAAAAA
AAAAAAAAwAhBGwAAAAAAAAAAAAAAAAAjBG0AAAAAAAAAAAAAAAAAjBC0AQAAAAAAAAAAAAAAADBC
0AYAAAAAAAAAAAAAAADACEEbAAAAAAAAAAAAAAAAACMEbQAAAAAAAAAAAAAAAACMELQBAAAAAAAA
AAAAAAAAMELQBgAAAAAAAAAAAAAAAMAIQRsAAAAAAAAAAAAAAAAAIwRtAAAAAAAAAAAAAAAAAIwQ
tAEAAAAAAAAAAAAAAAAwQtAGAAAAAAAAAAAAAAAAwAhBGwAAAAAAAAAAAAAAAAAjBG0AAAAAAAAA
AAAAAAAAjBC0AQAAAAAAAAAAAAAAADBC0AYAAAAAAAAAAAAAAADACEEbAAAAAAAAAAAAAAAAACME
bQAAAAAAAAAAAAAAAACMELQBAAAAAAAAAAAAAAAAMELQBgAAAAAAAAAAAAAAAMAIQRsAAAAAAAAA
AAAAAAAAIwRtAAAAAAAAAAAAAAAAAIwQtAEAAAAAAAAAAAAAAAAwQtAGAAAAAAAAAAAAAAAAwAhB
GwAAAAAAAAAAAAAAAAAjBG0AAAAAAAAAAAAAAAAAjBC0AQAAAAAAAAAAAAAAADBC0AYAAAAAAAAA
AAAAAADACEEbAAAAAAAAAAAAAAAAACMEbQAAAAAAAAAAAAAAAACMELQBAAAAAAAAAAAAAAAAMELQ
BgAAAAAAAAAAAAAAAMAIQRsAAAAAAAAAAAAAAAAAIwRtAAAAAAAAAAAAAAAAAIwQtAEAAAAAAAAA
AAAAAAAwQtAGAAAAAAAAAAAAAAAAwAhBGwAAAAAAAAAAAAAAAAAjBG0AAAAAAAAAAAAAAAAAjBC0
AQAAAAAAAAAAAAAAADBC0AYAAAAAAAAAAAAAAADACEEbAAAAAAAAAAAAAAAAACMEbQAAAAAAAAAA
AAAAAACMELQBAAAAAAAAAAAAAAAAMELQBgAAAAAAAAAAAAAAAMAIQRsAAAAAAAAAAAAAAAAAIwRt
AAAAAAAAAAAAAAAAAIwQtAEAAAAAAAAAAAAAAAAwQtAGAAAAAAAAAAAAAAAAwAhBGwAAAAAAAAAA
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
AADQQtAGAAAAAAAAAAAAAAAAQAtBGwAAAAAAAAAAAAAAAAAtBG0AAAAAAAAAAAAAAAAAtBC0AQAA
AAAAAAAAAAAAANBC0AYAAAAAAAAAAAAAAABAC0EbAAAAAAAAAAAAAAAAAC0EbQAAAAAAAAAAAAAA
AAC0ELQBAAAAAAAAAAAAAAAA0ELQBgAAAAAAAAAAAAAAAEALQRsAAAAAAAAAAAAAAAAALQRtAAAA
AAAAAAAAAAAAALQQtAEAAAAAAAAAAAAAAADQQtAGAAAAAAAAAAAAAAAAQAtBGwAAAAAAAAAAAAAA
AAAtBG0AAAAAAAAAAAAAAAAAtBC0AQAAAAAAAAAAAAAAANBC0AYAAAAAAAAAAAAAAABAC0EbAAAA
AAAAAAAAAAAAAC0EbQAAAAAAAAAAAAAAAAC0ELQBAAAAAAAAAAAAAAAA0ELQBgAAAAAAAAAAAAAA
AEALQRsAAAAAAAAAAAAAAAAALQRtAAAAAAAAAAAAAAAAALQQtAEAAAAAAAAAAAAAAADQQtAGAAAA
AAAAAAAAAAAAQAtBGwAAAAAAAAAAAAAAAAAtBG0AAAAAAAAAAAAAAAAAtBC0AQAAAAAAAAAAAAAA
ANBC0AYAAAAAAAAAAAAAAABAC0EbAAAAAAAAAAAAAAAAAC0EbQAAAAAAAAAAAAAAAAC0ELQBAAAA
AAAAAAAAAAAA0ELQBgAAAAAAAAAAAAAAAEALQRsAAAAAAAAAAAAAAAAALQRtAAAAAAAAAAAAAAAA
ALQQtAEAAAAAAAAAAAAAAADQQtAGAAAAAAAAAAAAAAAAQAtBGwAAAAAAAAAAAAAAAAAtBG0AAAAA
AAAAAAAAAAAAtBC0AQAAAAAAAAAAAAAAANBC0AYAAAAAAAAAAAAAAABAC0EbAAAAAAAAAAAAAAAA
AC0EbQAAAAAAAAAAAAAAAAC0ELQBAAAAAAAAAAAAAAAA0ELQBgAAAAAAAAAAAAAAAEALQRsAAAAA
AAAAAAAAAAAALQRtAAAAAAAAAAAAAAAAALQQtAEAAAAAAAAAAAAAAADQQtAGAAAAAAAAAAAAAAAA
QAtBGwAAAAAAAAAAAAAAAAAtBG0AAAAAAAAAAAAAAAAAtBC0AQAAAAAAAAAAAAAAANBC0AYAAAAA
AAAAAAAAAABAC0EbAAAAAAAAAAAAAAAAAC0EbQAAAAAAAAAAAAAAAAC0ELQBAAAAAAAAAAAAAAAA
0ELQBgAAAAAAAAAAAAAAAEALQRsAAAAAAAAAAAAAAAAALQRtAAAAAAAAAAAAAAAAALQQtAEAAAAA
AAAAAAAAAADQQtAGAAAAAAAAAAAAAAAAQAtBGwAAAAAAAAAAAAAAAAAtBG0AAAAAAAAAAAAAAAAA
tBC0AQAAAAAAAAAAAAAAANBC0AYAAAAAAAAAAAAAAABAC0EbAAAAAAAAAAAAAAAAAC0EbQAAAAAA
AAAAAAAAAAC0ELQBAAAAAAAAAAAAAAAA0ELQBgAAAAAAAAAAAAAAAEALQRsAAAAAAAAAAAAAAAAA
LQRtAAAAAAAAAAAAAAAAALQQtAEAAAAAAAAAAAAAAADQQtAGAAAAAAAAAAAAAAAAQAtBGwAAAAAA
AAAAAAAAAAAtBG0AAAAAAAAAAAAAAAAAtBC0AQAAAAAAAAAAAAAAANBC0AYAAAAAAAAAAAAAAABA
C0EbAAAAAAAAAAAAAAAAAC0EbQAAAAAAAAAAAAAAAAC0ELQBAAAAAAAAAAAAAAAA0ELQBgAAAAAA
AAAAAAAAAEALQRsAAAAAAAAAAAAAAAAALQRtAAAAAAAAAAAAAAAAALQQtAEAAAAAAAAAAAAAAADQ
QtAGAAAAAAAAAAAAAAAAQAtBGwAAAAAAAAAAAAAAAAAtBG0AAAAAAAAAAAAAAAAAtBC0AQAAAAAA
AAAAAAAAANBC0AYAAAAAAAAAAAAAAABAC0EbAAAAAAAAAAAAAAAAAC0EbQAAAAAAAAAAAAAAAAC0
ELQBAAAAAAAAAAAAAAAA0ELQBgAAAAAAAAAAAAAAAEALQRsAAAAAAAAAAAAAAAAALQRtAAAAAAAA
AAAAAAAAALQQtAEAAAAAAAAAAAAAAADQQtAGAAAAAAAAAAAAAAAAQAtBGwAAAAAAAAAAAAAAAAAt
BG0AAAAAAAAAAAAAAAAAtBC0AQAAAAAAAAAAAAAAANBC0AYAAAAAAAAAAAAAAABAC0EbAAAAAAAA
AAAAAAAAAC0EbQAAAAAAAAAAAAAAAAC0ELQBAAAAAAAAAAAAAAAA0ELQBgAAAAAAAAAAAAAAAEAL
QRsAAAAAAAAAAAAAAAAALQRtAAAAAAAAAAAAAAAAALQQtAEAAAAAAAAAAAAAAADQQtAGAAAAAAAA
AAAAAAAAQAtBGwAAAAAAAAAAAAAAAAAtBG0AAAAAAAAAAAAAAAAAtBC0AQAAAAAAAAAAAAAAANBC
0AYAAAAAAAAAAAAAAABAC0EbAAAAAAAAAAAAAAAAAC0EbQAAAAAAAAAAAAAAAAC0ELQBAAAAAAAA
AAAAAAAA0ELQBgAAAAAAAAAAAAAAAEALQRsAAAAAAAAAAAAAAAAALQRtAAAAAAAAAAAAAAAAALQQ
tAEAAAAAAAAAAAAAAADQQtAGAAAAAAAAAAAAAAAAQAtBGwAAAAAAAAAAAAAAAAAtBG0AAAAAAAAA
AAAAAAAAtBC0AQAAAAAAAAAAAAAAANBC0AYAAAAAAAAAAAAAAABAC0EbAAAAAAAAAAAAAAAAAC0E
bQAAAAAAAAAAAAAAAAC0ELQBAAAAAAAAAAAAAAAA0ELQBgAAAAAAAAAAAAAAAEALQRsAAAAAAAAA
AAAAAAAALQRtAAAAAAAAAAAAAAAAALQQtAEAAAAAAAAAAAAAAADQQtAGAAAAAAAAAAAAAAAAQAtB
GwAAAAAAAAAAAAAAAAAtBG0AAAAAAAAAAAAAAAAAtBC0AQAAAAAAAAAAAAAAANBC0AYAAAAAAAAA
AAAAAABAC0EbAAAAAAAAAAAAAAAAAC0EbQAAAAAAAAAAAAAAAAC0ELQBAAAAAAAAAAAAAAAA0ELQ
BgAAAAAAAAAAAAAAAEALQRsAAAAAAAAAAAAAAAAALQRtAAAAAAAAAAAAAAAAALQQtAEAAAAAAAAA
AAAAAADQQtAGAAAAAAAAAAAAAAAAQAtBGwAAAAAAAAAAAAAAAAAtBG0AAAAAAAAAAAAAAAAAtBC0
AQAAAAAAAAAAAAAAANBC0AYAAAAAAAAAAAAAAABAC0EbAAAAAAAAAAAAAAAAAC0EbQAAAAAAAAAA
AAAAAAC0ELQBAAAAAAAAAAAAAAAA0ELQBgAAAAAAAAAAAAAAAEALQRsAAAAAAAAAAAAAAAAALQRt
AAAAAAAAAAAAAAAAALQQtAEAAAAAAAAAAAAAAADQQtAGAAAAAAAAAAAAAAAAQAtBGwAAAAAAAAAA
AAAAAAAtBG0AAAAAAAAAAAAAAAAAtBC0AQAAAAAAAAAAAAAAANBC0AYAAAAAAAAAAAAAAABAC0Eb
AAAAAAAAAAAAAAAAAC0EbQAAAAAAAAAAAAAAAAC0ELQB/8e+HQsAAAAADPK3nsWu8ggAAAAAAAAA
AAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMA
AAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAA
AAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAA
AAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABY
CG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAA
AAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoA
AAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAA
AAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAA
AAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAA
FkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAA
AAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2
AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAA
AAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAA
AAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAA
gIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAA
AAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuh
DQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAA
AAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAA
AAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAA
AGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAA
AAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBC
aAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAA
AAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYA
AAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAA
AABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAA
AAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACw
ENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAA
AAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQB
AAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAA
AAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAA
AAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAA
LIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAA
AAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAht
AAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAA
AAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAA
AAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAA
AAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAA
AAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZC
GwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAA
AAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAA
AAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAA
AMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAA
AAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF
0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAA
AAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0A
AAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAA
AACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAA
AAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABg
IbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAA
AAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgD
AAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAA
AAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAA
AAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAA
WAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAA
AAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDa
AAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAA
AAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAA
AAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAA
ABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAA
AAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyE
NgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAA
AAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAA
AAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAA
AICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAA
AAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAAL
oQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAA
AAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsA
AAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAA
AABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAA
AAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADA
QmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAA
AAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAG
AAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAA
AAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAA
AAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAA
sBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAA
AAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0
AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAA
AAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAA
AAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAA
ACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAA
AAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgI
bQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAQOzbsQAAAADAIH/r
Wewqj1gIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAA
AAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAA
ALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAA
AAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAh
tAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAA
AAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMA
AAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAA
AAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAA
AAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABY
CG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAA
AAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoA
AAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAA
AAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAA
AAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAA
FkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAA
AAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2
AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAA
AAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAA
AAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAA
gIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAA
AAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuh
DQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAA
AAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAA
AAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAA
AGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAA
AAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBC
aAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAA
AAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYA
AAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAA
AABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAA
AAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACw
ENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAA
AAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQB
AAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAA
AAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAA
AAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAA
LIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAA
AAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAht
AAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAA
AAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAA
AAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAA
AAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAA
AAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZC
GwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAA
AAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAA
AAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAA
AMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAA
AAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF
0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAA
AAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0A
AAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAA
AACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAA
AAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABg
IbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAA
AAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgD
AAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAA
AAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAA
AAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAA
WAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAA
AAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDa
AAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAA
AAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAA
AAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAA
ABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAA
AAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyE
NgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAA
AAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAA
AAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAA
AICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAA
AAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAAL
oQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAA
AAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsA
AAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAA
AABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAA
AAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADA
QmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAA
AAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAG
AAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAA
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
AAAAAAAAAAAAkBC0AQAAAAAAAAAAAAAAAJAQtAEAAAAAAAAAAAAAAACQELQBAAAAAAAAAAAAAAAA
kBC0AQAAAAAAAAAAAAAAAJAQtAEAAAAAAAAAAAAAAACQELQBAAAAAAAAAAAAAAAAkBC0AQAAAAAA
AAAAAAAAAJAQtAEAAAAAAAAAAAAAAACQELQBAAAAAAAAAAAAAAAAkBC0AQAAAAAAAAAAAAAAAJAQ
tAEAAAAAAAAAAAAAAACQELQBAAAAAAAAAAAAAAAAkBC0AQAAAAAAAAAAAAAAAJAQtAEAAAAAAAAA
AAAAAACQELQBAAAAAAAAAAAAAAAAkBC0AQAAAAAAAAAAAAAAAJAQtAEAAAAAAAAAAAAAAACQELQB
AAAAAAAAAAAAAAAAkBC0AQAAAAAAAAAAAAAAAJAQtAEAAAAAAAAAAAAAAACQELQBAAAAAAAAAAAA
AAAAkBC0AQAAAAAAAAAAAAAAAJAQtAEAAAAAAAAAAAAAAACQELQBAAAAAAAAAAAAAAAAkBC0AQAA
AAAAAAAAAAAAAJAQtAEAAAAAAAAAAAAAAACQELQBAAAAAAAAAAAAAAAAkBC0AQAAAAAAAAAAAAAA
AJAQtAEAAAAAAAAAAAAAAACQELQBAAAAAAAAAAAAAAAAkBC0AQAAAAAAAAAAAAAAAJAQtAEAAAAA
AAAAAAAAAACQELQBAAAAAAAAAAAAAAAAkBC0AQAAAAAAAAAAAAAAAJAQtAEAAAAAAAAAAAAAAACQ
ELQBAAAAAAAAAAAAAAAAkBC0AQAAAAAAAAAAAAAAAJAQtAEAAAAAAAAAAAAAAACQELQBAAAAAAAA
AAAAAAAAkBC0AQAAAAAAAAAAAAAAAJAQtAEAAAAAAAAAAAAAAACQELQBAAAAAAAAAAAAAAAAkBC0
AQAAAAAAAAAAAAAAAJAQtAEAAAAAAAAAAAAAAACQELQBAAAAAAAAAAAAAAAAkBC0AQAAAAAAAAAA
AAAAAJAQtAEAAAAAAAAAAAAAAACQELQBAAAAAAAAAAAAAAAAkBC0AQAAAAAAAAAAAAAAAJAQtAEA
AAAAAAAAAAAAAACQELQBAAAAAAAAAAAAAAAAkBC0AQAAAAAAAAAAAAAAAJAQtAEAAAAAAAAAAAAA
AACQELQBAAAAAAAAAAAAAAAAkBC0AQAAAAAAAAAAAAAAAJAQtAEAAAAAAAAAAAAAAACQELQBAAAA
AAAAAAAAAAAAkBC0AQAAAAAAAAAAAAAAAJAQtAEAAAAAAAAAAAAAAACQELQBAAAAAAAAAAAAAAAA
kBC0AQAAAAAAAAAAAAAAAJAQtAEAAAAAAAAAAAAAAACQELQBAAAAAAAAAAAAAAAAkBC0AQAAAAAA
AAAAAAAAAJAQtAEAAAAAAAAAAAAAAACQELQBAAAAAAAAAAAAAAAAkBC0AQAAAAAAAAAAAAAAAJAQ
tAEAAAAAAAAAAAAAAACQELQBAAAAAAAAAAAAAAAAkBC0AQAAAAAAAAAAAAAAAJAQtAEAAAAAAAAA
AAAAAACQELQBAAAAAAAAAAAAAAAAkBC0AQAAAAAAAAAAAAAAAJAQtAEAAAAAAAAAAAAAAACQELQB
AAAAAAAAAAAAAAAAkBC0AQAAAAAAAAAAAAAAAJAQtAEAAAAAAAAAAAAAAACQELQBAAAAAAAAAAAA
AAAAkBC0AQAAAAAAAAAAAAAAAJAQtAEAAAAAAAAAAAAAAACQELQBAAAAAAAAAAAAAAAAkBC0AQAA
AAAAAAAAAAAAAJAQtAEAAAAAAAAAAAAAAACQELQBAAAAAAAAAAAAAAAAkBC0AQAAAAAAAAAAAAAA
AJAQtAEAAAAAAAAAAAAAAACQELQBAAAAAAAAAAAAAAAAkBC0AQAAAAAAAAAAAAAAAJAQtAEAAAAA
AAAAAAAAAACQELQBAAAAAAAAAAAAAAAAkBC0AQAAAAAAAAAAAAAAAJAQtAEAAAAAAAAAAAAAAACQ
ELQBAAAAAAAAAAAAAAAAkBC0AQAAAAAAAAAAAAAAAJAQtAEAAAAAAAAAAAAAAACQELQBAAAAAAAA
AAAAAAAAkBC0AQAAAAAAAAAAAAAAAJAQtAEAAAAAAAAAAAAAAACQELQBAAAAAAAAAAAAAAAAkBC0
AQAAAAAAAAAAAAAAAJAQtAEAAAAAAAAAAAAAAACQELQBAAAAAAAAAAAAAAAAkBC0AQAAAAAAAAAA
AAAAAJAQtAEAAAAAAAAAAAAAAACQELQBAAAAAAAAAAAAAAAAkBC0AQAAAAAAAAAAAAAAAJAQtAEA
AAAAAAAAAAAAAACQELQBAAAAAAAAAAAAAAAAkBC0AQAAAAAAAAAAAAAAAJAQtAEAAAAAAAAAAAAA
AACQELQBAAAAAAAAAAAAAAAAkBC0AQAAAAAAAAAAAAAAAJAQtAEAAAAAAAAAAAAAAACQELQBAAAA
AAAAAAAAAAAAkBC0AQAAAAAAAAAAAAAAAJAQtAEAAAAAAAAAAAAAAACQELQBAAAAAAAAAAAAAAAA
kBC0AQAAAAAAAAAAAAAAAJAQtAEAAAAAAAAAAAAAAACQELQBAAAAAAAAAAAAAAAAkBC0AQAAAAAA
AAAAAAAAAJAQtAEAAAAAAAAAAAAAAACQELQBAAAAAAAAAAAAAAAAkBC0AQAAAAAAAAAAAAAAAJAQ
tAEAAAAAAAAAAAAAAACQELQBAAAAAAAAAAAAAAAAkBC0AQAAAAAAAAAAAAAAAJAQtAEAAAAAAAAA
AAAAAACQELQBAAAAAAAAAAAAAAAAkBC0AQAAAAAAAAAAAAAAAJAQtAEAAAAAAAAAAAAAAACQELQB
AAAAAAAAAAAAAAAAkBC0AQAAAAAAAAAAAAAAAJAQtAEAAAAAAAAAAAAAAACQELQBAAAAAAAAAAAA
AAAAkBC0AQAAAAAAAAAAAAAAAJAQtAEAAAAAAAAAAAAAAACQELQBAAAAAAAAAAAAAAAAkBC0AQAA
AAAAAAAAAAAAAJAQtAEAAAAAAAAAAAAAAACQELQBAAAAAAAAAAAAAAAAkBC0AQAAAAAAAAAAAAAA
AJAQtAEAAAAAAAAAAAAAAACQELQBAAAAAAAAAAAAAAAAkBC0AQAAAAAAAAAAAAAAAJAQtAEAAAAA
AAAAAAAAAACQELQBAAAAAAAAAAAAAAAAkBC0AQAAAAAAAAAAAAAAAJAQtAEAAAAAAAAAAAAAAACQ
ELQBAAAAAAAAAAAAAAAAkBC0AQAAAAAAAAAAAAAAAJAQtAEAAAAAAAAAAAAAAACQELQBAAAAAAAA
AAAAAAAAkBC0AQAAAAAAAAAAAAAAAJAQtAEAAAAAAAAAAPxl344FAAAAAAb5W89iV3kEAAAAwEJo
AwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAA
AAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAA
AAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAA
AFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAA
AAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ
2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAA
AAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEA
AAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAA
AAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAA
AAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAs
hDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAA
AAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0A
AAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAA
AACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAA
AAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAA
C6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAA
AAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIb
AAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAA
AAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAA
AAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAA
wEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAA
AAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQ
BgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAA
AAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAA
AAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAA
ALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAA
AAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAh
tAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAA
AAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMA
AAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAA
AAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAA
AAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABY
CG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAA
AAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoA
AAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAA
AAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAA
AAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAA
FkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAA
AAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2
AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAA
AAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAA
AAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAA
gIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAA
AAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuh
DQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAA
AAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAA
AAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAA
AGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAA
AAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBC
aAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAA
AAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYA
AAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAA
AABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAA
AAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACw
ENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAA
AAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQB
AAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAA
AAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAA
AAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAA
LIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAA
AAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAht
AAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAA
AAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAA
AAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAA
AAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAA
AAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZC
GwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAA
AAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAA
AAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAA
AMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAA
AAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF
0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAA
AAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0A
AAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAA
AACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAA
AAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABg
IbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAA
AAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgD
AAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAA
AAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAA
AAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAA
WAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAA
AAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAQOzbsQAAAADA
IH/rWewqjwAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAA
AAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAA
AAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAA
YCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAA
AAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJo
AwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAA
AAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAA
AAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAA
AFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAA
AAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ
2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAA
AAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEA
AAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAA
AAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAA
AAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAs
hDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAA
AAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0A
AAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAA
AACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAA
AAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAA
C6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAA
AAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIb
AAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAA
AAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAA
AAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAA
wEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAA
AAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQ
BgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAA
AAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAA
AAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAA
ALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAA
AAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAh
tAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtACiBghwAACaFSURBVAAAAAAAAAAAAAAA
ACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAA
AAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgI
bQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAA
AAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAA
AAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAA
AAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAA
AAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAW
QhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAA
AAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYA
AAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAA
AADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAA
AAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACA
hdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAA
AAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6EN
AAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAA
AAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAA
AAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAA
YCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAA
AAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJo
AwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAA
AAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAA
AAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAA
AFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAA
AAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ
2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAA
AAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEA
AAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAA
AAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAA
AAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAs
hDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAA
AAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0A
AAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAA
AACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAA
AAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAA
C6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAA
AAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIb
AAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAA
AAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAA
AAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAA
wEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAA
AAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQ
BgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAA
AAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAA
AAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAA
ALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAA
AAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAh
tAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAA
AAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMA
AAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAA
AAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAA
AAAAAAAAAADAQmgDAAAAAACIfTsWAAAAABjkbz2LXeURAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF
0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAA
AAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0A
AAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAA
AACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAA
AAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABg
IbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAA
AAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgD
AAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAA
AAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAA
AAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAA
WAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAA
AAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDa
AAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAA
AAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAA
AAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAA
ABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAA
AAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyE
NgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAA
AAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAA
AAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAA
AICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAA
AAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAAL
oQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAA
AAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsA
AAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAA
AABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAA
AAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADA
QmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAA
AAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAG
AAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAA
AAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAA
AAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAA
sBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAA
AAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0
AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAA
AAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAA
AAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAA
ACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAA
AAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgI
bQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAA
AAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAA
AAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAA
AAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAA
AAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAW
QhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAA
AAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYA
AAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAA
AADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAA
AAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACA
hdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAA
AAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6EN
AAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAA
AAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAA
AAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAA
YCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAA
AAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJo
AwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAA
AAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAA
AAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAA
AFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAA
AAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ
2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAA
AAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEA
AAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAA
AAAWQhsAAAAAAAAAAAAAAAAAC6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAA
AAAAAAAAAABgIbQBAAAAAAAAAAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAs
hDYAAAAAAAAAAAAAAAAAFkIbAAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAA
AAAAAADAQmgDAAAAAAAAAAAAAAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0A
AAAAAAAAAAAAAAAALIQ2AAAAAAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAA
AACAhdAGAAAAAAAAAAAAAAAAwEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAA
AAAAAAAAAAAAWAhtAAAAAAAAAAAAAAAAACyENgAAAAAAAAAAAAAAAAAWQhsAAAAAAAAAAAAAAAAA
C6ENAAAAAAAAAAAAAAAAgIXQBgAAAAAAAAAAAAAAAMBCaAMAAAAAAAAAAAAAAABgIbQBAAAAAAAA
AAAAAAAAsBDaAAAAAAAAAAAAAAAAAFgIbQAAAAAAAAAAAAAAAAAshDYAAAAAAAAAAAAAAAAAFkIb
AAAAAAAAAAAAAAAAAAuhDQAAAAAAAAAAAAAAAICF0AYAAAAAAAAAAAAAAADAQmgDAAAAAAAAAAAA
AAAAYCG0AQAAAAAAAAAAAAAAALAQ2gAAAAAAAAAAAAAAAABYCG0AAAAAAAAAAAAAAAAALIQ2AAAA
AAAAAAAAAAAAABZCGwAAAAAAAAAAAAAAAAALoQ0AAAAAAAAAAAAAAACAhdAGAAAAAAAAAAAAAAAA
wEJoAwAAAAAAAAAAAAAAAGAhtAEAAAAAAAAAAAAAAACwENoAAAAAAAAAAAAAAAAAWAhtAAAAAAAA
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
CNoAAAAAAAAAAAAAAAAASAjaAAAAAAAAAAAAAAAAAEgI2gAAAAAAAAAAAAAAAABICNoAAAAAAAAA
AAAAAAAASAjaAAAAAAAAAAAAAAAAAEgI2gAAAAAAAAAAAAAAAABICNoAAAAAAAAAAAAAAAAASAja
AAAAAAAAAAAAAAAAAEgI2gAAAAAAAAAAAAAAAABICNoAAAAAAAAAAAAAAAAASAjaAAAAAAAAAAAA
AAAAAEgI2gAAAAAAAAAAAAAAAABICNoAAAAAAAAAAAAAAAAASAjaAAAAAAAAAAAAAAAAAEgI2gAA
AAAAAAAAAAAAAABICNoAAAAAAAAAAAAAAAAASAjaAAAAAAAAAAAAAAAAAEgI2gAAAAAAAAAAAAAA
AABICNoAAAAAAAAAAAAAAAAASAjaAAAAAAAAAAAAAAAAAEgI2gAAAAAAAAAAAAAAAABICNoAAAAA
AAAAAAAAAAAASAjaAAAAAAAAAAAAAAAAAEgI2gAAAAAAAAAAAAAAAABICNoAAAAAAAAAAAAAAAAA
SAjaAAAAAAAAAAAAAAAAAEgI2gAAAAAAAAAAAAAAAABICNoAAAAAAAAAAAAAAAAASAjaAAAAAAAA
AAAAAAAAAEgI2gAAAAAAAAAAAAAAAABICNoAAAAAAAAAAAAAAAAASAjaAAAAAAAAAAAAAAAAAEgI
2gAAAAAAAAAAAAAAAABICNoAAAAAAAAAAAAAAAAASAjaAAAAAAAAAAAAAAAAAEgI2gAAAAAAAAAA
AAAAAABICNoAAAAAAAAAAAAAAAAASAjaAAAAAAAAAAAAAAAAAEgI2gAAAAAAAAAAAAAAAABICNoA
AAAAAAAAAAAAAAAASAjaAAAAAAAAAAAAAAAAAEgI2gAAAAAAAAAAAAAAAABICNoAAAAAAAAAAAAA
AAAASAjaAAAAAAAAAAAAAAAAAEgI2gAAAAAAAAAAAAAAAABIfAMAAP//otUQpmPy4DcAAAAASUVO
RK5CYII=`
