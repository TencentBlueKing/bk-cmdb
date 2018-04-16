package transaction

import (
	"configcenter/src/framework/common"
	"container/list"
	"context"
)

// New create a new Transaction instance
func New(ctx context.Context) Transaction {

	tran := &transaction{tasks: list.New()}

	tran.key = TranKey(common.UUID())

	return tran
}
