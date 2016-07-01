package protocolo

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// FrequênciaPedido armazena os dados exigidos pelo Exército ao utilizar um
// estande de Tiro.
type FrequênciaPedido struct {
	Calibre           string `json:"calibre"`
	ArmaUtilizada     string `json:"armaUtilizada"`
	NúmeroSérie       string `json:"numeroSerie"`
	GuiaDeTráfego     string `json:"guiaTrafego"`
	QuantidadeMunição int    `json:"quantidadeMunicao"`

	// DataInício data e hora do início do treino de tiro no estande do clube.
	DataInício time.Time `json:"dataInicio"`

	// DataTérmino data e hora do término do treino de tiro no estande do clube.
	DataTérmino time.Time `json:"dataTermino"`
}

// FrequênciaPedidoCompleta é uma extensão do tipo FrequênciaPedido incluindo o
// CR enviado no endereço.
type FrequênciaPedidoCompleta struct {
	CR string
	FrequênciaPedido
}

// NovaFrequênciaPedidoCompleta inicializa o tipo FrequênciaPedidoCompleta a
// partir do CR e do tipo FrequênciaPedido.
func NovaFrequênciaPedidoCompleta(cr string, frequênciaPedido FrequênciaPedido) FrequênciaPedidoCompleta {
	return FrequênciaPedidoCompleta{
		CR:               cr,
		FrequênciaPedido: frequênciaPedido,
	}
}

// FrequênciaPendenteResposta armazena os dados que permitem ao Clube de Tiro
// confirmar a presença do Atirador.
type FrequênciaPendenteResposta struct {
	NúmeroControle NúmeroControle `json:"numeroControle"`
	Imagem         string         `json:"imagem"` // base64
}

// FrequênciaConfirmaçãoPedido armazena os dados necessários para confirmar a
// presença do Atirador no Clube de Tiro.
type FrequênciaConfirmaçãoPedido struct {
	Imagem string `json:"imagem"` // base64
}

// FrequênciaConfirmaçãoPedidoCompleta extende o tipo
// FrequênciaConfirmaçãoPedido incluindo o CR e o número de controle encontrados
// no endereço.
type FrequênciaConfirmaçãoPedidoCompleta struct {
	CR             string
	NúmeroControle NúmeroControle
	FrequênciaConfirmaçãoPedido
}

// NovaFrequênciaConfirmaçãoPedidoCompleta inicializa o tipo
// FrequênciaConfirmaçãoPedidoCompleta a partir do CR, número de controle e do
// tipo FrequênciaConfirmaçãoPedido.
func NovaFrequênciaConfirmaçãoPedidoCompleta(cr string, númeroControle NúmeroControle, frequênciaConfirmaçãoPedido FrequênciaConfirmaçãoPedido) FrequênciaConfirmaçãoPedidoCompleta {
	return FrequênciaConfirmaçãoPedidoCompleta{
		CR:                          cr,
		NúmeroControle:              númeroControle,
		FrequênciaConfirmaçãoPedido: frequênciaConfirmaçãoPedido,
	}
}

// NúmeroControle número gerado a partir do cadastro de uma frequência para
// comprovação da presença física do atirador no estande de tiro. Formado a
// partir do número de identificação da frequência com um número gerado
// aleatoriamente.
type NúmeroControle string

// NovoNúmeroControle gera um novo número de controle a partir do número de
// identificação da frequência com um número gerado aleatoriamente, chamado de
// controle.
func NovoNúmeroControle(id int64, controle int64) NúmeroControle {
	return NúmeroControle(fmt.Sprintf("%d-%d", id, controle))
}

// String retorna o formato texto do número de controle.
func (n NúmeroControle) String() string {
	return string(n)
}

// ID retorna o número de identificação da frequência contido no número de
// controle.
func (n NúmeroControle) ID() int64 {
	valor := string(n)
	partes := strings.Split(valor, "-")

	id, err := strconv.ParseInt(partes[0], 10, 64)
	if err != nil {
		return 0
	}

	return id
}

// Controle retorna o número aleatório contido no número de controle.
func (n NúmeroControle) Controle() int64 {
	valor := string(n)
	partes := strings.Split(valor, "-")

	if len(partes) != 2 {
		return 0
	}

	controle, err := strconv.ParseInt(partes[1], 10, 64)
	if err != nil {
		return 0
	}

	return controle
}
