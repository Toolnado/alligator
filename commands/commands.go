package commands

import (
	"errors"
	"fmt"
	"time"
)

type CommandName string

type CMD struct {
	Name  CommandName
	Key   string
	Value []byte
	TTL   time.Duration
}

const (
	SET_COMMAND    CommandName = "SET"
	GET_COMMAND    CommandName = "GET"
	DELETE_COMMAND CommandName = "DELETE"
	HAS_COMMAND    CommandName = "HAS"
)

var ErrorInvalidProtocolFormat = errors.New("invalid protocol format")
var ErrorInvalidCommand = errors.New("invalid command")

func ParseCommand(parts []string) (CMD, error) {
	cmd := CMD{}
	if len(parts) == 0 {
		return cmd, ErrorInvalidProtocolFormat
	}

	cmd.Name = CommandName(parts[0])

	switch cmd.Name {
	case SET_COMMAND:
		if len(parts) < 4 {
			return cmd, ErrorInvalidProtocolFormat
		}
		cmd.Key = parts[1]
		cmd.Value = []byte(parts[2])
		ttl, err := time.ParseDuration(parts[3])
		if err != nil {
			return cmd, fmt.Errorf("invalid ttl format: %s", err)
		}
		cmd.TTL = ttl

	case GET_COMMAND:
		if len(parts) < 2 {
			return cmd, ErrorInvalidProtocolFormat
		}
		cmd.Key = parts[1]

	case HAS_COMMAND:
		if len(parts) < 2 {
			return cmd, ErrorInvalidProtocolFormat
		}
		cmd.Key = parts[1]

	case DELETE_COMMAND:
		if len(parts) < 2 {
			return cmd, ErrorInvalidProtocolFormat
		}
		cmd.Key = parts[1]
	default:
		return cmd, ErrorInvalidCommand
	}

	return cmd, nil
}
