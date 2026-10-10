package luadoc

import (
	"fmt"
	"strings"
)

const (
	String  = "string"
	Number  = "number"
	Boolean = "boolean"
	Table   = "table"
	Any     = "any"
)

func List(of string) string {
	return of + "[]"
}

func Map(keys string, values string) string {
	return Table + "<" + keys + ", " + values + ">"
}

func TableLiteral(keysAndValues ...string) string {
	if len(keysAndValues)%2 != 0 {
		panic("keysAndValues must be even")
	}

	var pairs []string
	for i := 0; i < len(keysAndValues); i += 2 {
		key := keysAndValues[i]
		value := keysAndValues[i+1]
		pairs = append(pairs, fmt.Sprintf("%s: %s", key, value))
	}

	return "{" + strings.Join(pairs, ", ") + "}"
}

func Enum(members ...string) string {
	quoted := make([]string, len(members))
	for i, member := range members {
		quoted[i] = fmt.Sprintf(`'%s'`, member)
	}

	return strings.Join(quoted, " | ")
}
