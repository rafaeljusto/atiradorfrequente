package handler

import (
	"net/http"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/atirador"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
	"github.com/trajber/handy"
)

func init() {
	registrar("/frequencia/{cr}/{numeroControle}", func() handy.Handler { return &frequênciaAtiradorConfirmação{} })
}

type frequênciaAtiradorConfirmação struct {
	básico

	CR             string `urivar:"cr"`
	NúmeroControle int    `urivar:"numeroControle"`

	frequênciaConfirmaçãoPedido protocolo.FrequênciaConfirmaçãoPedidoCompleta `request:"put"`
}

func (f *frequênciaAtiradorConfirmação) Put() int {
	frequênciaConfirmaçãoPedidoCompleta := protocolo.NovaFrequênciaConfirmaçãoPedidoCompleta(f.CR, f.NúmeroControle, f.frequênciaConfirmaçãoPedido)
	serviçoAtirador := atirador.NovoServiço()

	if err := serviçoAtirador.ConfirmarFrequência(frequênciaConfirmaçãoPedidoCompleta); err != nil {
		return http.StatusInternalServerError
	}

	return http.StatusNoContent
}

func (f *frequênciaAtiradorConfirmação) Interceptors() handy.InterceptorChain {
	return nil
}
