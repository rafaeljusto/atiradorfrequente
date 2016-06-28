package simulador

import (
	"net/http"

	"github.com/trajber/handy"
)

// Handler simula as operações básicas de tratamento de uma requisição.
type Handler struct {
	handy.DefaultHandler

	SimulaRequisição        *http.Request
	SimulaResposta          http.ResponseWriter
	SimulaVariáveisEndereço handy.URIVars
	SimulaGet               func() int
	SimulaPost              func() int
	SimulaPut               func() int
	SimulaDelete            func() int
	SimulaPatch             func() int
	SimulaHead              func() int
}

// Req retorna a requisição HTTP.
func (h Handler) Req() *http.Request {
	return h.SimulaRequisição
}

// ResponseWriter retorna a estrutura que permite escrever respostas na conexão
// HTTP.
func (h Handler) ResponseWriter() http.ResponseWriter {
	return h.SimulaResposta
}

// URIVars retorna as variáveis de endereço da requisição HTTP.
func (h Handler) URIVars() handy.URIVars {
	return h.SimulaVariáveisEndereço
}

// Get trata requisições do tipo GET.
func (h Handler) Get() int {
	return h.SimulaGet()
}

// Post trata requisições do tipo POST.
func (h Handler) Post() int {
	return h.SimulaPost()
}

// Put trata requisições do tipo PUT.
func (h Handler) Put() int {
	return h.SimulaPut()
}

// Delete trata requisições do tipo DELETE.
func (h Handler) Delete() int {
	return h.SimulaDelete()
}

// Patch trata requisições do tipo PATCH.
func (h Handler) Patch() int {
	return h.SimulaPatch()
}

// Head trata requisições do tipo HEAD.
func (h Handler) Head() int {
	return h.SimulaHead()
}
