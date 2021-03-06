package servidor_test

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/bd"
	"github.com/rafaeljusto/atiradorfrequente/rest/config"
	"github.com/rafaeljusto/atiradorfrequente/rest/handler"
	"github.com/rafaeljusto/atiradorfrequente/rest/servidor"
	"github.com/rafaeljusto/atiradorfrequente/testes"
	"github.com/rafaeljusto/atiradorfrequente/testes/simulador"
	"github.com/registrobr/gostk/db"
	"github.com/registrobr/gostk/errors"
	"github.com/registrobr/gostk/log"
	"github.com/trajber/handy"
)

func TestIniciar(t *testing.T) {
	arquivoCertificado, err := ioutil.TempFile("", "atirador-frequente-")
	if err != nil {
		t.Fatalf("Erro ao criar o arquivo do certificado. Detalhes: %s", err)
	}
	defer arquivoCertificado.Close()
	arquivoCertificado.WriteString(certificado)

	arquivoChave, err := ioutil.TempFile("", "atirador-frequente-")
	if err != nil {
		t.Fatalf("Erro ao criar o arquivo da chave. Detalhes: %s", err)
	}
	defer arquivoChave.Close()
	arquivoChave.WriteString(chave)

	loggerOriginal := log.LocalLogger
	defer func() {
		log.LocalLogger = loggerOriginal
	}()

	var servidorLog simulador.ServidorLog
	syslog, err := servidorLog.Executar("localhost:0")
	if err != nil {
		t.Fatalf("Erro ao inicializar o servidor de log. Detalhes: %s", err)
	}
	defer syslog.Close()

	var endereçoServidor string

	cenários := []struct {
		descrição          string
		escuta             net.Listener
		configuração       config.Configuração
		conexãoBD          func(parâmetrosConexão db.ConnParams, txTempoEsgotado time.Duration) error
		fecharConexãoLog   func() error
		inicializar        func()
		finalizar          func()
		erroEsperado       error
		mensagensEsperadas *regexp.Regexp
	}{
		{
			descrição: "deve iniciar corretamente o servidor",
			escuta: func() net.Listener {
				escuta, err := net.Listen("tcp", "localhost:0")
				if err != nil {
					t.Fatalf("Erro ao inicializar o servidor. Detalhes: %s", err)
				}
				endereçoServidor = escuta.Addr().String()
				return escuta
			}(),
			configuração: func() config.Configuração {
				var c config.Configuração
				c.Servidor.Endereço = endereçoServidor
				c.Servidor.TLS.Habilitado = true
				c.Servidor.TLS.ArquivoCertificado = arquivoCertificado.Name()
				c.Servidor.TLS.ArquivoChave = arquivoChave.Name()
				c.Syslog.Endereço = syslog.Addr().String()
				c.Syslog.TempoEsgotadoConexão = 1 * time.Second
				return c
			}(),
			conexãoBD: func(parâmetrosConexão db.ConnParams, txTempoEsgotado time.Duration) error {
				bd.Conexão = simulador.BD{
					SimulaClose: func() error {
						return nil
					},
				}
				return nil
			},
			fecharConexãoLog: log.Close,
			erroEsperado:     errors.Errorf("accept tcp %s: use of closed network connection", endereçoServidor),
			mensagensEsperadas: regexp.MustCompile(`^.*Inicializando conexão com o servidor de log
.*Inicializando conexão com o banco de dados
.*Inicializando servidor
.*Erro ao iniciar o servidor\. Detalhes: .*use of closed network connection
$`),
		},
		{
			descrição: "deve detectar um erro ao inicializar a conexão com o servidor de log",
			escuta: func() net.Listener {
				escuta, err := net.Listen("tcp", "localhost:0")
				if err != nil {
					t.Fatalf("Erro ao inicializar o servidor. Detalhes: %s", err)
				}
				endereçoServidor = escuta.Addr().String()
				return escuta
			}(),
			configuração: func() config.Configuração {
				var c config.Configuração
				c.Servidor.Endereço = endereçoServidor
				c.Servidor.TLS.Habilitado = true
				c.Servidor.TLS.ArquivoCertificado = arquivoCertificado.Name()
				c.Servidor.TLS.ArquivoChave = arquivoChave.Name()
				c.Syslog.Endereço = "192.0.2.1:1234"
				c.Syslog.TempoEsgotadoConexão = 100 * time.Millisecond
				return c
			}(),
			conexãoBD: func(parâmetrosConexão db.ConnParams, txTempoEsgotado time.Duration) error {
				bd.Conexão = simulador.BD{
					SimulaClose: func() error {
						return nil
					},
				}
				return nil
			},
			fecharConexãoLog: log.Close,
			erroEsperado:     log.ErrDialTimeout,
			mensagensEsperadas: regexp.MustCompile(`^.*Inicializando conexão com o servidor de log
.*Erro ao conectar servidor de log. Detalhes: .*dial timeout
$`),
		},
		{
			descrição: "deve detectar um erro ao encerrar a conexão do log",
			escuta: func() net.Listener {
				escuta, err := net.Listen("tcp", "localhost:0")
				if err != nil {
					t.Fatalf("Erro ao inicializar o servidor. Detalhes: %s", err)
				}
				endereçoServidor = escuta.Addr().String()
				return escuta
			}(),
			configuração: func() config.Configuração {
				var c config.Configuração
				c.Servidor.Endereço = endereçoServidor
				c.Servidor.TLS.Habilitado = true
				c.Servidor.TLS.ArquivoCertificado = arquivoCertificado.Name()
				c.Servidor.TLS.ArquivoChave = arquivoChave.Name()
				c.Syslog.Endereço = syslog.Addr().String()
				c.Syslog.TempoEsgotadoConexão = 1 * time.Second
				return c
			}(),
			conexãoBD: func(parâmetrosConexão db.ConnParams, txTempoEsgotado time.Duration) error {
				bd.Conexão = simulador.BD{
					SimulaClose: func() error {
						return nil
					},
				}
				return nil
			},
			fecharConexãoLog: func() error {
				return fmt.Errorf("erro ao encerrar a conexão")
			},
			erroEsperado: errors.Errorf("accept tcp %s: use of closed network connection", endereçoServidor),
			mensagensEsperadas: regexp.MustCompile(`^.*Inicializando conexão com o servidor de log
.*Inicializando conexão com o banco de dados
.*Inicializando servidor
.*Erro ao iniciar o servidor\. Detalhes: .*use of closed network connection
.*Erro ao fechar a conexão do log. Detalhes: .*erro ao encerrar a conexão
$`),
		},
		{
			descrição: "deve detectar um erro ao iniciar a conexão com o banco de dados",
			escuta: func() net.Listener {
				escuta, err := net.Listen("tcp", "localhost:0")
				if err != nil {
					t.Fatalf("Erro ao inicializar o servidor. Detalhes: %s", err)
				}
				endereçoServidor = escuta.Addr().String()
				return escuta
			}(),
			configuração: func() config.Configuração {
				var c config.Configuração
				c.Servidor.Endereço = endereçoServidor
				c.Servidor.TLS.Habilitado = true
				c.Servidor.TLS.ArquivoCertificado = arquivoCertificado.Name()
				c.Servidor.TLS.ArquivoChave = arquivoChave.Name()
				c.Syslog.Endereço = syslog.Addr().String()
				c.Syslog.TempoEsgotadoConexão = 1 * time.Second
				return c
			}(),
			conexãoBD: func(parâmetrosConexão db.ConnParams, txTempoEsgotado time.Duration) error {
				return errors.Errorf("erro de conexão")
			},
			fecharConexãoLog: log.Close,
			erroEsperado:     errors.Errorf("accept tcp %s: use of closed network connection", endereçoServidor),
			mensagensEsperadas: regexp.MustCompile(`^.*Inicializando conexão com o servidor de log
.*Inicializando conexão com o banco de dados
.*Erro ao conectar o banco de dados. Detalhes: .*erro de conexão
.*Inicializando servidor
.*Erro ao iniciar o servidor\. Detalhes: .*use of closed network connection
$`),
		},
		{
			descrição: "deve detectar um erro ao encerrar a conexão do banco de dados",
			escuta: func() net.Listener {
				escuta, err := net.Listen("tcp", "localhost:0")
				if err != nil {
					t.Fatalf("Erro ao inicializar o servidor. Detalhes: %s", err)
				}
				endereçoServidor = escuta.Addr().String()
				return escuta
			}(),
			configuração: func() config.Configuração {
				var c config.Configuração
				c.Servidor.Endereço = endereçoServidor
				c.Servidor.TLS.Habilitado = true
				c.Servidor.TLS.ArquivoCertificado = arquivoCertificado.Name()
				c.Servidor.TLS.ArquivoChave = arquivoChave.Name()
				c.Syslog.Endereço = syslog.Addr().String()
				c.Syslog.TempoEsgotadoConexão = 1 * time.Second
				return c
			}(),
			conexãoBD: func(parâmetrosConexão db.ConnParams, txTempoEsgotado time.Duration) error {
				bd.Conexão = simulador.BD{
					SimulaClose: func() error {
						return errors.Errorf("erro na conexão com o banco de dados")
					},
				}
				return nil
			},
			fecharConexãoLog: log.Close,
			erroEsperado:     errors.Errorf("accept tcp %s: use of closed network connection", endereçoServidor),
			mensagensEsperadas: regexp.MustCompile(`^.*Inicializando conexão com o servidor de log
.*Inicializando conexão com o banco de dados
.*Inicializando servidor
.*Erro ao iniciar o servidor\. Detalhes: .*use of closed network connection
.*Erro ao fechar a conexão do banco de dados. Detalhes: .*erro na conexão com o banco de dados
$`),
		},
		{
			descrição: "deve detectar quando o sistema entra em pânico",
			escuta: func() net.Listener {
				escuta, err := net.Listen("tcp", "localhost:0")
				if err != nil {
					t.Fatalf("Erro ao inicializar o servidor. Detalhes: %s", err)
				}
				endereçoServidor = escuta.Addr().String()
				return escuta
			}(),
			configuração: func() config.Configuração {
				var c config.Configuração
				c.Servidor.Endereço = endereçoServidor
				c.Servidor.TLS.Habilitado = true
				c.Servidor.TLS.ArquivoCertificado = arquivoCertificado.Name()
				c.Servidor.TLS.ArquivoChave = arquivoChave.Name()
				c.Syslog.Endereço = syslog.Addr().String()
				c.Syslog.TempoEsgotadoConexão = 1 * time.Second
				return c
			}(),
			conexãoBD: func(parâmetrosConexão db.ConnParams, txTempoEsgotado time.Duration) error {
				bd.Conexão = simulador.BD{
					SimulaClose: func() error {
						return nil
					},
				}
				return nil
			},
			fecharConexãoLog: log.Close,
			inicializar: func() {
				handler.Rotas["/teste"] = handy.Constructor(func() handy.Handler {
					return &simulador.Handler{
						SimulaGet: func() int {
							panic("pânico no sistema")
						},
					}
				})
			},
			finalizar: func() {
				delete(handler.Rotas, "/teste")
			},
			erroEsperado: errors.Errorf("accept tcp %s: use of closed network connection", endereçoServidor),
			mensagensEsperadas: regexp.MustCompile(`^.*Inicializando conexão com o servidor de log
.*Inicializando conexão com o banco de dados
.*Inicializando servidor
.*Erro grave detectado. Detalhes: pânico no sistema
(.|\n)*
.*Erro ao iniciar o servidor\. Detalhes: .*use of closed network connection
$`),
		},
		{
			descrição: "deve detectar certificados HTTPS inválidos",
			escuta: func() net.Listener {
				escuta, err := net.Listen("tcp", "localhost:0")
				if err != nil {
					t.Fatalf("Erro ao inicializar o servidor. Detalhes: %s", err)
				}
				endereçoServidor = escuta.Addr().String()
				return escuta
			}(),
			configuração: func() config.Configuração {
				var c config.Configuração
				c.Servidor.Endereço = endereçoServidor
				c.Servidor.TLS.Habilitado = true
				c.Servidor.TLS.ArquivoCertificado = "/tmp/atiradorfrequente/nao-existo.crt"
				c.Servidor.TLS.ArquivoChave = "/tmp/atiradorfrequente/nao-existo.key"
				c.Syslog.Endereço = syslog.Addr().String()
				c.Syslog.TempoEsgotadoConexão = 1 * time.Second
				return c
			}(),
			conexãoBD: func(parâmetrosConexão db.ConnParams, txTempoEsgotado time.Duration) error {
				bd.Conexão = simulador.BD{
					SimulaClose: func() error {
						return nil
					},
				}
				return nil
			},
			fecharConexãoLog: log.Close,
			erroEsperado: &os.PathError{
				Op:   "open",
				Path: "/tmp/atiradorfrequente/nao-existo.crt",
				Err:  fmt.Errorf("no such file or directory"),
			},
			mensagensEsperadas: regexp.MustCompile(`^.*Inicializando conexão com o servidor de log
.*Inicializando conexão com o banco de dados
.*Inicializando servidor
.*Erro ao iniciar o servidor\. Detalhes: .*open /tmp/atiradorfrequente/nao-existo.crt: no such file or directory
$`),
		},
	}

	configuraçãoOriginal := config.Atual()
	defer func() {
		config.AtualizarConfiguração(configuraçãoOriginal)
	}()

	conexãoBDOriginal := bd.IniciarConexão
	defer func() {
		bd.IniciarConexão = conexãoBDOriginal
	}()

	fecharConexãoLogOriginal := log.Close
	defer func() {
		log.Close = fecharConexãoLogOriginal
	}()

	for i, cenário := range cenários {
		servidorLog.Limpar()
		config.AtualizarConfiguração(&cenário.configuração)

		bd.Conexão = nil
		bd.IniciarConexão = cenário.conexãoBD

		log.Close = cenário.fecharConexãoLog

		if cenário.inicializar != nil {
			cenário.inicializar()
		}

		canalRequisição := make(chan bool, 1)
		go func() {
			// aguarda para o servidor ser iniciado
			time.Sleep(10 * time.Millisecond)

			transporte := http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			}

			cliente := http.Client{
				Transport: &transporte,
				Timeout:   100 * time.Millisecond,
			}

			cliente.Get(fmt.Sprintf("https://%s/teste", cenário.escuta.Addr()))

			if cenário.escuta != nil {
				// fecha o escutador para desbloquear a função Iniciar
				cenário.escuta.Close()
			}

			close(canalRequisição)
		}()

		err := servidor.Iniciar(cenário.escuta)

		verificadorResultado := testes.NovoVerificadorResultados(cenário.descrição, i)
		verificadorResultado.DefinirEsperado(nil, cenário.erroEsperado)
		if err = verificadorResultado.VerificaResultado(nil, err); err != nil {
			t.Error(err)
		}

		// aguarda a rotina da requisição encerrar
		<-canalRequisição

		// aguarda as últimas mensagens serem escritas no log
		time.Sleep(10 * time.Millisecond)

		if !cenário.mensagensEsperadas.MatchString(servidorLog.Mensagens()) {
			t.Errorf("Item %d, “%s”: mensagem inesperada. Detalhes: %s",
				i, cenário.descrição, servidorLog.Mensagens())
		}

		if cenário.finalizar != nil {
			cenário.finalizar()
		}
	}
}

