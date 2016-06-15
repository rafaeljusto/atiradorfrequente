package handler

import "github.com/trajber/handy"

var (
	Rotas map[string]handy.Constructor
)

func registrar(rota string, handler handy.Constructor) {
	if Rotas == nil {
		Rotas = make(map[string]handy.Constructor)
	}

	Rotas[rota] = handler
}
