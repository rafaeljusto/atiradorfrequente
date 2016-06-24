package simulador

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

func (l Logger) Emerg(m ...interface{}) {
	l.SimulaEmerg(m...)
}

func (l Logger) Emergf(m string, a ...interface{}) {
	l.SimulaEmergf(m, a...)
}

func (l Logger) Alert(m ...interface{}) {
	l.SimulaAlert(m...)
}

func (l Logger) Alertf(m string, a ...interface{}) {
	l.SimulaAlertf(m, a...)
}

func (l Logger) Crit(m ...interface{}) {
	l.SimulaCrit(m...)
}

func (l Logger) Critf(m string, a ...interface{}) {
	l.SimulaCritf(m, a...)
}

func (l Logger) Error(e error) {
	l.SimulaError(e)
}

func (l Logger) Errorf(m string, a ...interface{}) {
	l.SimulaErrorf(m, a...)
}

func (l Logger) Warning(m ...interface{}) {
	l.SimulaWarning(m...)
}

func (l Logger) Warningf(m string, a ...interface{}) {
	l.SimulaWarningf(m, a...)
}

func (l Logger) Notice(m ...interface{}) {
	l.SimulaNotice(m...)
}

func (l Logger) Noticef(m string, a ...interface{}) {
	l.SimulaNoticef(m, a...)
}

func (l Logger) Info(m ...interface{}) {
	l.SimulaInfo(m...)
}

func (l Logger) Infof(m string, a ...interface{}) {
	l.SimulaInfof(m, a...)
}

func (l Logger) Debug(m ...interface{}) {
	l.SimulaDebug(m...)
}

func (l Logger) Debugf(m string, a ...interface{}) {
	l.SimulaDebugf(m, a...)
}

func (l Logger) SetCaller(n int) {
	l.SimulaSetCaller(n)
}
