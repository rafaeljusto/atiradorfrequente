package handler

import (
	"net/http"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/atirador"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
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
	NúmeroControle              int64                                 `urivar:"numeroControle"`
	FrequênciaConfirmaçãoPedido protocolo.FrequênciaConfirmaçãoPedido `request:"put"`
}

func (f *frequênciaAtiradorConfirmação) Put() int {
	serviçoAtirador := atirador.NovoServiço(f.Tx())
	frequênciaConfirmaçãoPedidoCompleta := protocolo.NovaFrequênciaConfirmaçãoPedidoCompleta(f.CR, f.NúmeroControle, f.FrequênciaConfirmaçãoPedido)

	if err := serviçoAtirador.ConfirmarFrequência(frequênciaConfirmaçãoPedidoCompleta); err != nil {
		return http.StatusInternalServerError
	}

	return http.StatusNoContent
}

func (f *frequênciaAtiradorConfirmação) Interceptors() handy.InterceptorChain {
	return criarCorrenteBásica(f).
		Chain(interceptador.NovoBD(f))
}
