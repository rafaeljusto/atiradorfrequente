package interceptador

import "github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"

// MensagensCompatível adiciona a possibilidade de retornar mensagens de erro
// nos handlers em um formato definido.
type MensagensCompatível struct {
	Mensagens protocolo.Mensagens `response:"all"`
}

// DefineMensagens define as mensagens que serão enviadas para o usuário na
// resposta.
func (m *MensagensCompatível) DefineMensagens(mensagens protocolo.Mensagens) {
	m.Mensagens = mensagens
}
