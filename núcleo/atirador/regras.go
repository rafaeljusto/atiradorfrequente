package atirador

import "github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"

func validarCR(cr string, frequência frequência) protocolo.Mensagens {
	if cr != frequência.CR {
		return protocolo.NovasMensagens(
			protocolo.NovaMensagemComValor(protocolo.MensagemCódigoCRInválido, cr),
		)
	}

	return nil
}

func validarNúmeroControle(númeroControle protocolo.NúmeroControle, frequência frequência) protocolo.Mensagens {
	if númeroControle.ID() != frequência.ID || númeroControle.Controle() != frequência.Controle {
		return protocolo.NovasMensagens(
			protocolo.NovaMensagemComValor(protocolo.MensagemCódigoNúmeroControleInválido, númeroControle.String()),
		)
	}

	return nil
}

func validarIntervaloMáximoConfirmação(frequência frequência) protocolo.Mensagens {
	// TODO(rafaeljusto): O intervalo deve ser configurável, mas não podemos
	// acessar diretamente a configuração deste ponto.
	return nil
}

// TODO(rafaeljusto): Validar imagem de confirmação. Podemos validar o formato.
