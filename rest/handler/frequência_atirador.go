package handler

import (
	"fmt"
	"net/http"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/atirador"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/erros"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
	"github.com/rafaeljusto/atiradorfrequente/rest/config"
	"github.com/rafaeljusto/atiradorfrequente/rest/interceptador"
	"github.com/trajber/handy"
)

func init() {
	registrar("/frequencia/{cr}", func() handy.Handler { return &frequênciaAtirador{} })
}

type frequênciaAtirador struct {
	básico
	interceptador.BDCompatível

	// TODO(rafaeljusto): Criar um tipo para o CR para padronizar a entrada
	CR                         string                                `urivar:"cr"`
	FrequênciaPedido           protocolo.FrequênciaPedido            `request:"post"`
	FrequênciaPendenteResposta *protocolo.FrequênciaPendenteResposta `response:"post"`
}

func (f *frequênciaAtirador) Post() int {
	if config.Atual() == nil {
		f.Logger().Crit("Não existe configuração definida para atender a requisição")
		return http.StatusInternalServerError
	}

	serviçoAtirador := atirador.NovoServiço(f.Tx(), config.Atual().Configuração)
	frequênciaPedidoCompleta := protocolo.NovaFrequênciaPedidoCompleta(f.CR, f.FrequênciaPedido)
	frequênciaPendenteResposta, err := serviçoAtirador.CadastrarFrequência(frequênciaPedidoCompleta)

	// TODO(rafaeljusto): verificar se o erro é do tipo protocolo.Mensagens, neste
	// caso retornar as mensagens para o usuário com o código HTTP BadRequest
	// (400)

	if err != nil {
		f.Logger().Error(erros.Novo(err))
		return http.StatusInternalServerError
	}

	f.FrequênciaPendenteResposta = &frequênciaPendenteResposta
	f.DefinirCabeçalho("Location", fmt.Sprintf("/frequencia-atirador/%s/%d", f.CR, f.FrequênciaPendenteResposta.NúmeroControle))
	return http.StatusCreated
}

func (f *frequênciaAtirador) Interceptors() handy.InterceptorChain {
	return criarCorrenteBásica(f).
		Chain(interceptador.NovoBD(f))
}
