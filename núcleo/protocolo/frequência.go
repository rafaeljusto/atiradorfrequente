package protocolo

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var númeroSérieFormato = regexp.MustCompile(`^([A-Z]+[0-9]+)?$`)
var númeroControleFormato = regexp.MustCompile(`^[0-9]+\-[0-9]+$`)

// FrequênciaPedido armazena os dados exigidos pelo Exército ao utilizar um
// estande de Tiro.
type FrequênciaPedido struct {
	Calibre           string `json:"calibre"`
	ArmaUtilizada     string `json:"armaUtilizada"`
	NúmeroSérie       string `json:"numeroSerie"`
	GuiaDeTráfego     int    `json:"guiaTrafego"`
	QuantidadeMunição int    `json:"quantidadeMunicao"`

	// DataInício data e hora do início do treino de tiro no estande do clube.
	DataInício time.Time `json:"dataInicio"`

	// DataTérmino data e hora do término do treino de tiro no estande do clube.
	DataTérmino time.Time `json:"dataTermino"`
}

// Normalizar padroniza o formato dos campos da requisição. Remove espaços e mantém alguns conteúdos
// em caixa alta.
func (f *FrequênciaPedido) Normalizar() {
	f.Calibre = strings.TrimSpace(f.Calibre)
	f.Calibre = strings.ToUpper(f.Calibre)

	f.ArmaUtilizada = strings.TrimSpace(f.ArmaUtilizada)
	f.ArmaUtilizada = strings.ToUpper(f.ArmaUtilizada)

	f.NúmeroSérie = strings.TrimSpace(f.NúmeroSérie)
	f.NúmeroSérie = strings.ToUpper(f.NúmeroSérie)
}

// Validar analisa se os dados informados possuem o formato correto e se os campos obrigatórios
// foram preenchidos.
func (f FrequênciaPedido) Validar() Mensagens {
	var mensagens Mensagens

	if f.Calibre == "" {
		mensagens = append(mensagens, NovaMensagemComCampo(MensagemCódigoCampoNãoPreenchido, "calibre", ""))
	}

	if f.ArmaUtilizada == "" {
		mensagens = append(mensagens, NovaMensagemComCampo(MensagemCódigoCampoNãoPreenchido, "armaUtilizada", ""))
	}

	if !númeroSérieFormato.MatchString(f.NúmeroSérie) {
		mensagens = append(mensagens, NovaMensagemComValor(MensagemCódigoNúmeroSérieInválido, f.NúmeroSérie))
	}

	if f.QuantidadeMunição == 0 {
		mensagens = append(mensagens, NovaMensagemComCampo(MensagemCódigoCampoNãoPreenchido, "quantidadeMunicao", "0"))
	}

	if f.DataInício.After(f.DataTérmino) || f.DataTérmino.After(time.Now()) {
		mensagens = append(mensagens, NovaMensagem(MensagemCódigoDatasPeríodoIncorreto))
	}

	return mensagens
}

// FrequênciaPedidoCompleta é uma extensão do tipo FrequênciaPedido incluindo o
// CR enviado no endereço.
type FrequênciaPedidoCompleta struct {
	CR int
	FrequênciaPedido
}

// NovaFrequênciaPedidoCompleta inicializa o tipo FrequênciaPedidoCompleta a
// partir do CR e do tipo FrequênciaPedido.
func NovaFrequênciaPedidoCompleta(cr int, frequênciaPedido FrequênciaPedido) FrequênciaPedidoCompleta {
	return FrequênciaPedidoCompleta{
		CR:               cr,
		FrequênciaPedido: frequênciaPedido,
	}
}

// FrequênciaPendenteResposta armazena os dados que permitem ao Clube de Tiro
// confirmar a presença do Atirador.
type FrequênciaPendenteResposta struct {
	NúmeroControle    NúmeroControle `json:"numeroControle"`
	CódigoVerificação string         `json:"codigoVerificacao"`
	Imagem            string         `json:"imagem"` // base64
}

// FrequênciaConfirmaçãoPedido armazena os dados necessários para confirmar a
// presença do Atirador no Clube de Tiro.
type FrequênciaConfirmaçãoPedido struct {
	Imagem string `json:"imagem"` // base64
}

