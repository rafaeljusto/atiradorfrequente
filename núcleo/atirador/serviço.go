package atirador

import (
	"github.com/rafaeljusto/atiradorfrequente/núcleo/bd"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/erros"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
)

// Serviço disponibiliza as ações que podem ser feitas relacionadas ao Atirador.
type Serviço interface {
	// CadastrarFrequência persiste em banco de dados as informações básicas
	// relacionados a visita do Atirador a um Clube de Tiro. Esta ação será
	// responsável por gerar o número de controle utilizado na confirmação da
	// frequência.
	CadastrarFrequência(protocolo.FrequênciaPedidoCompleta) (protocolo.FrequênciaPendenteResposta, error)

	// ConfirmarFrequência finaliza o cadastro da frequência, confirmando atraves
	// de uma imagem que o Atirador esta presente no Clube de Tiro.
	ConfirmarFrequência(protocolo.FrequênciaConfirmaçãoPedidoCompleta) error
}

// NovoServiço inicializa um serviço concreto do Atirador. Pode ser substituído
// em testes por simuladores, permitindo uma abstração da camada de serviços.
var NovoServiço = func(s *bd.SQLogger) Serviço {
	return serviço{
		sqlogger: s,
	}
}

type serviço struct {
	sqlogger *bd.SQLogger
}

func (s serviço) CadastrarFrequência(frequênciaPedidoCompleta protocolo.FrequênciaPedidoCompleta) (protocolo.FrequênciaPendenteResposta, error) {
	f := novaFrequência(frequênciaPedidoCompleta)

	dao := novaFrequênciaDAO(s.sqlogger)
	if err := dao.criar(&f); err != nil {
		return protocolo.FrequênciaPendenteResposta{}, erros.Novo(err)
	}

	// TODO(rafaeljusto): Criar imagem com o número de controle
	return protocolo.FrequênciaPendenteResposta{}, nil
}

func (s serviço) ConfirmarFrequência(frequênciaConfirmaçãoPedidoCompleta protocolo.FrequênciaConfirmaçãoPedidoCompleta) error {
	// TODO(rafaeljusto): Persistir imagem de confirmação
	return nil
}
