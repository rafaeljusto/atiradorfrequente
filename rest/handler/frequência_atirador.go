package handler

import (
	"fmt"
	"net/http"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/atirador"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
	"github.com/trajber/handy"
)

func init() {
	registrar("/frequencia/{cr}", func() handy.Handler { return &frequênciaAtirador{} })
}

type frequênciaAtirador struct {
	básico

	CR                 string                       `urivar:"cr"`
	FrequênciaPedido   protocolo.FrequênciaPedido   `request:"post"`
	FrequênciaResposta protocolo.FrequênciaResposta `response:"post"`
}

func (f *frequênciaAtirador) Post() int {
	frequênciaPedidoCompleta := protocolo.NovaFrequênciaPedidoCompleta(f.CR, f.FrequênciaPedido)
	serviçoAtirador := atirador.NovoServiço()

	var err error
	if f.FrequênciaResposta, err = serviçoAtirador.CadastrarFrequência(frequênciaPedidoCompleta); err != nil {
		return http.StatusInternalServerError
	}

	f.DefinirCabeçalho("Location", fmt.Sprintf("/frequencia-atirador/%s/%d", f.CR, f.FrequênciaResposta.NúmeroControle))
	return http.StatusCreated
}

func (f *frequênciaAtirador) Interceptors() handy.InterceptorChain {
	return nil
}
