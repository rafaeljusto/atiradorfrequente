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

	CR                         string                                `urivar:"cr"`
	FrequênciaPedido           protocolo.FrequênciaPedido            `request:"post"`
	FrequênciaPendenteResposta *protocolo.FrequênciaPendenteResposta `response:"post"`
}

func (f *frequênciaAtirador) Post() int {
	serviçoAtirador := atirador.NovoServiço()
	frequênciaPedidoCompleta := protocolo.NovaFrequênciaPedidoCompleta(f.CR, f.FrequênciaPedido)
	frequênciaPendenteResposta, err := serviçoAtirador.CadastrarFrequência(frequênciaPedidoCompleta)

	if err != nil {
		return http.StatusInternalServerError
	}

	f.FrequênciaPendenteResposta = &frequênciaPendenteResposta
	f.DefinirCabeçalho("Location", fmt.Sprintf("/frequencia-atirador/%s/%d", f.CR, f.FrequênciaPendenteResposta.NúmeroControle))
	return http.StatusCreated
}

func (f *frequênciaAtirador) Interceptors() handy.InterceptorChain {
	return nil
}
