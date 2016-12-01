// +build integração

package testes

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/rafaeljusto/atiradorfrequente/núcleo/protocolo"
	"github.com/rafaeljusto/atiradorfrequente/testes"
	"github.com/registrobr/gostk/errors"
	"golang.org/x/crypto/hkdf"
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

				códigoVerificação := gerarCódigoVerificação(frequênciaPendenteResposta, 380308, "abc123")

				return http.Header{
					"Content-Type": []string{"application/json; charset=utf-8"},
					"Location":     []string{"/frequencia/380308/" + frequênciaPendenteResposta.NúmeroControle.String() + "?verificacao=" + códigoVerificação},
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
		{
			descrição: "deve detectar quando campos obrigatórios não foram preenchidos",
			requisição: func() *http.Request {
				url := fmt.Sprintf("http://%s/frequencia/380308", endereçoServidor)
				r, err := http.NewRequest("POST", url, nil)
				if err != nil {
					t.Fatalf("Erro ao gerar a requisição. Detalhes: %s", err)
				}

				return r
			}(),
			códigoHTTPEsperado: http.StatusBadRequest,
			cabeçalhoEsperado: func(corpo []byte) (http.Header, error) {
				return http.Header{
					"Content-Type": []string{"application/json; charset=utf-8"},
				}, nil
			},
			corpoEsperado: func(corpo []byte) ([]byte, error) {
				mensagens := protocolo.NovasMensagens(
					protocolo.NovaMensagemComCampo(protocolo.MensagemCódigoCampoNãoPreenchido, "calibre", ""),
					protocolo.NovaMensagemComCampo(protocolo.MensagemCódigoCampoNãoPreenchido, "armaUtilizada", ""),
					protocolo.NovaMensagemComCampo(protocolo.MensagemCódigoCampoNãoPreenchido, "quantidadeMunicao", "0"),
				)

				corpoEsperado, err := json.Marshal(mensagens)
				if err != nil {
					return nil, errors.Errorf("Erro ao gerar os dados da resposta. Detalhes: %s", err)
				}

				return bytes.TrimSpace(corpoEsperado), nil
			},
		},
		{
			descrição: "deve detectar campos com dados inválidos",
			requisição: func() *http.Request {
				frequênciaPedido := protocolo.FrequênciaPedido{
					Calibre:           "calibre .380",
					ArmaUtilizada:     "arma do clube",
					NúmeroSérie:       "785671", // formato inválido
					GuiaDeTráfego:     762556223,
					QuantidadeMunição: 50,
					DataInício:        time.Now().Add(-30 * time.Minute),
					DataTérmino:       time.Now().Add(-40 * time.Minute), // término antes do inicio
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
			códigoHTTPEsperado: http.StatusBadRequest,
			cabeçalhoEsperado: func(corpo []byte) (http.Header, error) {
				return http.Header{
					"Content-Type": []string{"application/json; charset=utf-8"},
				}, nil
			},
			corpoEsperado: func(corpo []byte) ([]byte, error) {
				mensagens := protocolo.NovasMensagens(
					protocolo.NovaMensagemComCampo(protocolo.MensagemCódigoNúmeroSérieInválido, "", "785671"),
					protocolo.NovaMensagem(protocolo.MensagemCódigoDatasPeríodoIncorreto),
				)

				corpoEsperado, err := json.Marshal(mensagens)
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

				// copia campos variáveis da resposta definitiva para não causar
				// problemas. Infelizmente não temos como prever os valores destes
				// campos na resposta esperada.

				if data := resposta.Header.Get("Date"); data != "" {
					cabeçalhoEsperado.Set("Date", data)
				}
				if tamanhoConteúdo := resposta.Header.Get("Content-Length"); tamanhoConteúdo != "" {
					cabeçalhoEsperado.Set("Content-Length", tamanhoConteúdo)
				}

				// quando o Content-Type não esta definido, assume-se text/plain. Com
				// exceção de quando o código HTTP é 24 (NoContent).
				if tipoConteúdo := cabeçalhoEsperado.Get("Content-Type"); tipoConteúdo == "" && resposta.StatusCode != http.StatusNoContent {
					cabeçalhoEsperado.Set("Content-Type", "text/plain; charset=utf-8")
				}

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
				corpo = bytes.TrimSpace(corpo)

				if corpoEsperado == nil && corpo != nil {
					t.Errorf("Corpo inesperado na resposta.\n%s", string(corpo))

				} else if corpoEsperado != nil && corpo == nil {
					t.Error("Corpo inexistente na resposta")

				} else {
					verificadorResultado.DefinirEsperado(string(corpoEsperado), nil)
					if err = verificadorResultado.VerificaResultado(string(corpo), nil); err != nil {
						t.Error(err)
					}
				}
			}
		})
	}
}

