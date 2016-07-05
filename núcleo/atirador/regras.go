package atirador

import (
	"time"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
)

// validarCR garante que a frequência referente ao ID bate com o CR informado
// pelo usuário.
func validarCR(cr string, frequência frequência) protocolo.Mensagens {
	if cr != frequência.CR {
		return protocolo.NovasMensagens(
			protocolo.NovaMensagemComValor(protocolo.MensagemCódigoCRInválido, cr),
		)
	}

	return nil
}

// validarNúmeroControle garante que a frequência referente ao ID bate com o
// número de controle (ID + número aleatório).
func validarNúmeroControle(númeroControle protocolo.NúmeroControle, frequência frequência) protocolo.Mensagens {
	if númeroControle.ID() != frequência.ID || númeroControle.Controle() != frequência.Controle {
		return protocolo.NovasMensagens(
			protocolo.NovaMensagemComValor(protocolo.MensagemCódigoNúmeroControleInválido, númeroControle.String()),
		)
	}

	return nil
}

// validarIntervaloMáximoConfirmação verifique se o prazo máximo para envio da
// confirmação já expirou. No caso de expirado a mensagem de erro retornada
// informa qual foi a data limite.
func validarIntervaloMáximoConfirmação(frequência frequência, prazoConfirmação time.Duration) protocolo.Mensagens {
	if data := frequência.DataCriação.Add(prazoConfirmação); data.Before(time.Now()) {
		return protocolo.NovasMensagens(
			protocolo.NovaMensagemComValor(protocolo.MensagemCódigoPrazoConfirmaçãoExpirado, data.Format(time.RFC3339)),
		)
	}

	return nil
}

// TODO(rafaeljusto): Validar imagem de confirmação. Podemos validar o formato.
