package interceptador

import (
	"net/http"
	"strings"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
	"github.com/registrobr/gostk/log"
	"github.com/trajber/handy/interceptor"
)

type normalizador interface {
	Normalizar()
}

type validador interface {
	Validar() protocolo.Mensagens
}

type padronizador interface {
	Field(string, string) interface{}
	DefineMensagens(protocolo.Mensagens)
	Logger() log.Logger
	Req() *http.Request
}

// Padronizador normaliza e valida a requisição do usuário.
type Padronizador struct {
	interceptor.NopInterceptor
	handler padronizador
}

// NovoPadronizador cria um novo interceptador Padronizador.
func NovoPadronizador(p padronizador) *Padronizador {
	return &Padronizador{handler: p}
}

// Before faz um tratamento da requisição, padronizando e validando o formato
// dos campos.
func (p Padronizador) Before() int {
	p.handler.Logger().Debug("Interceptador Antes: Padronizador")

	método := strings.ToLower(p.handler.Req().Method)
	campo := p.handler.Field("request", método)
	if campo == nil {
		return 0
	}

	if n, ok := campo.(normalizador); ok {
		n.Normalizar()
	}

	if v, ok := campo.(validador); ok {
		mensagens := v.Validar()
		if mensagens != nil {
			p.handler.DefineMensagens(mensagens)
			return http.StatusBadRequest
		}
	}

	return 0
}
