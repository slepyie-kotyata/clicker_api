package format

import (
	"strconv"
)

func IntToString(u int) string {
	s := strconv.Itoa(u)
	return s
}

func StringToUint(s string) uint {
	u, _ := strconv.Atoi(s)
    return uint(u)
}