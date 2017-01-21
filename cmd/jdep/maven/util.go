package maven

import (
	"fmt"
	"strings"
)

func isProperty(value string) bool {
	return strings.HasPrefix(value, "${")
}

func varName(value string) string {
	if !isProperty(value) {
		return value
	}

	e := len(value) - 1
	return value[2:e]
}

func pprops(m map[string]string) {
	for k, v := range m {
		fmt.Printf("%-50s -> %#v\n", k, v)
	}
}
