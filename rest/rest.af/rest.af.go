package main

import (
	"fmt"

	"github.com/jpillora/overseer"
	"github.com/jpillora/overseer/fetcher"
	"github.com/rafaeljusto/atiradorfrequente/rest/config"
	"github.com/rafaeljusto/atiradorfrequente/rest/servidor"
	"github.com/urfave/cli"
)

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

	app.Action = func(c *cli.Context) {
		config.DefinirValoresPadrão()

		if arquivo := c.String("config"); arquivo != "" {
			if err := config.CarregarDeArquivo(arquivo); err != nil {
				fmt.Printf("Erro ao carregar o arquivo de configuração. Detalhes: %s", err)
				return
			}
		}

		if err := config.CarregarDeVariávelAmbiente(); err != nil {
			fmt.Printf("Erro ao carregar variáveis de ambiente. Detalhes: %s", err)
			return
		}

		// TODO(rafaeljusto): Mover o carregamento da configuração para dentro da
		// função executor. Assim todas as vezes em que um novo binário for
		// executado, o arquivo de configuração será recarregado. Só precisamos
		// descobrir como obter os valores necessários para iniciar o overseer que
		// hoje estão dentro da própria configuração.

		overseer.Run(overseer.Config{
			Program: executor,
			Address: config.Atual().Servidor.Endereço,
			Fetcher: &fetcher.HTTP{
				URL:      config.Atual().Binário.URL,
				Interval: config.Atual().Binário.TempoAtualização,
			},
		})
	}
}

func executor(estado overseer.State) {
	servidor.Iniciar(estado.Listener)
}
