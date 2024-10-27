package commands

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type CMD struct {
	raw []byte
}

type Command string

const (
	SetCommand    Command = "SET"
	GetCommand    Command = "GET"
	DeleteCommand Command = "DELETE"
	HasCommand    Command = "HAS"
)

var ErrorInvalidProtocolFormat = errors.New("invalid protocol format")
var ErrorInvalidCommand = errors.New("invalid command")

func New(raw []byte) CMD {
	return CMD{
		raw: raw,
	}
}

func (cmd CMD) Parse() (Message, error) {
	var (
		removeBytes = strings.TrimSuffix(string(cmd.raw), "\r\n")
		parts       = strings.Split(removeBytes, " ")
		name        = Command(parts[0])
		key         string
		value       []byte
		ttl         time.Duration
	)

	if len(parts) == 0 {
		return Message{}, ErrorInvalidProtocolFormat
	}

	switch name {
	case SetCommand:
		if len(parts) < 4 {
			return Message{}, ErrorInvalidProtocolFormat
		}
		key = parts[1]
		value = []byte(parts[2])
		latency, err := time.ParseDuration(parts[3])
		if err != nil {
			return Message{}, fmt.Errorf("invalid ttl format: %s", err)
		}
		ttl = latency

	case GetCommand:
		if len(parts) < 2 {
			return Message{}, ErrorInvalidProtocolFormat
		}
		key = parts[1]

	case HasCommand:
		if len(parts) < 2 {
			return Message{}, ErrorInvalidProtocolFormat
		}
		key = parts[1]

	case DeleteCommand:
		if len(parts) < 2 {
			return Message{}, ErrorInvalidProtocolFormat
		}
		key = parts[1]
	default:
		return Message{}, ErrorInvalidCommand
	}

	msg := newMessage(name, key, value, ttl)
	return msg, nil
}
