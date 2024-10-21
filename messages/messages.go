package messages

import "time"

type SetMessage struct {
	Key   string
	Value []byte
	TTL   time.Duration
}

type GetMessage struct {
	Key string
}
