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

	CR                         int                                   `urivar:"cr"`
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

	if err != nil {
		if mensagens, ok := err.(protocolo.Mensagens); ok {
			f.Mensagens = mensagens
			return http.StatusBadRequest
		}

		f.Logger().Error(erros.Novo(err))
		return http.StatusInternalServerError
	}

	f.FrequênciaPendenteResposta = &frequênciaPendenteResposta
	f.DefinirCabeçalho("Location", fmt.Sprintf("/frequencia-atirador/%d/%s", f.CR, f.FrequênciaPendenteResposta.NúmeroControle))
	return http.StatusCreated
}

func (f *frequênciaAtirador) Interceptors() handy.InterceptorChain {
	return criarCorrenteBásica(f).
		Chain(interceptador.NovoBD(f))
}
