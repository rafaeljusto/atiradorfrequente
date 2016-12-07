package handler

import (
	"net/http"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/atirador"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/erros"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
	"github.com/rafaeljusto/atiradorfrequente/rest/config"
	"github.com/rafaeljusto/atiradorfrequente/rest/interceptador"
	"github.com/registrobr/gostk/errors"
	"github.com/trajber/handy"
)

func init() {
	registrar("/frequencia/{cr}/{numeroControle}", func() handy.Handler { return &frequênciaAtiradorConfirmação{} })
}

type frequênciaAtiradorConfirmação struct {
	básico
	interceptador.BDCompatível

	CR                          int                                   `urivar:"cr"`
	NúmeroControle              protocolo.NúmeroControle              `urivar:"numeroControle"`
	CódigoVerificação           string                                `query:"verificacao"`
	FrequênciaConfirmaçãoPedido protocolo.FrequênciaConfirmaçãoPedido `request:"put"`
	FrequênciaResposta          *protocolo.FrequênciaResposta         `response:"get"`
}

func (f *frequênciaAtiradorConfirmação) Get() int {
	if config.Atual() == nil {
		f.Logger().Crit("Não existe configuração definida para atender a requisição")
		return http.StatusInternalServerError
	}

	serviçoAtirador := atirador.NovoServiço(f.Tx(), f.Logger(), config.Atual().Configuração)
	frequênciaResposta, err := serviçoAtirador.ObterFrequência(f.CR, f.NúmeroControle, f.CódigoVerificação)
	if err != nil {
		if errors.Equal(err, erros.NãoEncontrado) {
			return http.StatusNotFound
		}

		if mensagens, ok := err.(protocolo.Mensagens); ok {
			f.Mensagens = mensagens
			return http.StatusBadRequest
		}

		f.Logger().Error(erros.Novo(err))
		return http.StatusInternalServerError
	}

	f.FrequênciaResposta = &frequênciaResposta
	return http.StatusOK
}

func (f *frequênciaAtiradorConfirmação) Put() int {
	if config.Atual() == nil {
		f.Logger().Crit("Não existe configuração definida para atender a requisição")
		return http.StatusInternalServerError
	}

	serviçoAtirador := atirador.NovoServiço(f.Tx(), f.Logger(), config.Atual().Configuração)
	frequênciaConfirmaçãoPedidoCompleta := protocolo.NovaFrequênciaConfirmaçãoPedidoCompleta(f.CR, f.NúmeroControle, f.CódigoVerificação, f.FrequênciaConfirmaçãoPedido)

	if err := serviçoAtirador.ConfirmarFrequência(frequênciaConfirmaçãoPedidoCompleta); err != nil {
		if errors.Equal(err, erros.NãoEncontrado) {
			return http.StatusNotFound
		}

		if mensagens, ok := err.(protocolo.Mensagens); ok {
			f.Mensagens = mensagens
			return http.StatusBadRequest
		}

		f.Logger().Error(erros.Novo(err))
		return http.StatusInternalServerError
	}

	return http.StatusNoContent
}

func (f *frequênciaAtiradorConfirmação) Interceptors() handy.InterceptorChain {
	return criarCorrenteBásica(f).
		Chain(interceptador.NovoBD(f))
}
