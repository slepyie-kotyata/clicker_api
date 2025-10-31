package utils

import (
	"encoding/json"
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

func ToJSON(d any) ([]byte, error) {
	d_bytes, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	return d_bytes, nil
}