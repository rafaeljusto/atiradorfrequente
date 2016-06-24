package interceptador

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/registrobr/gostk/log"
)

func init() {
	// TODO(rafaeljusto): https://nishanths.svbtle.com/do-not-seed-the-global-random
	rand.Seed(time.Now().UTC().UnixNano())
}

type logger interface {
	EndereçoRemoto() net.IP
	DefineLogger(log.Logger)
	Logger() log.Logger
	Req() *http.Request
}

type Log struct {
	handler logger
}

func NovoLog(h logger) *Log {
	return &Log{handler: h}
}

func (l Log) Before() int {
	idRequisição := rand.Int31n(99999)
	identificador := fmt.Sprintf("%s %05d", l.handler.EndereçoRemoto(), idRequisição)
	l.handler.DefineLogger(log.NewLogger(identificador))

	requisição := l.handler.Req()
	l.handler.Logger().Infof("Requisicao %s %s", requisição.Method, requisição.RequestURI)
	return 0
}

func (l Log) After(status int) int {
	requisição := l.handler.Req()
	l.handler.Logger().Infof("Resposta %s %s %d %s", requisição.Method, requisição.RequestURI, status, http.StatusText(status))
	return status
}

type LogCompatível struct {
	logger log.Logger
}

func (l *LogCompatível) DefineLogger(logger log.Logger) {
	l.logger = logger
}

func (l LogCompatível) Logger() log.Logger {
	return l.logger
}
