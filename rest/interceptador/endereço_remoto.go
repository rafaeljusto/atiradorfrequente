package interceptador

import (
	"net"
	"net/http"
	"strings"

	"github.com/registrobr/gostk/log"
	"github.com/trajber/handy/interceptor"
)

type endereçoRemoto interface {
	EndereçoRemoto() net.IP
	DefineEndereçoRemoto(net.IP)
	Req() *http.Request
}

type EndereçoRemoto struct {
	interceptor.NopInterceptor
	handler endereçoRemoto
}

func NovoEndereçoRemoto(e endereçoRemoto) *EndereçoRemoto {
	return &EndereçoRemoto{handler: e}
}

func (r *EndereçoRemoto) Before() int {
	var endereçoCliente string

	xff := r.handler.Req().Header.Get("X-Forwarded-For")
	xff = strings.TrimSpace(xff)

	if len(xff) > 0 {
		xffParts := strings.Split(xff, ",")
		if len(xffParts) == 1 {
			endereçoCliente = strings.TrimSpace(xffParts[0])
		} else if len(xffParts) > 1 {
			endereçoCliente = strings.TrimSpace(xffParts[len(xffParts)-2])
		}

	} else {
		endereçoCliente = strings.TrimSpace(r.handler.Req().Header.Get("X-Real-IP"))
	}

	if len(endereçoCliente) > 0 {
		if endereço := net.ParseIP(endereçoCliente); endereço != nil {
			r.handler.DefineEndereçoRemoto(endereço)
			return 0
		}
	}

	endereçoCliente, _, err := net.SplitHostPort(r.handler.Req().RemoteAddr)

	if err != nil {
		log.Notice(err.Error())
		return http.StatusInternalServerError
	}

	r.handler.DefineEndereçoRemoto(net.ParseIP(strings.TrimSpace(endereçoCliente)))
	return 0
}

type EndereçoRemotoCompatível struct {
	endereçoRemoto net.IP
}

func (r *EndereçoRemotoCompatível) DefineEndereçoRemoto(a net.IP) {
	r.endereçoRemoto = a
}

func (r EndereçoRemotoCompatível) EndereçoRemoto() net.IP {
	return r.endereçoRemoto
}
