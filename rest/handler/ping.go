package handler

import (
	"net/http"
	"time"

	"github.com/rafaeljusto/atiradorfrequente/rest/interceptador"
	"github.com/registrobr/gostk/errors"
	"github.com/trajber/handy"
)

func init() {
	registrar("/ping", func() handy.Handler { return &ping{} })
}

type ping struct {
	básico
	interceptador.BDCompatível
}

func (p *ping) Get() int {
	resultado := p.Tx().QueryRow("SELECT NOW() AT TIME ZONE 'UTC'")

	var data time.Time
	if err := resultado.Scan(&data); err != nil {
		p.Logger().Error(errors.New(err))
		return http.StatusInternalServerError
	}

	return http.StatusNoContent
}

func (p *ping) Interceptors() handy.InterceptorChain {
	return criarCorrenteBásica(p).
		Chain(interceptador.NovoBD(p))
}
