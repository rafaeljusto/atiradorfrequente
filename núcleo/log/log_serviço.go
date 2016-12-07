// log armazena estruturas de log específicas para o projeto.
package log

// Serviço descreve a interface de log disponível nas camadas de serviço do
// sistema. Este log será injetado em cada serviço para permitir mensagens de
// testes e de informações no núcleo do sistema.
type Serviço interface {
	Debug(m ...interface{})
	Debugf(m string, a ...interface{})
	Info(m ...interface{})
	Infof(s string, a ...interface{})
}
