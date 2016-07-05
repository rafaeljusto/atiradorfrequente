package interceptador_test

import (
	"fmt"
	"net"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/bd"
	"github.com/rafaeljusto/atiradorfrequente/rest/config"
	"github.com/rafaeljusto/atiradorfrequente/rest/interceptador"
	"github.com/rafaeljusto/atiradorfrequente/testes"
	"github.com/rafaeljusto/atiradorfrequente/testes/simulador"
	"github.com/registrobr/gostk/db"
	"github.com/registrobr/gostk/errors"
	"github.com/registrobr/gostk/log"
)

func TestBD_Before(t *testing.T) {
	cenários := []struct {
		descrição          string
		iniciarConexão     func(db.ConnParams, time.Duration) error
		conexão            bd.BD
		configuração       *config.Configuração
		endereçoRemoto     net.IP
		logger             log.Logger
		códigoHTTPEsperado int
		sqloggerEsperado   *bd.SQLogger
	}{
		{
			descrição: "deve se conectar corretamente ao banco de dados",
			iniciarConexão: func(parâmetrosConexão db.ConnParams, txTempoEsgotado time.Duration) error {
				if parâmetrosConexão.Host != "127.0.0.1" {
					t.Errorf("endereço do banco de dados inesperado: %s", parâmetrosConexão.Host)
				}

				if parâmetrosConexão.DatabaseName != "teste" {
					t.Errorf("nome do banco de dados inesperado: %s", parâmetrosConexão.DatabaseName)
				}

				if parâmetrosConexão.Username != "usuario" {
					t.Errorf("usuário do banco de dados inesperado: %s", parâmetrosConexão.Username)
				}

				if parâmetrosConexão.Password != "senha" {
					t.Errorf("senha do banco de dados inesperada: %s", parâmetrosConexão.Password)
				}

				if parâmetrosConexão.ConnectTimeout != 2*time.Second {
					t.Errorf("tempo esgotado de conexão do banco de dados inesperado: %s", parâmetrosConexão.ConnectTimeout)
				}

				if parâmetrosConexão.StatementTimeout != 10*time.Second {
					t.Errorf("tempo esgotado de comando do banco de dados inesperado: %s", parâmetrosConexão.StatementTimeout)
				}

				if txTempoEsgotado != 1*time.Second {
					t.Errorf("tempo esgotado de transação do banco de dados inesperado: %s", txTempoEsgotado)
				}

				if parâmetrosConexão.MaxIdleConnections != 16 {
					t.Errorf("número máximo de conexões inativas do banco de dados inesperado: %d", parâmetrosConexão.MaxIdleConnections)
				}

				if parâmetrosConexão.MaxOpenConnections != 32 {
					t.Errorf("número máximo de conexões abertas do banco de dados inesperado: %d", parâmetrosConexão.MaxOpenConnections)
				}

				bd.Conexão = simulador.BD{
					SimulaBegin: func() (bd.Tx, error) {
						return simulador.Tx{}, nil
					},
				}

				return nil
			},
			configuração: func() *config.Configuração {
				configuração := new(config.Configuração)
				configuração.BancoDados.Endereço = "127.0.0.1"
				configuração.BancoDados.Nome = "teste"
				configuração.BancoDados.Usuário = "usuario"
				configuração.BancoDados.Senha = "senha"
				configuração.BancoDados.TempoEsgotadoConexão = 2 * time.Second
				configuração.BancoDados.TempoEsgotadoComando = 10 * time.Second
				configuração.BancoDados.TempoEsgotadoTransação = 1 * time.Second
				configuração.BancoDados.MáximoNúmeroConexõesInativas = 16
				configuração.BancoDados.MáximoNúmeroConexõesAbertas = 32
				return configuração
			}(),
			endereçoRemoto: net.ParseIP("192.168.1.1"),
			logger: simulador.Logger{
				SimulaDebug: func(m ...interface{}) {
					mensagem := fmt.Sprint(m...)
					if mensagem != "Interceptador Antes: BD" {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
			},
			sqloggerEsperado: bd.NovoSQLogger(simulador.Tx{}, net.ParseIP("192.168.1.1")),
		},
		{
			descrição:      "deve detectar quando a configuração não foi inicializada",
			endereçoRemoto: net.ParseIP("192.168.1.1"),
			logger: simulador.Logger{
				SimulaCrit: func(m ...interface{}) {
					mensagem := fmt.Sprint(m...)
					if mensagem != "Não existe configuração definida para iniciar a conexão com o banco de dados" {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
				SimulaDebug: func(m ...interface{}) {
					mensagem := fmt.Sprint(m...)
					if mensagem != "Interceptador Antes: BD" {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
			},
			códigoHTTPEsperado: http.StatusInternalServerError,
		},
		{
			descrição: "deve ignorar quando já existe uma conexão aberta",
			conexão: simulador.BD{
				SimulaBegin: func() (bd.Tx, error) {
					return simulador.Tx{}, nil
				},
			},
			configuração: func() *config.Configuração {
				configuração := new(config.Configuração)
				configuração.BancoDados.Endereço = "127.0.0.1"
				configuração.BancoDados.Nome = "teste"
				configuração.BancoDados.Usuário = "usuario"
				configuração.BancoDados.Senha = "senha"
				configuração.BancoDados.TempoEsgotadoConexão = 2 * time.Second
				configuração.BancoDados.TempoEsgotadoComando = 10 * time.Second
				configuração.BancoDados.TempoEsgotadoTransação = 1 * time.Second
				configuração.BancoDados.MáximoNúmeroConexõesInativas = 16
				configuração.BancoDados.MáximoNúmeroConexõesAbertas = 32
				return configuração
			}(),
			endereçoRemoto: net.ParseIP("192.168.1.1"),
			logger: simulador.Logger{
				SimulaDebug: func(m ...interface{}) {
					mensagem := fmt.Sprint(m...)
					if mensagem != "Interceptador Antes: BD" {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
			},
			sqloggerEsperado: bd.NovoSQLogger(simulador.Tx{}, net.ParseIP("192.168.1.1")),
		},
		{
			descrição: "deve detectar um erro ao iniciar uma conexão",
			iniciarConexão: func(parâmetrosConexão db.ConnParams, txTempoEsgotado time.Duration) error {
				return errors.Errorf("erro de conexão")
			},
			configuração: func() *config.Configuração {
				configuração := new(config.Configuração)
				configuração.BancoDados.Endereço = "127.0.0.1"
				configuração.BancoDados.Nome = "teste"
				configuração.BancoDados.Usuário = "usuario"
				configuração.BancoDados.Senha = "senha"
				configuração.BancoDados.TempoEsgotadoConexão = 2 * time.Second
				configuração.BancoDados.TempoEsgotadoComando = 10 * time.Second
				configuração.BancoDados.TempoEsgotadoTransação = 1 * time.Second
				configuração.BancoDados.MáximoNúmeroConexõesInativas = 16
				configuração.BancoDados.MáximoNúmeroConexõesAbertas = 32
				return configuração
			}(),
			endereçoRemoto: net.ParseIP("192.168.1.1"),
			logger: simulador.Logger{
				SimulaCritf: func(m string, a ...interface{}) {
					mensagem := fmt.Sprintf(m, a...)
					if !regexp.MustCompile(`Erro ao conectar o banco de dados. Detalhes: .*:[0-9]+: erro de conexão`).MatchString(mensagem) {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
				SimulaDebug: func(m ...interface{}) {
					mensagem := fmt.Sprint(m...)
					if mensagem != "Interceptador Antes: BD" {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
			},
			códigoHTTPEsperado: http.StatusInternalServerError,
		},
		{
			descrição: "deve detectar um erro ao iniciar uma transação",
			iniciarConexão: func(parâmetrosConexão db.ConnParams, txTempoEsgotado time.Duration) error {
				bd.Conexão = simulador.BD{
					SimulaBegin: func() (bd.Tx, error) {
						return nil, fmt.Errorf("erro na transação")
					},
				}

				return nil
			},
			configuração: func() *config.Configuração {
				configuração := new(config.Configuração)
				configuração.BancoDados.Endereço = "127.0.0.1"
				configuração.BancoDados.Nome = "teste"
				configuração.BancoDados.Usuário = "usuario"
				configuração.BancoDados.Senha = "senha"
				configuração.BancoDados.TempoEsgotadoConexão = 2 * time.Second
				configuração.BancoDados.TempoEsgotadoComando = 10 * time.Second
				configuração.BancoDados.TempoEsgotadoTransação = 1 * time.Second
				configuração.BancoDados.MáximoNúmeroConexõesInativas = 16
				configuração.BancoDados.MáximoNúmeroConexõesAbertas = 32
				return configuração
			}(),
			endereçoRemoto: net.ParseIP("192.168.1.1"),
			logger: simulador.Logger{
				SimulaErrorf: func(m string, a ...interface{}) {
					mensagem := fmt.Sprintf(m, a...)
					if !regexp.MustCompile(`Erro ao iniciar uma transação no banco de dados. Detalhes: .*:[0-9]+: erro na transação`).MatchString(mensagem) {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
				SimulaDebug: func(m ...interface{}) {
					mensagem := fmt.Sprint(m...)
					if mensagem != "Interceptador Antes: BD" {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
			},
			códigoHTTPEsperado: http.StatusInternalServerError,
		},
	}

	configuraçãoOriginal := config.Atual()
	defer func() {
		config.AtualizarConfiguração(configuraçãoOriginal)
	}()

	iniciarConexãoOriginal := bd.IniciarConexão
	defer func() {
		bd.IniciarConexão = iniciarConexãoOriginal
	}()

	conexãoOriginal := bd.Conexão
	defer func() {
		bd.Conexão = conexãoOriginal
	}()

	for i, cenário := range cenários {
		config.AtualizarConfiguração(cenário.configuração)
		bd.IniciarConexão = cenário.iniciarConexão
		bd.Conexão = cenário.conexão

		requisição, err := http.NewRequest("GET", "/teste", nil)
		if err != nil {
			t.Fatal(err)
		}

		handler := &bdSimulado{}
		handler.SimulaRequisição = requisição
		handler.DefineEndereçoRemoto(cenário.endereçoRemoto)
		handler.DefineLogger(cenário.logger)

		bd := interceptador.NovoBD(handler)
		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)

		verificadorResultado.DefinirEsperado(cenário.códigoHTTPEsperado, nil)
		if err := verificadorResultado.VerificaResultado(bd.Before(), nil); err != nil {
			t.Error(err)
		}

		verificadorResultado.DefinirEsperado(cenário.sqloggerEsperado, nil)
		if err := verificadorResultado.VerificaResultado(handler.Tx(), nil); err != nil {
			t.Error(err)
		}
	}
}

func TestBD_After(t *testing.T) {
	cenários := []struct {
		descrição          string
		códigoHTTP         int
		conexão            bd.BD
		logger             log.Logger
		códigoHTTPEsperado int
	}{
		{
			descrição:  "deve detectar uma transação não inicializada",
			códigoHTTP: http.StatusNoContent,
			conexão: simulador.BD{
				SimulaBegin: func() (bd.Tx, error) {
					return nil, fmt.Errorf("erro na transação")
				},
			},
			logger: simulador.Logger{
				SimulaErrorf: func() func(m string, a ...interface{}) {
					i := 0
					return func(m string, a ...interface{}) {
						mensagem := fmt.Sprintf(m, a...)

						if i == 0 {
							if !regexp.MustCompile(`Erro ao iniciar uma transação no banco de dados. Detalhes: .*:[0-9]+: erro na transação`).MatchString(mensagem) {
								t.Errorf("mensagem inesperada: %s", mensagem)
							}
						} else {
							if mensagem != "Transação não inicializada detectada" {
								t.Errorf("mensagem inesperada: %s", mensagem)
							}
						}

						i++
					}
				}(),
				SimulaDebug: func() func(m ...interface{}) {
					i := 0

					return func(m ...interface{}) {
						mensagem := fmt.Sprint(m...)

						if i == 0 {
							if mensagem != "Interceptador Antes: BD" {
								t.Errorf("mensagem inesperada: %s", mensagem)
							}
						} else {
							if mensagem != "Interceptador Depois: BD" {
								t.Errorf("mensagem inesperada: %s", mensagem)
							}
						}

						i++
					}
				}(),
			},
			códigoHTTPEsperado: http.StatusNoContent,
		},
		{
			descrição:  "deve confirmar uma transação corretamente",
			códigoHTTP: http.StatusNoContent,
			conexão: simulador.BD{
				SimulaBegin: func() (bd.Tx, error) {
					return simulador.Tx{
						SimulaCommit: func() error {
							return nil
						},
					}, nil
				},
			},
			logger: simulador.Logger{
				SimulaDebug: func() func(m ...interface{}) {
					i := 0

					return func(m ...interface{}) {
						mensagem := fmt.Sprint(m...)

						if i == 0 {
							if mensagem != "Interceptador Antes: BD" {
								t.Errorf("mensagem inesperada: %s", mensagem)
							}
						} else {
							if mensagem != "Interceptador Depois: BD" {
								t.Errorf("mensagem inesperada: %s", mensagem)
							}
						}

						i++
					}
				}(),
			},
			códigoHTTPEsperado: http.StatusNoContent,
		},
		{
			descrição:  "deve detectar um erro ao confirmar uma transação",
			códigoHTTP: http.StatusNoContent,
			conexão: simulador.BD{
				SimulaBegin: func() (bd.Tx, error) {
					return simulador.Tx{
						SimulaCommit: func() error {
							return fmt.Errorf("erro ao confirmar")
						},
					}, nil
				},
			},
			logger: simulador.Logger{
				SimulaErrorf: func(m string, a ...interface{}) {
					mensagem := fmt.Sprintf(m, a...)
					if !regexp.MustCompile(`Erro ao confirmar uma transação. Detalhes: .*:[0-9]+: erro ao confirmar`).MatchString(mensagem) {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
				SimulaDebug: func() func(m ...interface{}) {
					i := 0

					return func(m ...interface{}) {
						mensagem := fmt.Sprint(m...)

						if i == 0 {
							if mensagem != "Interceptador Antes: BD" {
								t.Errorf("mensagem inesperada: %s", mensagem)
							}
						} else {
							if mensagem != "Interceptador Depois: BD" {
								t.Errorf("mensagem inesperada: %s", mensagem)
							}
						}

						i++
					}
				}(),
			},
			códigoHTTPEsperado: http.StatusInternalServerError,
		},
		{
			descrição:  "deve desfazer uma transação corretamente",
			códigoHTTP: http.StatusBadRequest,
			conexão: simulador.BD{
				SimulaBegin: func() (bd.Tx, error) {
					return simulador.Tx{
						SimulaRollback: func() error {
							return nil
						},
					}, nil
				},
			},
			logger: simulador.Logger{
				SimulaDebug: func() func(m ...interface{}) {
					i := 0

					return func(m ...interface{}) {
						mensagem := fmt.Sprint(m...)

						if i == 0 {
							if mensagem != "Interceptador Antes: BD" {
								t.Errorf("mensagem inesperada: %s", mensagem)
							}
						} else {
							if mensagem != "Interceptador Depois: BD" {
								t.Errorf("mensagem inesperada: %s", mensagem)
							}
						}

						i++
					}
				}(),
			},
			códigoHTTPEsperado: http.StatusBadRequest,
		},
		{
			descrição:  "deve detectar um erro ao desfazer uma transação",
			códigoHTTP: http.StatusBadRequest,
			conexão: simulador.BD{
				SimulaBegin: func() (bd.Tx, error) {
					return simulador.Tx{
						SimulaRollback: func() error {
							return fmt.Errorf("erro ao desfazer")
						},
					}, nil
				},
			},
			logger: simulador.Logger{
				SimulaErrorf: func(m string, a ...interface{}) {
					mensagem := fmt.Sprintf(m, a...)
					if !regexp.MustCompile(`Erro ao desfazer uma transação. Detalhes: .*:[0-9]+: erro ao desfazer`).MatchString(mensagem) {
						t.Errorf("mensagem inesperada: %s", mensagem)
					}
				},
				SimulaDebug: func() func(m ...interface{}) {
					i := 0

					return func(m ...interface{}) {
						mensagem := fmt.Sprint(m...)

						if i == 0 {
							if mensagem != "Interceptador Antes: BD" {
								t.Errorf("mensagem inesperada: %s", mensagem)
							}
						} else {
							if mensagem != "Interceptador Depois: BD" {
								t.Errorf("mensagem inesperada: %s", mensagem)
							}
						}

						i++
					}
				}(),
			},
			códigoHTTPEsperado: http.StatusBadRequest,
		},
	}

	conexãoOriginal := bd.Conexão
	defer func() {
		bd.Conexão = conexãoOriginal
	}()

	for i, cenário := range cenários {
		bd.Conexão = cenário.conexão

		requisição, err := http.NewRequest("GET", "/teste", nil)
		if err != nil {
			t.Fatal(err)
		}

		handler := &bdSimulado{}
		handler.SimulaRequisição = requisição
		handler.DefineLogger(cenário.logger)

		bd := interceptador.NovoBD(handler)
		bd.Before()

		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)

		verificadorResultado.DefinirEsperado(cenário.códigoHTTPEsperado, nil)
		if err := verificadorResultado.VerificaResultado(bd.After(cenário.códigoHTTP), nil); err != nil {
			t.Error(err)
		}
	}
}

type bdSimulado struct {
	interceptador.EndereçoRemotoCompatível
	interceptador.LogCompatível
	interceptador.BDCompatível
	simulador.Handler
}
