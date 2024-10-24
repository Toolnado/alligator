package commands

import (
	"errors"
	"fmt"
	"time"

	"github.com/Toolnado/alligator/messages"
)

type Command string

const (
	SET_COMMAND    Command = "SET"
	GET_COMMAND    Command = "GET"
	DELETE_COMMAND Command = "DELETE"
	HAS_COMMAND    Command = "HAS"
)

var ErrorInvalidLength = errors.New("invalid parts len")
var ErrorInvalidCommand = errors.New("invalid command")

// SET key value 1000ms
func ParseSetCommand(parts []string) (messages.FullCommandMessage, error) {
	if len(parts) != 3 {
		return nil, ErrorInvalidLength
	}
	ttl, err := time.ParseDuration(parts[2])
	if err != nil {
		return nil, fmt.Errorf("invalid ttl value: %s", err)
	}

	return &messages.Message{
		Key:   parts[0],
		Value: []byte(parts[1]),
		TTL:   ttl,
	}, nil
}

func ParseGetCommand(parts []string) (messages.CommandMessage, error) {
	return defaultHandler(parts)
}

func ParseDeleteCommand(parts []string) (messages.CommandMessage, error) {
	return defaultHandler(parts)
}

func ParseHasCommand(parts []string) (messages.CommandMessage, error) {
	return defaultHandler(parts)
}

func defaultHandler(parts []string) (messages.CommandMessage, error) {
	if len(parts) != 1 {
		return nil, ErrorInvalidLength
	}

	return &messages.Message{
		Key: parts[0],
	}, nil
}
