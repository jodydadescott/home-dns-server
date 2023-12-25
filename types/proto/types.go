package proto

import (
	"strings"
)

// Proto is the protocol type. Currently only UDP and TCP are supported.
type Proto string

const (
	Empty   Proto = ""
	UDP     Proto = "udp"
	TCP     Proto = "tcp"
	Invalid Proto = "INVALID"
)

// NewFromString returns enum value from string
func NewFromString(input string) Proto {

	switch strings.ToLower(input) {

	case string(UDP):
		return UDP

	case string(TCP):
		return TCP

	case "":
		return Empty

	}

	return Invalid
}
