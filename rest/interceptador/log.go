package interceptador

import (
	"fmt"
	"net"
	"net/http"

	"github.com/rafaeljusto/atiradorfrequente/núcleo/randômico"
	"github.com/registrobr/gostk/log"
)

type logger interface {
	EndereçoRemoto() net.IP
	DefineLogger(log.Logger)
	Logger() log.Logger
	Req() *http.Request
}

// Log disponibiliza ao handler uma estrutura de log contextualizada para a
// requisição, de forma a permitir identificar as mensagens de log de um
// usuário.
type Log struct {
	handler logger
}

// NovoLog cria um novo interceptador Log.
func NovoLog(h logger) *Log {
	return &Log{handler: h}
}

// Before inicializa uma estrutura de log contextualizada, utilizando como
// identificador o endereço IP remoto e um número aleatório. Existe uma pequena
// chance de colisão de identificadores caso gere um número aleatório repetido
// para o mesmo endereço IP remoto. Ao inicializar, adiciona informações da
// requisição no log.
func (l Log) Before() int {
	idRequisição := randômico.FonteRandômica.Int31n(99999)
	identificador := fmt.Sprintf("%s %05d", l.handler.EndereçoRemoto(), idRequisição)
	l.handler.DefineLogger(log.NewLogger(identificador))

	requisição := l.handler.Req()
	l.handler.Logger().Infof("Requisicao %s %s", requisição.Method, requisição.RequestURI)
	return 0
}

// After adiciona informações da resposta no log.
func (l Log) After(status int) int {
	requisição := l.handler.Req()
	l.handler.Logger().Infof("Resposta %s %s %d %s", requisição.Method, requisição.RequestURI, status, http.StatusText(status))
	return status
}

// LogCompatível implementa os métodos que serão utilizados pelo handler para
// acessar o log criado por este interceptador.
type LogCompatível struct {
	logger log.Logger
}

// DefineLogger defile o logger que será utilizado pelo handler.
func (l *LogCompatível) DefineLogger(logger log.Logger) {
	l.logger = logger
}

// Logger obtém o logger que será utilizado pelo handler.
func (l LogCompatível) Logger() log.Logger {
	// TODO(rafaeljusto): Se o logger estiver indefinido devemos ter um plano B?
	return l.logger
}