func TestObterFrequência(t *testing.T) {
	cenários := []struct {
		descrição          string
		requisição         *http.Request
		códigoHTTPEsperado int
		cabeçalhoEsperado  func(corpo []byte) (http.Header, error)
		corpoEsperado      func(corpo []byte) ([]byte, error)
	}{
		{
			descrição: "deve obter corretamente uma frequência",
			requisição: func() *http.Request {
				url := fmt.Sprintf("http://%s/frequencia/380308/1-1234?verificacao=EdapMpB39KQejKKqTq1CsU3MBb4EKzkckFM3eZqrvhd6", endereçoServidor)
				r, err := http.NewRequest("GET", url, nil)
				if err != nil {
					t.Fatalf("Erro ao gerar a requisição. Detalhes: %s", err)
				}

				return r
			}(),
			códigoHTTPEsperado: http.StatusOK,
			cabeçalhoEsperado: func(corpo []byte) (http.Header, error) {
				return http.Header{
					"Content-Type": []string{"application/json; charset=utf-8"},
				}, nil
			},
			corpoEsperado: func(corpo []byte) ([]byte, error) {
				var frequênciaResposta protocolo.FrequênciaResposta
				if err := json.Unmarshal(corpo, &frequênciaResposta); err != nil {
					return nil, errors.Errorf("Erro ao interpretar o corpo da resposta. Detalhes: %s", err)
				}

				corpoEsperado, err := json.Marshal(frequênciaResposta)
				if err != nil {
					return nil, errors.Errorf("Erro ao gerar os dados da resposta. Detalhes: %s", err)
				}

				return bytes.TrimSpace(corpoEsperado), nil
			},
		},
		{
			descrição: "deve detectar quando uma frequência não existe",
			requisição: func() *http.Request {
				url := fmt.Sprintf("http://%s/frequencia/380308/9999-1234?verificacao=EdapMpB39KQejKKqTq1CsU3MBb4EKzkckFM3eZqrvhd6", endereçoServidor)
				r, err := http.NewRequest("GET", url, nil)
				if err != nil {
					t.Fatalf("Erro ao gerar a requisição. Detalhes: %s", err)
				}

				return r
			}(),
			códigoHTTPEsperado: http.StatusNotFound,
			cabeçalhoEsperado: func(corpo []byte) (http.Header, error) {
				return make(http.Header), nil
			},
			corpoEsperado: func(corpo []byte) ([]byte, error) {
				return nil, nil
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

				// copia campos variáveis da resposta definitiva para não causar
				// problemas. Infelizmente não temos como prever os valores destes
				// campos na resposta esperada.

				if data := resposta.Header.Get("Date"); data != "" {
					cabeçalhoEsperado.Set("Date", data)
				}
				if tamanhoConteúdo := resposta.Header.Get("Content-Length"); tamanhoConteúdo != "" {
					cabeçalhoEsperado.Set("Content-Length", tamanhoConteúdo)
				}

				// quando o Content-Type não esta definido, assume-se text/plain. Com
				// exceção de quando o código HTTP é 24 (NoContent).
				if tipoConteúdo := cabeçalhoEsperado.Get("Content-Type"); tipoConteúdo == "" && resposta.StatusCode != http.StatusNoContent {
					cabeçalhoEsperado.Set("Content-Type", "text/plain; charset=utf-8")
				}

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
				corpo = bytes.TrimSpace(corpo)

				if corpoEsperado == nil && corpo != nil {
					t.Errorf("Corpo inesperado na resposta.\n%s", string(corpo))

				} else if corpoEsperado != nil && corpo == nil {
					t.Error("Corpo inexistente na resposta")

				} else {
					verificadorResultado.DefinirEsperado(string(corpoEsperado), nil)
					if err = verificadorResultado.VerificaResultado(string(corpo), nil); err != nil {
						t.Error(err)
					}
				}
			}
		})
	}
}

