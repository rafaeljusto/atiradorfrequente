package atirador

import "github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"

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
	// TODO(rafaeljusto): Persistir na base de dados
	// TODO(rafaeljusto): Criar imagem com o número de controle
	return protocolo.FrequênciaResposta{}, nil
}

func (s serviço) ConfirmarFrequência(frequênciaConfirmaçãoPedidoCompleta protocolo.FrequênciaConfirmaçãoPedidoCompleta) error {
	// TODO(rafaeljusto): Persistir imagem de confirmação
	return nil
}
