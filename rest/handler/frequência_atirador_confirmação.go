package handler

import (
	"net/http"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/atirador"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/erros"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
	"github.com/rafaeljusto/atiradorfrequente/rest/config"
	"github.com/rafaeljusto/atiradorfrequente/rest/interceptador"
	"github.com/trajber/handy"
)

func init() {
	registrar("/frequencia/{cr}/{numeroControle}", func() handy.Handler { return &frequênciaAtiradorConfirmação{} })
}

type frequênciaAtiradorConfirmação struct {
	básico
	interceptador.BDCompatível

	CR                          string                                `urivar:"cr"`
	NúmeroControle              protocolo.NúmeroControle              `urivar:"numeroControle"`
	FrequênciaConfirmaçãoPedido protocolo.FrequênciaConfirmaçãoPedido `request:"put"`
}

func (f *frequênciaAtiradorConfirmação) Put() int {
	if config.Atual() == nil {
		f.Logger().Crit("Não existe configuração definida para atender a requisição")
		return http.StatusInternalServerError
	}

	serviçoAtirador := atirador.NovoServiço(f.Tx(), config.Atual().Configuração)
	frequênciaConfirmaçãoPedidoCompleta := protocolo.NovaFrequênciaConfirmaçãoPedidoCompleta(f.CR, f.NúmeroControle, f.FrequênciaConfirmaçãoPedido)

	if err := serviçoAtirador.ConfirmarFrequência(frequênciaConfirmaçãoPedidoCompleta); err != nil {
		// TODO(rafaeljusto): verificar se o erro é do tipo protocolo.Mensagens, neste
		// caso retornar as mensagens para o usuário com o código HTTP BadRequest
		// (400)

		f.Logger().Error(erros.Novo(err))
		return http.StatusInternalServerError
	}

	return http.StatusNoContent
}

func (f *frequênciaAtiradorConfirmação) Interceptors() handy.InterceptorChain {
	return criarCorrenteBásica(f).
		Chain(interceptador.NovoBD(f))
}