// Normalizar padroniza o formato dos campos da requisição. Remove espaços, mas
// mantém a caixa alta ou baixa para respeitar o padrão do base64.
func (f *FrequênciaConfirmaçãoPedido) Normalizar() {
	f.Imagem = strings.TrimSpace(f.Imagem)
}

// Validar verifica se a imagem enviada na confirmação possuí um formato
// correto.
func (f FrequênciaConfirmaçãoPedido) Validar() Mensagens {
	imagem, err := base64.StdEncoding.DecodeString(f.Imagem)
	if err != nil {
		return NovasMensagens(NovaMensagemComCampo(MensagemCódigoImagemBase64Inválido, "imagem", f.Imagem))
	}

	leitorImagem := bytes.NewReader(imagem)
	if _, _, err = image.Decode(leitorImagem); err != nil {
		return NovasMensagens(NovaMensagemComCampo(MensagemCódigoImagemFormatoInválido, "imagem", f.Imagem))
	}

	return nil
}

// FrequênciaConfirmaçãoPedidoCompleta extende o tipo
// FrequênciaConfirmaçãoPedido incluindo o CR e o número de controle encontrados
// no endereço.
type FrequênciaConfirmaçãoPedidoCompleta struct {
	CR             int
	NúmeroControle NúmeroControle
	FrequênciaConfirmaçãoPedido
}

// NovaFrequênciaConfirmaçãoPedidoCompleta inicializa o tipo
// FrequênciaConfirmaçãoPedidoCompleta a partir do CR, número de controle e do
// tipo FrequênciaConfirmaçãoPedido.
func NovaFrequênciaConfirmaçãoPedidoCompleta(cr int, númeroControle NúmeroControle, frequênciaConfirmaçãoPedido FrequênciaConfirmaçãoPedido) FrequênciaConfirmaçãoPedidoCompleta {
	return FrequênciaConfirmaçãoPedidoCompleta{
		CR:                          cr,
		NúmeroControle:              númeroControle,
		FrequênciaConfirmaçãoPedido: frequênciaConfirmaçãoPedido,
	}
}

// NúmeroControle número gerado a partir do cadastro de uma frequência para
// comprovação da presença física do atirador no estande de tiro. Formado a
// partir do número de identificação da frequência com um número gerado
// aleatoriamente.
type NúmeroControle string

// NovoNúmeroControle gera um novo número de controle a partir do número de
// identificação da frequência com um número gerado aleatoriamente, chamado de
// controle.
func NovoNúmeroControle(id, controle int64) NúmeroControle {
	return NúmeroControle(fmt.Sprintf("%d-%d", id, controle))
}

// String retorna o formato texto do número de controle.
func (n NúmeroControle) String() string {
	return string(n)
}

// ID retorna o número de identificação da frequência contido no número de
// controle.
func (n NúmeroControle) ID() int64 {
	valor := string(n)
	partes := strings.Split(valor, "-")

	id, err := strconv.ParseInt(partes[0], 10, 64)
	if err != nil {
		return 0
	}

	return id
}

// Controle retorna o número aleatório contido no número de controle.
func (n NúmeroControle) Controle() int64 {
	valor := string(n)
	partes := strings.Split(valor, "-")

	if len(partes) != 2 {
		return 0
	}

	controle, err := strconv.ParseInt(partes[1], 10, 64)
	if err != nil {
		return 0
	}

	return controle
}

// Normalizar padroniza o formato do número de controle.
func (n *NúmeroControle) Normalizar() {
	if n == nil {
		return
	}

	texto := n.String()
	texto = strings.TrimSpace(texto)
	*n = NúmeroControle(texto)
}

// Validar verifica o formato do número de controle.
func (n NúmeroControle) Validar() Mensagens {
	var mensagens Mensagens
	if !númeroControleFormato.MatchString(n.String()) {
		mensagens = append(mensagens, NovaMensagemComValor(MensagemCódigoNúmeroControleInválido, n.String()))
	}
	return mensagens
}

// UnmarshalText converte um texto no tipo NúmeroControle, normalizando e
// validando o valor convertido.
func (n *NúmeroControle) UnmarshalText(texto []byte) error {
	*n = NúmeroControle(texto)
	n.Normalizar()
	return n.Validar().Expor()
}
