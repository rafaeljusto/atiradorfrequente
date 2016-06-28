package simulador

import "github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"

// ServiçoAtirador simula o serviço que representa um Atirador. Muito útil para
// simular as camadas de serviços em testes unitários.
type ServiçoAtirador struct {
	SimulaCadastrarFrequência func(protocolo.FrequênciaPedidoCompleta) (protocolo.FrequênciaPendenteResposta, error)
	SimulaConfirmarFrequência func(protocolo.FrequênciaConfirmaçãoPedidoCompleta) error
}

// CadastrarFrequência persiste em banco de dados as informações básicas
// relacionados a visita do Atirador a um Clube de Tiro. Esta ação será
// responsável por gerar o número de controle utilizado na confirmação da
// frequência.
func (s ServiçoAtirador) CadastrarFrequência(frequênciaPedidoCompleta protocolo.FrequênciaPedidoCompleta) (protocolo.FrequênciaPendenteResposta, error) {
	return s.SimulaCadastrarFrequência(frequênciaPedidoCompleta)
}

// ConfirmarFrequência finaliza o cadastro da frequência, confirmando atraves de
// uma imagem que o Atirador esta presente no Clube de Tiro.
func (s ServiçoAtirador) ConfirmarFrequência(frequênciaConfirmaçãoPedidoCompleta protocolo.FrequênciaConfirmaçãoPedidoCompleta) error {
	return s.SimulaConfirmarFrequência(frequênciaConfirmaçãoPedidoCompleta)
}
