package handler

import "github.com/trajber/handy"

var (
	// Rotas armazena todos os endereços registrados para este servidor REST e
	// seus respectivos tratadores.
	Rotas map[string]handy.Constructor
)

func registrar(rota string, handler handy.Constructor) {
	if Rotas == nil {
		Rotas = make(map[string]handy.Constructor)
	}

	Rotas[rota] = handler
}
