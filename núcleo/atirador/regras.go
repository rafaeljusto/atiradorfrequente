package atirador

import (
	"strconv"
	"time"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
)

// validarCR garante que a frequência referente ao ID bate com o CR informado
// pelo usuário.
func validarCR(cr int, frequência frequência) protocolo.Mensagens {
	if cr != frequência.CR {
		return protocolo.NovasMensagens(
			protocolo.NovaMensagemComValor(protocolo.MensagemCódigoCRInválido, strconv.Itoa(cr)),
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

// validarDuraçãoTreino garante que o tempo de treino não seja algo anormal.
func validarDuraçãoTreino(frequência frequência, duraçãoMáximaTreino time.Duration) protocolo.Mensagens {
	if intervalo := frequência.DataTérmino.Sub(frequência.DataInício); intervalo > duraçãoMáximaTreino {
		return protocolo.NovasMensagens(
			protocolo.NovaMensagem(protocolo.MensagemCódigoTreinoMuitoLongo),
		)
	}

	return nil
}

// validarTempoMáximoParaCadastro garante que o cadastro não seja feito após
// muito tempo do treino.
func validarTempoMáximoParaCadastro(frequência frequência, tempoMáximaCadastro time.Duration) protocolo.Mensagens {
	if intervalo := time.Now().UTC().Sub(frequência.DataTérmino); intervalo > tempoMáximaCadastro {
		return protocolo.NovasMensagens(
			protocolo.NovaMensagem(protocolo.MensagemCódigoTempoMáximaCadastroExcedido),
		)
	}

	return nil
}

// validarCódigoVerificação analisa se o código de verificação informado é o
// correto. Calculamos o código correto utilizando os dados da frequência com
// uma chave simétrica global.
func validarCódigoVerificação(frequência frequência, chaveCódigoVerificação, códigoVerificação string) protocolo.Mensagens {
	if frequência.gerarCódigoVerificação(chaveCódigoVerificação) != códigoVerificação {
		return protocolo.NovasMensagens(
			protocolo.NovaMensagemComValor(protocolo.MensagemCódigoVerificaçãoInválida, códigoVerificação),
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
			protocolo.NovaMensagem(protocolo.MensagemCódigoPrazoConfirmaçãoExpirado),
		)
	}

	return nil
}

// validarImagemConfirmação evita que a imagem gerada com o número de controle
// seja a mesma utilizada na confirmação da frequência do atirador.
func validarImagemConfirmação(frequência frequência, imagem string) protocolo.Mensagens {
	if frequência.ImagemNúmeroControle == imagem {
		return protocolo.NovasMensagens(
			protocolo.NovaMensagemComValor(protocolo.MensagemCódigoImagemNãoAceita, imagem),
		)
	}

	return nil
}

// validarEstadoFrequência verifica se a frequência já foi confirmada.
func validarEstadoFrequência(frequência frequência) protocolo.Mensagens {
	if !frequência.DataConfirmação.IsZero() && frequência.ImagemConfirmação != "" {
		return protocolo.NovasMensagens(
			protocolo.NovaMensagem(protocolo.MensagemCódigoFrequênciaJáConfirmada),
		)
	}

	return nil
}
