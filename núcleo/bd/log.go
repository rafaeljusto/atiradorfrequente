package bd

import "time"

// SQLogger armazena além dos dados da transação do banco de dados, referências
// para rastrear todas as alterações do usuário nesta transação.
type SQLogger struct {
	sqler
	Log Log
}

// NovoSQLogger gera um novo SQLogger com os dados da transação.
func NovoSQLogger(s sqler) *SQLogger {
	return &SQLogger{
		sqler: s,
	}
}

// Log armazena os dados para rastreamento de todas as modificações do usuário.
type Log struct {
	ID          uint64
	DataCriação time.Time
}
