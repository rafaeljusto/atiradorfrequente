package simulador

import (
	"net/http"

	"github.com/trajber/handy"
)

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

func (h Handler) Req() *http.Request {
	return h.SimulaRequisição
}

func (h Handler) ResponseWriter() http.ResponseWriter {
	return h.SimulaResposta
}

func (h Handler) URIVars() handy.URIVars {
	return h.SimulaVariáveisEndereço
}

func (h Handler) Get() int {
	return h.SimulaGet()
}

func (h Handler) Post() int {
	return h.SimulaPost()
}

func (h Handler) Put() int {
	return h.SimulaPut()
}

func (h Handler) Delete() int {
	return h.SimulaDelete()
}

func (h Handler) Patch() int {
	return h.SimulaPatch()
}

func (h Handler) Head() int {
	return h.SimulaHead()
}