const certificado = `-----BEGIN CERTIFICATE-----
MIIC+TCCAeGgAwIBAgIQQeiyKgAaawHjIOaFGmYFujANBgkqhkiG9w0BAQsFADAS
MRAwDgYDVQQKEwdBY21lIENvMB4XDTE2MDgwMjExMDgwNFoXDTE3MDgwMjExMDgw
NFowEjEQMA4GA1UEChMHQWNtZSBDbzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCC
AQoCggEBALu/3579jVlIkO+Z2XHs8pSNzlDbDurcHmT9tLzPzub7ufaHUij+MR1I
W62MFTkSZNlycPjD4KF1ChhKhRZojYFZXGOCPv0gDKgkZ0nTHRCQvS4uaY+mEvE+
vqEb+S4wzdO6ZPCPHeqCVTkJHRwKqJfydFPsQGYzHPQEOTw+Q0IQPh1Ba7VRVG+o
RN4yzvz7QlVFYPo2OTJDl7TPl/Zuw6+cvjMWLao3AieWvaJSPCvDxw9P7ugz9KzU
DvzZybBwYzvM1RcjEaN2rkHzEQxOtGgNLQ2ZEosu8OITgWx5lgVXjmxrnClLVzPc
7LhiP98j9d149RFkE830Y4vwjG1hx5sCAwEAAaNLMEkwDgYDVR0PAQH/BAQDAgWg
MBMGA1UdJQQMMAoGCCsGAQUFBwMBMAwGA1UdEwEB/wQCMAAwFAYDVR0RBA0wC4IJ
bG9jYWxob3N0MA0GCSqGSIb3DQEBCwUAA4IBAQBa0z57uKQRZupao6TCt2vLarLj
UUJtyfjAw4VMxEBjLa3CkCB2EruLZPdaWgDee0Qd1bmo1sycyQKbYAmIY62EgX7o
6xIM8guAs2/kw0AkvjybPHlxK1UzcQDIryKNr/4vg8jnQAMayZjz9VDEzj4FNrgL
SnN0YPjRi3AO5XEH5PO+cT2qdwONnyJ1+cqptA9Vw2+4N+iUi+QkK2JuyrNksp6E
qXykACDeeyELMCWocOi0qwGgQyT3CUgbqrNk6fjjq7D7dncrJSiV63YSqfkyxne/
7kswiN8P15QyciSDVasMg/Ce5OHRVJK3KJF4/oVzcrspKmK/Z1jxdSDUQFH0
-----END CERTIFICATE-----`

