package base58id

import (
	"errors"
	"strconv"
	"strings"
)

type BrokerOption func(b *Broker) error

func WithCapacity(capacity int) BrokerOption {
	return func(b *Broker) error {
		b.capacity = capacity
		return nil
	}
}
func WithInstanceID(instanceID int) BrokerOption {
	return func(b *Broker) error {
		id := strconv.Itoa(instanceID)
		if strings.Contains(id, "0") {
			return errors.New("your instance ID contained a zero which is not allowed")
		}
		b.instanceID = id
		b.shortest = false
		return nil
	}
}
