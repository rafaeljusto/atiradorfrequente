package protocolo

import (
	"fmt"
	"strings"
)

const (
	// MensagemCódigoParâmetroInválido indica que uma variável do endereço possui
	// um formato incorreto e não pode ser atribuída a variável correspondente.
	MensagemCódigoParâmetroInválido MensagemCódigo = "parametro-invalido"

	// MensagemCódigoNúmeroControleInválido número de controle informado não é
	// válido.
	MensagemCódigoNúmeroControleInválido MensagemCódigo = "numero-controle-invalido"

	// MensagemCódigoCRInválido CR informado não é válido.
	MensagemCódigoCRInválido MensagemCódigo = "cr-invalido"

	// MensagemCódigoPrazoConfirmaçãoExpirado prazo limite para envio da
	// confirmação de sequência expirado.
	MensagemCódigoPrazoConfirmaçãoExpirado MensagemCódigo = "prazo-confirmacao-expirado"

	// MensagemCódigoDatasPeríodoIncorreto datas informadas de ínico e fim não tem coerência.
	MensagemCódigoDatasPeríodoIncorreto = "datas-periodo-incorreto"

	// MensagemCódigoNúmeroSérieInválido número de série informado não é válido.
	MensagemCódigoNúmeroSérieInválido = "numero-serie-invalido"

	// MensagemCódigoCampoNãoPreenchido indica um campo que possue preenchimento obrigatório.
	MensagemCódigoCampoNãoPreenchido = "campo-nao-preenchido"

	// MensagemCódigoImagemBase64Inválido imagem enviada na confirmação possuí um
	// base64 inválido.
	MensagemCódigoImagemBase64Inválido = "imagem-base64-invalido"

	// MensagemCódigoImagemFormatoInválido imagem enviada na confirmação possuí um
	// formato inválido ou não suportado. Ao gerar este erro a imagem já foi
	// extraída corretamente de um base64.
	MensagemCódigoImagemFormatoInválido = "imagem-formato-invalido"

	// MensagemCódigoImagemNãoAceita imagem não foi aceita por alguma análise da
	// imagem enviada. Geralmente ocorre quando se envia a imagem de número de
	// controle na confirmação da frequência do atirador.
	MensagemCódigoImagemNãoAceita = "imagem-nao-aceita"
)

// MensagemCódigo tipo que define as possíveis mensagens a serem retornadas. A
// utilização de códigos permite a tradução para diferentes idiomas e facilita
// aintegração com outros sistemas.
type MensagemCódigo string

// Mensagem armazena todas as informações necessárias para localizar ao que se
// refere uma mensagem do sistema.
type Mensagem struct {
	Código MensagemCódigo `json:"codigo"`
	Campo  string         `json:"campo,omitempty"`
	Valor  string         `json:"valor,omitempty"`
	Texto  string         `json:"texto,omitempty"`
}

// NovaMensagem cria uma nova mensagem somente com o código. Em alguns casos o
// código da mensagem é o suficiente para identificar um caso.
func NovaMensagem(código MensagemCódigo) Mensagem {
	return Mensagem{
		Código: código,
	}
}

// NovaMensagemComValor cria uma nova mensagem com um código e um valor. Existem
// cenários onde um valor gerou uma mensagem sem ter um campo associado.
func NovaMensagemComValor(código MensagemCódigo, valor string) Mensagem {
	return Mensagem{
		Código: código,
		Valor:  valor,
	}
}

// NovaMensagemComCampo cria uma mensagem com código, campo e valor. Normalmente
// utilizado para indicar problema em algum campo enviado pelo usuário.
func NovaMensagemComCampo(código MensagemCódigo, campo, valor string) Mensagem {
	return Mensagem{
		Código: código,
		Campo:  campo,
		Valor:  valor,
	}
}

// String converte a mensagem para o formato texto. Possui inteligência para
// tratar os casos em que não existe o campo e/ou o valor.
func (m Mensagem) String() string {
	if m.Campo != "" {
		return fmt.Sprintf("Código de erro “%s” referente ao campo “%s” com valor “%s”",
			m.Código, m.Campo, m.Valor)

	} else if m.Valor != "" {
		return fmt.Sprintf("Código de erro “%s” referente ao valor “%s”",
			m.Código, m.Valor)

	} else {
		return fmt.Sprintf("Código de erro “%s”", m.Código)
	}
}

// Error converte a mensagem para o formato de erro. Possui o mesmo
// comportamento do método String.
func (m Mensagem) Error() string {
	return m.String()
}

// Mensagens encapsula múltiplas mensagens.
type Mensagens []Mensagem

// NovasMensagens cria um objeto já populado com mensagens.
func NovasMensagens(mensagens ...Mensagem) Mensagens {
	return Mensagens(mensagens)
}

// String converte as mensagens para o formato texto.
func (m Mensagens) String() string {
	if len(m) == 0 {
		return ""
	}

	mensagens := make([]string, len(m))
	for i, msg := range m {
		mensagens[i] = msg.String()
	}
	return "Mensagens:\n\t* " + strings.Join(mensagens, "\n\t* ")
}

// Error converte as mensagens para o formato de erro. Possui o mesmo
// comportamento do método String.
func (m Mensagens) Error() string {
	return m.String()
}

// JuntarMensagens concatena grupos de mensagens para facilitar a usabilidade.
func JuntarMensagens(m ...Mensagens) Mensagens {
	var mensagensJuntas Mensagens
	for _, mensagens := range m {
		mensagensJuntas = append(mensagensJuntas, mensagens...)
	}
	return mensagensJuntas
}

// TODO(rafaeljusto): Criar método de tradução do código para uma mensagem.
// Embora tenhamos somente o idioma português, devemos evitar de deixar as
// mensagens escritas diretamente no código fonte.
