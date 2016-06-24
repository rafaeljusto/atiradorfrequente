package bd

import "time"

type SQLogger struct {
	sqler
	Log Log
}

func NovoSQLogger(s sqler) *SQLogger {
	return &SQLogger{
		sqler: s,
	}
}

type Log struct {
	ID          uint64
	DataCriação time.Time
}
