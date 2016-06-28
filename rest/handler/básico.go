package handler

import (
	"net"
	"net/http"

	"github.com/rafaeljusto/atiradorfrequente/rest/interceptador"
	"github.com/registrobr/gostk/log"
	"github.com/trajber/handy"
)

type básico struct {
	handy.DefaultHandler
	interceptador.CabeçalhoCompatível
	interceptador.EndereçoRemotoCompatível
	interceptador.LogCompatível
}

type correnteBásica interface {
	Req() *http.Request
	DefineEndereçoRemoto(net.IP)
	EndereçoRemoto() net.IP
	EndereçoProxy() net.IP
	DefineEndereçoProxy(net.IP)
	DefineLogger(log.Logger)
	Logger() log.Logger
}

func criarCorrenteBásica(c correnteBásica) handy.InterceptorChain {
	return handy.NewInterceptorChain().
		Chain(interceptador.NovoEndereçoRemoto(c)).
		Chain(interceptador.NovoLog(c))
}
