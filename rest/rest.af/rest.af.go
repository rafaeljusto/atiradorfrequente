package main

import (
	"fmt"
	"os"

	"github.com/jpillora/overseer"
	"github.com/jpillora/overseer/fetcher"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/erros"
	"github.com/rafaeljusto/atiradorfrequente/rest/config"
	"github.com/rafaeljusto/atiradorfrequente/rest/servidor"
	"github.com/urfave/cli"
)

const (
	códigoSaídaAplicação códigoSaída = 1
)

type códigoSaída int

func (c códigoSaída) Código() int {
	return int(c)
}

func main() {
	app := cli.NewApp()
	app.Name = "rest.af"
	app.Usage = "Serviços REST que controlam a frequência dos atiradores nos estandes de tiro"
	app.Author = "Rafael Dantas Justo"
	app.Version = config.Versão

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config,c",
			EnvVar: "AF_REST_CONFIG",
			Usage:  "arquivo de configuração",
		},
	}

	app.Action = cli.ActionFunc(func(c *cli.Context) error {
		config.DefinirValoresPadrão()

		if arquivo := c.String("config"); arquivo != "" {
			if err := config.CarregarDeArquivo(arquivo); err != nil {
				return erros.Novo(err)
			}
		}

		if err := config.CarregarDeVariávelAmbiente(); err != nil {
			return erros.Novo(err)
		}

		// TODO(rafaeljusto): Mover o carregamento da configuração para dentro da
		// função executor. Assim todas as vezes em que um novo binário for
		// executado, o arquivo de configuração será recarregado. Só precisamos
		// descobrir como obter os valores necessários para iniciar o overseer que
		// hoje estão dentro da própria configuração.

		err := overseer.RunErr(overseer.Config{
			Program: executor,
			Address: config.Atual().Servidor.Endereço,
			Fetcher: &fetcher.HTTP{
				URL:      config.Atual().Binário.URL,
				Interval: config.Atual().Binário.TempoAtualização,
			},
		})

		return erros.Novo(err)
	})

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao executar a aplicação. Detalhes: %s\n", err)
		os.Exit(códigoSaídaAplicação.Código())
	}
}

func executor(estado overseer.State) {
	servidor.Iniciar(estado.Listener)
}
