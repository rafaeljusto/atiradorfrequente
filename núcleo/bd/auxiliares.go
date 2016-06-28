package bd

import (
	"bytes"
	"strconv"
)

// MarcadoresPSQL gera os argumentos de um comando SQL para PostgresSQL com a
// notação "$1,$...,$n".
func MarcadoresPSQL(n int) string {
	return MarcadoresPSQLComInício(1, n)
}

// MarcadoresPSQLComInício gera os argumentos de um comando SQL para PostgresSQL
// com a notação "$início,$...,$n".
func MarcadoresPSQLComInício(início, n int) string {
	if n <= 0 || início <= 0 {
		return ""
	}

	buf := bytes.NewBufferString("")
	for i := início; i <= n; i++ {
		buf.WriteString("$")
		buf.WriteString(strconv.Itoa(i))
		buf.WriteString(",")
	}

	// remove a última vírgula
	if buf.Len() > 1 {
		buf.Truncate(buf.Len() - 1)
	}

	return buf.String()
}
