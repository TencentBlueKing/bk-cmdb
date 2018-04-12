package publisher

import (
	"configcenter/src/framework/core/types"
)

// SenderKey the sender key
type SenderKey string

// MapSender the sender map container
type MapSender map[SenderKey]Sender

// Sender is the interface that must be implemented by erver Sender.
type Sender interface {
	Description() string
	Put(event types.Event) error
	BatchPut(events []types.Event) error
}

// Publisher  is the interface that must be implemented by erver Publiser.
type Publisher interface {

	// RegisterCustom register the custom sender , return the sender key.
	RegisterCustom(sender Sender) (string, error)

	// GetCustomSender get the custom sender by the sender key, return the custom sender.
	GetCustomSender(senderKey string) (Sender, error)
}
