package handler

import (
	"github.com/rafaeljusto/cr/rest/interceptador"
	"github.com/trajber/handy"
)

type básico struct {
	handy.DefaultHandler
	interceptador.CabeçalhoCompatível
}
