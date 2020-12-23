package tools

import (
	"strconv"
	"strings"
)

func ParseRemoteAddress(address string) (string, int) {
	parts := strings.Split(address, ":")

	if len(parts) != 2 {
		return "", 0
	}

	ip := parts[0]
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return ip, 0
	}

	return ip, port
}
