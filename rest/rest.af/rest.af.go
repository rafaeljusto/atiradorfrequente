package main

import (
	"fmt"
	"os"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/erros"
	"github.com/rafaeljusto/atiradorfrequente/rest/config"
	"github.com/rafaeljusto/atiradorfrequente/rest/servidor"
	"github.com/rafaeljusto/overseer"
	"github.com/rafaeljusto/overseer/fetcher"
	"github.com/urfave/cli"
)

// teste define o modo de execução sem criar um sub-processo e permitindo uma
// execução única.
var teste = false

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
				fmt.Fprintf(os.Stderr, "Erro ao carregar o arquivo de configuração. Detalhes: %s\n", erros.Novo(err))
				return nil
			}
		}

		if err := config.CarregarDeVariávelAmbiente(); err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao carregar as variáveis de ambiente. Detalhes: %s\n", erros.Novo(err))
			return nil
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
			Test: teste,
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao executar a aplicação. Detalhes: %s\n", erros.Novo(err))
		}

		return nil
	})

	// não verificamos o erro de retorno aqui, pois por padrão a biblioteca já
	// encerra a aplicação em caso de erro. A única situação em que seria
	// interessante analisar o erro seria no caso de configurar argumentos
	// repetidos ou inválidos, mas isto pode ser resolvido no ambiente de
	// desenvolvimento.
	app.Run(os.Args)
}

func executor(estado overseer.State) {
	servidor.Iniciar(estado.Listener)
}
