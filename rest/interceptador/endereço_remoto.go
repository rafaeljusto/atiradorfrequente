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

// EndereçoRemoto disponibiliza ao handler o endereço do cliente. Possuí a
// capacidade de tratar endereços enviados via proxy tanto com o cabeçalhos HTTP
// X-Forwarded-For quanto com o X-Real-IP.
type EndereçoRemoto struct {
	interceptor.NopInterceptor
	handler endereçoRemoto
}

// NovoEndereçoRemoto cria um novo interceptador EndereçoRemoto.
func NovoEndereçoRemoto(e endereçoRemoto) *EndereçoRemoto {
	return &EndereçoRemoto{handler: e}
}

// Before interpreta a conexão e os cabeçalhos HTTP para identificar o endereço
// IP do cliente. A prioridade é dada na seguinte ordem: Cabeçalhos HTTP
// X-Forwarded-For, X-Real-IP e IP da conexão.
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

// EndereçoRemotoCompatível implementa os métodos que serão utilizados pelo
// handler para acessar o endereço IP remoto armazenado por este interceptador.
type EndereçoRemotoCompatível struct {
	endereçoRemoto net.IP
}

// DefineEndereçoRemoto armazena o endereço IP remoto.
func (r *EndereçoRemotoCompatível) DefineEndereçoRemoto(a net.IP) {
	r.endereçoRemoto = a
}

// EndereçoRemoto obtém o endereço IP remoto.
func (r EndereçoRemotoCompatível) EndereçoRemoto() net.IP {
	return r.endereçoRemoto
}
