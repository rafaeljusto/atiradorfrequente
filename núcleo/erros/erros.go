package erros

import (
	"database/sql"

	"github.com/registrobr/gostk/errors"
)

var (
	NãoEncontrado = errors.Errorf("Objeto não encontrado")
	NãoAtualizado = errors.Errorf("Objeto não atualizado devido a problema de versões")
)

func Novo(err error) error {
	if err == sql.ErrNoRows {
		return NãoEncontrado
	}

	return errors.New(err)
}
