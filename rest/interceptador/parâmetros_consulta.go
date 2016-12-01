package interceptador

import (
	"net/http"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
	"github.com/registrobr/gostk/log"
	"github.com/trajber/handy/interceptor"
)

type parâmetrosConsulta interface {
	Field(string, string) interface{}
	Logger() log.Logger
	Req() *http.Request
	DefineMensagens(protocolo.Mensagens)
}

// ParâmetrosConsulta preenche os atributos do handler que se referenciam a
// parâmetros da consulta (query string) através da tag "query".
type ParâmetrosConsulta struct {
	interceptor.NopInterceptor
	handler parâmetrosConsulta
}

// NovoParâmetrosConsulta cria um novo interceptador de parâmetros de consulta.
func NovoParâmetrosConsulta(h parâmetrosConsulta) *ParâmetrosConsulta {
	return &ParâmetrosConsulta{handler: h}
}

// Before percorre os parâmetros de consulta e preenche nos atributos
// correspondentes do handler. Caso ocorra algum erro ao preencher um atributo
// uma mensagem é definida para alertar o usuário e detalhes serão escritos no
// log.
func (p *ParâmetrosConsulta) Before() int {
	p.handler.Logger().Debug("Interceptador Antes: Parâmetros Consulta")

	if p.handler.Req().Form == nil {
		p.handler.Req().ParseMultipartForm(32 << 20) // 32 MB
	}

	for chave, valores := range p.handler.Req().Form {
		if len(valores) == 0 {
			continue
		}

		campo := p.handler.Field("query", chave)
		if campo == nil {
			continue
		}

		if mensagens := defineValorCampo(campo, chave, valores[0], p.handler.Logger()); mensagens != nil {
			p.handler.DefineMensagens(mensagens)
			return http.StatusBadRequest
		}
	}

	return 0
}
