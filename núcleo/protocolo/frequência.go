package protocolo

import "time"

type FrequênciaPedido struct {
	Calibre           string    `json:"calibre"`
	ArmaUtilizada     string    `json:"armaUtilizada"`
	QuantidadeMunição int       `json:"quantidadeMunicao"`
	HorárioInício     time.Time `json:"horarioInicio"`
	HorárioTérmino    time.Time `json:"horarioTermino"`
}

type FrequênciaPedidoCompleta struct {
	CR string
	FrequênciaPedido
}

func NovaFrequênciaPedidoCompleta(cr string, frequênciaPedido FrequênciaPedido) FrequênciaPedidoCompleta {
	return FrequênciaPedidoCompleta{
		CR:               cr,
		FrequênciaPedido: frequênciaPedido,
	}
}

type FrequênciaResposta struct {
	NúmeroControle int    `json:"numeroControle"`
	Imagem         string `json:"imagem"` // base64
}

type FrequênciaConfirmaçãoPedido struct {
	Imagem string `json:"imagem"` // base64
}

type FrequênciaConfirmaçãoPedidoCompleta struct {
	CR             string
	NúmeroControle int
	FrequênciaConfirmaçãoPedido
}

func NovaFrequênciaConfirmaçãoPedidoCompleta(cr string, númeroControle int, frequênciaConfirmaçãoPedido FrequênciaConfirmaçãoPedido) FrequênciaConfirmaçãoPedidoCompleta {
	return FrequênciaConfirmaçãoPedidoCompleta{
		CR:                          cr,
		NúmeroControle:              númeroControle,
		FrequênciaConfirmaçãoPedido: frequênciaConfirmaçãoPedido,
	}
}
