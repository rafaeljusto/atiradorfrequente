package erros_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/erros"
	"github.com/rafaeljusto/atiradorfrequente/testes"
	"github.com/registrobr/gostk/errors"
)

func TestNovo(t *testing.T) {
	cenários := []struct {
		descrição    string
		erro         error
		erroEsperado error
	}{
		{
			descrição:    "deve detectar um erro de resultados não encontrados no banco de dados",
			erro:         sql.ErrNoRows,
			erroEsperado: erros.NãoEncontrado,
		},
		{
			descrição:    "deve encapsular um erro não específico",
			erro:         fmt.Errorf("erro genérico"),
			erroEsperado: errors.Errorf("erro genérico"),
		},
	}

	for i, cenário := range cenários {
		err := erros.Novo(cenário.erro)

		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(nil, cenário.erroEsperado)
		if err := verificadorResultado.VerificaResultado(nil, err); err != nil {
			t.Error(err)
		}
	}
}
