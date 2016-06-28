package bd

import (
	"bytes"
	"strconv"
)

// Marcadores gera os argumentos de um comando SQL para PostgresSQL com a
// notação "$1,$...,$n".
func Marcadores(n int) string {
	return MarcadoresComInicio(1, n)
}

// MarcadoresComInicio gera os argumentos de um comando SQL para PostgresSQL com
// a notação "$inicio,$...,$n".
func MarcadoresComInicio(inicio, n int) string {
	if n <= 0 || inicio <= 0 {
		return ""
	}

	buf := bytes.NewBufferString("")
	for i := inicio; i <= n; i++ {
		buf.WriteString("$")
		buf.WriteString(strconv.Itoa(i))
		buf.WriteString(",")
	}

	// remove a última virgula
	if buf.Len() > 1 {
		buf.Truncate(buf.Len() - 1)
	}

	return buf.String()
}
