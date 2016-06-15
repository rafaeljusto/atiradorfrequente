package handler

import (
	"github.com/rafaeljusto/atiradorfrequente/rest/interceptador"
	"github.com/trajber/handy"
)

type básico struct {
	handy.DefaultHandler
	interceptador.CabeçalhoCompatível
}
