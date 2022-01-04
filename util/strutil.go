package util

import (
	"fmt"
	"strings"
)

func Replace(s string, args []string) string {
	for i, v := range args {
		s = strings.ReplaceAll(s, fmt.Sprintf("{%d}", i), v)
	}
	return s
}
