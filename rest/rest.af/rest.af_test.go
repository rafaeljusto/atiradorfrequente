package main

import (
	"os"
	"regexp"
	"testing"

	"github.com/rafaeljusto/atiradorfrequente/rest/config"
)

func Test_main(t *testing.T) {
	cenários := []struct {
		descrição            string
		argumentos           []string
		variáveisAmbiente    map[string]string
		configuraçãoEsperada *config.Configuração
		mensagensEsperadas   *regexp.Regexp
	}{}

	for _, cenário := range cenários {
		os.Args = cenário.argumentos

		os.Clearenv()
		for chave, valor := range cenário.variáveisAmbiente {
			os.Setenv(chave, valor)
		}

		main()

		// TODO
	}
}
