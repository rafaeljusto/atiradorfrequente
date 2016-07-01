package protocolo_test

import (
	"reflect"
	"testing"

	"fmt"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
	"github.com/rafaeljusto/atiradorfrequente/testes"
)

func TestNovaMensagem(t *testing.T) {
	esperado := protocolo.Mensagem{
		Código: protocolo.MensagemCódigoParâmetroInválido,
	}

	mensagem := protocolo.NovaMensagem(protocolo.MensagemCódigoParâmetroInválido)

	if !reflect.DeepEqual(esperado, mensagem) {
		t.Errorf("Contrutor de mensagem inválido.\n%v", testes.Diff(esperado, mensagem))
	}
}

func TestNovaMensagemComValor(t *testing.T) {
	esperado := protocolo.Mensagem{
		Código: protocolo.MensagemCódigoParâmetroInválido,
		Valor:  "valor",
	}

	mensagem := protocolo.NovaMensagemComValor(protocolo.MensagemCódigoParâmetroInválido, "valor")

	if !reflect.DeepEqual(esperado, mensagem) {
		t.Errorf("Contrutor de mensagem com valor inválido.\n%v", testes.Diff(esperado, mensagem))
	}
}

func TestNovaMensagemComCampo(t *testing.T) {
	esperado := protocolo.Mensagem{
		Código: protocolo.MensagemCódigoParâmetroInválido,
		Campo:  "campo",
		Valor:  "valor",
	}

	mensagem := protocolo.NovaMensagemComCampo(protocolo.MensagemCódigoParâmetroInválido, "campo", "valor")

	if !reflect.DeepEqual(esperado, mensagem) {
		t.Errorf("Contrutor de mensagem com campo inválido.\n%v", testes.Diff(esperado, mensagem))
	}
}

func TestMensagem_String(t *testing.T) {
	cenários := []struct {
		descrição string
		mensagem  protocolo.Mensagem
		esperado  string
	}{
		{
			descrição: "deve gerar uma mensagem com código, campo e valor",
			mensagem:  protocolo.NovaMensagemComCampo(protocolo.MensagemCódigoParâmetroInválido, "campo", "valor"),
			esperado: fmt.Sprintf("Código de erro “%s” referente ao campo “campo” com valor “valor”",
				protocolo.MensagemCódigoParâmetroInválido),
		},
		{
			descrição: "deve gerar uma mensagem com código e valor",
			mensagem:  protocolo.NovaMensagemComValor(protocolo.MensagemCódigoParâmetroInválido, "valor"),
			esperado: fmt.Sprintf("Código de erro “%s” referente ao valor “valor”",
				protocolo.MensagemCódigoParâmetroInválido),
		},
		{
			descrição: "deve gerar uma mensagem só com código",
			mensagem:  protocolo.NovaMensagem(protocolo.MensagemCódigoParâmetroInválido),
			esperado:  fmt.Sprintf("Código de erro “%s”", protocolo.MensagemCódigoParâmetroInválido),
		},
	}

	for i, cenário := range cenários {
		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(cenário.esperado, nil)
		if err := verificadorResultado.VerificaResultado(cenário.mensagem.String(), nil); err != nil {
			t.Error(err)
		}
	}
}

func TestMensagem_Error(t *testing.T) {
	cenários := []struct {
		descrição string
		mensagem  protocolo.Mensagem
		esperado  string
	}{
		{
			descrição: "deve gerar uma mensagem com código, campo e valor",
			mensagem:  protocolo.NovaMensagemComCampo(protocolo.MensagemCódigoParâmetroInválido, "campo", "valor"),
			esperado: fmt.Sprintf("Código de erro “%s” referente ao campo “campo” com valor “valor”",
				protocolo.MensagemCódigoParâmetroInválido),
		},
		{
			descrição: "deve gerar uma mensagem com código e valor",
			mensagem:  protocolo.NovaMensagemComValor(protocolo.MensagemCódigoParâmetroInválido, "valor"),
			esperado: fmt.Sprintf("Código de erro “%s” referente ao valor “valor”",
				protocolo.MensagemCódigoParâmetroInválido),
		},
		{
			descrição: "deve gerar uma mensagem só com código",
			mensagem:  protocolo.NovaMensagem(protocolo.MensagemCódigoParâmetroInválido),
			esperado:  fmt.Sprintf("Código de erro “%s”", protocolo.MensagemCódigoParâmetroInválido),
		},
	}

	for i, cenário := range cenários {
		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(cenário.esperado, nil)
		if err := verificadorResultado.VerificaResultado(cenário.mensagem.Error(), nil); err != nil {
			t.Error(err)
		}
	}
}

