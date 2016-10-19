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
	"github.com/rafaeljusto/atiradorfrequente/testes"
	"github.com/registrobr/gostk/errors"
)

var endereçoServidor string

func TestCriaçãoDeFrequência(t *testing.T) {
	cenários := []struct {
		descrição          string
		requisição         *http.Request
		códigoHTTPEsperado int
		cabeçalhoEsperado  func(corpo []byte) (http.Header, error)
		corpoEsperado      func(corpo []byte) ([]byte, error)
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
			cabeçalhoEsperado: func(corpo []byte) (http.Header, error) {
				var frequênciaPendenteResposta protocolo.FrequênciaPendenteResposta
				if err := json.Unmarshal(corpo, &frequênciaPendenteResposta); err != nil {
					return nil, errors.Errorf("Erro ao interpretar o corpo da resposta. Detalhes: %s", err)
				}

				return http.Header{
					"Content-Type": []string{"application/json; charset=utf-8"},
					"Location":     []string{"/frequencia/380308/" + frequênciaPendenteResposta.NúmeroControle.String()},
				}, nil
			},
			corpoEsperado: func(corpo []byte) ([]byte, error) {
				var frequênciaPendenteResposta protocolo.FrequênciaPendenteResposta
				if err := json.Unmarshal(corpo, &frequênciaPendenteResposta); err != nil {
					return nil, errors.Errorf("Erro ao interpretar o corpo da resposta. Detalhes: %s", err)
				}

				corpoEsperado, err := json.Marshal(frequênciaPendenteResposta)
				if err != nil {
					return nil, errors.Errorf("Erro ao gerar os dados da resposta. Detalhes: %s", err)
				}

				return bytes.TrimSpace(corpoEsperado), nil
			},
		},
	}

	for _, cenário := range cenários {
		t.Run(cenário.descrição, func(t *testing.T) {
			var cliente http.Client

			resposta, err := cliente.Do(cenário.requisição)
			if err != nil {
				t.Fatalf("Erro inesperado ao enviar a requisição. Detalhes: %s", err)
			}
			defer resposta.Body.Close()

			corpo, err := ioutil.ReadAll(resposta.Body)
			if err != nil {
				t.Fatalf("Erro inesperado ao ler o corpo da resposta. Detalhes: %s", err)
			}

			var verificadorResultado testes.VerificadorResultados

			verificadorResultado.DefinirEsperado(cenário.códigoHTTPEsperado, nil)
			if err = verificadorResultado.VerificaResultado(resposta.StatusCode, nil); err != nil {
				t.Error(err)
			}

			if cenário.cabeçalhoEsperado != nil {
				cabeçalhoEsperado, err := cenário.cabeçalhoEsperado(corpo)
				if err != nil {
					t.Fatal(err)
				}

				// como a data pode ser variável, sempre copia a da resposta definitiva
				// para não causar problemas
				cabeçalhoEsperado.Set("Date", resposta.Header.Get("Date"))

				verificadorResultado.DefinirEsperado(cabeçalhoEsperado, nil)
				if err = verificadorResultado.VerificaResultado(resposta.Header, nil); err != nil {
					t.Error(err)
				}
			}

			if cenário.corpoEsperado != nil {
				corpoEsperado, err := cenário.corpoEsperado(corpo)
				if err != nil {
					t.Fatal(err)
				}

				if corpoEsperado == nil && corpo != nil {
					t.Errorf("Corpo inesperado na resposta.\n%s", string(corpo))

				} else if corpoEsperado != nil && corpo == nil {
					t.Error("Corpo inexistente na resposta")

				} else {
					corpo = bytes.TrimSpace(corpo)
					verificadorResultado.DefinirEsperado(string(corpoEsperado), nil)
					if err = verificadorResultado.VerificaResultado(string(corpo), nil); err != nil {
						t.Error(err)
					}
				}
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
		código = m.Run()

	case <-time.Tick(20 * time.Second):
		close(abortarInício)
		fmt.Println("Tempo esgotado aguardando o servidor iniciar")
		código = 3
	}
}
