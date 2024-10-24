package messages

import "time"

type CommandMessage interface {
	ItemKey() string
}

type FullCommandMessage interface {
	CommandMessage
	ItemValue() []byte
	ItemTTL() time.Duration
}

type Message struct {
	Key   string
	Value []byte
	TTL   time.Duration
}

func (m Message) ItemKey() string {
	return m.Key
}
func (m Message) ItemValue() []byte {
	return m.Value
}
func (m Message) ItemTTL() time.Duration {
	return m.TTL
}
