package config

import "time"

// Configuração define os valores configuráveis referentes a regras de negócio e
// políticas nos serviços.
type Configuração struct {
	Atirador struct {
		// PrazoConfirmação define o tempo máximo permitido para confirmar uma
		// frequência a partir do momento de sua criação.
		PrazoConfirmação time.Duration `yaml:"prazo confirmacao"`
	} `yaml:"atirador"`
}
