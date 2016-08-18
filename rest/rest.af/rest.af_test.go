package main

import (
	"os"
	"testing"
)

func Test_main(t *testing.T) {
	cenários := []struct {
		descrição         string
		argumentos        []string
		variáveisAmbiente map[string]string
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
