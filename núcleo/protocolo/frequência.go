package protocolo

import "time"

// FrequênciaPedido armazena os dados exigidos pelo Exército ao utilizar um
// estande de Tiro.
type FrequênciaPedido struct {
	Calibre           string    `json:"calibre"`
	ArmaUtilizada     string    `json:"armaUtilizada"`
	NúmeroSérie       string    `json:"numeroSerie"`
	GuiaDeTráfego     string    `json:"guiaTrafego"`
	QuantidadeMunição int       `json:"quantidadeMunicao"`
	HorárioInício     time.Time `json:"horarioInicio"`
	HorárioTérmino    time.Time `json:"horarioTermino"`
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
	NúmeroControle int64  `json:"numeroControle"`
	Imagem         string `json:"imagem"` // base64
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
	NúmeroControle int64
	FrequênciaConfirmaçãoPedido
}

// NovaFrequênciaConfirmaçãoPedidoCompleta inicializa o tipo
// FrequênciaConfirmaçãoPedidoCompleta a partir do CR, número de controle e do
// tipo FrequênciaConfirmaçãoPedido.
func NovaFrequênciaConfirmaçãoPedidoCompleta(cr string, númeroControle int64, frequênciaConfirmaçãoPedido FrequênciaConfirmaçãoPedido) FrequênciaConfirmaçãoPedidoCompleta {
	return FrequênciaConfirmaçãoPedidoCompleta{
		CR:                          cr,
		NúmeroControle:              númeroControle,
		FrequênciaConfirmaçãoPedido: frequênciaConfirmaçãoPedido,
	}
}
