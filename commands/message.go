package commands

import "time"

type Message struct {
	cmd   Command
	key   string
	value []byte
	ttl   time.Duration
}

func (m Message) Command() Command {
	return m.cmd
}
func (m Message) Key() string {
	return m.key
}
func (m Message) Value() []byte {
	return m.value
}
func (m Message) TTL() time.Duration {
	return m.ttl
}

func newMessage(cmd Command, key string, value []byte, ttl time.Duration) Message {
	return Message{
		cmd:   cmd,
		key:   key,
		value: value,
		ttl:   ttl,
	}
}
