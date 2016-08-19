package simulador

import (
	"bytes"
	golog "log"
	"net"
	"strings"
	"sync"

	gostklog "github.com/registrobr/gostk/log"
)

// Logger simula uma estrutura de log.
type Logger struct {
	SimulaEmerg     func(m ...interface{})
	SimulaEmergf    func(m string, a ...interface{})
	SimulaAlert     func(m ...interface{})
	SimulaAlertf    func(m string, a ...interface{})
	SimulaCrit      func(m ...interface{})
	SimulaCritf     func(m string, a ...interface{})
	SimulaError     func(e error)
	SimulaErrorf    func(m string, a ...interface{})
	SimulaWarning   func(m ...interface{})
	SimulaWarningf  func(m string, a ...interface{})
	SimulaNotice    func(m ...interface{})
	SimulaNoticef   func(m string, a ...interface{})
	SimulaInfo      func(m ...interface{})
	SimulaInfof     func(m string, a ...interface{})
	SimulaDebug     func(m ...interface{})
	SimulaDebugf    func(m string, a ...interface{})
	SimulaSetCaller func(n int)
}

// Emerg escreve mensagens em nível de emergência.
func (l Logger) Emerg(m ...interface{}) {
	l.SimulaEmerg(m...)
}

// Emergf escreve mensagens formatadas em nível de emergência.
func (l Logger) Emergf(m string, a ...interface{}) {
	l.SimulaEmergf(m, a...)
}

// Alert escreve mensagens em nível de alerta.
func (l Logger) Alert(m ...interface{}) {
	l.SimulaAlert(m...)
}

// Alertf escreve mensagens formatadas em nível de alerta.
func (l Logger) Alertf(m string, a ...interface{}) {
	l.SimulaAlertf(m, a...)
}

// Crit escreve mensagens em nível crítico.
func (l Logger) Crit(m ...interface{}) {
	l.SimulaCrit(m...)
}

// Critf escreve mensagens formatadas em nível de alerta.
func (l Logger) Critf(m string, a ...interface{}) {
	l.SimulaCritf(m, a...)
}

// Error converte erros em mensagens.
func (l Logger) Error(e error) {
	l.SimulaError(e)
}

// Errorf escreve mensagens formatadas em nível de erro.
func (l Logger) Errorf(m string, a ...interface{}) {
	l.SimulaErrorf(m, a...)
}

// Warning escreve mensagens em nível de atenção.
func (l Logger) Warning(m ...interface{}) {
	l.SimulaWarning(m...)
}

// Warningf escreve mensagens formatadas em nível de atenção.
func (l Logger) Warningf(m string, a ...interface{}) {
	l.SimulaWarningf(m, a...)
}

// Notice escreve mensagens em nível de notícia.
func (l Logger) Notice(m ...interface{}) {
	l.SimulaNotice(m...)
}

// Noticef escreve mensagens formatadas em nível de notícia.
func (l Logger) Noticef(m string, a ...interface{}) {
	l.SimulaNoticef(m, a...)
}

// Info escreve mensagens em nível de informação.
func (l Logger) Info(m ...interface{}) {
	l.SimulaInfo(m...)
}

// Infof escreve mensagens formatadas em nível de informação.
func (l Logger) Infof(m string, a ...interface{}) {
	l.SimulaInfof(m, a...)
}

// Debug escreve mensagens em nível de desenvolvimento.
func (l Logger) Debug(m ...interface{}) {
	l.SimulaDebug(m...)
}

// Debugf escreve mensagens formatadas em nível de desenvolvimento.
func (l Logger) Debugf(m string, a ...interface{}) {
	l.SimulaDebugf(m, a...)
}

// SetCaller define quantos níveis a estrutura logger deve subir na pilha de
// chamadas de funções para obter o local de onde foi invocado o log. Este valor
// deve ser modificado somente em casos muito específicos.
func (l Logger) SetCaller(n int) {
	l.SimulaSetCaller(n)
}

// ServidorLog simula um servidor de log remoto escutando em uma porta TCP.
type ServidorLog struct {
	mensagens logMensagens
}

// Executar inicia o servidor de log retornando a escuta. Para encerrar o
// servidor basta fechar a escuta.
func (s *ServidorLog) Executar(endereço string) (net.Listener, error) {
	gostklog.LocalLogger = golog.New(&s.mensagens, "", golog.Lshortfile)

	syslog, err := net.Listen("tcp", endereço)
	if err != nil {
		return nil, err
	}

	go func(l net.Listener) {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}

			go func(conn net.Conn) {
				defer conn.Close()

				for {
					var buffer [1024]byte
					n, err := conn.Read(buffer[:])
					if err != nil {
						break
					}

					linhas := strings.Split(string(buffer[:n]), "\n")
					for _, linha := range linhas {
						linha = strings.TrimSpace(linha)
						if linha == "" {
							continue
						}

						if i := strings.Index(linha, "[]"); i != -1 && len(linha)-i > 3 {
							s.mensagens.WriteString(linha[i+3:])
						} else {
							s.mensagens.WriteString(linha)
						}

						s.mensagens.WriteString("\n")
					}
				}
			}(conn)
		}
	}(syslog)

	return syslog, nil
}

// Mensagens retorna as mensagens obtidas pelo servidor de logs.
func (s *ServidorLog) Mensagens() string {
	s.mensagens.Lock()
	defer s.mensagens.Unlock()

	return s.mensagens.Buffer.String()
}

// Limpar remove as mensagens obtidas pelo servidor de logs.
func (s *ServidorLog) Limpar() {
	s.mensagens.Lock()
	defer s.mensagens.Unlock()
	s.mensagens.Buffer.Reset()
}

// logMensagens armazena os logs garantindo a concorrência.
type logMensagens struct {
	sync.Mutex
	bytes.Buffer
}

// Write escreve a mensagem de log tratando problemas de concorrência.
func (l *logMensagens) Write(p []byte) (int, error) {
	l.Lock()
	defer l.Unlock()
	return l.Buffer.Write(p)
}

// WriteString escreve a mensagem de log tratando problemas de concorrência.
func (l *logMensagens) WriteString(s string) (int, error) {
	l.Lock()
	defer l.Unlock()
	return l.Buffer.WriteString(s)
}
