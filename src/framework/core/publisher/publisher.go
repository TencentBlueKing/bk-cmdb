package publisher

import (
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/log"
)

type publisher struct {
	senders MapSender
}

// RegisterCustom register the custom sender , return the sender key.
func (cli *publisher) RegisterCustom(sender Sender) (string, error) {

	log.Infof("register a custom sender(%s)", sender.Description())
	key := common.UUID()
	cli.senders[SenderKey(key)] = sender
	return key, nil
}

// GetCustomSender get the custom sender by the sender key, return the custom sender.
func (cli *publisher) GetCustomSender(senderKey string) (Sender, error) {
	return nil, nil
}