func TestConfirmaçãoDeFrequência(t *testing.T) {
	cenários := []struct {
		descrição          string
		requisição         *http.Request
		códigoHTTPEsperado int
		cabeçalhoEsperado  func(corpo []byte) (http.Header, error)
		corpoEsperado      func(corpo []byte) ([]byte, error)
	}{
		{
			descrição: "deve confirmar corretamente uma frequência",
			requisição: func() *http.Request {
				frequênciaConfirmaçãoPedido := protocolo.FrequênciaConfirmaçãoPedido{
					Imagem: "iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKCAYAAACNMs+9AAAATElEQVR4XpWPCwoAIAhDswPrWTyxMcKEBn2EwYPG0yQi2st0gJllO5kHRlUNBIwsrkzjbnN3AdPqc5lXV/gMNqYt7cn9Vvr+NT0k7xkVBY1RndW3lwAAAABJRU5ErkJggg==",
				}

				corpo, err := json.Marshal(frequênciaConfirmaçãoPedido)
				if err != nil {
					t.Fatalf("Erro ao gerar os dados da requisição. Detalhes: %s", err)
				}

				url := fmt.Sprintf("http://%s/frequencia/380308/1-1234?verificacao=EdapMpB39KQejKKqTq1CsU3MBb4EKzkckFM3eZqrvhd6", endereçoServidor)
				r, err := http.NewRequest("PUT", url, bytes.NewReader(corpo))
				if err != nil {
					t.Fatalf("Erro ao gerar a requisição. Detalhes: %s", err)
				}

				return r
			}(),
			códigoHTTPEsperado: http.StatusNoContent,
			cabeçalhoEsperado: func(corpo []byte) (http.Header, error) {
				return make(http.Header), nil
			},
			corpoEsperado: func(corpo []byte) ([]byte, error) {
				return nil, nil
			},
		},
		{
			descrição: "deve identificar uma frequência inexistente",
			requisição: func() *http.Request {
				frequênciaConfirmaçãoPedido := protocolo.FrequênciaConfirmaçãoPedido{
					Imagem: "iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKCAYAAACNMs+9AAAATElEQVR4XpWPCwoAIAhDswPrWTyxMcKEBn2EwYPG0yQi2st0gJllO5kHRlUNBIwsrkzjbnN3AdPqc5lXV/gMNqYt7cn9Vvr+NT0k7xkVBY1RndW3lwAAAABJRU5ErkJggg==",
				}

				corpo, err := json.Marshal(frequênciaConfirmaçãoPedido)
				if err != nil {
					t.Fatalf("Erro ao gerar os dados da requisição. Detalhes: %s", err)
				}

				url := fmt.Sprintf("http://%s/frequencia/380308/9999-1234?verificacao=EdapMpB39KQejKKqTq1CsU3MBb4EKzkckFM3eZqrvhd6", endereçoServidor)
				r, err := http.NewRequest("PUT", url, bytes.NewReader(corpo))
				if err != nil {
					t.Fatalf("Erro ao gerar a requisição. Detalhes: %s", err)
				}

				return r
			}(),
			códigoHTTPEsperado: http.StatusNotFound,
			cabeçalhoEsperado: func(corpo []byte) (http.Header, error) {
				return make(http.Header), nil
			},
			corpoEsperado: func(corpo []byte) ([]byte, error) {
				return nil, nil
			},
		},
		{
			descrição: "deve detectar quando a imagem é inválida",
			requisição: func() *http.Request {
				frequênciaConfirmaçãoPedido := protocolo.FrequênciaConfirmaçãoPedido{
					Imagem: "@@@",
				}

				corpo, err := json.Marshal(frequênciaConfirmaçãoPedido)
				if err != nil {
					t.Fatalf("Erro ao gerar os dados da requisição. Detalhes: %s", err)
				}

				url := fmt.Sprintf("http://%s/frequencia/380308/1-1234?verificacao=EdapMpB39KQejKKqTq1CsU3MBb4EKzkckFM3eZqrvhd6", endereçoServidor)
				r, err := http.NewRequest("PUT", url, bytes.NewReader(corpo))
				if err != nil {
					t.Fatalf("Erro ao gerar a requisição. Detalhes: %s", err)
				}

				return r
			}(),
			códigoHTTPEsperado: http.StatusBadRequest,
			cabeçalhoEsperado: func(corpo []byte) (http.Header, error) {
				return http.Header{
					"Content-Type": []string{"application/json; charset=utf-8"},
				}, nil
			},
			corpoEsperado: func(corpo []byte) ([]byte, error) {
				mensagens := protocolo.NovasMensagens(
					protocolo.NovaMensagemComCampo(protocolo.MensagemCódigoImagemBase64Inválido, "imagem", "@@@"),
				)

				corpoEsperado, err := json.Marshal(mensagens)
				if err != nil {
					return nil, errors.Errorf("Erro ao gerar os dados da resposta. Detalhes: %s", err)
				}

				return bytes.TrimSpace(corpoEsperado), nil
			},
		},
		{
			descrição: "deve detectar quando o formato da imagem não é suportado",
			requisição: func() *http.Request {
				frequênciaConfirmaçãoPedido := protocolo.FrequênciaConfirmaçãoPedido{
					Imagem: "aXNzbyDDqSB1bSB0ZXN0ZQo=",
				}

				corpo, err := json.Marshal(frequênciaConfirmaçãoPedido)
				if err != nil {
					t.Fatalf("Erro ao gerar os dados da requisição. Detalhes: %s", err)
				}

				url := fmt.Sprintf("http://%s/frequencia/380308/1-1234?verificacao=EdapMpB39KQejKKqTq1CsU3MBb4EKzkckFM3eZqrvhd6", endereçoServidor)
				r, err := http.NewRequest("PUT", url, bytes.NewReader(corpo))
				if err != nil {
					t.Fatalf("Erro ao gerar a requisição. Detalhes: %s", err)
				}

				return r
			}(),
			códigoHTTPEsperado: http.StatusBadRequest,
			cabeçalhoEsperado: func(corpo []byte) (http.Header, error) {
				return http.Header{
					"Content-Type": []string{"application/json; charset=utf-8"},
				}, nil
			},
			corpoEsperado: func(corpo []byte) ([]byte, error) {
				mensagens := protocolo.NovasMensagens(
					protocolo.NovaMensagemComCampo(protocolo.MensagemCódigoImagemFormatoInválido, "imagem", "aXNzbyDDqSB1bSB0ZXN0ZQo="),
				)

				corpoEsperado, err := json.Marshal(mensagens)
				if err != nil {
					return nil, errors.Errorf("Erro ao gerar os dados da resposta. Detalhes: %s", err)
				}

				return bytes.TrimSpace(corpoEsperado), nil
			},
		},
		{
			descrição: "deve detectar quando a confirmação já está expirada",
			requisição: func() *http.Request {
				frequênciaConfirmaçãoPedido := protocolo.FrequênciaConfirmaçãoPedido{
					Imagem: "iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKCAYAAACNMs+9AAAATElEQVR4XpWPCwoAIAhDswPrWTyxMcKEBn2EwYPG0yQi2st0gJllO5kHRlUNBIwsrkzjbnN3AdPqc5lXV/gMNqYt7cn9Vvr+NT0k7xkVBY1RndW3lwAAAABJRU5ErkJggg==",
				}

				corpo, err := json.Marshal(frequênciaConfirmaçãoPedido)
				if err != nil {
					t.Fatalf("Erro ao gerar os dados da requisição. Detalhes: %s", err)
				}

				url := fmt.Sprintf("http://%s/frequencia/923714/2-7344?verificacao=A9W4z9dyPMQzGisYDn1zQLHoDiXSUtQc4vFQFByFdj1P", endereçoServidor)
				r, err := http.NewRequest("PUT", url, bytes.NewReader(corpo))
				if err != nil {
					t.Fatalf("Erro ao gerar a requisição. Detalhes: %s", err)
				}

				return r
			}(),
			códigoHTTPEsperado: http.StatusBadRequest,
			cabeçalhoEsperado: func(corpo []byte) (http.Header, error) {
				return http.Header{
					"Content-Type": []string{"application/json; charset=utf-8"},
				}, nil
			},
			corpoEsperado: func(corpo []byte) ([]byte, error) {
				mensagens := protocolo.NovasMensagens(
					protocolo.NovaMensagem(protocolo.MensagemCódigoPrazoConfirmaçãoExpirado),
				)

				corpoEsperado, err := json.Marshal(mensagens)
				if err != nil {
					return nil, errors.Errorf("Erro ao gerar os dados da resposta. Detalhes: %s", err)
				}

				return bytes.TrimSpace(corpoEsperado), nil
			},
		},
		{
			descrição: "deve detectar quando a confirmação já está confirmada",
			requisição: func() *http.Request {
				frequênciaConfirmaçãoPedido := protocolo.FrequênciaConfirmaçãoPedido{
					Imagem: "iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKCAYAAACNMs+9AAAATElEQVR4XpWPCwoAIAhDswPrWTyxMcKEBn2EwYPG0yQi2st0gJllO5kHRlUNBIwsrkzjbnN3AdPqc5lXV/gMNqYt7cn9Vvr+NT0k7xkVBY1RndW3lwAAAABJRU5ErkJggg==",
				}

				corpo, err := json.Marshal(frequênciaConfirmaçãoPedido)
				if err != nil {
					t.Fatalf("Erro ao gerar os dados da requisição. Detalhes: %s", err)
				}

				url := fmt.Sprintf("http://%s/frequencia/114239/3-1246?verificacao=DqkQsDhqiYpRGPgq9LDzWQvLNqTxWnxURF47R99ztEdm", endereçoServidor)
				r, err := http.NewRequest("PUT", url, bytes.NewReader(corpo))
				if err != nil {
					t.Fatalf("Erro ao gerar a requisição. Detalhes: %s", err)
				}

				return r
			}(),
			códigoHTTPEsperado: http.StatusBadRequest,
			cabeçalhoEsperado: func(corpo []byte) (http.Header, error) {
				return http.Header{
					"Content-Type": []string{"application/json; charset=utf-8"},
				}, nil
			},
			corpoEsperado: func(corpo []byte) ([]byte, error) {
				mensagens := protocolo.NovasMensagens(
					protocolo.NovaMensagem(protocolo.MensagemCódigoFrequênciaJáConfirmada),
				)

				corpoEsperado, err := json.Marshal(mensagens)
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

				// copia campos variáveis da resposta definitiva para não causar
				// problemas. Infelizmente não temos como prever os valores destes
				// campos na resposta esperada.

				if data := resposta.Header.Get("Date"); data != "" {
					cabeçalhoEsperado.Set("Date", data)
				}
				if tamanhoConteúdo := resposta.Header.Get("Content-Length"); tamanhoConteúdo != "" {
					cabeçalhoEsperado.Set("Content-Length", tamanhoConteúdo)
				}

				// quando o Content-Type não esta definido, assume-se text/plain. Com
				// exceção de quando o código HTTP é 24 (NoContent).
				if tipoConteúdo := cabeçalhoEsperado.Get("Content-Type"); tipoConteúdo == "" && resposta.StatusCode != http.StatusNoContent {
					cabeçalhoEsperado.Set("Content-Type", "text/plain; charset=utf-8")
				}

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
				corpo = bytes.TrimSpace(corpo)

				if corpoEsperado == nil && corpo != nil {
					t.Errorf("Corpo inesperado na resposta.\n%s", string(corpo))

				} else if corpoEsperado != nil && corpo == nil {
					t.Error("Corpo inexistente na resposta")

				} else {
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

	defer pararServidorREST()

	var err error
	if endereçoServidor, err = iniciarServidorREST(); err != nil {
		fmt.Println(err)
		código = 1
		return
	}

	código = m.Run()
}

func gerarCódigoVerificação(f protocolo.FrequênciaPendenteResposta, cr int, chave string) string {
	buffer := new(bytes.Buffer)

	// o erro retornado nesta escrita é ignorado, pois o tipo bytes.Buffer não
	// gera erro no método Write
	binary.Write(buffer, binary.LittleEndian, int64(f.NúmeroControle.ID()))

	derivaçãoChave := make([]byte, 32)
	funçãoDerivação := hkdf.New(sha256.New, []byte(chave), nil, buffer.Bytes())

	// o erro retornado nesta leitura é ignorado, pois como estamos lendo a
	// quantidade total de bytes o cenário de erro nunca será atingido
	io.ReadFull(funçãoDerivação, derivaçãoChave)

	mensagem := fmt.Sprintf("%010d %d %d", f.NúmeroControle.ID(), cr, f.NúmeroControle.Controle())

	// o erro retornado é ignorado, pois o método Write do SHA256 não gera erro
	mac := hmac.New(sha256.New, derivaçãoChave)
	mac.Write([]byte(mensagem))
	mensagemCodificada := base58.Encode(mac.Sum(nil))

	if tamanho := 44 - len(mensagemCodificada); tamanho > 0 {
		mensagemCodificada += strings.Repeat("o", tamanho)
	}
	return mensagemCodificada
}
