// +build gofuzz

package atirador

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/erikstmartin/go-testdb"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/bd"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/config"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
	"golang.org/x/image/font/basicfont"
)

func init() {
	gob.Register(protocolo.FrequênciaPedidoCompleta{})

	imagemLogoExtraída, err := base64.StdEncoding.DecodeString(imagemLogoPNG)

	if err != nil {
		panic(fmt.Errorf("Erro ao extrair a imagem de teste do logo. Detalhes: %s", err))
	}

	imagemLogoBuffer := bytes.NewBuffer(imagemLogoExtraída)
	if imagemLogo, _, err = image.Decode(imagemLogoBuffer); err != nil {
		panic(fmt.Errorf("Erro ao interpretar imagem. Detalhes: %s", err))
	}
}

// FuzzCadastrarFrequência é utilizado pela ferramenta go-fuzz, responsável por
// testar o cadastro de frequência com dados aleatórios.
func FuzzCadastrarFrequência(dados []byte) int {
	dadosExtraídos, err := base64.StdEncoding.DecodeString(string(dados))
	if err != nil {
		return -1
	}
	decodificador := gob.NewDecoder(bytes.NewReader(dadosExtraídos))

	var frequênciaPedidoCompleta protocolo.FrequênciaPedidoCompleta
	if decodificador.Decode(&frequênciaPedidoCompleta) != nil {
		return -1
	}

	conexão, err := sql.Open("testdb", "")
	if err != nil {
		panic(err)
	}

	testdb.StubQuery(frequênciaCriaçãoComando, testdb.RowsFromSlice([]string{"id"}, [][]driver.Value{{1}}))
	testdb.StubExec(frequênciaAtualizaçãoComando, testdb.NewResult(1, nil, 1, nil))
	testdb.StubExec(frequênciaLogCriaçãoComando, testdb.NewResult(1, nil, 1, nil))

	logCriaçãoComando := `INSERT INTO log (id, data_criacao, endereco_remoto) VALUES (DEFAULT, $1, $2) RETURNING id`
	testdb.StubQuery(logCriaçãoComando, testdb.RowsFromSlice([]string{"id"}, [][]driver.Value{{1}}))

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

	serviço := NovoServiço(bd.NovoSQLogger(conexão, nil), configuração)
	if _, err := serviço.CadastrarFrequência(frequênciaPedidoCompleta); err != nil {
		if _, ok := err.(protocolo.Mensagens); !ok {
			panic(err)
		}

		return 0
	}

	return 1
}

// FuzzConfirmarFrequência é utilizado pela ferramenta go-fuzz, responsável por
// testar a confirmar a frequência com dados aleatórios.
func FuzzConfirmarFrequência(dados []byte) int {
	dadosExtraídos, err := base64.StdEncoding.DecodeString(string(dados))
	if err != nil {
		return -1
	}
	decodificador := gob.NewDecoder(bytes.NewReader(dadosExtraídos))

	var frequênciaConfirmaçãoPedidoCompleta protocolo.FrequênciaConfirmaçãoPedidoCompleta
	if decodificador.Decode(&frequênciaConfirmaçãoPedidoCompleta) != nil {
		return -1
	}

	conexão, err := sql.Open("testdb", "")
	if err != nil {
		panic(err)
	}

	testdb.StubQuery(frequênciaResgateComando, testdb.RowsFromSlice(frequênciaResgateCampos, [][]driver.Value{
		{
			frequênciaConfirmaçãoPedidoCompleta.NúmeroControle.ID(),
			frequênciaConfirmaçãoPedidoCompleta.NúmeroControle.Controle(),
			frequênciaConfirmaçãoPedidoCompleta.CR, ".380", "Arma Clube", "ZA785671", 762556223, 50,
			time.Now().Add(-1 * time.Hour), time.Now().Add(-10 * time.Minute), time.Now(), time.Time{}, time.Time{},
			`TWFuIGlzIGRpc3Rpbmd1aXNoZWQsIG5vdCBvbmx5IGJ5IGhpcyByZWFzb24sIGJ1dCBieSB0aGlz
IHNpbmd1bGFyIHBhc3Npb24gZnJvbSBvdGhlciBhbmltYWxzLCB3aGljaCBpcyBhIGx1c3Qgb2Yg
dGhlIG1pbmQsIHRoYXQgYnkgYSBwZXJzZXZlcmFuY2Ugb2YgZGVsaWdodCBpbiB0aGUgY29udGlu
dWVkIGFuZCBpbmRlZmF0aWdhYmxlIGdlbmVyYXRpb24gb2Yga25vd2xlZGdlLCBleGNlZWRzIHRo
ZSBzaG9ydCB2ZWhlbWVuY2Ugb2YgYW55IGNhcm5hbCBwbGVhc3VyZS4=`, "", 0,
		},
	}))

	testdb.StubExec(frequênciaAtualizaçãoComando, testdb.NewResult(1, nil, 1, nil))
	testdb.StubExec(frequênciaLogCriaçãoComando, testdb.NewResult(1, nil, 1, nil))

	logCriaçãoComando := `INSERT INTO log (id, data_criacao, endereco_remoto) VALUES (DEFAULT, $1, $2) RETURNING id`
	testdb.StubQuery(logCriaçãoComando, testdb.RowsFromSlice([]string{"id"}, [][]driver.Value{{1}}))

	var configuração config.Configuração
	configuração.Atirador.PrazoConfirmação = 20 * time.Minute

	serviço := NovoServiço(bd.NovoSQLogger(conexão, nil), configuração)
	if err := serviço.ConfirmarFrequência(frequênciaConfirmaçãoPedidoCompleta); err != nil {
		if _, ok := err.(protocolo.Mensagens); !ok {
			panic(err)
		}

		return 0
	}

	return 1
}

var imagemLogo image.Image

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
