// +build integração

package testes

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/docker/ctx"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/project/options"
)

var projeto project.APIProject

// iniciarServidorREST executa o servidor REST e todas as suas dependências
// utilizando containers.
func iniciarServidorREST() (endereçoServidor string, err error) {
	// descarta os logs da biblioteca de containers
	logrus.SetOutput(ioutil.Discard)

	projeto, err = docker.NewProject(&ctx.Context{
		Context: project.Context{
			ComposeFiles: []string{"docker-compose.yml"},
			ProjectName:  "atiradorfrequente",
		},
	}, nil)

	if err != nil {
		err = fmt.Errorf("Erro ao inicializar o projeto para testes. Detalhes: %s\n", err)
		return
	}

	err = projeto.Up(context.Background(), options.Up{
		Create: options.Create{
			ForceBuild: true,
		},
	})

	if err != nil {
		err = fmt.Errorf("Erro ao executar o projeto para testes. Detalhes: %s\n", err)
		return
	}

	abortarInício := make(chan bool)
	servidorRodando := make(chan bool)

	go func() {
		for {
			select {
			case <-time.Tick(100 * time.Millisecond):
				break
			case <-abortarInício:
				return
			}

			if endereçoServidor == "" {
				endereçoServidor, err = projeto.Port(context.Background(), 1, "tcp", "restaf", "80")
				if err != nil || endereçoServidor == " " {
					continue
				}
			}

			url := fmt.Sprintf("http://%s/ping", endereçoServidor)
			if resposta, err := http.Get(url); err != nil || resposta.StatusCode != http.StatusNoContent {
				continue
			}

			close(servidorRodando)
			return
		}
	}()

	select {
	case <-servidorRodando:
		break

	case <-time.Tick(20 * time.Second):
		close(abortarInício)
		err = fmt.Errorf("Tempo esgotado aguardando o servidor iniciar")
	}

	return
}

// pararServidorREST para a execução do servidor REST e todos os seus
// containers.
func pararServidorREST() error {
	if projeto == nil {
		return nil
	}

	err := projeto.Down(context.Background(), options.Down{
		RemoveVolume:  true,
		RemoveOrphans: true,
	})

	if err != nil {
		return fmt.Errorf("Erro ao finalizar o projeto de testes. Detalhes: %s\n", err)
	}

	return nil
}
