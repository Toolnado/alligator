package interfaces

import "time"

type Cache interface {
	Set(string, []byte, time.Duration) error
	Get(string) ([]byte, error)
	Delete(string) error
	Has(string) bool
}