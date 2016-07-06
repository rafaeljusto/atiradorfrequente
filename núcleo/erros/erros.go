// Package erros adiciona uma camada extra a estrutura de erros do gostk
// tratando alguns casos excepcionais.
package erros

import (
	"database/sql"

	"github.com/registrobr/gostk/errors"
)

var (
	// NãoEncontrado erro utilizado quando não se encontra um objeto na base de
	// dados.
	NãoEncontrado = errors.Errorf("Objeto não encontrado")

	// NãoAtualizado erro utilizado quando a versão do objeto atualizado não é
	// mais a mesma, ou o objeto já foi removido.
	NãoAtualizado = errors.Errorf("Objeto não atualizado devido a problema de versões")

	// ObjetoIndefinido erro utilizado quando se tenta manipular um objeto não
	// inicializado.
	ObjetoIndefinido = errors.Errorf("Objeto indefinido")
)

// Novo cria um novo erro tratando casos de erros de baixo nível específicos,
// como o caso de dados não encontrados no banco de dados.
func Novo(err error) error {
	if err == sql.ErrNoRows {
		return NãoEncontrado
	}

	return errors.NewWithFollowUp(err, 2)
}
