package simulador

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
