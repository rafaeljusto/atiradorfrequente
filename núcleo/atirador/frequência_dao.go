package atirador

import (
	"fmt"
	"strings"
	"time"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/bd"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/erros"
)

type frequênciaDAO interface {
	criar(*frequência) error
	atualizar(*frequência) error
	resgatar(id int64) (frequência, error)
}

var novaFrequênciaDAO = func(sqlogger *bd.SQLogger) frequênciaDAO {
	return frequênciaDAOImpl{sqlogger: sqlogger}
}

type frequênciaDAOImpl struct {
	sqlogger *bd.SQLogger
}

func (f frequênciaDAOImpl) criar(frequência *frequência) error {
	frequência.DataCriação = time.Now().UTC()
	frequência.revisão = 0

	resultado, err := f.sqlogger.Exec(frequênciaCriaçãoComando,
		frequência.Controle,
		frequência.CR,
		frequência.Calibre,
		frequência.ArmaUtilizada,
		frequência.NúmeroSérie,
		frequência.GuiaDeTráfego,
		frequência.QuantidadeMunição,
		frequência.DataInício.UTC(),
		frequência.DataTérmino.UTC(),
		frequência.DataCriação.UTC(),
		frequência.revisão,
	)

	if err != nil {
		return erros.Novo(err)
	}

	if frequência.ID, err = resultado.LastInsertId(); err != nil {
		return erros.Novo(err)
	}

	// TODO(rafaeljusto): inserir na tabela log

	return nil
}

func (f frequênciaDAOImpl) atualizar(frequência *frequência) error {
	frequência.DataAtualização = time.Now().UTC()
	frequência.revisão++

	resultado, err := f.sqlogger.Exec(frequênciaAtualizaçãoComando,
		frequência.DataAtualização.UTC(),
		frequência.DataConfirmação.UTC(),
		frequência.revisão,
		frequência.ImagemNúmeroControle,
		frequência.ImagemConfirmação,
		frequência.ID,
		frequência.revisão-1,
	)

	if err != nil {
		return erros.Novo(err)
	}

	atualizados, err := resultado.RowsAffected()

	if err != nil {
		return erros.Novo(err)
	}

	if atualizados != 1 {
		return erros.NãoAtualizado
	}

	// TODO(rafaeljusto): inserir na tabela log

	return nil
}

func (f frequênciaDAOImpl) resgatar(id int64) (frequência, error) {
	resultado := f.sqlogger.QueryRow(frequênciaResgateComando, id)

	var freq frequência
	err := resultado.Scan(
		&freq.ID,
		&freq.Controle,
		&freq.CR,
		&freq.Calibre,
		&freq.ArmaUtilizada,
		&freq.NúmeroSérie,
		&freq.GuiaDeTráfego,
		&freq.QuantidadeMunição,
		&freq.DataInício,
		&freq.DataTérmino,
		&freq.DataCriação,
		&freq.DataAtualização,
		&freq.DataConfirmação,
		&freq.ImagemNúmeroControle,
		&freq.ImagemConfirmação,
		&freq.revisão,
	)

	return freq, erros.Novo(err)
}

var (
	frequênciaTabela = "frequencia_atirador"

	frequênciaCriaçãoCampos = []string{
		"id",
		"controle",
		"cr",
		"calibre",
		"arma_utilizada",
		"numero_serie",
		"guia_de_trafego",
		"quantidade_municao",
		"data_inicio",
		"data_termino",
		"data_criacao",
		"revisao",
	}
	frequênciaCriaçãoCamposTexto = strings.Join(frequênciaCriaçãoCampos, ", ")
	frequênciaCriaçãoComando     = fmt.Sprintf(`INSERT INTO %s (%s) VALUES (DEFAULT, %s)`,
		frequênciaTabela, frequênciaCriaçãoCamposTexto, bd.MarcadoresPSQL(len(frequênciaCriaçãoCampos)-1))

	frequênciaAtualizaçãoComando = fmt.Sprintf(`UPDATE %s SET
	data_atualizacao = $1,
	data_confirmacao = $2,
	revisao = $3,
	imagem_numero_controle = $4,
	imagem_confirmacao = $5
	WHERE id = $6 AND revisao = $7`, frequênciaTabela)

	frequênciaResgateCampos = []string{
		"id",
		"controle",
		"cr",
		"calibre",
		"arma_utilizada",
		"numero_serie",
		"guia_de_trafego",
		"quantidade_municao",
		"data_inicio",
		"data_termino",
		"data_criacao",
		"data_atualizacao",
		"data_confirmacao",
		"imagem_numero_controle",
		"imagem_confirmacao",
		"revisao",
	}
	frequênciaResgateCamposTexto = strings.Join(frequênciaResgateCampos, ", ")
	frequênciaResgateComando     = fmt.Sprintf(`SELECT %s FROM %s WHERE id = $1`,
		frequênciaTabela, frequênciaResgateCamposTexto)
)
