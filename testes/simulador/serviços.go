package simulador

import "github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"

// ServiçoAtirador simula o serviço que representa um Atirador. Muito útil para
// simular as camadas de serviços em testes unitários.
type ServiçoAtirador struct {
	SimulaCadastrarFrequência func(protocolo.FrequênciaPedidoCompleta) (protocolo.FrequênciaPendenteResposta, error)
	SimulaObterFrequência     func(cr int, númeroControle protocolo.NúmeroControle, códigoVerificação string) (protocolo.FrequênciaResposta, error)
	SimulaConfirmarFrequência func(protocolo.FrequênciaConfirmaçãoPedidoCompleta) error
}

// CadastrarFrequência persiste em banco de dados as informações básicas
// relacionados a visita do Atirador a um Clube de Tiro. Esta ação será
// responsável por gerar o número de controle utilizado na confirmação da
// frequência.
func (s ServiçoAtirador) CadastrarFrequência(frequênciaPedidoCompleta protocolo.FrequênciaPedidoCompleta) (protocolo.FrequênciaPendenteResposta, error) {
	return s.SimulaCadastrarFrequência(frequênciaPedidoCompleta)
}

// ObterFrequência retorna a frequência relacionada ao CR e número de controle
// informados. O código de verificação deve bater com o informado no momento da
// criação para que a informação seja liberada.
func (s ServiçoAtirador) ObterFrequência(cr int, númeroControle protocolo.NúmeroControle, códigoVerificação string) (protocolo.FrequênciaResposta, error) {
	return s.SimulaObterFrequência(cr, númeroControle, códigoVerificação)
}

// ConfirmarFrequência finaliza o cadastro da frequência, confirmando atraves de
// uma imagem que o Atirador esta presente no Clube de Tiro.
func (s ServiçoAtirador) ConfirmarFrequência(frequênciaConfirmaçãoPedidoCompleta protocolo.FrequênciaConfirmaçãoPedidoCompleta) error {
	return s.SimulaConfirmarFrequência(frequênciaConfirmaçãoPedidoCompleta)
}
