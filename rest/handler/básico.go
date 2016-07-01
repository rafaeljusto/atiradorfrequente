package handler

import (
	"net"
	"net/http"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
	"github.com/rafaeljusto/atiradorfrequente/rest/interceptador"
	"github.com/registrobr/gostk/log"
	"github.com/trajber/handy"
	"github.com/trajber/handy/interceptor"
)

type básico struct {
	handy.DefaultHandler
	interceptador.CabeçalhoCompatível
	interceptador.EndereçoRemotoCompatível
	interceptador.LogCompatível
	interceptador.MensagensCompatível
	interceptor.IntrospectorCompliant
}

type correnteBásica interface {
	Req() *http.Request
	DefineEndereçoRemoto(net.IP)
	EndereçoRemoto() net.IP
	EndereçoProxy() net.IP
	DefineEndereçoProxy(net.IP)
	DefineLogger(log.Logger)
	Logger() log.Logger
	URIVars() handy.URIVars
	Field(tag, valor string) interface{}
	SetFields(interceptor.StructFields)
	DefineMensagens(protocolo.Mensagens)
}

func criarCorrenteBásica(c correnteBásica) handy.InterceptorChain {
	return handy.NewInterceptorChain().
		Chain(interceptador.NovoEndereçoRemoto(c)).
		Chain(interceptador.NovoLog(c)).
		Chain(interceptor.NewIntrospector(c)).
		Chain(interceptador.NovaVariáveisEndereço(c))
}
