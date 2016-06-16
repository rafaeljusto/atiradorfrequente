package interceptador

import "net/http"

// CabeçalhoCompatível armazena os dados do cabeçalho a serem retornados em uma
// resposta HTTP do servidor REST.
type CabeçalhoCompatível struct {
	Cabeçalho http.Header `response:"header"`
}

// DefinirCabeçalho define um valor para um cabeçalho HTTP, substituíndo algum
// valor que já exista.
func (c *CabeçalhoCompatível) DefinirCabeçalho(chave, valor string) {
	if c.Cabeçalho == nil {
		c.Cabeçalho = make(http.Header)
	}

	c.Cabeçalho.Set(chave, valor)
}

// AdicionarCabeçalho adiciona um valor no cabeçalho HTTP, podendo ser agregado
// a algum valor já existente com a mesma chave.
func (c *CabeçalhoCompatível) AdicionarCabeçalho(chave, valor string) {
	if c.Cabeçalho == nil {
		c.Cabeçalho = make(http.Header)
	}

	c.Cabeçalho.Add(chave, valor)
}
