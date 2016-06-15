package interceptador

import "net/http"

type CabeçalhoCompatível struct {
	Cabeçalho http.Header `response:"header"`
}

func (c *CabeçalhoCompatível) DefinirCabeçalho(chave, valor string) {
	if c.Cabeçalho == nil {
		c.Cabeçalho = make(http.Header)
	}

	c.Cabeçalho.Set(chave, valor)
}

func (c *CabeçalhoCompatível) AdicionarCabeçalho(chave, valor string) {
	if c.Cabeçalho == nil {
		c.Cabeçalho = make(http.Header)
	}

	c.Cabeçalho.Add(chave, valor)
}
