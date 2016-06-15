package atirador

import "github.com/rafaeljusto/cr/núcleo/protocolo"

type Serviço interface {
	CadastrarFrequência(protocolo.FrequênciaPedidoCompleta) (protocolo.FrequênciaResposta, error)
	ConfirmarFrequência(protocolo.FrequênciaConfirmaçãoPedidoCompleta) error
}

var NovoServiço = func() Serviço {
	return serviço{}
}

type serviço struct {
}

func (s serviço) CadastrarFrequência(frequênciaPedidoCompleta protocolo.FrequênciaPedidoCompleta) (protocolo.FrequênciaResposta, error) {
	return protocolo.FrequênciaResposta{}, nil
}

func (s serviço) ConfirmarFrequência(frequênciaConfirmaçãoPedidoCompleta protocolo.FrequênciaConfirmaçãoPedidoCompleta) error {
	return nil
}