func TestNovasMensagens(t *testing.T) {
	esperado := protocolo.Mensagens{
		protocolo.NovaMensagem(protocolo.MensagemCódigoParâmetroInválido),
		protocolo.NovaMensagemComValor(protocolo.MensagemCódigoParâmetroInválido, "valor"),
		protocolo.NovaMensagemComCampo(protocolo.MensagemCódigoParâmetroInválido, "campo", "valor"),
	}

	mensagens := protocolo.NovasMensagens(
		protocolo.NovaMensagem(protocolo.MensagemCódigoParâmetroInválido),
		protocolo.NovaMensagemComValor(protocolo.MensagemCódigoParâmetroInválido, "valor"),
		protocolo.NovaMensagemComCampo(protocolo.MensagemCódigoParâmetroInválido, "campo", "valor"),
	)

	if !reflect.DeepEqual(esperado, mensagens) {
		t.Errorf("Contrutor de mensagens inválido.\n%v", testes.Diff(esperado, mensagens))
	}
}

func TestMensagens_String(t *testing.T) {
	cenários := []struct {
		descrição string
		mensagens protocolo.Mensagens
		esperado  string
	}{
		{
			descrição: "deve gerar uma mensagem com várias linhas",
			mensagens: protocolo.NovasMensagens(
				protocolo.NovaMensagemComCampo(protocolo.MensagemCódigoParâmetroInválido, "campo", "valor"),
				protocolo.NovaMensagemComValor(protocolo.MensagemCódigoParâmetroInválido, "valor"),
				protocolo.NovaMensagem(protocolo.MensagemCódigoParâmetroInválido),
			),
			esperado: fmt.Sprintf(`Mensagens:
	* Código de erro “%s” referente ao campo “campo” com valor “valor”
	* Código de erro “%s” referente ao valor “valor”
	* Código de erro “%s”`,
				protocolo.MensagemCódigoParâmetroInválido,
				protocolo.MensagemCódigoParâmetroInválido,
				protocolo.MensagemCódigoParâmetroInválido),
		},
		{
			descrição: "deve detectar quando não existem mensagens",
			mensagens: protocolo.NovasMensagens(),
		},
	}

	for i, cenário := range cenários {
		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(cenário.esperado, nil)
		if err := verificadorResultado.VerificaResultado(cenário.mensagens.String(), nil); err != nil {
			t.Error(err)
		}
	}
}

func TestMensagens_Error(t *testing.T) {
	cenários := []struct {
		descrição string
		mensagens protocolo.Mensagens
		esperado  string
	}{
		{
			descrição: "deve gerar uma mensagem com várias linhas",
			mensagens: protocolo.NovasMensagens(
				protocolo.NovaMensagemComCampo(protocolo.MensagemCódigoParâmetroInválido, "campo", "valor"),
				protocolo.NovaMensagemComValor(protocolo.MensagemCódigoParâmetroInválido, "valor"),
				protocolo.NovaMensagem(protocolo.MensagemCódigoParâmetroInválido),
			),
			esperado: fmt.Sprintf(`Mensagens:
	* Código de erro “%s” referente ao campo “campo” com valor “valor”
	* Código de erro “%s” referente ao valor “valor”
	* Código de erro “%s”`,
				protocolo.MensagemCódigoParâmetroInválido,
				protocolo.MensagemCódigoParâmetroInválido,
				protocolo.MensagemCódigoParâmetroInválido),
		},
		{
			descrição: "deve detectar quando não existem mensagens",
			mensagens: protocolo.NovasMensagens(),
		},
	}

	for i, cenário := range cenários {
		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(cenário.esperado, nil)
		if err := verificadorResultado.VerificaResultado(cenário.mensagens.Error(), nil); err != nil {
			t.Error(err)
		}
	}
}
