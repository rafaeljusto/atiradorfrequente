package atirador

import (
	"fmt"
	"strings"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/bd"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/erros"
)

type frequênciaLogDAO interface {
	criar(frequência, bd.AçãoLog) error
}

var novaFrequênciaLogDAO = func(sqlogger *bd.SQLogger) frequênciaLogDAO {
	return frequênciaLogDAOImpl{sqlogger: sqlogger}
}

type frequênciaLogDAOImpl struct {
	sqlogger *bd.SQLogger
}

func (f frequênciaLogDAOImpl) criar(frequência frequência, ação bd.AçãoLog) error {
	if err := f.sqlogger.Gerar(); err != nil {
		return erros.Novo(err)
	}

	_, err := f.sqlogger.Exec(frequênciaLogCriaçãoComando,
		f.sqlogger.Log.ID,
		ação,
		frequência.ID,
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
		frequência.DataAtualização.UTC(),
		frequência.DataConfirmação.UTC(),
		frequência.ImagemNúmeroControle,
		frequência.ImagemConfirmação,
		frequência.revisão,
	)

	return erros.Novo(err)
}

var (
	frequênciaLogTabela = "frequencia_atirador_log"

	frequênciaLogCriaçãoCampos = []string{
		"id",
		"id_log",
		"acao",
		"id_frequencia_atirador",
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
	frequênciaLogCriaçãoCamposTexto = strings.Join(frequênciaLogCriaçãoCampos, ", ")
	frequênciaLogCriaçãoComando     = fmt.Sprintf(`INSERT INTO %s (%s) VALUES (DEFAULT, %s)`,
		frequênciaLogTabela, frequênciaLogCriaçãoCamposTexto, bd.MarcadoresPSQL(len(frequênciaLogCriaçãoCampos)-1))
)
