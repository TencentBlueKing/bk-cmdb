package sender

import (
	"configcenter/src/framework/core/types"
)

type CustomSender struct {
}

func (cli *CustomSender) Description() string {
	return "custom_sender"
}
func (cli *CustomSender) Put(event types.Event) error {
	return nil
}

func (cli *CustomSender) BatchPut(events []types.Event) error {
	return nil
}
