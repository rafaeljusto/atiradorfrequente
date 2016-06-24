package simulador

import "github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"

type ServiçoAtirador struct {
	SimulaCadastrarFrequência func(protocolo.FrequênciaPedidoCompleta) (protocolo.FrequênciaPendenteResposta, error)
	SimulaConfirmarFrequência func(protocolo.FrequênciaConfirmaçãoPedidoCompleta) error
}

func (s ServiçoAtirador) CadastrarFrequência(frequênciaPedidoCompleta protocolo.FrequênciaPedidoCompleta) (protocolo.FrequênciaPendenteResposta, error) {
	return s.SimulaCadastrarFrequência(frequênciaPedidoCompleta)
}

func (s ServiçoAtirador) ConfirmarFrequência(frequênciaConfirmaçãoPedidoCompleta protocolo.FrequênciaConfirmaçãoPedidoCompleta) error {
	return s.SimulaConfirmarFrequência(frequênciaConfirmaçãoPedidoCompleta)
}
