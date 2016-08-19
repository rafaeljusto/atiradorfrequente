package simulador_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"regexp"

	"log/syslog"

	"github.com/rafaeljusto/atiradorfrequente/testes/simulador"
	"github.com/registrobr/gostk/log"
)

func TestLogger(t *testing.T) {
	var logSimulado simulador.Logger
	var métodosSimulados []string

	estruturaBDSimulado := reflect.TypeOf(logSimulado)
	for i := 0; i < estruturaBDSimulado.NumField(); i++ {
		// trata somente funções como argumentos, ignorando atributos simples
		if !strings.HasPrefix(estruturaBDSimulado.Field(i).Type.String(), "func (") {
			continue
		}

		métodosSimulados = append(métodosSimulados, estruturaBDSimulado.Field(i).Name)
	}

	visitou := func(métodoSimulado string) {
		for i := len(métodosSimulados) - 1; i >= 0; i-- {
			if métodosSimulados[i] == métodoSimulado {
				métodosSimulados = append(métodosSimulados[:i], métodosSimulados[i+1:]...)
				break
			}
		}
	}

	logSimulado.SimulaEmerg = func(m ...interface{}) {
		visitou("SimulaEmerg")
	}

	logSimulado.SimulaEmergf = func(m string, a ...interface{}) {
		visitou("SimulaEmergf")
	}

	logSimulado.SimulaAlert = func(m ...interface{}) {
		visitou("SimulaAlert")
	}

	logSimulado.SimulaAlertf = func(m string, a ...interface{}) {
		visitou("SimulaAlertf")
	}

	logSimulado.SimulaCrit = func(m ...interface{}) {
		visitou("SimulaCrit")
	}

	logSimulado.SimulaCritf = func(m string, a ...interface{}) {
		visitou("SimulaCritf")
	}

	logSimulado.SimulaError = func(e error) {
		visitou("SimulaError")
	}

	logSimulado.SimulaErrorf = func(m string, a ...interface{}) {
		visitou("SimulaErrorf")
	}

	logSimulado.SimulaWarning = func(m ...interface{}) {
		visitou("SimulaWarning")
	}

	logSimulado.SimulaWarningf = func(m string, a ...interface{}) {
		visitou("SimulaWarningf")
	}

	logSimulado.SimulaNotice = func(m ...interface{}) {
		visitou("SimulaNotice")
	}

	logSimulado.SimulaNoticef = func(m string, a ...interface{}) {
		visitou("SimulaNoticef")
	}

	logSimulado.SimulaInfo = func(m ...interface{}) {
		visitou("SimulaInfo")
	}

	logSimulado.SimulaInfof = func(m string, a ...interface{}) {
		visitou("SimulaInfof")
	}

	logSimulado.SimulaDebug = func(m ...interface{}) {
		visitou("SimulaDebug")
	}

	logSimulado.SimulaDebugf = func(m string, a ...interface{}) {
		visitou("SimulaDebugf")
	}

	logSimulado.SimulaSetCaller = func(n int) {
		visitou("SimulaSetCaller")
	}

	logSimulado.Emerg("")
	logSimulado.Emergf("")
	logSimulado.Alert("")
	logSimulado.Alertf("")
	logSimulado.Crit("")
	logSimulado.Critf("")
	logSimulado.Error(fmt.Errorf(""))
	logSimulado.Errorf("")
	logSimulado.Warning("")
	logSimulado.Warningf("")
	logSimulado.Notice("")
	logSimulado.Noticef("")
	logSimulado.Info("")
	logSimulado.Infof("")
	logSimulado.Debug("")
	logSimulado.Debugf("")
	logSimulado.SetCaller(1)

	if len(métodosSimulados) > 0 {
		t.Errorf("métodos %#v não foram chamados", métodosSimulados)
	}
}

func TestServidorLog(t *testing.T) {
	var servidorLog simulador.ServidorLog
	escuta, err := servidorLog.Executar("localhost:0")
	if err != nil {
		t.Fatalf("Erro ao inicializar o servidor de log. Detalhes: %s", err)
	}
	defer escuta.Close()

	log.Debug("Teste1")

	if err = log.Dial("tcp", escuta.Addr().String(), "teste", 2*time.Second); err != nil {
		t.Fatalf("Erro ao conectar-se ao servidor de log. Detalhes: %s", err)
	}
	defer log.Close()

	log.Debug("Teste2\nTeste3")

	s, err := syslog.Dial("tcp", escuta.Addr().String(), syslog.LOG_INFO|syslog.LOG_LOCAL0, "teste")
	if err != nil {
		t.Fatalf("Erro ao conectar-se ao servidor de log. Detalhes: %s", err)
	}

	s.Debug("Teste4")

	// aguarda o servidor escrever a mensagem
	time.Sleep(10 * time.Millisecond)

	mensagensEsperadas := regexp.MustCompile(`.* Teste1
.* Teste2
.* Teste3
.* Teste4`)

	if !mensagensEsperadas.MatchString(servidorLog.Mensagens()) {
		t.Errorf("Mensagens inesperadas: %s", servidorLog.Mensagens())
	}

	servidorLog.Limpar()

	if servidorLog.Mensagens() != "" {
		t.Errorf("Mensagens inesperadas: %s", servidorLog.Mensagens())
	}
}

func TestServidorLog_endereçoInválido(t *testing.T) {
	var servidorLog simulador.ServidorLog
	if _, err := servidorLog.Executar("xxxxx"); err == nil {
		t.Fatal("Servidor de log não detectou um endereço inválido")
	}
}
