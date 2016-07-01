package interceptador_test

import (
	"reflect"
	"testing"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
	"github.com/rafaeljusto/atiradorfrequente/rest/interceptador"
	"github.com/rafaeljusto/atiradorfrequente/testes"
)

func TestMensagensCompatível_DefineMensagens(t *testing.T) {
	esperado := protocolo.Mensagens{
		protocolo.NovaMensagem(protocolo.MensagemCódigoParâmetroInválido),
		protocolo.NovaMensagemComValor(protocolo.MensagemCódigoParâmetroInválido, "valor"),
		protocolo.NovaMensagemComCampo(protocolo.MensagemCódigoParâmetroInválido, "campo", "valor"),
	}

	var mensagens interceptador.MensagensCompatível
	mensagens.DefineMensagens(protocolo.NovasMensagens(
		protocolo.NovaMensagem(protocolo.MensagemCódigoParâmetroInválido),
		protocolo.NovaMensagemComValor(protocolo.MensagemCódigoParâmetroInválido, "valor"),
		protocolo.NovaMensagemComCampo(protocolo.MensagemCódigoParâmetroInválido, "campo", "valor"),
	))

	if !reflect.DeepEqual(esperado, mensagens.Mensagens) {
		t.Errorf("Mensagens não definidas corretamente.\n%v", testes.Diff(esperado, mensagens))
	}
}
