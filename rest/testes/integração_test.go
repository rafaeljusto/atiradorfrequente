// +build integração

package testes

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/docker/ctx"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/project/options"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
)

var endereçoServidor string

func TestCriaçãoDeFrequência(t *testing.T) {
	cenários := []struct {
		descrição          string
		requisição         *http.Request
		códigoHTTPEsperado int
		corpoEsperado      string
	}{
		{
			descrição: "deve criar corretamente uma frequência",
			requisição: func() *http.Request {
				frequênciaPedido := protocolo.FrequênciaPedido{
					Calibre:           "calibre .380",
					ArmaUtilizada:     "arma do clube",
					NúmeroSérie:       "za785671",
					GuiaDeTráfego:     762556223,
					QuantidadeMunição: 50,
					DataInício:        time.Now().Add(-30 * time.Minute),
					DataTérmino:       time.Now().Add(-10 * time.Minute),
				}

				corpo, err := json.Marshal(frequênciaPedido)
				if err != nil {
					t.Fatalf("Erro ao gerar os dados da requisição. Detalhes: %s", err)
				}

				url := fmt.Sprintf("http://%s/frequencia/380308", endereçoServidor)
				r, err := http.NewRequest("POST", url, bytes.NewReader(corpo))
				if err != nil {
					t.Fatalf("Erro ao gerar a requisição. Detalhes: %s", err)
				}

				return r
			}(),
			códigoHTTPEsperado: http.StatusCreated,
		},
	}

	for _, cenário := range cenários {
		t.Run(cenário.descrição, func(t *testing.T) {
			var cliente http.Client

			resposta, err := cliente.Do(cenário.requisição)
			if err != nil {
				t.Fatalf("Erro inesperado ao enviar a requisição. Detalhes: %s", err)
			}

			if resposta.StatusCode != cenário.códigoHTTPEsperado {
				t.Errorf("Código HTTP inesperado. Esperava “%d” e obteve “%d”", cenário.códigoHTTPEsperado, resposta.StatusCode)
			}
		})
	}
}

func TestMain(m *testing.M) {
	flag.Parse()

	código := 0
	defer func() {
		os.Exit(código)
	}()

	// descarta os logs da biblioteca de containers
	logrus.SetOutput(ioutil.Discard)

	projeto, err := docker.NewProject(&ctx.Context{
		Context: project.Context{
			ComposeFiles: []string{"docker-compose.yml"},
			ProjectName:  "atiradorfrequente",
		},
	}, nil)

	if err != nil {
		fmt.Printf("Erro ao inicializar o projeto para testes. Detalhes: %s\n", err)
		código = 1
		return
	}

	defer func() {
		err = projeto.Down(context.Background(), options.Down{
			RemoveVolume:  true,
			RemoveImages:  options.ImageType("all"),
			RemoveOrphans: true,
		})

		if err != nil {
			fmt.Printf("Erro ao finalizar o projeto de testes. Detalhes: %s\n", err)
		}
	}()

	err = projeto.Up(context.Background(), options.Up{
		Create: options.Create{
			ForceRecreate: true,
			ForceBuild:    true,
		},
	})

	if err != nil {
		fmt.Printf("Erro ao executar o projeto para testes. Detalhes: %s\n", err)
		código = 2
		return
	}

	// temos que aguardar todos as dependências serem executadas antes de rodar o
	// teste principal
	time.Sleep(10 * time.Second)

	endereçoServidor, err = projeto.Port(context.Background(), 1, "tcp", "restaf", "80")
	if err != nil {
		fmt.Printf("Erro ao obter informações da porta do servidor rest.af. Detalhes: %s\n", err)
		código = 3
		return

	} else if endereçoServidor == "" {
		fmt.Println("Não foi possível obter o endereço do servidor")
		código = 4
		return
	}

	código = m.Run()
}