const chave = `-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAu7/fnv2NWUiQ75nZcezylI3OUNsO6tweZP20vM/O5vu59odS
KP4xHUhbrYwVORJk2XJw+MPgoXUKGEqFFmiNgVlcY4I+/SAMqCRnSdMdEJC9Li5p
j6YS8T6+oRv5LjDN07pk8I8d6oJVOQkdHAqol/J0U+xAZjMc9AQ5PD5DQhA+HUFr
tVFUb6hE3jLO/PtCVUVg+jY5MkOXtM+X9m7Dr5y+MxYtqjcCJ5a9olI8K8PHD0/u
6DP0rNQO/NnJsHBjO8zVFyMRo3auQfMRDE60aA0tDZkSiy7w4hOBbHmWBVeObGuc
KUtXM9zsuGI/3yP13Xj1EWQTzfRji/CMbWHHmwIDAQABAoIBAD5wIgs64V4W7vRv
4surdEUJH9rt7vkWORl28jt0lKdmgcLF4AH3/xdw7+Q4WPqA7n8OOxnP8o1fYfsQ
FVCNdrnUoRAKya3ekbb9XEhF6D2RFQkdsEdwgI4wQq9LoUPGQT0vmNATmGxb6cGt
ETw5IzZdEGi0gfo6918DZJFvV2jJcZb0ovuqLUI9ocJLvDZgRhgF9H8y/T+t6v34
T0/GuE0U2nblIDzFcDJRDQ/ck2XCXURMT7z6ET0oUYAyK/WckfQ8NZJqgwwTFmNs
ZjZ7P1sI3vz1Co1vLebo4C05a1SgjvrsJALwRy391V6+v3E05SS0VVoer7WJ8BMt
aK+vlvECgYEA8OeVxAkxjO/UG60oUH4+98HXKVoO3DB7ZKqzgNIxSxMVIlmPU765
UFCzIgIZl/7GREl/WDLGDiez5OlpK8qgDguNT/yVpyoNUBz9IdrP+U6ZgKTm05uw
ewua74DkT48AjN/UUr54pXUvFGUapKtGxkWZ5qB+a/x/SB8ZSgcgEakCgYEAx4Oc
8y7CWp9bFcbOWX6TvT8tX9Fl85SXSy2vKwy5gYpkcnjgEiX3KJSpzWTWs6EZtlQ8
3Ha1fRlE38EPAdtsnayWXukmBhlyXA8yvCW+aydzQBdSKLr4SK+SStVaYA5L+Dek
0ZjpAo2I0WZ1kL/VvCX6cm2Di3YOQOqb/dmJ4aMCgYAbTtUuTLB+Pm132bAZN8Zh
hWqjeF743NIP/j2s26bU0Mvzgd16a8NL9Gnp7/0AutO0x/QUhmTnE98TktXmLejo
zqxtJb+9HEo4C6EyJkCvDRbfe1HjKOHfgNhGUAERd69jSLgjzQ2WC+uTT0au5e92
6Eri1syd5xhyj3vpZVdgSQKBgES3Q9tOA4qK0ChT7MZOHjxUAiC6Uk7uop02Atrk
6w9+xtHWZ/ZYNSQ477LaREhh+CUgJkYYbLHFfj9CkxSkqmg0BSZzTrFTGlwyr9q1
dTwavksYvSdiHhmKvuwfR51Fz0ySfaXi8H38mV7l1yAfslG3EudOaLwj0QzywP9R
aXfZAoGAdugSb2ppeb+4Z29gUeRwSTxjJpTAkjI0VO/yPfVPRrh0+TC5FYCivMYN
6zrHmC6GZFiJxA0NBySnTMG6kuCprOgoYdYgvlL8WU5TL8luARxE3MU6rNo/K4qM
xzSQ9bVHfA4fHL+um+w20K6xwMHMxnnif6xgpqbAIWVBo1JzlAc=
-----END RSA PRIVATE KEY-----`
