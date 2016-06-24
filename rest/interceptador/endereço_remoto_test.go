package interceptador_test

import (
	"net"
	"net/http"
	"testing"

	"github.com/rafaeljusto/atiradorfrequente/rest/interceptador"
	"github.com/rafaeljusto/atiradorfrequente/testes"
	"github.com/rafaeljusto/atiradorfrequente/testes/simulador"
)

func TestEndereçoRemoto_Before(t *testing.T) {
	cenários := []struct {
		descrição              string
		endereçoRemoto         string
		cabeçalho              map[string]string
		códigoHTTPEsperado     int
		endereçoRemotoEsperado net.IP
	}{
		{
			descrição: "deve interpretar corretamente quando o cabeçalho X-Forwarded-For possui somente um endereço",
			cabeçalho: map[string]string{
				"X-Forwarded-For": "  192.168.1.1  ",
			},
			endereçoRemotoEsperado: net.ParseIP("192.168.1.1"),
		},
		{
			descrição: "deve interpretar corretamente quando o cabeçalho X-Forwarded-For possui múltiplos endereços",
			cabeçalho: map[string]string{
				"X-Forwarded-For": "  192.168.1.2  ,  192.168.1.3  ",
			},
			endereçoRemotoEsperado: net.ParseIP("192.168.1.2"),
		},
		{
			descrição:      "deve ignorar quando o X-Forwarded-For possuir um endereço IP inválido",
			endereçoRemoto: "  192.168.1.4:123  ",
			cabeçalho: map[string]string{
				"X-Forwarded-For": "  X.X.X.X  ",
			},
			endereçoRemotoEsperado: net.ParseIP("192.168.1.4"),
		},
		{
			descrição: "deve interpretar corretamente quando o cabeçalho X-Real-IP esta definido",
			cabeçalho: map[string]string{
				"X-Real-IP": "  192.168.1.5  ",
			},
			endereçoRemotoEsperado: net.ParseIP("192.168.1.5"),
		},
		{
			descrição:      "deve ignorar quando o X-Real-IP possuir um endereço IP inválido",
			endereçoRemoto: "  192.168.1.6:123  ",
			cabeçalho: map[string]string{
				"X-Real-IP": "  X.X.X.X  ",
			},
			endereçoRemotoEsperado: net.ParseIP("192.168.1.6"),
		},
		{
			descrição:          "deve detectar um endereço IP inválido",
			endereçoRemoto:     "  192.168.1.6  ",
			códigoHTTPEsperado: http.StatusInternalServerError,
		},
	}

	for i, cenário := range cenários {
		requisição, err := http.NewRequest("GET", "/teste", nil)
		if err != nil {
			t.Fatal(err)
		}

		requisição.RemoteAddr = cenário.endereçoRemoto
		for chave, valor := range cenário.cabeçalho {
			requisição.Header.Set(chave, valor)
		}

		handler := &endereçoRemotoSimulado{}
		handler.SimulaRequisição = requisição

		endereçoRemoto := interceptador.NovoEndereçoRemoto(handler)
		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)

		verificadorResultado.DefinirEsperado(cenário.códigoHTTPEsperado, nil)
		if err := verificadorResultado.VerificaResultado(endereçoRemoto.Before(), nil); err != nil {
			t.Error(err)
		}

		verificadorResultado.DefinirEsperado(cenário.endereçoRemotoEsperado, nil)
		if err := verificadorResultado.VerificaResultado(handler.EndereçoRemoto(), nil); err != nil {
			t.Error(err)
		}
	}
}

type endereçoRemotoSimulado struct {
	interceptador.EndereçoRemotoCompatível
	simulador.Handler
}
