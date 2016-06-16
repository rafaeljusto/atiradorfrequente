package atirador

import (
	"time"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
)

type frequência struct {
	ID                   uint64
	Controle             int64
	CR                   string
	Calibre              string
	ArmaUtilizada        string
	QuantidadeMunição    int
	HorárioInício        time.Time
	HorárioTérmino       time.Time
	ImagemNúmeroControle string
	ImagemConfirmação    string
	DataCriação          time.Time
	DataConfirmação      time.Time
}

func novaFrequência(frequênciaPedidoCompleta protocolo.FrequênciaPedidoCompleta) frequência {
	return frequência{
		Controle:          origemRandômica.Int63(),
		CR:                frequênciaPedidoCompleta.CR,
		Calibre:           frequênciaPedidoCompleta.Calibre,
		ArmaUtilizada:     frequênciaPedidoCompleta.ArmaUtilizada,
		QuantidadeMunição: frequênciaPedidoCompleta.QuantidadeMunição,
		HorárioInício:     frequênciaPedidoCompleta.HorárioInício,
		HorárioTérmino:    frequênciaPedidoCompleta.HorárioTérmino,
		DataCriação:       time.Now().UTC(),
	}
}
