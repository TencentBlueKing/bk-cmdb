package input

import (
	"configcenter/src/framework/core/types"
	"time"
)

type transaction struct {
	isTiming bool
	duration time.Duration
	savers   []types.Saver
}

// AddSaver add a saver
func (cli *transaction) AddSaver(saver types.Saver) {

}

// Execute execute this transaction
func (cli *transaction) Execute() error {

	return nil
}
