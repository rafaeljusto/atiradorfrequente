package atirador

import (
	"time"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/randômico"
)

type frequência struct {
	ID                   int64
	Controle             int64
	CR                   string
	Calibre              string
	ArmaUtilizada        string
	NúmeroSérie          string
	GuiaDeTráfego        string
	QuantidadeMunição    int
	DataInício           time.Time
	DataTérmino          time.Time
	DataCriação          time.Time
	DataAtualização      time.Time
	DataConfirmação      time.Time
	ImagemNúmeroControle string
	ImagemConfirmação    string

	// revisão utilizado para o controle de versão do objeto na base de dados,
	// minimizando problemas de concorrência quando 2 transações alteram o mesmo
	// objeto.
	revisão int
}

func (f *frequência) confirmar(frequênciaConfirmaçãoPedidoCompleta protocolo.FrequênciaConfirmaçãoPedidoCompleta) {
	f.DataConfirmação = time.Now().UTC()
	f.ImagemConfirmação = frequênciaConfirmaçãoPedidoCompleta.Imagem
}

func novaFrequência(frequênciaPedidoCompleta protocolo.FrequênciaPedidoCompleta) frequência {
	return frequência{
		Controle:          randômico.FonteRandômica.Int63(),
		CR:                frequênciaPedidoCompleta.CR,
		Calibre:           frequênciaPedidoCompleta.Calibre,
		ArmaUtilizada:     frequênciaPedidoCompleta.ArmaUtilizada,
		NúmeroSérie:       frequênciaPedidoCompleta.NúmeroSérie,
		GuiaDeTráfego:     frequênciaPedidoCompleta.GuiaDeTráfego,
		QuantidadeMunição: frequênciaPedidoCompleta.QuantidadeMunição,
		DataInício:        frequênciaPedidoCompleta.DataInício,
		DataTérmino:       frequênciaPedidoCompleta.DataTérmino,
	}
}
