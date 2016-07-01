package interceptador

import (
	"encoding"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
	"github.com/registrobr/gostk/log"
	"github.com/trajber/handy"
	"github.com/trajber/handy/interceptor"
)

type variáveisEndereço interface {
	Logger() log.Logger
	URIVars() handy.URIVars
	Field(tag, valor string) interface{}
	DefineMensagens(protocolo.Mensagens)
}

// VariáveisEndereço preenche os atributos do handler que se referenciam a
// parâmetros do endereço através da tag "urivar".
type VariáveisEndereço struct {
	interceptor.NopInterceptor
	handler variáveisEndereço
}

// NovaVariáveisEndereço cria um novo interceptador de váriaveis de endereço.
func NovaVariáveisEndereço(v variáveisEndereço) *VariáveisEndereço {
	return &VariáveisEndereço{handler: v}
}

// Before percorre as variáveis de endereço e preenche nos atributos
// correspondentes do handler. Caso ocorra algum erro ao preencher um atributo
// uma mensagem definida para alertar o usuário e detalhes serão escritos no
// log.
func (v *VariáveisEndereço) Before() int {
	v.handler.Logger().Debug("Interceptador Antes: Variáveis Endereço")

	for nomeCampo, valor := range v.handler.URIVars() {
		campo := v.handler.Field("urivar", nomeCampo)
		if campo == nil {
			v.handler.Logger().Warningf("Tentando definir um valor no campo “%s” que não existe", nomeCampo)
			continue
		}

		if mensagens := defineValorCampo(campo, nomeCampo, valor, v.handler.Logger()); mensagens != nil {
			v.handler.DefineMensagens(mensagens)
			return http.StatusBadRequest
		}
	}

	return 0
}

func defineValorCampo(ptr interface{}, nomeCampo, valor string, logger log.Logger) protocolo.Mensagens {
	switch f := ptr.(type) {
	case *string:
		*f = valor

	case *bool:
		lower := strings.ToLower(valor)
		*f = lower == "true"

	case *int, *int8, *int16, *int32, *int64:
		return defineValorCampoInt(ptr, nomeCampo, valor, logger)

	case *uint, *uint8, *uint16, *uint32, *uint64:
		return defineValorCampoUint(ptr, nomeCampo, valor, logger)

	case *float32, *float64:
		return defineValorCampoFloat(ptr, nomeCampo, valor, logger)

	default:
		return defineValorCampoUnmarshal(ptr, nomeCampo, valor, logger)
	}

	return nil
}

func defineValorCampoInt(ptr interface{}, nomeCampo, valor string, logger log.Logger) protocolo.Mensagens {
	n, err := strconv.ParseInt(valor, 10, 64)
	if err != nil {
		logger.Error(err)
		return protocolo.NovasMensagens(
			protocolo.NovaMensagemComCampo(protocolo.MensagemCódigoParâmetroInválido, nomeCampo, valor),
		)
	}

	v := reflect.ValueOf(ptr)
	v.Elem().SetInt(n)
	return nil
}

func defineValorCampoUint(ptr interface{}, nomeCampo, valor string, logger log.Logger) protocolo.Mensagens {
	n, err := strconv.ParseUint(valor, 10, 64)
	if err != nil {
		logger.Error(err)
		return protocolo.NovasMensagens(
			protocolo.NovaMensagemComCampo(protocolo.MensagemCódigoParâmetroInválido, nomeCampo, valor),
		)
	}

	v := reflect.ValueOf(ptr)
	v.Elem().SetUint(n)
	return nil
}

func defineValorCampoFloat(ptr interface{}, nomeCampo, valor string, logger log.Logger) protocolo.Mensagens {
	n, err := strconv.ParseFloat(valor, 64)
	if err != nil {
		logger.Error(err)
		return protocolo.NovasMensagens(
			protocolo.NovaMensagemComCampo(protocolo.MensagemCódigoParâmetroInválido, nomeCampo, valor),
		)
	}

	v := reflect.ValueOf(ptr)
	v.Elem().SetFloat(n)
	return nil
}

func defineValorCampoUnmarshal(ptr interface{}, nomeCampo, valor string, logger log.Logger) protocolo.Mensagens {
	u, ok := ptr.(encoding.TextUnmarshaler)
	if !ok {
		logger.Errorf("tipo de valor não suportado: %#v", ptr)
		return protocolo.NovasMensagens(
			protocolo.NovaMensagemComCampo(protocolo.MensagemCódigoParâmetroInválido, nomeCampo, valor),
		)
	}

	if err := u.UnmarshalText([]byte(valor)); err != nil {
		logger.Error(err)

		if mensagens, ok := err.(protocolo.Mensagens); ok && mensagens != nil {
			return mensagens
		}

		return protocolo.NovasMensagens(
			protocolo.NovaMensagemComCampo(protocolo.MensagemCódigoParâmetroInválido, nomeCampo, valor),
		)
	}

	return nil
}
